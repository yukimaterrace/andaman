package flow

import (
	"time"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/config"
)

// Builder is a builder for flow
type Builder struct {
	broker           broker.Broker
	ordererFactory   broker.OrdererFactory
	pricerTradePairs []broker.TradePair
	initialTradeMode TradeMode
	pricerFactory    PricerFactory
	traderFactory    TraderFactory
	recorderFactory  RecorderFactory
}

// Broker sets broker in builder
func (builder *Builder) Broker(broker broker.Broker) *Builder {
	builder.broker = broker
	return builder
}

// OrdererFactory sets orderer factory in builder
func (builder *Builder) OrdererFactory(ordererFactory broker.OrdererFactory) *Builder {
	builder.ordererFactory = ordererFactory
	return builder
}

// PricerTradePairs sets pricer trade pairs in builder
func (builder *Builder) PricerTradePairs(pricerTradePairs []broker.TradePair) *Builder {
	builder.pricerTradePairs = pricerTradePairs
	return builder
}

// InitialTradeMode sets initial trade mode in builder
func (builder *Builder) InitialTradeMode(initialTradeMode TradeMode) *Builder {
	builder.initialTradeMode = initialTradeMode
	return builder
}

// PricerFactory sets pricer factory in builder
func (builder *Builder) PricerFactory(pricerFactory PricerFactory) *Builder {
	builder.pricerFactory = pricerFactory
	return builder
}

// TraderFactory sets trader factory in builder
func (builder *Builder) TraderFactory(traderFactory TraderFactory) *Builder {
	builder.traderFactory = traderFactory
	return builder
}

// RecorderFactory sets recorder factory in builder
func (builder *Builder) RecorderFactory(recorderFactory RecorderFactory) *Builder {
	builder.recorderFactory = recorderFactory
	return builder
}

// Build builds flow
func (builder *Builder) Build() *Flow {
	pricer := builder.pricerFactory.create(builder.broker, builder.pricerTradePairs)
	trader := builder.traderFactory.create(builder.broker, builder.ordererFactory)
	recorder := builder.recorderFactory.create()

	recordWorker := &recordWorker{
		recorder: recorder,
		ch:       make(chan interface{}, config.FlowChanCap),
		ticker:   time.NewTicker(time.Minute), // write ever 1 minute
	}

	tradeWorker := &tradeWorker{
		trader:       trader,
		recordWorker: recordWorker,
		mode:         builder.initialTradeMode,
		ch:           make(chan interface{}, config.FlowChanCap),
	}

	priceWorker := &priceWorker{
		pricer:            pricer,
		tradeWorker:       tradeWorker,
		createPriceResult: make(chan *createPriceResult, 1),
		ch:                make(chan interface{}, config.FlowChanCap),
		init:              false,
	}

	return &Flow{
		priceWorker:      priceWorker,
		tradeWorker:      tradeWorker,
		recordWorker:     recordWorker,
		priceWorkerDone:  make(chan bool, 1),
		tradeWorkerDone:  make(chan bool, 1),
		recordWorkerDone: make(chan bool, 1),
	}
}
