package pricer

import (
	"log"
	"math"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/trader"
)

type oandaPrice struct {
	candlesMap map[broker.TradePair]*broker.OandaCandles
	pricesMap  map[broker.TradePair]*broker.OandaClientPrice
	priceTime  int64
}

func newOandaPrice(
	candlesMap map[broker.TradePair]*broker.OandaCandles,
	pricesMap map[broker.TradePair]*broker.OandaClientPrice,
	priceTime int64,
) *oandaPrice {

	return &oandaPrice{
		candlesMap: candlesMap,
		pricesMap:  pricesMap,
		priceTime:  priceTime,
	}
}

func (oandaPrice *oandaPrice) TradePairs() []broker.TradePair {
	tradePairs := make([]broker.TradePair, len(oandaPrice.candlesMap))

	i := 0
	for pair := range oandaPrice.candlesMap {
		tradePairs[i] = pair
		i++
	}

	return tradePairs
}

func (oandaPrice *oandaPrice) Price(tradePair broker.TradePair) broker.Price {
	price, ok := oandaPrice.pricesMap[tradePair]

	if !ok || len(price.Bids) == 0 || len(price.Asks) == 0 {
		log.Panicf("no price exists for %v\n", tradePair)
	}

	return price
}

func (oandaPrice *oandaPrice) Time() int64 {
	return oandaPrice.priceTime
}

func (oandaPrice *oandaPrice) calculate(tradePair broker.TradePair, length int) *trader.Frame {
	candles, ok := oandaPrice.candlesMap[tradePair]
	if !ok {
		log.Panicf("no candle exists for %v\n", tradePair)
	}

	size := len(candles.Candles)
	if size < length {
		log.Panicf("not enough candles size to calculate frame: size %d\n", size)
	}

	mids := make([]*broker.OandaCandleStickData, length)
	for i, j := 0, size-length; j < size; i, j = i+1, j+1 {
		mids[i] = &candles.Candles[j].Mid
	}

	o := mids[0].O
	c := mids[length-1].C

	h := -math.MaxFloat64
	l := math.MaxFloat64
	for _, mid := range mids {
		if h < mid.H {
			h = mid.H
		}

		if l > mid.L {
			l = mid.L
		}
	}

	return &trader.Frame{o, h, l, c}
}
