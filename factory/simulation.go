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

	paramSet0 := flow.FrameTradeParamSet0{
		TradeDirectionLong: true,
	}

	paramSet1 := flow.FrameTradeParamSet1{
		SmallFrameLength: 30,
		LargeFrameLength: 60,
	}

	paramSet2 := flow.FrameTradeParamSet2{
		PipsGapForCreateOrder: 10.0,
	}

	paramSet3 := flow.FrameTradeParamSet3{
		PipsForStopLoss: -100.0,
	}

	paramSet4 := flow.FrameTradeParamSet4{
		PipsForAdditionalOrder: -5.0,
	}

	paramSet5 := flow.FrameTradeParamSet5{
		TimeForProfit1: 40,
		TimeForProfit2: 60,
		TimeForProfit3: 80,
	}

	paramSet6 := flow.FrameTradeParamSet6{
		PipsForProfit1: 20.0,
		PipsForProfit2: 10.0,
		PipsForProfit3: 5.0,
	}

	longTradeParam := flow.FrameTradeParam{
		FrameTradeParamSet0: paramSet0,
		FrameTradeParamSet1: paramSet1,
		FrameTradeParamSet2: paramSet2,
		FrameTradeParamSet3: paramSet3,
		FrameTradeParamSet4: paramSet4,
		FrameTradeParamSet5: paramSet5,
		FrameTradeParamSet6: paramSet6,
	}

	shortTradeParam := longTradeParam
	shortTradeParam.TradeDirectionLong = false

	longTradeAlgorithm := flow.NewFrameTradeAlgorithm(&longTradeParam)
	shortTradeAlgorithm := flow.NewFrameTradeAlgorithm(&shortTradeParam)

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
