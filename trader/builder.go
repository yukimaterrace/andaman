package trader

import (
	"encoding/json"
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/service"
	"yukimaterrace/andaman/util"
)

// Builder is a builder for trader
type Builder struct {
	tradeSetName   string
	tradeRunType   model.TradeRunType
	broker         broker.Broker
	ordererFactory broker.OrdererFactory
	parallel       int
}

// NewBuilder is a constructor for trader builder
func NewBuilder() *Builder {
	return &Builder{}
}

// TradeSet sets trade set in builder
func (builder *Builder) TradeSet(name string) *Builder {
	builder.tradeSetName = name
	return builder
}

// TradeRunType sets trade run type in builder
func (builder *Builder) TradeRunType(_type model.TradeRunType) *Builder {
	builder.tradeRunType = _type
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
	tradeSet, err := service.GetTradeSetDetail(builder.tradeSetName, tradeParamObjectCreator)
	if err != nil {
		panic(err)
	}

	orderer := builder.ordererFactory.Create(builder.broker)
	orderPartitionAggregator := newOrderPartitionAggregator()

	var tradeRunners []*tradeRunner
	for _, tradeConfiguration := range tradeSet.Configurations {
		algorithm, err := createTradeAlgorithm(tradeConfiguration.Algorithm)
		if err != nil {
			panic(err)
		}

		tradeConfigurationKey := tradeConfiguration.Key()
		orderAggregator := newOrderAggregator(orderer, tradeConfigurationKey, orderPartitionAggregator)

		tradeRunners = append(tradeRunners, &tradeRunner{
			tradeConfigurationKey: tradeConfigurationKey,
			tradeConfiguration:    tradeConfiguration,
			algorithm:             algorithm,
			orderAggregator:       orderAggregator,
			done:                  make(chan *combinedOrders, 1),
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

// BuildTradeRun is a method to build trade run
func (builder *Builder) BuildTradeRun() *model.TradeRun {
	tradeSet, err := service.AddTradeRun(builder.tradeSetName, builder.tradeRunType)
	if err != nil {
		panic(err)
	}
	return tradeSet
}

func tradeParamObjectCreator(_type model.TradeAlgorithmType, param string) (interface{}, error) {
	switch _type {
	case model.Frame:
		p := FrameTradeParam{}
		if err := json.Unmarshal([]byte(param), &param); err != nil {
			return nil, err
		}
		return &p, nil

	default:
		return nil, util.ErrWrongType
	}
}

func createTradeAlgorithm(algorithm model.TradeAlgorithmDetail) (TradeAlgorithm, error) {
	switch algorithm.Type {
	case model.Frame:
		param, ok := algorithm.ParamObject.(*FrameTradeParam)
		if !ok {
			return nil, util.ErrWrongType
		}
		return NewFrameTradeAlgorithm(param), nil

	default:
		return nil, util.ErrWrongType
	}
}
