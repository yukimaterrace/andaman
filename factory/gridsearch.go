package factory

import (
	"time"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/gridsearch"
	"yukimaterrace/andaman/pricer"
	"yukimaterrace/andaman/recorder"
	"yukimaterrace/andaman/trader"
)

// CreateGridSearchFlow is a factory method to create grid search flow
func CreateGridSearchFlow() *flow.Flow {
	paramsForGridSearch := gridsearch.FrameTradeParamsForGridSearch()

	timezones := []*trader.TradableTimeZone{
		trader.CreateTokyoAMTimeZone(),
		trader.CreateTokyoPMTimeZone(),
		trader.CreateLondonAMTimeZone(),
		trader.CreateLondonPMTimeZone(),
		trader.CreateNewYorkAMTimeZone(),
		trader.CreateNewYorkPMTimeZone(),
	}

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

	tradeBuilder := trader.NewSimpleTraderBuilder()
	for i, timezone := range timezones {
		for j, param := range paramsForGridSearch {
			algorithm := trader.NewFrameTradeAlgorithm(param)
			paritionID := trader.PartitionID(j + i*len(paramsForGridSearch))

			tradeBuilder.TradableTimeZone(paritionID, timezone)
			for _, tradePair := range pricerTradePairs {
				tradeBuilder.Trade(paritionID, tradePair, algorithm)
			}
		}
	}

	tradeBuilder.Parallel(2)

	start := time.Date(2020, time.June, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(2020, time.July, 31, 23, 59, 59, 0, time.Local)

	flow := flow.NewFlowBuilder().
		Broker(broker.NewSimpleSimulationBroker()).
		OrdererFactory(broker.NewSimpleSimulationOrdererFactory()).
		PricerTradePairs(pricerTradePairs).
		InitialTradeMode(flow.Trade).
		PricerFactory(pricer.NewOandaSimulationPricerFactory(start, end)).
		TraderFactory(trader.NewSimpleTraderFactory(tradeBuilder)).
		RecorderFactory(recorder.NewSimpleTradePairSummaryRecorderFactory(tradeBuilder)).
		Build()

	return flow
}
