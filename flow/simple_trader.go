package flow

import (
	"log"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/util"
)

// SimpleTrader is a struct for simple trader
type SimpleTrader struct {
	broker       broker.Broker
	tradeRunners []*simpleTradeRunner
	executor     *simpleTradeRunnersExecutor
}

func (trader *SimpleTrader) trade(material tradeMaterial, mode TradeMode) (recordMaterial, bool) {
	if simulationBroker, ok := trader.broker.(broker.SimulationBroker); ok {
		if price, ok := material.(broker.PriceExtractor); ok {
			simulationBroker.Update(price)
		} else {
			panic(util.ErrWrongType)
		}
	}

	if mode == Watch {
		return nil, false
	}

	trader.executor.run(material, mode)

	ordersMap := map[broker.AccountID]*combinedOrders{}
	for _, runner := range trader.tradeRunners {
		combinedOrders := <-runner.done
		if combinedOrders != nil && (len(combinedOrders.createdOrders) > 0 || len(combinedOrders.closedOrders) > 0) {
			if _, ok := ordersMap[runner.accountID]; ok {
				panic("duplicate runner for an accountID detected")
			}

			ordersMap[runner.accountID] = combinedOrders
		}
	}

	return accountCombinedOrders(ordersMap), len(ordersMap) > 0
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
		if i == parallel-1 {
			runnerGroups[i] = runners[i*count:]
		} else {
			runnerGroups[i] = runners[i*count : (i+1)*count]
		}
	}

	return &simpleTradeRunnersExecutor{
		runners:       runners,
		runnersGroups: runnerGroups,
		parallel:      parallel,
	}
}

func (executor *simpleTradeRunnersExecutor) run(material tradeMaterial, mode TradeMode) {
	if executor.parallel == 0 {
		for _, runner := range executor.runners {
			runner.run(material, mode)
		}
	} else {
		for i := 0; i < executor.parallel; i++ {
			go func(runners []*simpleTradeRunner) {
				for _, runner := range runners {
					runner.run(material, mode)
				}
			}(executor.runnersGroups[i])
		}
	}
}

// SimpleTradeAlgorithm is an interface for simple trade algorithm
type SimpleTradeAlgorithm interface {
	initialTrade(material tradeMaterial, aggregator *orderAggregator, tradePair broker.TradePair)
	proceedTrade(material tradeMaterial, aggregator *orderAggregator, openOrders []broker.OpenOrder, tradePair broker.TradePair)
	tradeParamLoader
}

type simpleTradeRunner struct {
	accountID         broker.AccountID
	tradableTimeZone  *TradableTimeZone
	algorithmMap      map[broker.TradePair]SimpleTradeAlgorithm
	orderer           broker.Orderer
	done              chan *combinedOrders
	openOrdersExisted bool
}

func (runner *simpleTradeRunner) run(material tradeMaterial, mode TradeMode) {
	tradable := runner.tradableTimeZone.OK(material)

	if !runner.openOrdersExisted && !tradable {
		runner.done <- nil
		return
	}

	res := <-runner.orderer.OpenOrders(runner.accountID)
	if res.Err != nil {
		log.Println(res.Err.Error())
		runner.done <- nil
		return
	}

	if len(res.OpenOrders) > 0 {
		runner.openOrdersExisted = true
	} else {
		runner.openOrdersExisted = false
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

type accountCombinedOrders map[broker.AccountID]*combinedOrders

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
	tradableTimeZoneMap map[broker.AccountID]*TradableTimeZone
	algorithmMap        map[broker.AccountID]map[broker.TradePair]SimpleTradeAlgorithm
	broker              broker.Broker
	ordererFactory      broker.OrdererFactory
	parallel            int
}

// NewSimpleTraderBuilder is a constructor for simple trader builder
func NewSimpleTraderBuilder() *SimpleTraderBuilder {
	return &SimpleTraderBuilder{
		tradableTimeZoneMap: map[broker.AccountID]*TradableTimeZone{},
		algorithmMap:        map[broker.AccountID]map[broker.TradePair]SimpleTradeAlgorithm{},
	}
}

// TradableTimeZone is a method to add a tradable time zone
func (builder *SimpleTraderBuilder) TradableTimeZone(accountID broker.AccountID, tradableTimeZone *TradableTimeZone) *SimpleTraderBuilder {
	if _, ok := builder.tradableTimeZoneMap[accountID]; ok {
		panic("duplicate tradable time zone for an account ID detected")
	}

	builder.tradableTimeZoneMap[accountID] = tradableTimeZone
	return builder
}

// Trade is a method to add a trade piece in builder
func (builder *SimpleTraderBuilder) Trade(accountID broker.AccountID, tradePair broker.TradePair, algorithm SimpleTradeAlgorithm) *SimpleTraderBuilder {
	_, ok := builder.algorithmMap[accountID]
	if !ok {
		builder.algorithmMap[accountID] = map[broker.TradePair]SimpleTradeAlgorithm{}
	}

	builder.algorithmMap[accountID][tradePair] = algorithm
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
	tradeRunners := []*simpleTradeRunner{}
	for accountID, algorithmMap := range builder.algorithmMap {
		tradableTimeZone, ok := builder.tradableTimeZoneMap[accountID]
		if !ok {
			panic("no tradable time zone specified")
		}

		tradeRunners = append(tradeRunners, &simpleTradeRunner{
			accountID:        accountID,
			tradableTimeZone: tradableTimeZone,
			algorithmMap:     algorithmMap,
			orderer:          builder.ordererFactory.Create(builder.broker),
			done:             make(chan *combinedOrders, 1),
		})
	}

	executor := newSimpleTradeRunnersExecutor(tradeRunners, builder.parallel)

	return &SimpleTrader{
		broker:       builder.broker,
		tradeRunners: tradeRunners,
		executor:     executor,
	}
}

// BuildTradeSpecs builds trade specs
func (builder *SimpleTraderBuilder) BuildTradeSpecs() *TradeSpecs {
	paramLoaders := map[broker.AccountID]map[broker.TradePair]tradeParamLoader{}

	for accountID, algorithmMap := range builder.algorithmMap {
		if _, ok := paramLoaders[accountID]; !ok {
			paramLoaders[accountID] = map[broker.TradePair]tradeParamLoader{}
		}

		for tradePair, algorithm := range algorithmMap {
			paramLoaders[accountID][tradePair] = algorithm
		}
	}

	return &TradeSpecs{
		timeZones:    builder.tradableTimeZoneMap,
		paramLoaders: paramLoaders,
	}
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
