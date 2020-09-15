package factory

import (
	"time"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/gridsearch"
)

// CreateGridSearchFlow is a factory method to create grid search flow
func CreateGridSearchFlow() *flow.Flow {
	paramsForGridSearch := gridsearch.FrameTradeParamsForGridSearch()

	timezones := []*flow.TradableTimeZone{
		flow.CreateTokyoAMTimeZone(),
		flow.CreateTokyoPMTimeZone(),
		flow.CreateLondonAMTimeZone(),
		flow.CreateLondonPMTimeZone(),
		flow.CreateNewYorkAMTimeZone(),
		flow.CreateNewYorkPMTimeZone(),
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

	tradeBuilder := flow.NewSimpleTraderBuilder()
	for i, timezone := range timezones {
		for j, param := range paramsForGridSearch {
			algorithm := flow.NewFrameTradeAlgorithm(param)
			paritionID := flow.PartitionID(j + i*len(paramsForGridSearch))

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
		PricerFactory(flow.NewOandaSimulationPricerFactory(start, end)).
		TraderFactory(flow.NewSimpleTraderFactory(tradeBuilder)).
		RecorderFactory(flow.NewSimpleTradePairSummaryRecorderFactory(tradeBuilder)).
		Build()

	return flow
}
