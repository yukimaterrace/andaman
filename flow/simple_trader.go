package flow

import (
	"log"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/config"
	"yukimaterrace/andaman/util"
)

// SimpleTrader is a struct for simple trader
type SimpleTrader struct {
	orderer                  broker.Orderer
	orderPartitionAggregator *orderPartitionAggregator
	tradeRunners             []*simpleTradeRunner
	executor                 *simpleTradeRunnersExecutor
}

func (trader *SimpleTrader) trade(material tradeMaterial, mode TradeMode) (recordMaterial, bool) {
	if simulationBroker, ok := trader.orderer.(broker.SimulationBroker); ok {
		if price, ok := material.(broker.PriceExtractor); ok {
			simulationBroker.Update(price)
		} else {
			panic(util.ErrWrongType)
		}
	}

	if mode == Watch {
		return nil, false
	}

	openOrdersResult := <-trader.orderer.OpenOrders()
	if openOrdersResult.Err != nil {
		return nil, false
	}

	partitionedOpenOrders := <-trader.orderPartitionAggregator.partitionedOpenOrders(openOrdersResult.OpenOrders)
	trader.executor.run(material, partitionedOpenOrders, mode)

	ordersMap := map[PartitionID]*combinedOrders{}
	for _, runner := range trader.tradeRunners {
		combinedOrders := <-runner.done
		if combinedOrders != nil && (len(combinedOrders.createdOrders) > 0 || len(combinedOrders.closedOrders) > 0) {
			if _, ok := ordersMap[runner.partitionID]; ok {
				panic("duplicate runner for an partitionID detected")
			}

			ordersMap[runner.partitionID] = combinedOrders
		}
	}

	return partitionCombinedOrders(ordersMap), len(ordersMap) > 0
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

func (executor *simpleTradeRunnersExecutor) run(material tradeMaterial, partitionedOpenOrders partitionedOpenOrders, mode TradeMode) {
	if executor.parallel == 0 {
		for _, runner := range executor.runners {
			runner.run(material, partitionedOpenOrders, mode)
		}
	} else {
		for i := 0; i < executor.parallel; i++ {
			go func(runners []*simpleTradeRunner) {
				for _, runner := range runners {
					runner.run(material, partitionedOpenOrders, mode)
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
	partitionID            PartitionID
	tradableTimeZone       *TradableTimeZone
	algorithmMap           map[broker.TradePair]SimpleTradeAlgorithm
	orderAggregatorFactory *orderAggregatorFactory
	done                   chan *combinedOrders
	openOrdersExisted      bool
}

func (runner *simpleTradeRunner) run(material tradeMaterial, partitionedOpenOrders partitionedOpenOrders, mode TradeMode) {
	timeExtractor, ok := material.(broker.TimeExtractor)
	if !ok {
		panic(util.ErrWrongType)
	}

	tradable := runner.tradableTimeZone.OK(timeExtractor)

	if !runner.openOrdersExisted && !tradable {
		runner.done <- nil
		return
	}

	openOrdersMap := partitionedOpenOrders[runner.partitionID]

	if len(openOrdersMap) > 0 {
		runner.openOrdersExisted = true
	} else {
		runner.openOrdersExisted = false
	}

	aggregator := runner.orderAggregatorFactory.create(runner.partitionID)
	for tradePair, algorithm := range runner.algorithmMap {
		openOrders, ok := openOrdersMap[tradePair]
		if !ok {
			if mode != Terminate && tradable {
				algorithm.initialTrade(material, aggregator, tradePair)
			}
		} else {
			algorithm.proceedTrade(material, aggregator, openOrders, tradePair)
		}
	}

	runner.done <- aggregator.reduce()
}

type (
	combinedOrders struct {
		createdOrders []broker.CreatedOrder
		closedOrders  []broker.ClosedOrder
	}

	partitionCombinedOrders map[PartitionID]*combinedOrders
)

type orderAggregator struct {
	broker.Orderer
	partitionID              PartitionID
	orderPartitionAggregator *orderPartitionAggregator
	createOrderDone          []<-chan *broker.CreateOrderResult
	closeOrderDone           []<-chan *broker.CloseOrderResult
}

func newOrderAggregator(orderer broker.Orderer, partitionID PartitionID, orderPartitionAggregator *orderPartitionAggregator) *orderAggregator {
	return &orderAggregator{
		Orderer:                  orderer,
		partitionID:              partitionID,
		orderPartitionAggregator: orderPartitionAggregator,
		createOrderDone:          []<-chan *broker.CreateOrderResult{},
		closeOrderDone:           []<-chan *broker.CloseOrderResult{},
	}
}

type orderAggregatorFactory struct {
	orderer                  broker.Orderer
	orderPartitionAggregator *orderPartitionAggregator
}

func newOrderAggregatorFactory(orderer broker.Orderer, orderPartitionAggregator *orderPartitionAggregator) *orderAggregatorFactory {
	return &orderAggregatorFactory{
		orderer:                  orderer,
		orderPartitionAggregator: orderPartitionAggregator,
	}
}

func (factory *orderAggregatorFactory) create(partitionID PartitionID) *orderAggregator {
	return newOrderAggregator(factory.orderer, partitionID, factory.orderPartitionAggregator)
}

func (aggregator *orderAggregator) createOrder(tradePair broker.TradePair, units float64, isLong bool) {
	result := aggregator.CreateOrder(tradePair, units, isLong)
	aggregator.createOrderDone = append(aggregator.createOrderDone, result)
}

func (aggregator *orderAggregator) closeOrder(orderID broker.OrderID) {
	result := aggregator.CloseOrder(orderID)
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
			aggregator.orderPartitionAggregator.put(createdOrder.CreatedOrder.OrderID(), aggregator.partitionID)
		}
	}

	closedOrders := []broker.ClosedOrder{}
	for _, done := range aggregator.closeOrderDone {
		closedOrder := <-done
		if closedOrder.Err != nil {
			log.Println(closedOrder.Err.Error())
		} else {
			closedOrders = append(closedOrders, closedOrder.ClosedOrder)
			aggregator.orderPartitionAggregator.deelete(closedOrder.ClosedOrder.OrderID())
		}
	}

	return &combinedOrders{
		createdOrders: createdOrders,
		closedOrders:  closedOrders,
	}
}

// PartitionID is a definition for parition id
type PartitionID int

type partitionedOpenOrders map[PartitionID]map[broker.TradePair][]broker.OpenOrder

type orderPartitionAggregator struct {
	ch                chan interface{}
	orderPartitionMap map[broker.OrderID]PartitionID
}

func newOrderPartitionAggregator() *orderPartitionAggregator {
	return &orderPartitionAggregator{
		ch:                make(chan interface{}, config.OrderPartitionAggregatorChanCap),
		orderPartitionMap: map[broker.OrderID]PartitionID{},
	}
}

type (
	putOrderPartitionRequest struct {
		orderID     broker.OrderID
		partitionID PartitionID
	}

	deleteOrderPartitionRequest struct {
		orderID broker.OrderID
	}

	partitionedOpenOrdersRequest struct {
		openOrders []broker.OpenOrder
		done       chan<- partitionedOpenOrders
	}
)

func (aggregator *orderPartitionAggregator) put(orderID broker.OrderID, partitionID PartitionID) {
	aggregator.ch <- &putOrderPartitionRequest{
		orderID:     orderID,
		partitionID: partitionID,
	}
}

func (aggregator *orderPartitionAggregator) deelete(orderID broker.OrderID) {
	aggregator.ch <- &deleteOrderPartitionRequest{
		orderID: orderID,
	}
}

func (aggregator *orderPartitionAggregator) partitionedOpenOrders(openOrders []broker.OpenOrder) <-chan partitionedOpenOrders {
	done := make(chan partitionedOpenOrders, 1)
	aggregator.ch <- &partitionedOpenOrdersRequest{
		openOrders: openOrders,
		done:       done,
	}

	return done
}

func (aggregator *orderPartitionAggregator) run() {
	go func() {
		for {
			aggregator.work()
		}
	}()
}

func (aggregator *orderPartitionAggregator) work() {
	request := <-aggregator.ch

	switch req := request.(type) {
	case *putOrderPartitionRequest:
		if _, ok := aggregator.orderPartitionMap[req.orderID]; ok {
			panic("duplicate order id detected")
		}
		aggregator.orderPartitionMap[req.orderID] = req.partitionID

	case *deleteOrderPartitionRequest:
		if _, ok := aggregator.orderPartitionMap[req.orderID]; !ok {
			panic("no order id detected")
		}
		delete(aggregator.orderPartitionMap, req.orderID)

	case *partitionedOpenOrdersRequest:
		result := map[PartitionID]map[broker.TradePair][]broker.OpenOrder{}

		for _, openOrder := range req.openOrders {
			partitionID, ok := aggregator.orderPartitionMap[openOrder.OrderID()]
			if !ok {
				panic("no order id detected")
			}

			if _, ok := result[partitionID]; !ok {
				result[partitionID] = map[broker.TradePair][]broker.OpenOrder{}
			}

			if _, ok := result[partitionID][openOrder.TradePair()]; !ok {
				result[partitionID][openOrder.TradePair()] = []broker.OpenOrder{}
			}
			openOrders := result[partitionID][openOrder.TradePair()]

			result[partitionID][openOrder.TradePair()] = append(openOrders, openOrder)
		}

		req.done <- result
	}
}

type (
	tradeParamLoader interface {
		paramCsvHeader() []string
		paramCsvValue() []string
	}

	// TradableTimeZone is a struct for tradable time zone
	TradableTimeZone struct {
		Name string
		OK   func(timeExtractor broker.TimeExtractor) bool
	}

	tradableTimeZones map[PartitionID]*TradableTimeZone

	keyPartitionIDTradePair struct {
		partitionID PartitionID
		tradePair   broker.TradePair
	}

	tradeSpecs struct {
		timeZones    tradableTimeZones
		paramLoaders map[keyPartitionIDTradePair]tradeParamLoader
	}
)

// SimpleTraderBuilder is a builder for simple trader
type SimpleTraderBuilder struct {
	tradableTimeZoneMap map[PartitionID]*TradableTimeZone
	algorithmMap        map[PartitionID]map[broker.TradePair]SimpleTradeAlgorithm
	broker              broker.Broker
	ordererFactory      broker.OrdererFactory
	parallel            int
}

// NewSimpleTraderBuilder is a constructor for simple trader builder
func NewSimpleTraderBuilder() *SimpleTraderBuilder {
	return &SimpleTraderBuilder{
		tradableTimeZoneMap: map[PartitionID]*TradableTimeZone{},
		algorithmMap:        map[PartitionID]map[broker.TradePair]SimpleTradeAlgorithm{},
	}
}

// TradableTimeZone is a method to add a tradable time zone
func (builder *SimpleTraderBuilder) TradableTimeZone(partitionID PartitionID, tradableTimeZone *TradableTimeZone) *SimpleTraderBuilder {
	if _, ok := builder.tradableTimeZoneMap[partitionID]; ok {
		panic("duplicate tradable time zone for an partition ID detected")
	}

	builder.tradableTimeZoneMap[partitionID] = tradableTimeZone
	return builder
}

// Trade is a method to add a trade piece in builder
func (builder *SimpleTraderBuilder) Trade(partitionID PartitionID, tradePair broker.TradePair, algorithm SimpleTradeAlgorithm) *SimpleTraderBuilder {
	_, ok := builder.algorithmMap[partitionID]
	if !ok {
		builder.algorithmMap[partitionID] = map[broker.TradePair]SimpleTradeAlgorithm{}
	}

	builder.algorithmMap[partitionID][tradePair] = algorithm
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
	orderer := builder.ordererFactory.Create(builder.broker)
	orderPartitionAggregator := newOrderPartitionAggregator()
	orderAggregatorFactory := newOrderAggregatorFactory(orderer, orderPartitionAggregator)

	tradeRunners := []*simpleTradeRunner{}
	for partitionID, algorithmMap := range builder.algorithmMap {
		tradableTimeZone, ok := builder.tradableTimeZoneMap[partitionID]
		if !ok {
			panic("no tradable time zone specified")
		}

		tradeRunners = append(tradeRunners, &simpleTradeRunner{
			partitionID:            partitionID,
			tradableTimeZone:       tradableTimeZone,
			algorithmMap:           algorithmMap,
			orderAggregatorFactory: orderAggregatorFactory,
			done:                   make(chan *combinedOrders, 1),
		})
	}

	executor := newSimpleTradeRunnersExecutor(tradeRunners, builder.parallel)

	trader := &SimpleTrader{
		orderer:                  orderer,
		orderPartitionAggregator: orderPartitionAggregator,
		tradeRunners:             tradeRunners,
		executor:                 executor,
	}

	orderPartitionAggregator.run()
	return trader
}

func (builder *SimpleTraderBuilder) buildTradeSpecs() *tradeSpecs {
	paramLoaders := map[keyPartitionIDTradePair]tradeParamLoader{}

	for partitionID, algorithmMap := range builder.algorithmMap {
		for tradePair, algorithm := range algorithmMap {
			paramLoaders[keyPartitionIDTradePair{partitionID, tradePair}] = algorithm
		}
	}

	return &tradeSpecs{
		timeZones:    builder.tradableTimeZoneMap,
		paramLoaders: paramLoaders,
	}
}

func (builder *SimpleTraderBuilder) buildTradableTimeZones() tradableTimeZones {
	return builder.tradableTimeZoneMap
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
