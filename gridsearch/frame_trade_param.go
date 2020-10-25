package gridsearch

import (
	"log"
	"reflect"
	"yukimaterrace/andaman/trader"
	"yukimaterrace/andaman/util"
)

type frameTradeParam trader.FrameTradeParam

func (param frameTradeParam) partialParamSetSize() int {
	return reflect.ValueOf(param).NumField()
}

func (param frameTradeParam) partialParamSetValueSize(paramSetID partialParamSetID) int {
	value := reflect.ValueOf(param).Field(int(paramSetID)).Interface()

	switch value.(type) {
	case trader.FrameTradeParamSet0:
		return 2
	case trader.FrameTradeParamSet1:
		return 2
	case trader.FrameTradeParamSet2:
		return 3
	case trader.FrameTradeParamSet3:
		return 1
	case trader.FrameTradeParamSet4:
		return 3
	case trader.FrameTradeParamSet5:
		return 2
	case trader.FrameTradeParamSet6:
		return 2
	default:
		panicForValue()
	}

	return -1
}

func (param frameTradeParam) getPartialParamSet(paramSetID partialParamSetID, paramSetValueID partialParamSetValueID) partialParamSet {
	value := reflect.ValueOf(param).Field(int(paramSetID)).Interface()

	switch value.(type) {
	case trader.FrameTradeParamSet0:
		switch paramSetValueID {
		case 0:
			return trader.FrameTradeParamSet0{
				TradeDirectionLong: true,
			}
		case 1:
			return trader.FrameTradeParamSet0{
				TradeDirectionLong: false,
			}
		default:
			panicForValue()
		}

	case trader.FrameTradeParamSet1:
		switch paramSetValueID {
		case 0:
			return trader.FrameTradeParamSet1{
				SmallFrameLength: 30,
				LargeFrameLength: 60,
			}
		case 1:
			return trader.FrameTradeParamSet1{
				SmallFrameLength: 60,
				LargeFrameLength: 120,
			}
		default:
			panicForValue()
		}

	case trader.FrameTradeParamSet2:
		switch paramSetValueID {
		case 0:
			return trader.FrameTradeParamSet2{
				PipsGapForCreateOrder: 10,
			}
		case 1:
			return trader.FrameTradeParamSet2{
				PipsGapForCreateOrder: 20,
			}
		case 2:
			return trader.FrameTradeParamSet2{
				PipsGapForCreateOrder: 30,
			}
		default:
			panicForValue()
		}

	case trader.FrameTradeParamSet3:
		switch paramSetValueID {
		case 0:
			return trader.FrameTradeParamSet3{
				PipsForStopLoss: -100,
			}
		default:
			panicForValue()
		}

	case trader.FrameTradeParamSet4:
		switch paramSetValueID {
		case 0:
			return trader.FrameTradeParamSet4{
				PipsForAdditionalOrder: -5,
			}
		case 1:
			return trader.FrameTradeParamSet4{
				PipsForAdditionalOrder: -10,
			}
		case 2:
			return trader.FrameTradeParamSet4{
				PipsForAdditionalOrder: -20,
			}
		default:
			panicForValue()
		}

	case trader.FrameTradeParamSet5:
		switch paramSetValueID {
		case 0:
			return trader.FrameTradeParamSet5{
				TimeForProfit1: 40,
				TimeForProfit2: 80,
				TimeForProfit3: 120,
			}
		case 1:
			return trader.FrameTradeParamSet5{
				TimeForProfit1: 20,
				TimeForProfit2: 40,
				TimeForProfit3: 60,
			}
		default:
			panicForValue()
		}

	case trader.FrameTradeParamSet6:
		switch paramSetValueID {
		case 0:
			return trader.FrameTradeParamSet6{
				PipsForProfit1: 20,
				PipsForProfit2: 10,
				PipsForProfit3: 5,
			}
		case 1:
			return trader.FrameTradeParamSet6{
				PipsForProfit1: 30,
				PipsForProfit2: 15,
				PipsForProfit3: 5,
			}
		default:
			panicForValue()
		}

	default:
		panicForValue()
	}

	return nil
}

// CreateParam is a method to create param
func (param frameTradeParam) createParam(paramSets []partialParamSet) param {
	if len(paramSets) < param.partialParamSetSize() {
		panicForValue()
	}

	p := frameTradeParam{}

	for _, paramSet := range paramSets {
		switch v := paramSet.(type) {
		case trader.FrameTradeParamSet0:
			p.FrameTradeParamSet0 = v
		case trader.FrameTradeParamSet1:
			p.FrameTradeParamSet1 = v
		case trader.FrameTradeParamSet2:
			p.FrameTradeParamSet2 = v
		case trader.FrameTradeParamSet3:
			p.FrameTradeParamSet3 = v
		case trader.FrameTradeParamSet4:
			p.FrameTradeParamSet4 = v
		case trader.FrameTradeParamSet5:
			p.FrameTradeParamSet5 = v
		case trader.FrameTradeParamSet6:
			p.FrameTradeParamSet6 = v
		default:
			panicForValue()
		}
	}

	return p
}

func panicForValue() {
	panic("invalid value")
}

// FrameTradeParamsForGridSearch is a method to get frame trade params for grid search
func FrameTradeParamsForGridSearch() []*trader.FrameTradeParam {
	params := paramsForGridSearch(frameTradeParam{})

	frameTradeParams := make([]*trader.FrameTradeParam, len(params))
	for i, param := range params {
		p, ok := param.(frameTradeParam)
		if !ok {
			panic(util.ErrWrongType)
		}
		frameTradeParam := trader.FrameTradeParam(p)
		frameTradeParams[i] = &frameTradeParam
	}

	log.Printf("%d grids for search\n", len(frameTradeParams))

	return frameTradeParams
}
