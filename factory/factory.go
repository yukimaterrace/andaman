package factory

import (
	"time"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
)

// CreateSimulationApp is a factory method to create simulation app
func CreateSimulationApp() *flow.Flow {
	pricerTradePairs := []broker.TradePair{
		broker.GbpUsd,
		broker.EurUsd,
		broker.AudUsd,
		broker.UsdJpy,
		broker.AudJpy,
		broker.GbpAud,
		broker.EurAud,
		broker.EurGbp,
		broker.GbpJpy,
		broker.EurJpy,
	}

	longTradeParam := &flow.FrameTradeParam{
		TradeDirectionLong:     true,
		SmallFrameLength:       30,
		LargeFrameLength:       60,
		PipsGapForCreateOrder:  10.0,
		PipsForAdditionalOrder: 5.0,
		PipsForStopLoss:        100.0,
		TimeForProfit1:         40,
		TimeForProfit2:         60,
		TimeForProfit3:         80,
		PipsForProfit1:         20.0,
		PipsForProfit2:         10.0,
		PipsForProfit3:         5.0,
	}

	// shortTradeParam := &(*longTradeParam)
	// shortTradeParam.TradeDirectionLong = false

	longTradeAlgorithm := flow.NewFrameTradeAlgorithm(longTradeParam)
	// shortTradeAlhorithm := flow.NewFrameTradeAlgorithm(shortTradeParam)

	tradeBuilder := flow.NewSimpleTraderBuilder().
		TradableTimeZone(0, flow.CreateWeekdayTradableTimeZone())

	for _, tradePair := range pricerTradePairs {
		tradeBuilder = tradeBuilder.Trade(0, tradePair, longTradeAlgorithm)
	}

	tradeBuilder.Parallel(1)

	start := time.Date(2020, time.July, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(2020, time.July, 23, 23, 59, 59, 0, time.Local)

	flow := flow.NewFlowBuilder().
		Broker(broker.NewSimpleSimulationBroker()).
		OrdererFactory(broker.NewSimpleSimulationOrdererFactory()).
		PricerTradePairs(pricerTradePairs).
		InitialTradeMode(flow.Trade).
		PricerFactory(flow.NewOandaSimulationPricerFactory(start, end)).
		TraderFactory(flow.NewSimpleTraderFactory(tradeBuilder)).
		RecorderFactory(flow.NewSimpleRecorderFactory(tradeBuilder)).
		Build()

	return flow
}
