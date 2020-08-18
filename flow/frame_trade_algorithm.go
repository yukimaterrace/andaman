package flow

import (
	"log"
	"os"
	"strconv"
	"yukimaterrace/andaman/broker"
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

type paramInt int

func (i paramInt) String() string {
	return strconv.FormatInt(int64(i), 10)
}

type paramFloat64 float64

func (f paramFloat64) String() string {
	return strconv.FormatFloat(float64(f), 'f', 6, 64)
}

// FrameTradeParam is parameters for frame trade algorithm
type FrameTradeParam struct {
	SmallFrameLength       int
	LargeFrameLength       int
	PipsGapForCreateOrder  float64
	PipsForStopLoss        float64
	PipsForAdditionalOrder float64
	TimeForProfit1         int
	TimeForProfit2         int
	TimeForProfit3         int
	PipsForProfit1         float64
	PipsForProfit2         float64
	PipsForProfit3         float64
	PipsForProfit4         float64
}

func (param *FrameTradeParam) csvHeader() []string {
	return []string{
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
		"pipsForProfit4",
	}
}

func (param *FrameTradeParam) csvValue() []string {
	return []string{
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
		paramFloat64(param.PipsForProfit4).String(),
	}
}

// FrameTradeAlgorithm is a struct for frame trade algorithm
type FrameTradeAlgorithm struct {
	*FrameTradeParam
	units float64
}

// NewFrameTradeAlgorithm is a constructor for new frame trade algorithm
func NewFrameTradeAlgorithm(param *FrameTradeParam) *FrameTradeAlgorithm {
	unitsEnv := os.Getenv("UNITS")
	if unitsEnv == "" {
		log.Panicln("UNITS has not been set")
	}

	units, err := strconv.ParseFloat(unitsEnv, 64)
	if err != nil {
		log.Panicln(err.Error())
	}

	return &FrameTradeAlgorithm{
		FrameTradeParam: param,
		units:           units,
	}
}

func (algorithm *FrameTradeAlgorithm) initialTrade(
	material tradeMaterial, orderer broker.Orderer, accountID broker.AccountID, tradePair broker.TradePair,
) *simpleTradeAlgorithmResult {
	calculator, ok := material.(frameCalculator)
	if !ok {
		panic("wrong type has been passed")
	}

	price := calculator.Price(tradePair)
	spread := spread(price)

	smallFrame := calculator.calculate(tradePair, algorithm.SmallFrameLength)
	largeFrame := calculator.calculate(tradePair, algorithm.LargeFrameLength)

	gapCond := smallFrame.h-smallFrame.l > algorithm.PipsGapForCreateOrder*tradePair.PricePerPip()+spread

	longCond := smallFrame.h == largeFrame.h && smallFrame.l > largeFrame.l && smallFrame.l == smallFrame.c

	if longCond && gapCond {
		return &simpleTradeAlgorithmResult{
			createOrderDone: []<-chan *broker.CreateOrderResult{
				orderer.CreateOrder(accountID, tradePair, algorithm.units, true),
			},

			closeOrderDone: []<-chan *broker.CloseOrderResult{},
		}
	}

	shortCond := smallFrame.l == largeFrame.l && smallFrame.h < largeFrame.h && smallFrame.h == smallFrame.c

	if shortCond && gapCond {
		return &simpleTradeAlgorithmResult{
			createOrderDone: []<-chan *broker.CreateOrderResult{
				orderer.CreateOrder(accountID, tradePair, algorithm.units, false),
			},

			closeOrderDone: []<-chan *broker.CloseOrderResult{},
		}
	}

	return nil
}

func (algorithm *FrameTradeAlgorithm) proceedTrade(
	material tradeMaterial, orderer broker.Orderer, openOrders []broker.OpenOrder, accountID broker.AccountID, tradePair broker.TradePair,
) *simpleTradeAlgorithmResult {

	return nil
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
