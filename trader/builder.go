package trader

import "yukimaterrace/andaman/broker"

// Builder is a builder for trader
type Builder struct {
	tradableTimeZoneMap map[PartitionID]*TradableTimeZone
	algorithmMap        map[PartitionID]map[broker.TradePair]TradeAlgorithm
	broker              broker.Broker
	ordererFactory      broker.OrdererFactory
	parallel            int
}

// NewBuilder is a constructor for trader builder
func NewBuilder() *Builder {
	return &Builder{
		tradableTimeZoneMap: map[PartitionID]*TradableTimeZone{},
		algorithmMap:        map[PartitionID]map[broker.TradePair]TradeAlgorithm{},
	}
}

// TradableTimeZone is a method to add a tradable time zone
func (builder *Builder) TradableTimeZone(partitionID PartitionID, tradableTimeZone *TradableTimeZone) *Builder {
	if _, ok := builder.tradableTimeZoneMap[partitionID]; ok {
		panic("duplicate tradable time zone for an partition ID detected")
	}

	builder.tradableTimeZoneMap[partitionID] = tradableTimeZone
	return builder
}

// Trade is a method to add a trade piece in builder
func (builder *Builder) Trade(partitionID PartitionID, tradePair broker.TradePair, algorithm TradeAlgorithm) *Builder {
	_, ok := builder.algorithmMap[partitionID]
	if !ok {
		builder.algorithmMap[partitionID] = map[broker.TradePair]TradeAlgorithm{}
	}

	builder.algorithmMap[partitionID][tradePair] = algorithm
	return builder
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

// Parallel sets parallel parameter in builder
func (builder *Builder) Parallel(paralle int) *Builder {
	builder.parallel = paralle
	return builder
}

// Build builds simple trader
func (builder *Builder) Build() *Trader {
	orderer := builder.ordererFactory.Create(builder.broker)
	orderPartitionAggregator := newOrderPartitionAggregator()
	orderAggregatorFactory := newOrderAggregatorFactory(orderer, orderPartitionAggregator)

	var tradeRunners []*tradeRunner
	for partitionID, algorithmMap := range builder.algorithmMap {
		tradableTimeZone, ok := builder.tradableTimeZoneMap[partitionID]
		if !ok {
			panic("no tradable time zone specified")
		}

		tradeRunners = append(tradeRunners, &tradeRunner{
			partitionID:            partitionID,
			tradableTimeZone:       tradableTimeZone,
			algorithmMap:           algorithmMap,
			orderAggregatorFactory: orderAggregatorFactory,
			done:                   make(chan *combinedOrders, 1),
		})
	}

	executor := newTradeRunnersExecutor(tradeRunners, builder.parallel)

	trader := &Trader{
		orderer:                  orderer,
		orderPartitionAggregator: orderPartitionAggregator,
		tradeRunners:             tradeRunners,
		executor:                 executor,
	}

	orderPartitionAggregator.run()
	return trader
}

// BuildTradeSpecs is a method to build trade specs
func (builder *Builder) BuildTradeSpecs() *TradeSpecs {
	paramLoaders := map[KeyPartitionIDTradePair]TradeParamLoader{}

	for partitionID, algorithmMap := range builder.algorithmMap {
		for tradePair, algorithm := range algorithmMap {
			paramLoaders[KeyPartitionIDTradePair{partitionID, tradePair}] = algorithm
		}
	}

	return &TradeSpecs{
		TimeZones:    builder.tradableTimeZoneMap,
		ParamLoaders: paramLoaders,
	}
}

// BuildTradableTimeZones is a method to build tradable timezones
func (builder *Builder) BuildTradableTimeZones() TradableTimeZones {
	return builder.tradableTimeZoneMap
}
