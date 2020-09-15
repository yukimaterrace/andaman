package factory

import (
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/gridsearch"
)

// CreateGridSearchFlow is a factory method to create grid search flow
func CreateGridSearchFlow() *flow.Flow {
	paramsForGridSearch := gridsearch.ParamsForGridSearch(flow.FrameTradeParam{})

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

	tokyoAM := flow.CreateTokyoAMTimeZone()
	tokyoPM := flow.CreateTokyoPMTimeZone()
	londonAM := flow.CreateLondonAMTimeZone()
	londonPM := flow.CreateLondonPMTimeZone()
	newyorkAM := flow.CreateNewYorkAMTimeZone()
	newyorkPM := flow.CreateNewYorkPMTimeZone()

	return nil
}
