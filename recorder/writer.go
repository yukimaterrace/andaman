package recorder

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/trader"
	"yukimaterrace/andaman/util"
)

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

type (
	tradeSummary struct {
		realizedProfit float64
		tradeCount     int
	}

	keyTradePairTradableTimeZone struct {
		tradePair            broker.TradePair
		tradableTimeZoneName string
	}

	tradePairSummaryWriter struct {
		tradeSpecs      *trader.TradeSpecs
		recordDir       string
		tradeSummaryMap map[trader.KeyPartitionIDTradePair]*tradeSummary
		writerMap       map[keyTradePairTradableTimeZone]*csv.Writer
		files           []*os.File
	}
)

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
