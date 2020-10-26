package trader

import (
	"log"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/config"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/util"
)

// Trader is a struct for trader
type Trader struct {
	orderer                  broker.Orderer
	orderPartitionAggregator *orderPartitionAggregator
	tradeRunners             []*tradeRunner
	executor                 *tradeRunnersExecutor
}

// Trade is a method to trade
func (trader *Trader) Trade(material flow.TradeMaterial, mode flow.TradeMode) (flow.RecordMaterial, bool) {
	if simulationBroker, ok := trader.orderer.(broker.SimulationBroker); ok {
		if price, ok := material.(broker.PriceExtractor); ok {
			simulationBroker.Update(price)
		} else {
			panic(util.ErrWrongType)
		}
	}

	if mode == flow.Watch {
		return nil, false
	}

	openOrdersResult := <-trader.orderer.OpenOrders()
	if openOrdersResult.Err != nil {
		return nil, false
	}

	partitionedOpenOrders := <-trader.orderPartitionAggregator.partitionedOpenOrders(openOrdersResult.OpenOrders)
	trader.executor.run(material, partitionedOpenOrders, mode)

	var partitionCombineOrders []*PartitionCombinedOrder
	for _, runner := range trader.tradeRunners {
		combinedOrders := <-runner.done
		if combinedOrders != nil && (len(combinedOrders.CreatedOrders) > 0 || len(combinedOrders.ClosedOrders) > 0) {
			partitionCombinedOrder := &PartitionCombinedOrder{
				runner.tradeConfiguration,
				combinedOrders,
			}

			partitionCombineOrders = append(partitionCombineOrders, partitionCombinedOrder)
		}
	}

	recordMaterial := RecordMaterial{
		TradeMaterial:           material,
		PartitionCombinedOrders: partitionCombineOrders,
	}

	return recordMaterial, len(recordMaterial.PartitionCombinedOrders) > 0
}

type tradeRunnersExecutor struct {
	runners       []*tradeRunner
	runnersGroups [][]*tradeRunner
	parallel      int
}

func newTradeRunnersExecutor(runners []*tradeRunner, parallel int) *tradeRunnersExecutor {
	if parallel == 0 {
		return &tradeRunnersExecutor{
			runners:       runners,
			runnersGroups: nil,
			parallel:      0,
		}
	}

	if parallel > len(runners) {
		parallel = len(runners)
	}

	runnerGroups := make([][]*tradeRunner, parallel)
	count := len(runners) / parallel
	for i := 0; i < parallel; i++ {
		if i == parallel-1 {
			runnerGroups[i] = runners[i*count:]
		} else {
			runnerGroups[i] = runners[i*count : (i+1)*count]
		}
	}

	return &tradeRunnersExecutor{
		runners:       runners,
		runnersGroups: runnerGroups,
		parallel:      parallel,
	}
}

func (executor *tradeRunnersExecutor) run(material flow.TradeMaterial, partitionedOpenOrders partitionedOpenOrders, mode flow.TradeMode) {
	if executor.parallel == 0 {
		for _, runner := range executor.runners {
			runner.run(material, partitionedOpenOrders, mode)
		}
	} else {
		for i := 0; i < executor.parallel; i++ {
			go func(runners []*tradeRunner) {
				for _, runner := range runners {
					runner.run(material, partitionedOpenOrders, mode)
				}
			}(executor.runnersGroups[i])
		}
	}
}

// TradeAlgorithm is an interface for trade algorithm
type TradeAlgorithm interface {
	initialTrade(material flow.TradeMaterial, aggregator *orderAggregator, tradePair model.TradePair)
	proceedTrade(material flow.TradeMaterial, aggregator *orderAggregator, openOrders []broker.OpenOrder, tradePair model.TradePair)
}

type tradeRunner struct {
	tradeConfigurationKey model.TradeConfigurationKey
	tradeConfiguration    *model.TradeConfigurationDetail
	algorithm             TradeAlgorithm
	orderAggregator       *orderAggregator
	done                  chan *combinedOrders
	openOrdersExisted     bool
}

func (runner *tradeRunner) run(material flow.TradeMaterial, partitionedOpenOrders partitionedOpenOrders, mode flow.TradeMode) {
	timeExtractor, ok := material.(broker.TimeExtractor)
	if !ok {
		panic(util.ErrWrongType)
	}

	tradable := runner.tradeConfiguration.Timezone.OK(timeExtractor.Time())

	if !runner.openOrdersExisted && !tradable {
		runner.done <- nil
		return
	}

	openOrders := partitionedOpenOrders[runner.tradeConfigurationKey]

	if len(openOrders) > 0 {
		runner.openOrdersExisted = true
	} else {
		runner.openOrdersExisted = false
	}

	tradePair := runner.tradeConfiguration.TradePair
	if !runner.openOrdersExisted {
		if mode != flow.Terminate && tradable {
			runner.algorithm.initialTrade(material, runner.orderAggregator, tradePair)
		}
	} else {
		runner.algorithm.proceedTrade(material, runner.orderAggregator, openOrders, tradePair)
	}

	runner.done <- runner.orderAggregator.reduce()
}

