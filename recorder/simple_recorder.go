package recorder

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/trader"
	"yukimaterrace/andaman/util"
)

type simpleRecorder struct {
	writer   writer
	orderMap map[trader.PartitionID]map[broker.OrderID]*completableOrder
}

func newSimpleRecorder(writer writer) *simpleRecorder {
	return &simpleRecorder{
		writer:   writer,
		orderMap: map[trader.PartitionID]map[broker.OrderID]*completableOrder{},
	}
}

func (recorder *simpleRecorder) Record(material flow.RecordMaterial) {
	partitionCombinedOrders, ok := material.(trader.PartitionCombinedOrders)
	if !ok {
		panic(util.ErrWrongType)
	}

	for partitionID, combinedOrders := range partitionCombinedOrders {
		if _, ok := recorder.orderMap[partitionID]; !ok {
			recorder.orderMap[partitionID] = map[broker.OrderID]*completableOrder{}
		}

		completableOrderMap := recorder.orderMap[partitionID]
		for _, createdOrder := range combinedOrders.CreatedOrders {
			if _, ok := completableOrderMap[createdOrder.OrderID()]; ok {
				panic("duplicate order id for created order detected")
			}

			completableOrderMap[createdOrder.OrderID()] = &completableOrder{
				createdOrder: createdOrder,
				closedOrder:  nil,
			}
		}

		for _, closedOrder := range combinedOrders.ClosedOrders {
			completableOrder, ok := completableOrderMap[closedOrder.OrderID()]
			if !ok {
				panic("no created order for the order id exist")
			}

			completableOrder.closedOrder = closedOrder
		}
	}
}

func (recorder *simpleRecorder) Write() {
	identifiedCompletableOrders := recorder.flush(true)
	recorder.writer.write(identifiedCompletableOrders)
}

func (recorder *simpleRecorder) Close() {
	identifiedCompletableOrders := recorder.flush(false)

	recorder.writer.write(identifiedCompletableOrders)
	recorder.writer.close()
}

func (recorder *simpleRecorder) flush(onlyCompleted bool) identifiedCompletableOrders {
	orders := identifiedCompletableOrders{}

	for partitionID, orderMap := range recorder.orderMap {
		var closedOrderIDs []broker.OrderID

		for orderID, order := range orderMap {
			if onlyCompleted && order.closedOrder == nil {
				continue
			}

			orders = append(orders, &identifiedCompletableOrder{
				partitionID: partitionID,
				tradePair:   order.createdOrder.TradePair(),
				order:       order,
			})

			if order.closedOrder != nil {
				closedOrderIDs = append(closedOrderIDs, orderID)
			}
		}

		for _, orderID := range closedOrderIDs {
			delete(orderMap, orderID)
		}
	}

	sort.Sort(orders)
	return orders
}

type (
	completableOrder struct {
		createdOrder broker.CreatedOrder
		closedOrder  broker.ClosedOrder
	}

	identifiedCompletableOrder struct {
		partitionID trader.PartitionID
		tradePair   broker.TradePair
		order       *completableOrder
	}
)

func (identifiedCompletableOrder *identifiedCompletableOrder) csvHeaders() []string {
	return []string{
		"orderID",
		"tradePair",
		"timeAtOpen",
		"priceAtOpen",
		"units",
		"isLong",
		"timeAtClose",
		"priceAtClose",
		"realizedPL",
	}
}

func (identifiedCompletableOrder *identifiedCompletableOrder) csvValues() []string {
	created := identifiedCompletableOrder.order.createdOrder
	closed := identifiedCompletableOrder.order.closedOrder

	csv := []string{
		strconv.FormatInt(int64(created.OrderID()), 10),
		string(created.TradePair()),
		strconv.FormatInt(created.TimeAtOpen(), 10),
		strconv.FormatFloat(created.PriceAtOpen(), 'f', 6, 64),
		strconv.FormatFloat(created.Units(), 'f', 8, 64),
		strconv.FormatBool(created.IsLong()),
	}

	if closed == nil {
		return append(csv, "not closed", "not closed", "0")
	}

	return append(csv,
		strconv.FormatInt(closed.TimeAtClose(), 10),
		strconv.FormatFloat(closed.PriceAtClose(), 'f', 6, 64),
		strconv.FormatFloat(closed.RealizedProfit(), 'f', 6, 64),
	)
}

type identifiedCompletableOrders []*identifiedCompletableOrder

func (orders identifiedCompletableOrders) Len() int {
	return len(orders)
}

func (orders identifiedCompletableOrders) Less(i, j int) bool {
	return orders[i].order.createdOrder.OrderID() < orders[j].order.createdOrder.OrderID()
}

func (orders identifiedCompletableOrders) Swap(i, j int) {
	order := orders[i]
	orders[i] = orders[j]
	orders[j] = order
}

type writer interface {
	write(orders identifiedCompletableOrders)
	close()
}

type simpleWriter struct {
	recordDir         string
	tradableTimeZones trader.TradableTimeZones
	csvWriterMap      map[trader.KeyPartitionIDTradePair]*csv.Writer
	files             []*os.File
}

func newSimpleWriter(tradableTimeZones trader.TradableTimeZones) *simpleWriter {
	return &simpleWriter{
		recordDir:         util.GetEnv("RECORD_DIR"),
		tradableTimeZones: tradableTimeZones,
		csvWriterMap:      map[trader.KeyPartitionIDTradePair]*csv.Writer{},
		files:             []*os.File{},
	}
}

