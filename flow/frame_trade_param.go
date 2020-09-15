package flow

import (
	"reflect"
	"yukimaterrace/andaman/gridsearch"
)

type (
	// FrameTradeParamSet0 is a definition for frame trade param set 0
	FrameTradeParamSet0 struct {
		TradeDirectionLong bool
	}

	// FrameTradeParamSet1 is a definition for frame trade param set 1
	FrameTradeParamSet1 struct {
		SmallFrameLength int
		LargeFrameLength int
	}

	// FrameTradeParamSet2 is a definition for frame trade param set 2
	FrameTradeParamSet2 struct {
		PipsGapForCreateOrder float64
	}

	// FrameTradeParamSet3 is a definition for frame trade param set 3
	FrameTradeParamSet3 struct {
		PipsForStopLoss float64
	}

	// FrameTradeParamSet4 is a definition for frame trade param set 4
	FrameTradeParamSet4 struct {
		PipsForAdditionalOrder float64
	}

	// FrameTradeParamSet5 is a definition for frame trade param set 5
	FrameTradeParamSet5 struct {
		TimeForProfit1 int64
		TimeForProfit2 int64
		TimeForProfit3 int64
	}

	// FrameTradeParamSet6 is a definition for frame trade param set 6
	FrameTradeParamSet6 struct {
		PipsForProfit1 float64
		PipsForProfit2 float64
		PipsForProfit3 float64
	}
)

// FrameTradeParam is parameters for frame trade algorithm
type FrameTradeParam struct {
	FrameTradeParamSet0
	FrameTradeParamSet1
	FrameTradeParamSet2
	FrameTradeParamSet3
	FrameTradeParamSet4
	FrameTradeParamSet5
	FrameTradeParamSet6
}

func (param *FrameTradeParam) csvHeader() []string {
	return []string{
		"tradeDirectionLong",
		"smallFrameLength",
		"largeFrameLength",
		"pipsGapForCreateOrder",
		"pipsForStopLoss",
		"pipsForAdditionalOrder",
		"timeForProfit1",
		"timeForProfit2",
		"timeForProfit3",
		"pipsForProfit1",
		"pipsForProfit2",
		"pipsForProfit3",
	}
}

func (param *FrameTradeParam) csvValue() []string {
	return []string{
		paramBool(param.TradeDirectionLong).String(),
		paramInt(param.SmallFrameLength).String(),
		paramInt(param.LargeFrameLength).String(),
		paramFloat64(param.PipsGapForCreateOrder).String(),
		paramFloat64(param.PipsForStopLoss).String(),
		paramFloat64(param.PipsForAdditionalOrder).String(),
		paramInt(param.TimeForProfit1).String(),
		paramInt(param.TimeForProfit2).String(),
		paramInt(param.TimeForProfit3).String(),
		paramFloat64(param.PipsForProfit1).String(),
		paramFloat64(param.PipsForProfit2).String(),
		paramFloat64(param.PipsForProfit3).String(),
	}
}

// PartialParamSetSize is a method to get partial param set size
func (param FrameTradeParam) PartialParamSetSize() int {
	return reflect.ValueOf(param).NumField()
}

// PartialParamSetValueSize is a method to get partial param set value size
func (param FrameTradeParam) PartialParamSetValueSize(partialParamSetID gridsearch.PartialParamSetID) int {
	value := reflect.ValueOf(param).Field(int(partialParamSetID)).Interface()

	switch value.(type) {
	case FrameTradeParamSet0:
		return 2
	case FrameTradeParamSet1:
		return 2
	case FrameTradeParamSet2:
		return 3
	case FrameTradeParamSet3:
		return 1
	case FrameTradeParamSet4:
		return 3
	case FrameTradeParamSet5:
		return 2
	case FrameTradeParamSet6:
		return 2
	default:
		panicForValue()
	}

	return -1
}

// GetPartialParamSet is a method to get partial param set
func (param FrameTradeParam) GetPartialParamSet(
	partialParamSetID gridsearch.PartialParamSetID, partialParamSetValueID gridsearch.PartialParamSetValueID) gridsearch.PartialParamSet {
	value := reflect.ValueOf(param).Field(int(partialParamSetID)).Interface()

	switch value.(type) {
	case FrameTradeParamSet0:
		switch partialParamSetValueID {
		case 0:
			return true
		case 1:
			return false
		default:
			panicForValue()
		}

	case FrameTradeParamSet1:
		switch partialParamSetValueID {
		case 0:
			return FrameTradeParamSet1{
				SmallFrameLength: 30,
				LargeFrameLength: 60,
			}
		case 1:
			return FrameTradeParamSet1{
				SmallFrameLength: 60,
				LargeFrameLength: 120,
			}
		default:
			panicForValue()
		}

	case FrameTradeParamSet2:
		switch partialParamSetValueID {
		case 0:
			return FrameTradeParamSet2{
				PipsGapForCreateOrder: 10,
			}
		case 1:
			return FrameTradeParamSet2{
				PipsGapForCreateOrder: 20,
			}
		case 2:
			return FrameTradeParamSet2{
				PipsGapForCreateOrder: 30,
			}
		default:
			panicForValue()
		}

	case FrameTradeParamSet3:
		switch partialParamSetValueID {
		case 0:
			return FrameTradeParamSet3{
				PipsForStopLoss: -100,
			}
		default:
			panicForValue()
		}

	case FrameTradeParamSet4:
		switch partialParamSetValueID {
		case 0:
			return FrameTradeParamSet4{
				PipsForAdditionalOrder: -5,
			}
		case 1:
			return FrameTradeParamSet4{
				PipsForAdditionalOrder: -10,
			}
		case 2:
			return FrameTradeParamSet4{
				PipsForAdditionalOrder: -20,
			}
		default:
			panicForValue()
		}

	case FrameTradeParamSet5:
		switch partialParamSetValueID {
		case 0:
			return FrameTradeParamSet5{
				TimeForProfit1: 40,
				TimeForProfit2: 80,
				TimeForProfit3: 120,
			}
		case 1:
			return FrameTradeParamSet5{
				TimeForProfit1: 20,
				TimeForProfit2: 40,
				TimeForProfit3: 60,
			}
		default:
			panicForValue()
		}

	case FrameTradeParamSet6:
		switch partialParamSetValueID {
		case 0:
			return FrameTradeParamSet6{
				PipsForProfit1: 20,
				PipsForProfit2: 10,
				PipsForProfit3: 5,
			}
		case 1:
			return FrameTradeParamSet6{
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
func (param FrameTradeParam) CreateParam(partialParamSet ...gridsearch.PartialParamSet) gridsearch.Param {
	if len(partialParamSet) < param.PartialParamSetSize() {
		panicForValue()
	}

	p := FrameTradeParam{}

	for _, paramSet := range partialParamSet {
		switch v := paramSet.(type) {
		case FrameTradeParamSet0:
			p.FrameTradeParamSet0 = v
		case FrameTradeParamSet1:
			p.FrameTradeParamSet1 = v
		case FrameTradeParamSet2:
			p.FrameTradeParamSet2 = v
		case FrameTradeParamSet3:
			p.FrameTradeParamSet3 = v
		case FrameTradeParamSet4:
			p.FrameTradeParamSet4 = v
		case FrameTradeParamSet5:
			p.FrameTradeParamSet5 = v
		case FrameTradeParamSet6:
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
