package trade

import (
	"log"
	"yukimaterrace/andaman/indicator"
	"yukimaterrace/andaman/market"
	"yukimaterrace/andaman/recorder"
)

// Trader is an interface for trader
type Trader interface {
	Start()
}

type routine struct {
	instrument market.Instrument
	market.Pricer
	market.Orderer
	recorder.Recorder
	tradeStrategy
	indicators  []indicator.Indicator
	granularity market.Granularity
	priceCount  int
}

type tradeStrategy interface {
	requireInidicators() []indicator.Indicator
	requireGranularity() market.Granularity
	requirePriceCount() int
	processOrder(orderer market.Orderer, tradePrice *market.TradePrice, indicatorValues []*indicator.Value) (*market.MadeOrder, []*market.ClosedOrder)
}

func (routine *routine) Start() {
	routine.indicators = routine.requireInidicators()
	routine.granularity = routine.requireGranularity()
	routine.priceCount = routine.requirePriceCount()

	go routine.run()
}

func (routine *routine) run() {
	for {
		pricesChan := routine.Prices(routine.instrument, routine.granularity, routine.priceCount)

		pricesStat := <-pricesChan
		if pricesStat.Err != nil {
			log.Fatal(pricesStat.Err)
			continue
		}
		prices := pricesStat.PriceSequence

		tradePriceChan := routine.TradePrice(routine.instrument)

		values := make([]*indicator.Value, 0)
		for _, indicator := range routine.indicators {
			values = append(values, indicator.Indicate(prices))
		}

		tradePriceStat := <-tradePriceChan
		if tradePriceStat.Err != nil {
			log.Fatal(tradePriceStat.Err)
			continue
		}
		tradePrice := tradePriceStat.TradePrice

		madeOrder, closedOrders := routine.processOrder(routine.Orderer, tradePrice, values)
		material := &recorder.Material{
			Instrument:  routine.instrument,
			Prices:      prices,
			Indicators:  values,
			MadeOrder:   madeOrder,
			ClosedOrder: closedOrders,
		}

		routine.Record(material)
	}
}