func (writer *simpleWriter) write(orders identifiedCompletableOrders) {
	if len(orders) == 0 {
		return
	}

	for _, order := range orders {
		key := trader.KeyPartitionIDTradePair{PartitionID: order.partitionID, TradePair: order.tradePair}
		if _, ok := writer.csvWriterMap[key]; !ok {
			tradableTimeZone, ok := writer.tradableTimeZones[key.PartitionID]
			if !ok {
				panic("no tradable time zone specified")
			}

			path := fmt.Sprintf("%s/%s_%s_%d.csv", writer.recordDir, string(order.tradePair), tradableTimeZone.Name, order.partitionID)

			file, err := os.Create(path)
			if err != nil {
				panic(err)
			}

			writer.files = append(writer.files, file)

			csvWriter := csv.NewWriter(file)
			if err := csvWriter.Write(order.csvHeaders()); err != nil {
				panic(err)
			}

			writer.csvWriterMap[key] = csvWriter
		}

		csvWriter := writer.csvWriterMap[key]
		if err := csvWriter.Write(order.csvValues()); err != nil {
			log.Println(err.Error())
		}
	}
}

func (writer *simpleWriter) close() {
	for _, csvWriter := range writer.csvWriterMap {
		csvWriter.Flush()
	}

	for _, file := range writer.files {
		if err := file.Close(); err != nil {
			log.Println(err.Error())
		}
	}
}

// SimpleRecorderFactory is a factory for simple trader using simple recorder
type SimpleRecorderFactory struct {
	builder *trader.SimpleTraderBuilder
}

// NewSimpleRecorderFactory is a constructor for simple recorder factory
func NewSimpleRecorderFactory(builder *trader.SimpleTraderBuilder) *SimpleRecorderFactory {
	return &SimpleRecorderFactory{builder}
}

// Create is a factory method to create recorder
func (factory *SimpleRecorderFactory) Create() flow.Recorder {
	tradableTimeZones := factory.builder.BuildTradableTimeZones()
	return newSimpleRecorder(newSimpleWriter(tradableTimeZones))
}

type tradeSummary struct {
	realizedProfit float64
	tradeCount     int
}

type keyTradePairTradableTimeZone struct {
	tradePair            broker.TradePair
	tradableTimeZoneName string
}

type tradePairSummaryWriter struct {
	tradeSpecs      *trader.TradeSpecs
	recordDir       string
	tradeSummaryMap map[trader.KeyPartitionIDTradePair]*tradeSummary
	writerMap       map[keyTradePairTradableTimeZone]*csv.Writer
	files           []*os.File
}

func newTradePairSummaryWriter(tradeSpecs *trader.TradeSpecs) *tradePairSummaryWriter {
	return &tradePairSummaryWriter{
		tradeSpecs:      tradeSpecs,
		recordDir:       util.GetEnv("RECORD_DIR"),
		tradeSummaryMap: map[trader.KeyPartitionIDTradePair]*tradeSummary{},
		writerMap:       map[keyTradePairTradableTimeZone]*csv.Writer{},
		files:           []*os.File{},
	}
}

func (writer *tradePairSummaryWriter) write(orders identifiedCompletableOrders) {
	for _, order := range orders {
		closedOrder := order.order.closedOrder

		if closedOrder != nil {
			key := trader.KeyPartitionIDTradePair{PartitionID: order.partitionID, TradePair: order.tradePair}
			if _, ok := writer.tradeSummaryMap[key]; !ok {
				writer.tradeSummaryMap[key] = &tradeSummary{}
			}

			tradeSummary := writer.tradeSummaryMap[key]
			tradeSummary.realizedProfit += closedOrder.RealizedProfit()
			tradeSummary.tradeCount++
		}
	}
}

func (writer *tradePairSummaryWriter) close() {
	for key, tradeSummary := range writer.tradeSummaryMap {
		tradableTimeZone, ok := writer.tradeSpecs.TimeZones[key.PartitionID]
		if !ok {
			panic("no tradable time zone specified")
		}

		paramLoader, ok := writer.tradeSpecs.ParamLoaders[key]
		if !ok {
			panic("no param loader specified")
		}

		_key := keyTradePairTradableTimeZone{key.TradePair, tradableTimeZone.Name}
		if _, ok := writer.writerMap[_key]; !ok {
			path := fmt.Sprintf("%s/%s_%s.csv", writer.recordDir, string(key.TradePair), tradableTimeZone.Name)
			file, err := os.Create(path)
			if err != nil {
				panic(err)
			}

			writer.files = append(writer.files, file)

			csvWriter := csv.NewWriter(file)
			csvWriter.Write(paramLoader.ParamCsvHeader())

			writer.writerMap[_key] = csvWriter
		}

		realizedProfit := strconv.FormatFloat(tradeSummary.realizedProfit, 'f', 6, 64)
		tradeCount := strconv.FormatInt(int64(tradeSummary.tradeCount), 10)

		writer.writerMap[_key].Write(append(paramLoader.ParamCsvValue(), realizedProfit, tradeCount))
	}

	for _, csvWriter := range writer.writerMap {
		csvWriter.Flush()
	}

	for _, file := range writer.files {
		file.Close()
	}
}

// SimpleTradePairSummaryRecorderFactory is a factory for simple trader using trade summary recorder
type SimpleTradePairSummaryRecorderFactory struct {
	builder *trader.SimpleTraderBuilder
}

// NewSimpleTradePairSummaryRecorderFactory is a constructor for simple trade pair summary recorder
func NewSimpleTradePairSummaryRecorderFactory(builder *trader.SimpleTraderBuilder) *SimpleTradePairSummaryRecorderFactory {
	return &SimpleTradePairSummaryRecorderFactory{builder}
}

func (factory *SimpleTradePairSummaryRecorderFactory) Create() flow.Recorder {
	tradeSpecs := factory.builder.BuildTradeSpecs()
	return newSimpleRecorder(newTradePairSummaryWriter(tradeSpecs))
}
