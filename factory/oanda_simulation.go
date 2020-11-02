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
func CreateSimulationFlow(tradeSetName string, writeInterval time.Duration, start int, end int) *flow.Flow {
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
		TradeRunType(model.OandaSimulation).
		Parallel(1)

	_start := time.Unix(int64(start), 0)
	_end := time.Unix(int64(end), 0)

	flow := flow.NewFlowBuilder().
		Broker(broker.NewSimpleSimulationBroker()).
		OrdererFactory(broker.NewSimpleSimulationOrdererFactory()).
		PricerTradePairs(pricerTradePairs).
		InitialTradeMode(flow.Trade).
		PricerFactory(pricer.NewOandaSimulationPricerFactory(_start, _end)).
		TraderFactory(trader.NewFactory(tradeBuilder)).
		RecorderFactory(recorder.NewFactory(tradeBuilder)).
		WriteInterval(writeInterval).
		Build()

	return flow
}
