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
		results = append(results, accountOrders...)
	}

	return results, len(results) > 0
}

// SimpleTradeAlgorithm is an interface for simple trade algorithm
type SimpleTradeAlgorithm interface {
	initialTrade(material tradeMaterial, orderer broker.Orderer, accountID broker.AccountID, tradePair broker.TradePair) (*accountOrders, error)
	proceedTrade(material tradeMaterial, orderer broker.Orderer, openOrders []broker.OpenOrder, accountID broker.AccountID, tradePair broker.TradePair) (*accountOrders, error)
}

type simpleTradeRunner struct {
	accountID    broker.AccountID
	algorithmMap map[broker.TradePair]SimpleTradeAlgorithm
	orderer      broker.Orderer
	done         chan []*accountOrders
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

	results := make([]*accountOrders, 0)
	for tradePair, openOrders := range openOrdersMap {
		algorithm, ok := runner.algorithmMap[tradePair]
		if !ok {
			log.Printf("no algorithm registered for %v\n", tradePair)
		}

		var accountOrders *accountOrders
		var err error
		if len(openOrders) == 0 {
			if mode != Terminate && tradable {
				accountOrders, err = algorithm.initialTrade(material, runner.orderer, runner.accountID, tradePair)
			}
		} else {
			accountOrders, err = algorithm.proceedTrade(material, runner.orderer, openOrders, runner.accountID, tradePair)
		}

		if err != nil {
			log.Println(err.Error())
		} else if accountOrders != nil {
			results = append(results, accountOrders)
		}
	}

	for _, result := range results {
		for _, createdOrder := range result.createdOrders {
			createdOrder.createdOrder = <-createdOrder.done
		}

		for _, closedOrder := range result.closedOrders {
			closedOrder.closedOrder = <-closedOrder.done
		}
	}

	runner.done <- results
}

type accountCreatedOrder struct {
	accountID    broker.AccountID
	createdOrder broker.CreatedOrder
	done         <-chan broker.CreatedOrder
}

type accountClosedOrder struct {
	accountID   broker.AccountID
	closedOrder broker.ClosedOrder
	done        <-chan broker.ClosedOrder
}

type accountOrders struct {
	createdOrders []*accountCreatedOrder
	closedOrders  []*accountClosedOrder
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
			done:         make(chan []*accountOrders, 1),
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
