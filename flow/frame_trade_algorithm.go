package flow

import "yukimaterrace/andaman/broker"

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
