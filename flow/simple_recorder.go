package flow

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/util"
)

type simpleRecorder struct {
	writer   writer
	orderMap map[broker.AccountID]map[broker.OrderID]*completableOrder
}

func newSimpleRecorder(writer writer) *simpleRecorder {
	return &simpleRecorder{
		writer:   writer,
		orderMap: map[broker.AccountID]map[broker.OrderID]*completableOrder{},
	}
}

func (recorder *simpleRecorder) record(material recordMaterial) {
	accountCombinedOrders, ok := material.(accountCombinedOrders)
	if !ok {
		panic(util.ErrWrongType)
	}

	for accountID, combinedOrders := range accountCombinedOrders {
		if _, ok := recorder.orderMap[accountID]; !ok {
			recorder.orderMap[accountID] = map[broker.OrderID]*completableOrder{}
		}

		completableOrderMap := recorder.orderMap[accountID]
		for _, createdOrder := range combinedOrders.createdOrders {
			if _, ok := completableOrderMap[createdOrder.OrderID()]; ok {
				panic("duplicate order id for created order detected")
			}

			completableOrderMap[createdOrder.OrderID()] = &completableOrder{
				createdOrder: createdOrder,
				closedOrder:  nil,
			}
		}

		for _, closedOrder := range combinedOrders.closedOrders {
			completableOrder, ok := completableOrderMap[closedOrder.OrderID()]
			if !ok {
				panic("no created order for the order id exist")
			}

			completableOrder.closedOrder = closedOrder
		}
	}
}

func (recorder *simpleRecorder) write() {
	identifiedCompletableOrders := recorder.flush(true)
	recorder.writer.write(identifiedCompletableOrders)
}

func (recorder *simpleRecorder) close() {
	identifiedCompletableOrders := recorder.flush(false)

	recorder.writer.write(identifiedCompletableOrders)
	recorder.writer.close()
}

func (recorder *simpleRecorder) flush(onlyCompleted bool) identifiedCompletableOrders {
	identifiedCompletableOrders := []*identifiedCompletableOrder{}

	for accountID, orderMap := range recorder.orderMap {
		closedOrderIDs := []broker.OrderID{}

		for orderID, order := range orderMap {
			if onlyCompleted && order.closedOrder == nil {
				continue
			}

			identifiedCompletableOrders = append(identifiedCompletableOrders, &identifiedCompletableOrder{
				accountID: accountID,
				tradePair: order.createdOrder.TradePair(),
				order:     order,
			})

			if order.closedOrder != nil {
				closedOrderIDs = append(closedOrderIDs, orderID)
			}
		}

		for _, orderID := range closedOrderIDs {
			delete(orderMap, orderID)
		}
	}

	return identifiedCompletableOrders
}

type completableOrder struct {
	createdOrder broker.CreatedOrder
	closedOrder  broker.ClosedOrder
}

type identifiedCompletableOrder struct {
	accountID broker.AccountID
	tradePair broker.TradePair
	order     *completableOrder
}

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
		strconv.FormatInt(int64(created.TimeAtOpen()), 10),
		strconv.FormatFloat(created.PriceAtOpen(), 'f', 6, 64),
		strconv.FormatFloat(created.Units(), 'f', 8, 64),
		strconv.FormatBool(created.IsLong()),
	}

	if closed != nil {
		return append(csv, "not closed", "not closed", "0")
	}

	return append(csv,
		strconv.FormatInt(int64(closed.TimeAtClose()), 10),
		strconv.FormatFloat(closed.PriceAtClose(), 'f', 6, 64),
		strconv.FormatFloat(closed.RealizedProfit(), 'f', 6, 64),
	)
}

type identifiedCompletableOrders []*identifiedCompletableOrder

type writer interface {
	write(orders identifiedCompletableOrders)
	close()
}

type simpleWriter struct {
	recordDir    string
	csvWriterMap map[broker.TradePair]map[broker.AccountID]*csv.Writer
	files        []*os.File
}

func newSimpleWriter() *simpleWriter {
	return &simpleWriter{
		recordDir:    util.GetEnv("RECORD_DIR"),
		csvWriterMap: map[broker.TradePair]map[broker.AccountID]*csv.Writer{},
		files:        []*os.File{},
	}
}

func (writer *simpleWriter) write(orders identifiedCompletableOrders) {
	if len(orders) == 0 {
		return
	}

	for _, order := range orders {
		if _, ok := writer.csvWriterMap[order.tradePair]; !ok {
			writer.csvWriterMap[order.tradePair] = map[broker.AccountID]*csv.Writer{}
		}

		accountCsvWriterMap := writer.csvWriterMap[order.tradePair]
		if _, ok := accountCsvWriterMap[order.accountID]; !ok {
			path := fmt.Sprintf("%s/%s_%s", writer.recordDir, string(order.tradePair), string(order.accountID))

			file, err := os.Create(path)
			if err != nil {
				panic(err)
			}

			writer.files = append(writer.files, file)

			csvWriter := csv.NewWriter(file)
			if err := csvWriter.Write(order.csvHeaders()); err != nil {
				panic(err)
			}

			accountCsvWriterMap[order.accountID] = csvWriter
		}

		csvWriter := accountCsvWriterMap[order.accountID]
		if err := csvWriter.Write(order.csvValues()); err != nil {
			log.Println(err.Error())
		}
	}
}

func (writer *simpleWriter) close() {
	for _, accountCsvWriterMap := range writer.csvWriterMap {
		for _, csvWrite := range accountCsvWriterMap {
			csvWrite.Flush()
		}
	}

	for _, file := range writer.files {
		if err := file.Close(); err != nil {
			log.Println(err.Error())
		}
	}
}

// SimpleRecorderFactory is a factory for simple recorder
type SimpleRecorderFactory struct{}

func (factory *SimpleRecorderFactory) create() recorder {
	return newSimpleRecorder(newSimpleWriter())
}
