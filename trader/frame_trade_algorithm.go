package trader

import (
	"log"
	"strconv"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/util"
)

type frameCalculator interface {
	broker.PriceExtractor
	Calculate(tradePair model.TradePair, length int) *Frame
}

// Frame is a struct for frame
type Frame struct {
	O float64
	H float64
	L float64
	C float64
}

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

// FrameTradeAlgorithm is a struct for frame trade algorithm
type FrameTradeAlgorithm struct {
	*FrameTradeParam
	units float64
}

// NewFrameTradeAlgorithm is a constructor for new frame trade algorithm
func NewFrameTradeAlgorithm(param *FrameTradeParam) *FrameTradeAlgorithm {
	unitsEnv := util.GetEnv("UNITS")

	units, err := strconv.ParseFloat(unitsEnv, 64)
	if err != nil {
		log.Panicln(err.Error())
	}

	return &FrameTradeAlgorithm{
		FrameTradeParam: param,
		units:           units,
	}
}

func (algorithm *FrameTradeAlgorithm) initialTrade(material flow.TradeMaterial, agggregator *orderAggregator, tradePair model.TradePair) {
	calculator, ok := material.(frameCalculator)
	if !ok {
		panic(model.ErrWrongType)
	}

	price := calculator.Price(tradePair)
	spread := spread(price)

	smallFrame := calculator.Calculate(tradePair, algorithm.SmallFrameLength)
	largeFrame := calculator.Calculate(tradePair, algorithm.LargeFrameLength)

	gapCond := smallFrame.H-smallFrame.L > algorithm.PipsGapForCreateOrder*tradePair.PricePerPip()+spread

	if algorithm.TradeDirectionLong {
		longCond := smallFrame.H == largeFrame.H && smallFrame.L > largeFrame.L && smallFrame.L == smallFrame.C

		if longCond && gapCond {
			agggregator.createOrder(tradePair, algorithm.units, true)
		}
	} else {
		shortCond := smallFrame.L == largeFrame.L && smallFrame.H < largeFrame.H && smallFrame.H == smallFrame.C

		if shortCond && gapCond {
			agggregator.createOrder(tradePair, algorithm.units, false)
		}
	}
}

func (algorithm *FrameTradeAlgorithm) proceedTrade(material flow.TradeMaterial, agggregator *orderAggregator, openOrders []broker.OpenOrder, tradePair model.TradePair) {
	calculator, ok := material.(frameCalculator)
	if !ok {
		panic(model.ErrWrongType)
	}

	if len(openOrders) == 0 {
		panic("no open orders has been passed")
	}

	currentPrice := calculator.Price(tradePair)

	initialOrder := openOrders[0]
	lastOrder := initialOrder
	totalProfitPips := 0.0
	for _, openOrder := range openOrders {
		if initialOrder.TimeAtOpen() > openOrder.TimeAtOpen() {
			initialOrder = openOrder
		}

		if lastOrder.TimeAtOpen() < openOrder.TimeAtOpen() {
			lastOrder = openOrder
		}

		totalProfitPips += profit(openOrder, currentPrice)
	}

	totalProfitPips /= tradePair.PricePerPip()
	tradeTime := calculator.Time() - initialOrder.TimeAtOpen()

	takeProfitCond1 := tradeTime <= algorithm.TimeForProfit1*60 && totalProfitPips >= algorithm.PipsForProfit1
	takeProfitCond2 := algorithm.TimeForProfit1*60 < tradeTime && tradeTime <= algorithm.TimeForProfit2*60 && totalProfitPips >= algorithm.PipsForProfit2
	takeProfitCond3 := algorithm.TimeForProfit2*60 < tradeTime && tradeTime <= algorithm.TimeForProfit3*60 && totalProfitPips >= algorithm.PipsForProfit3
	takeProfitCond4 := algorithm.TimeForProfit3*60 < tradeTime && totalProfitPips >= 0

	stopLossCond := totalProfitPips < algorithm.PipsForStopLoss

	if takeProfitCond1 || takeProfitCond2 || takeProfitCond3 || takeProfitCond4 || stopLossCond {
		for _, openOrder := range openOrders {
			agggregator.closeOrder(openOrder.OrderID())
		}
		return
	}

	lastOrderProfitPips := profit(lastOrder, currentPrice) / tradePair.PricePerPip()

	if lastOrderProfitPips < algorithm.PipsForAdditionalOrder-spread(currentPrice)/tradePair.PricePerPip() {
		agggregator.createOrder(tradePair, algorithm.units, lastOrder.IsLong())
	}
}

func profit(openOrder broker.OpenOrder, price broker.Price) float64 {
	if openOrder.IsLong() {
		return price.Bid() - openOrder.PriceAtOpen()
	}
	return openOrder.PriceAtOpen() - price.Ask()
}

func spread(price broker.Price) float64 {
	return price.Ask() - price.Bid()
}
