package factory

import (
	"yukimaterrace/andaman/gridsearch"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/service"
	"yukimaterrace/andaman/trader"
)

const (
	// SimulationTradeSetName is a name for simulation trade set
	SimulationTradeSetName = "frame_simulation"

	// GridSearchTradeSetName is a name for grid search trade set
	GridSearchTradeSetName = "grid_search"
)

// AddSimulationTradeSet is a method to add simulation trade set
func AddSimulationTradeSet() {
	paramSet0 := trader.FrameTradeParamSet0{
		TradeDirectionLong: true,
	}

	paramSet1 := trader.FrameTradeParamSet1{
		SmallFrameLength: 30,
		LargeFrameLength: 60,
	}

	paramSet2 := trader.FrameTradeParamSet2{
		PipsGapForCreateOrder: 10.0,
	}

	paramSet3 := trader.FrameTradeParamSet3{
		PipsForStopLoss: -100.0,
	}

	paramSet4 := trader.FrameTradeParamSet4{
		PipsForAdditionalOrder: -5.0,
	}

	paramSet5 := trader.FrameTradeParamSet5{
		TimeForProfit1: 40,
		TimeForProfit2: 60,
		TimeForProfit3: 80,
	}

	paramSet6 := trader.FrameTradeParamSet6{
		PipsForProfit1: 20.0,
		PipsForProfit2: 10.0,
		PipsForProfit3: 5.0,
	}

	longTradeParam := trader.FrameTradeParam{
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

	frameTradeAlgorithmForLong := model.TradeAlgorithmParam{
		Type:           model.Frame,
		Param:          &longTradeParam,
		TradeDirection: model.Long,
	}

	frameTradeAlgorithmForShort := model.TradeAlgorithmParam{
		Type:           model.Frame,
		Param:          &shortTradeParam,
		TradeDirection: model.Short,
	}

	timezoneIterator := model.TimezoneIterator{}
	tradePairIterator := model.TradePairIterator{}
	algorithmParams := []*model.TradeAlgorithmParam{
		&frameTradeAlgorithmForLong,
		&frameTradeAlgorithmForShort,
	}

	var configurationParams []*model.TradeConfigurationParam
	for timezoneIterator.Next() {
		timezone := timezoneIterator.Value()

		for tradePairIterator.Next() {
			tradePair := tradePairIterator.Value()

			for _, algorithmParam := range algorithmParams {
				configurationParam := model.TradeConfigurationParam{
					TradePair:      tradePair,
					Timezone:       timezone,
					AlgorithmParam: algorithmParam,
				}

				configurationParams = append(configurationParams, &configurationParam)
			}
		}
	}

	tradeSetParam := model.TradeSetParam{
		Name:                SimulationTradeSetName,
		Type:                model.Simulation,
		ConfigurationParams: configurationParams,
	}

	if err := service.AddTradeSet(&tradeSetParam); err != nil {
		panic(err)
	}
}

// AddGridSearchTradeSet is a method to add grid search trade set
func AddGridSearchTradeSet() {
	params := gridsearch.FrameTradeParamsForGridSearch()

	var algorithmParams []*model.TradeAlgorithmParam
	for _, param := range params {
		var tradeDirection model.TradeDirection
		if param.TradeDirectionLong {
			tradeDirection = model.Long
		} else {
			tradeDirection = model.Short
		}

		algorithmParam := model.TradeAlgorithmParam{
			Type:           model.Frame,
			Param:          param,
			TradeDirection: tradeDirection,
		}
		algorithmParams = append(algorithmParams, &algorithmParam)
	}

	timezoneIterator := model.TimezoneIterator{}
	tradePairIterator := model.TradePairIterator{}

	var configurationParams []*model.TradeConfigurationParam
	for timezoneIterator.Next() {
		timezone := timezoneIterator.Value()

		for tradePairIterator.Next() {
			tradePair := tradePairIterator.Value()

			for _, algorithmParam := range algorithmParams {
				configurationParam := model.TradeConfigurationParam{
					TradePair:      tradePair,
					Timezone:       timezone,
					AlgorithmParam: algorithmParam,
				}

				configurationParams = append(configurationParams, &configurationParam)
			}
		}
	}

	tradeSetParam := model.TradeSetParam{
		Name:                GridSearchTradeSetName,
		Type:                model.GridSearch,
		ConfigurationParams: configurationParams,
	}

	if err := service.AddTradeSet(&tradeSetParam); err != nil {
		panic(err)
	}
}
