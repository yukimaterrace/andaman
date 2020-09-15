package flow

import (
	"log"
	"strconv"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/util"
)

type frameCalculator interface {
	broker.PriceExtractor
	calculate(tradePair broker.TradePair, length int) *frame
}

type frame struct {
	o float64
	h float64
	l float64
	c float64
}

type paramBool bool

func (b paramBool) String() string {
	return strconv.FormatBool(bool(b))
}

type paramInt int

func (i paramInt) String() string {
	return strconv.FormatInt(int64(i), 10)
}

type paramFloat64 float64

func (f paramFloat64) String() string {
	return strconv.FormatFloat(float64(f), 'f', 6, 64)
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

func (algorithm *FrameTradeAlgorithm) initialTrade(material tradeMaterial, agggregator *orderAggregator, tradePair broker.TradePair) {
	calculator, ok := material.(frameCalculator)
	if !ok {
		panic(util.ErrWrongType)
	}

	price := calculator.Price(tradePair)
	spread := spread(price)

	smallFrame := calculator.calculate(tradePair, algorithm.SmallFrameLength)
	largeFrame := calculator.calculate(tradePair, algorithm.LargeFrameLength)

	gapCond := smallFrame.h-smallFrame.l > algorithm.PipsGapForCreateOrder*tradePair.PricePerPip()+spread

	if algorithm.TradeDirectionLong {
		longCond := smallFrame.h == largeFrame.h && smallFrame.l > largeFrame.l && smallFrame.l == smallFrame.c

		if longCond && gapCond {
			agggregator.createOrder(tradePair, algorithm.units, true)
		}
	} else {
		shortCond := smallFrame.l == largeFrame.l && smallFrame.h < largeFrame.h && smallFrame.h == smallFrame.c

		if shortCond && gapCond {
			agggregator.createOrder(tradePair, algorithm.units, false)
		}
	}
}

func (algorithm *FrameTradeAlgorithm) proceedTrade(material tradeMaterial, agggregator *orderAggregator, openOrders []broker.OpenOrder, tradePair broker.TradePair) {
	calculator, ok := material.(frameCalculator)
	if !ok {
		panic(util.ErrWrongType)
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

func (algorithm *FrameTradeAlgorithm) paramCsvHeader() []string {
	return algorithm.FrameTradeParam.csvHeader()
}

func (algorithm *FrameTradeAlgorithm) paramCsvValue() []string {
	return algorithm.FrameTradeParam.csvValue()
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
