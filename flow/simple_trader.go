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
	executor     *simpleTradeRunnersExecutor
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

	trader.executor.run(material, mode, tradable)

	ordersMap := map[broker.AccountID]*combinedOrders{}
	for _, runner := range trader.tradeRunners {
		combinedOrders := <-runner.done
		if len(combinedOrders.createdOrders) > 0 || len(combinedOrders.closedOrders) > 0 {
			if _, ok := ordersMap[runner.accountID]; ok {
				panic("duplicate runner for an accountID detected")
			}

			ordersMap[runner.accountID] = combinedOrders
		}
	}

	return &accountCombinedOrders{ordersMap: ordersMap}, len(ordersMap) > 0
}

type simpleTradeRunnersExecutor struct {
	runners       []*simpleTradeRunner
	runnersGroups [][]*simpleTradeRunner
	parallel      int
}

func newSimpleTradeRunnersExecutor(runners []*simpleTradeRunner, parallel int) *simpleTradeRunnersExecutor {
	if parallel == 0 {
		return &simpleTradeRunnersExecutor{
			runners:       runners,
			runnersGroups: nil,
			parallel:      0,
		}
	}

	if parallel > len(runners) {
		parallel = len(runners)
	}

	runnerGroups := make([][]*simpleTradeRunner, parallel)
	count := len(runners) / parallel
	for i := 0; i < parallel; i++ {
		runnerGroups[i] = runners[i*count : (i+1)*count]
	}

	return &simpleTradeRunnersExecutor{
		runners:       runners,
		runnersGroups: runnerGroups,
		parallel:      parallel,
	}
}

func (executor *simpleTradeRunnersExecutor) run(material tradeMaterial, mode TradeMode, tradable bool) {
	if executor.parallel == 0 {
		for _, runner := range executor.runners {
			runner.run(material, mode, tradable)
		}
	} else {
		for i := 0; i < executor.parallel; i++ {
			go func(runners []*simpleTradeRunner) {
				for _, runner := range runners {
					runner.run(material, mode, tradable)
				}
			}(executor.runnersGroups[i])
		}
	}
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
	done         chan *combinedOrders
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

type combinedOrders struct {
	createdOrders []broker.CreatedOrder
	closedOrders  []broker.ClosedOrder
}

type accountCombinedOrders struct {
	ordersMap map[broker.AccountID]*combinedOrders
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

func (aggregator *orderAggregator) reduce() *combinedOrders {
	createdOrders := []broker.CreatedOrder{}
	for _, done := range aggregator.createOrderDone {
		createdOrder := <-done
		if createdOrder.Err != nil {
			log.Println(createdOrder.Err.Error())
		} else {
			createdOrders = append(createdOrders, createdOrder.CreatedOrder)
		}
	}

	closedOrders := []broker.ClosedOrder{}
	for _, done := range aggregator.closeOrderDone {
		closedOrder := <-done
		if closedOrder.Err != nil {
			log.Println(closedOrder.Err.Error())
		} else {
			closedOrders = append(closedOrders, closedOrder.ClosedOrder)
		}
	}

	return &combinedOrders{
		createdOrders: createdOrders,
		closedOrders:  closedOrders,
	}
}

// SimpleTraderBuilder is a builder for simple trader
type SimpleTraderBuilder struct {
	configs        []*simpleTraderConfig
	tradableFunc   TradableFunc
	broker         broker.Broker
	ordererFactory broker.OrdererFactory
	parallel       int
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

// Parallel sets parallel parameter in builder
func (builder *SimpleTraderBuilder) Parallel(paralle int) *SimpleTraderBuilder {
	builder.parallel = paralle
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
			done:         make(chan *combinedOrders, 1),
		})
	}

	executor := newSimpleTradeRunnersExecutor(tradeRunners, builder.parallel)

	return &SimpleTrader{
		broker:       builder.broker,
		tradableFunc: builder.tradableFunc,
		tradeRunners: tradeRunners,
		executor:     executor,
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
