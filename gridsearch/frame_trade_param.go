package gridsearch

import (
	"log"
	"reflect"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/util"
)

type frameTradeParam flow.FrameTradeParam

func (param frameTradeParam) partialParamSetSize() int {
	return reflect.ValueOf(param).NumField()
}

func (param frameTradeParam) partialParamSetValueSize(paramSetID partialParamSetID) int {
	value := reflect.ValueOf(param).Field(int(paramSetID)).Interface()

	switch value.(type) {
	case flow.FrameTradeParamSet0:
		return 2
	case flow.FrameTradeParamSet1:
		return 2
	case flow.FrameTradeParamSet2:
		return 3
	case flow.FrameTradeParamSet3:
		return 1
	case flow.FrameTradeParamSet4:
		return 3
	case flow.FrameTradeParamSet5:
		return 2
	case flow.FrameTradeParamSet6:
		return 2
	default:
		panicForValue()
	}

	return -1
}

func (param frameTradeParam) getPartialParamSet(paramSetID partialParamSetID, paramSetValueID partialParamSetValueID) partialParamSet {
	value := reflect.ValueOf(param).Field(int(paramSetID)).Interface()

	switch value.(type) {
	case flow.FrameTradeParamSet0:
		switch paramSetValueID {
		case 0:
			return flow.FrameTradeParamSet0{
				TradeDirectionLong: true,
			}
		case 1:
			return flow.FrameTradeParamSet0{
				TradeDirectionLong: false,
			}
		default:
			panicForValue()
		}

	case flow.FrameTradeParamSet1:
		switch paramSetValueID {
		case 0:
			return flow.FrameTradeParamSet1{
				SmallFrameLength: 30,
				LargeFrameLength: 60,
			}
		case 1:
			return flow.FrameTradeParamSet1{
				SmallFrameLength: 60,
				LargeFrameLength: 120,
			}
		default:
			panicForValue()
		}

	case flow.FrameTradeParamSet2:
		switch paramSetValueID {
		case 0:
			return flow.FrameTradeParamSet2{
				PipsGapForCreateOrder: 10,
			}
		case 1:
			return flow.FrameTradeParamSet2{
				PipsGapForCreateOrder: 20,
			}
		case 2:
			return flow.FrameTradeParamSet2{
				PipsGapForCreateOrder: 30,
			}
		default:
			panicForValue()
		}

	case flow.FrameTradeParamSet3:
		switch paramSetValueID {
		case 0:
			return flow.FrameTradeParamSet3{
				PipsForStopLoss: -100,
			}
		default:
			panicForValue()
		}

	case flow.FrameTradeParamSet4:
		switch paramSetValueID {
		case 0:
			return flow.FrameTradeParamSet4{
				PipsForAdditionalOrder: -5,
			}
		case 1:
			return flow.FrameTradeParamSet4{
				PipsForAdditionalOrder: -10,
			}
		case 2:
			return flow.FrameTradeParamSet4{
				PipsForAdditionalOrder: -20,
			}
		default:
			panicForValue()
		}

	case flow.FrameTradeParamSet5:
		switch paramSetValueID {
		case 0:
			return flow.FrameTradeParamSet5{
				TimeForProfit1: 40,
				TimeForProfit2: 80,
				TimeForProfit3: 120,
			}
		case 1:
			return flow.FrameTradeParamSet5{
				TimeForProfit1: 20,
				TimeForProfit2: 40,
				TimeForProfit3: 60,
			}
		default:
			panicForValue()
		}

	case flow.FrameTradeParamSet6:
		switch paramSetValueID {
		case 0:
			return flow.FrameTradeParamSet6{
				PipsForProfit1: 20,
				PipsForProfit2: 10,
				PipsForProfit3: 5,
			}
		case 1:
			return flow.FrameTradeParamSet6{
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
		case flow.FrameTradeParamSet0:
			p.FrameTradeParamSet0 = v
		case flow.FrameTradeParamSet1:
			p.FrameTradeParamSet1 = v
		case flow.FrameTradeParamSet2:
			p.FrameTradeParamSet2 = v
		case flow.FrameTradeParamSet3:
			p.FrameTradeParamSet3 = v
		case flow.FrameTradeParamSet4:
			p.FrameTradeParamSet4 = v
		case flow.FrameTradeParamSet5:
			p.FrameTradeParamSet5 = v
		case flow.FrameTradeParamSet6:
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
func FrameTradeParamsForGridSearch() []*flow.FrameTradeParam {
	params := paramsForGridSearch(frameTradeParam{})

	frameTradeParams := make([]*flow.FrameTradeParam, len(params))
	for i, param := range params {
		p, ok := param.(frameTradeParam)
		if !ok {
			panic(util.ErrWrongType)
		}
		frameTradeParam := flow.FrameTradeParam(p)
		frameTradeParams[i] = &frameTradeParam
	}

	log.Printf("%d grids for search\n", len(frameTradeParams))

	return frameTradeParams
}
