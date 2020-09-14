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
		PipsForAdditionalOrder: -5.0,
		PipsForStopLoss:        -100.0,
		TimeForProfit1:         40,
		TimeForProfit2:         60,
		TimeForProfit3:         80,
		PipsForProfit1:         20.0,
		PipsForProfit2:         10.0,
		PipsForProfit3:         5.0,
	}

	shortTradeParam := &flow.FrameTradeParam{
		TradeDirectionLong:     false,
		SmallFrameLength:       30,
		LargeFrameLength:       60,
		PipsGapForCreateOrder:  10.0,
		PipsForAdditionalOrder: -5.0,
		PipsForStopLoss:        -100.0,
		TimeForProfit1:         40,
		TimeForProfit2:         60,
		TimeForProfit3:         80,
		PipsForProfit1:         20.0,
		PipsForProfit2:         10.0,
		PipsForProfit3:         5.0,
	}

	longTradeAlgorithm := flow.NewFrameTradeAlgorithm(longTradeParam)
	shortTradeAlgorithm := flow.NewFrameTradeAlgorithm(shortTradeParam)

	tokyoAM := flow.CreateTokyoAMTimeZone()
	tokyoPM := flow.CreateTokyoPMTimeZone()
	londonAM := flow.CreateLondonAMTimeZone()
	londonPM := flow.CreateLondonPMTimeZone()
	newyorkAM := flow.CreateNewYorkAMTimeZone()
	newyorkPM := flow.CreateNewYorkPMTimeZone()

	tradeBuilder := flow.NewSimpleTraderBuilder().
		TradableTimeZone(0, tokyoAM).
		TradableTimeZone(1, tokyoAM).
		TradableTimeZone(2, tokyoPM).
		TradableTimeZone(3, tokyoPM).
		TradableTimeZone(4, londonAM).
		TradableTimeZone(5, londonAM).
		TradableTimeZone(6, londonPM).
		TradableTimeZone(7, londonPM).
		TradableTimeZone(8, newyorkAM).
		TradableTimeZone(9, newyorkAM).
		TradableTimeZone(10, newyorkPM).
		TradableTimeZone(11, newyorkPM)

	for _, tradePair := range pricerTradePairs {
		for i := 0; i < 12; i++ {
			if i%2 == 0 {
				tradeBuilder.Trade(flow.PartitionID(i), tradePair, longTradeAlgorithm)
			} else {
				tradeBuilder.Trade(flow.PartitionID(i), tradePair, shortTradeAlgorithm)
			}
		}
	}

	tradeBuilder.Parallel(1)

	start := time.Date(2020, time.June, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(2020, time.July, 23, 23, 59, 59, 0, time.Local)

	flow := flow.NewFlowBuilder().
		Broker(broker.NewSimpleSimulationBroker()).
		OrdererFactory(broker.NewSimpleSimulationOrdererFactory()).
		PricerTradePairs(pricerTradePairs).
		InitialTradeMode(flow.Trade).
		PricerFactory(flow.NewOandaSimulationPricerFactory(start, end)).
		TraderFactory(flow.NewSimpleTraderFactory(tradeBuilder)).
		RecorderFactory(flow.NewSimpleTradePairSummaryRecorderFactory(tradeBuilder)).
		Build()

	return flow
}
