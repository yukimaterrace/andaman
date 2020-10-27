package factory

import (
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/trader"
)

const simulationTradeSetName = "Frame Simulation"

// CreateSimulationTradeSet is a factory method to create simulation trade set
func CreateSimulationTradeSet() {
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

	var configurationParams []*model.TradeConfigurationParam
	for timezone := model.TokyoAM; timezone <= model.NewYorkPM; timezone++ {
		for tradePair := model.GbpUsd; tradePair <= model.EurGbp; tradePair++ {
			configurationForLong := model.TradeConfigurationParam{
				TradePair:      tradePair,
				Timezone:       timezone,
				AlgorithmParam: &frameTradeAlgorithmForLong,
			}

			configurationForShort := model.TradeConfigurationParam{
				TradePair:      tradePair,
				Timezone:       timezone,
				AlgorithmParam: &frameTradeAlgorithmForShort,
			}

			configurationParams = append(configurationParams, &configurationForLong, &configurationForShort)
		}
	}

	tradeSetParam := model.TradeSetParam{
		Name:                simulationTradeSetName,
		Type:                model.Simulation,
		ConfigurationParams: configurationParams,
	}

	if err := model.AddTradeSet(&tradeSetParam); err != nil {
		panic(err)
	}
}
