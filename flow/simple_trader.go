package flow

import (
	"log"
	"yukimaterrace/andaman/broker"
)

// SimpleTrader is a struct for simple trader
type SimpleTrader struct {
	broker       broker.Broker
	tradableFunc TradableFunc
	tradeRunners []*simpleTradeRunner
}

func (trader *SimpleTrader) trade(material tradeMaterial, mode TradeMode) (recordMaterial, bool) {
	if simulationBroker, ok := trader.broker.(broker.SimulationBroker); ok {
		if price, ok := material.(broker.PriceExtractor); ok {
			simulationBroker.Update(price)
		} else {
			panic("wrong type has been passed")
		}
	}

	if mode == Watch {
		return nil, false
	}

	tradable := trader.tradableFunc(material)

	for _, runner := range trader.tradeRunners {
		go func(runner *simpleTradeRunner) {
			runner.run(material, mode, tradable)
		}(runner)
	}

	results := make([]*accountOrders, 0)
	for _, runner := range trader.tradeRunners {
		accountOrders := <-runner.done
		if len(accountOrders.createdOrders) > 0 || len(accountOrders.closedOrders) > 0 {
			results = append(results, accountOrders)
		}
	}

	return results, len(results) > 0
}

// SimpleTradeAlgorithm is an interface for simple trade algorithm
type SimpleTradeAlgorithm interface {
	initialTrade(material tradeMaterial, aggregator *orderAggregator, tradePair broker.TradePair)

	proceedTrade(material tradeMaterial, aggregator *orderAggregator, openOrders []broker.OpenOrder, tradePair broker.TradePair)
}

type simpleTradeRunner struct {
	accountID    broker.AccountID
	algorithmMap map[broker.TradePair]SimpleTradeAlgorithm
	orderer      broker.Orderer
	done         chan *accountOrders
}

func (runner *simpleTradeRunner) run(material tradeMaterial, mode TradeMode, tradable bool) {
	res := <-runner.orderer.OpenOrders(runner.accountID)
	if res.Err != nil {
		log.Println(res.Err.Error())
		runner.done <- nil
		return
	}

	openOrdersMap := map[broker.TradePair][]broker.OpenOrder{}
	for _, openOrder := range res.OpenOrders {
		openOrders, ok := openOrdersMap[openOrder.TradePair()]
		if !ok {
			openOrders = make([]broker.OpenOrder, 0)
		}

		openOrdersMap[openOrder.TradePair()] = append(openOrders, openOrder)
	}

	aggregator := newOrderAggregator(runner.orderer, runner.accountID)
	for tradePair, openOrders := range openOrdersMap {
		algorithm, ok := runner.algorithmMap[tradePair]
		if !ok {
			log.Panicf("no algorithm registered for %v\n", tradePair)
		}

		if len(openOrders) == 0 {
			if mode != Terminate && tradable {
				algorithm.initialTrade(material, aggregator, tradePair)
			}
		} else {
			algorithm.proceedTrade(material, aggregator, openOrders, tradePair)
		}
	}

	runner.done <- aggregator.reduce()
}

type accountCreatedOrder struct {
	accountID    broker.AccountID
	createdOrder broker.CreatedOrder
}

type accountClosedOrder struct {
	accountID   broker.AccountID
	closedOrder broker.ClosedOrder
}

type accountOrders struct {
	createdOrders []*accountCreatedOrder
	closedOrders  []*accountClosedOrder
}

type orderAggregator struct {
	broker.Orderer
	accountID       broker.AccountID
	createOrderDone []<-chan *broker.CreateOrderResult
	closeOrderDone  []<-chan *broker.CloseOrderResult
}

func newOrderAggregator(orderer broker.Orderer, accountID broker.AccountID) *orderAggregator {
	return &orderAggregator{
		Orderer:         orderer,
		accountID:       accountID,
		createOrderDone: []<-chan *broker.CreateOrderResult{},
		closeOrderDone:  []<-chan *broker.CloseOrderResult{},
	}
}

func (aggregator *orderAggregator) createOrder(tradePair broker.TradePair, units float64, isLong bool) {
	result := aggregator.CreateOrder(aggregator.accountID, tradePair, units, isLong)
	aggregator.createOrderDone = append(aggregator.createOrderDone, result)
}