type (
	combinedOrders struct {
		CreatedOrders []broker.CreatedOrder
		ClosedOrders  []broker.ClosedOrder
	}

	// PartitionCombinedOrder is a definition for partition combined order
	PartitionCombinedOrder struct {
		tradeConfiguration *model.TradeConfigurationDetail
		*combinedOrders
	}

	// RecordMaterial is a definition for concrete record material
	RecordMaterial struct {
		TradeMaterial           flow.TradeMaterial
		PartitionCombinedOrders []*PartitionCombinedOrder
	}
)

type orderAggregator struct {
	broker.Orderer
	tradeConfigurationKey    model.TradeConfigurationKey
	orderPartitionAggregator *orderPartitionAggregator
	createOrderDone          []<-chan *broker.CreateOrderResult
	closeOrderDone           []<-chan *broker.CloseOrderResult
}

func newOrderAggregator(orderer broker.Orderer, tradeConfigurationKey model.TradeConfigurationKey, orderPartitionAggregator *orderPartitionAggregator) *orderAggregator {
	return &orderAggregator{
		Orderer:                  orderer,
		tradeConfigurationKey:    tradeConfigurationKey,
		orderPartitionAggregator: orderPartitionAggregator,
	}
}

func (aggregator *orderAggregator) createOrder(tradePair model.TradePair, units float64, isLong bool) {
	result := aggregator.CreateOrder(tradePair, units, isLong)
	aggregator.createOrderDone = append(aggregator.createOrderDone, result)
}

func (aggregator *orderAggregator) closeOrder(orderID broker.OrderID) {
	result := aggregator.CloseOrder(orderID)
	aggregator.closeOrderDone = append(aggregator.closeOrderDone, result)
}

func (aggregator *orderAggregator) reduce() *combinedOrders {
	var createdOrders []broker.CreatedOrder
	for _, done := range aggregator.createOrderDone {
		createdOrder := <-done
		if createdOrder.Err != nil {
			log.Println(createdOrder.Err.Error())
		} else {
			createdOrders = append(createdOrders, createdOrder.CreatedOrder)
			aggregator.orderPartitionAggregator.put(createdOrder.CreatedOrder.OrderID(), aggregator.tradeConfigurationKey)
		}
	}

	var closedOrders []broker.ClosedOrder
	for _, done := range aggregator.closeOrderDone {
		closedOrder := <-done
		if closedOrder.Err != nil {
			log.Println(closedOrder.Err.Error())
		} else {
			closedOrders = append(closedOrders, closedOrder.ClosedOrder)
			aggregator.orderPartitionAggregator.delete(closedOrder.ClosedOrder.OrderID())
		}
	}

	aggregator.createOrderDone = nil
	aggregator.closeOrderDone = nil

	return &combinedOrders{
		CreatedOrders: createdOrders,
		ClosedOrders:  closedOrders,
	}
}

type (
	partitionedOpenOrders map[model.TradeConfigurationKey][]broker.OpenOrder

	orderPartitionAggregator struct {
		ch                chan interface{}
		orderPartitionMap map[broker.OrderID]model.TradeConfigurationKey
	}
)

func newOrderPartitionAggregator() *orderPartitionAggregator {
	return &orderPartitionAggregator{
		ch:                make(chan interface{}, config.OrderPartitionAggregatorChanCap),
		orderPartitionMap: map[broker.OrderID]model.TradeConfigurationKey{},
	}
}

type (
	putOrderPartitionRequest struct {
		orderID               broker.OrderID
		tradeConfigurationkey model.TradeConfigurationKey
	}

	deleteOrderPartitionRequest struct {
		orderID broker.OrderID
	}

	partitionedOpenOrdersRequest struct {
		openOrders []broker.OpenOrder
		done       chan<- partitionedOpenOrders
	}
)

func (aggregator *orderPartitionAggregator) put(orderID broker.OrderID, tradeConfigurationKey model.TradeConfigurationKey) {
	aggregator.ch <- &putOrderPartitionRequest{
		orderID:               orderID,
		tradeConfigurationkey: tradeConfigurationKey,
	}
}

func (aggregator *orderPartitionAggregator) delete(orderID broker.OrderID) {
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
		aggregator.orderPartitionMap[req.orderID] = req.tradeConfigurationkey

	case *deleteOrderPartitionRequest:
		if _, ok := aggregator.orderPartitionMap[req.orderID]; !ok {
			panic("no order id detected")
		}
		delete(aggregator.orderPartitionMap, req.orderID)

	case *partitionedOpenOrdersRequest:
		result := map[model.TradeConfigurationKey][]broker.OpenOrder{}

		for _, openOrder := range req.openOrders {
			key, ok := aggregator.orderPartitionMap[openOrder.OrderID()]
			if !ok {
				panic("no order id detected")
			}

			openOrders := result[key]
			result[key] = append(openOrders, openOrder)
		}

		req.done <- result
	}
}

// Factory is a factory for simple trader
type Factory struct {
	builder *Builder
}

// NewFactory is a constructor for simple trader factory
func NewFactory(builder *Builder) *Factory {
	return &Factory{
		builder: builder,
	}
}

// Create is a factory method to create trader
func (factory *Factory) Create(broker broker.Broker, ordererFactory broker.OrdererFactory) flow.Trader {
	return factory.builder.Broker(broker).OrdererFactory(ordererFactory).Build()
}
