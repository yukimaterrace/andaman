package trade

import (
	"log"
	"yukimaterrace/andaman/indicate"
	"yukimaterrace/andaman/market"
	"yukimaterrace/andaman/record"
)

// Trader is an interface for trader
type Trader interface {
	Start()
}

type routine struct {
	instrument market.Instrument
	market.Pricer
	market.Orderer
	record.Recorder
	tradeStrategy
	indicators  []indicate.Indicator
	granularity market.Granularity
	priceCount  int
}

func newRoutine(instrument market.Instrument, market market.Market, recorder record.Recorder, tradeStrategy tradeStrategy) *routine {
	return &routine{
		instrument:    instrument,
		Pricer:        market,
		Orderer:       market,
		Recorder:      recorder,
		tradeStrategy: tradeStrategy,
	}
}

type tradeStrategy interface {
	requireIndicators() []indicate.Indicator
	requireGranularity() market.Granularity
	requirePriceCount() int
	processOrder(orderer market.Orderer, tradePrice *market.TradePrice, indicatorValues []*indicate.Value) (*market.MadeOrder, []*market.ClosedOrder)
}

func (routine *routine) Start() {
	routine.indicators = routine.requireIndicators()
	routine.granularity = routine.requireGranularity()
	routine.priceCount = routine.requirePriceCount()

	go routine.run()
}

func (routine *routine) run() {
	for {
		pricesChan := routine.Prices(routine.instrument, routine.granularity, routine.priceCount)

		pricesStat := <-pricesChan
		if pricesStat.Err != nil {
			log.Println(pricesStat.Err)
			continue
		}
		prices := pricesStat.PriceSequence

		tradePriceChan := routine.TradePrice(routine.instrument)

		values := make([]*indicate.Value, 0)
		for _, indicator := range routine.indicators {
			values = append(values, indicate.Indicate(prices))
		}

		tradePriceStat := <-tradePriceChan
		if tradePriceStat.Err != nil {
			log.Println(tradePriceStat.Err)
			continue
		}
		tradePrice := tradePriceStat.TradePrice

		madeOrder, closedOrders := routine.processOrder(routine.Orderer, tradePrice, values)
		material := &record.Material{
			Instrument:  routine.instrument,
			Prices:      prices,
			Indicators:  values,
			MadeOrder:   madeOrder,
			ClosedOrder: closedOrders,
		}

		routine.Record(material)
	}
}