func (aggregator *orderAggregator) closeOrder(orderID broker.OrderID) {
	result := aggregator.CloseOrder(aggregator.accountID, orderID)
	aggregator.closeOrderDone = append(aggregator.closeOrderDone, result)
}

func (aggregator *orderAggregator) reduce() *accountOrders {
	accountCreatedOrders := []*accountCreatedOrder{}
	for _, done := range aggregator.createOrderDone {
		createdOrder := <-done
		if createdOrder.Err != nil {
			log.Println(createdOrder.Err.Error())
		} else {
			accountCreatedOrders = append(accountCreatedOrders, &accountCreatedOrder{
				accountID:    aggregator.accountID,
				createdOrder: createdOrder.CreatedOrder,
			})
		}
	}

	accountClosedOrders := []*accountClosedOrder{}
	for _, done := range aggregator.closeOrderDone {
		closedOrder := <-done
		if closedOrder.Err != nil {
			log.Println(closedOrder.Err.Error())
		} else {
			accountClosedOrders = append(accountClosedOrders, &accountClosedOrder{
				accountID:   aggregator.accountID,
				closedOrder: closedOrder.ClosedOrder,
			})
		}
	}

	return &accountOrders{
		createdOrders: accountCreatedOrders,
		closedOrders:  accountClosedOrders,
	}
}

// SimpleTraderBuilder is a builder for simple trader
type SimpleTraderBuilder struct {
	configs        []*simpleTraderConfig
	tradableFunc   TradableFunc
	broker         broker.Broker
	ordererFactory broker.OrdererFactory
}

// Trade is a method to add a trade piece in builder
func (builder *SimpleTraderBuilder) Trade(accountID broker.AccountID, tradePair broker.TradePair, algorithm SimpleTradeAlgorithm) *SimpleTraderBuilder {
	builder.configs = append(builder.configs, &simpleTraderConfig{
		accountID: accountID,
		tradePair: tradePair,
		algorithm: algorithm,
	})
	return builder
}

// TradableFunc sets tradablefunc in builder
func (builder *SimpleTraderBuilder) TradableFunc(tradableFunc TradableFunc) *SimpleTraderBuilder {
	builder.tradableFunc = tradableFunc
	return builder
}

// Broker sets broker in builder
func (builder *SimpleTraderBuilder) Broker(broker broker.Broker) *SimpleTraderBuilder {
	builder.broker = broker
	return builder
}

// OrdererFactory sets orderer factory in builder
func (builder *SimpleTraderBuilder) OrdererFactory(ordererFactory broker.OrdererFactory) *SimpleTraderBuilder {
	builder.ordererFactory = ordererFactory
	return builder
}

// Build builds simple trader
func (builder *SimpleTraderBuilder) Build() *SimpleTrader {
	configMap := make(map[broker.AccountID]map[broker.TradePair]SimpleTradeAlgorithm)

	for _, config := range builder.configs {
		algorithmMap, ok := configMap[config.accountID]
		if !ok {
			algorithmMap = make(map[broker.TradePair]SimpleTradeAlgorithm)
			configMap[config.accountID] = algorithmMap
		}

		algorithmMap[config.tradePair] = config.algorithm
	}

	tradeRunners := make([]*simpleTradeRunner, 0)
	for accountID, algorithmMap := range configMap {
		tradeRunners = append(tradeRunners, &simpleTradeRunner{
			accountID:    accountID,
			algorithmMap: algorithmMap,
			orderer:      builder.ordererFactory.Create(builder.broker),
			done:         make(chan *accountOrders, 1),
		})
	}

	return &SimpleTrader{
		broker:       builder.broker,
		tradableFunc: builder.tradableFunc,
		tradeRunners: tradeRunners,
	}
}

type simpleTraderConfig struct {
	accountID broker.AccountID
	tradePair broker.TradePair
	algorithm SimpleTradeAlgorithm
}

// SimpleTraderFactory is a factory for simple trader
type SimpleTraderFactory struct {
	builder *SimpleTraderBuilder
}

// NewSimpleTraderFactory is a constructor for simple trader factory
func NewSimpleTraderFactory(builder *SimpleTraderBuilder) *SimpleTraderFactory {
	return &SimpleTraderFactory{
		builder: builder,
	}
}

func (factory *SimpleTraderFactory) create(broker broker.Broker, ordererFactory broker.OrdererFactory) trader {
	return factory.builder.Broker(broker).OrdererFactory(ordererFactory).Build()
}
