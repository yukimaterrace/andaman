package factory

import (
	"time"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/pricer"
	"yukimaterrace/andaman/recorder"
	"yukimaterrace/andaman/trader"
)

// CreateSimulationFlow is a factory method to create simulation app
func CreateSimulationFlow(tradeSetName string, writeInterval time.Duration) *flow.Flow {
	pricerTradePairs := []model.TradePair{
		model.GbpUsd,
		model.EurUsd,
		model.AudUsd,
		model.UsdJpy,
		model.AudJpy,
		model.GbpAud,
		model.EurAud,
		model.EurGbp,
		model.GbpJpy,
		model.EurJpy,
	}

	tradeBuilder := trader.NewBuilder().
		TradeSet(tradeSetName).
		Parallel(1)

	start := time.Date(2020, time.July, 1, 0, 0, 0, 0, time.Local)
	end := time.Date(2020, time.August, 31, 23, 59, 59, 0, time.Local)

	flow := flow.NewFlowBuilder().
		Broker(broker.NewSimpleSimulationBroker()).
		OrdererFactory(broker.NewSimpleSimulationOrdererFactory()).
		PricerTradePairs(pricerTradePairs).
		InitialTradeMode(flow.Trade).
		PricerFactory(pricer.NewOandaSimulationPricerFactory(start, end)).
		TraderFactory(trader.NewFactory(tradeBuilder)).
		RecorderFactory(recorder.NewFactory(tradeBuilder)).
		WriteInterval(writeInterval).
		Build()

	return flow
}
