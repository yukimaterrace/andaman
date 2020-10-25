package recorder

import (
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/trader"
)

// SimpleRecorderFactory is a factory for recorder using simple writer
type SimpleRecorderFactory struct {
	builder *trader.Builder
}

// NewSimpleRecorderFactory is a constructor for simple recorder factory
func NewSimpleRecorderFactory(builder *trader.Builder) *SimpleRecorderFactory {
	return &SimpleRecorderFactory{builder}
}

// Create is a factory method to create recorder
func (factory *SimpleRecorderFactory) Create() flow.Recorder {
	tradableTimeZones := factory.builder.BuildTradableTimeZones()
	return newRecorder(newSimpleWriter(tradableTimeZones))
}

// TradePairSummaryRecorderFactory is a factory for simple trader using trade summary recorder
type TradePairSummaryRecorderFactory struct {
	builder *trader.Builder
}

// NewTradePairSummaryRecorderFactory is a constructor for simple trade pair summary recorder
func NewTradePairSummaryRecorderFactory(builder *trader.Builder) *TradePairSummaryRecorderFactory {
	return &TradePairSummaryRecorderFactory{builder}
}

// Create is a factory method to create trade pair summary recorder factory
func (factory *TradePairSummaryRecorderFactory) Create() flow.Recorder {
	tradeSpecs := factory.builder.BuildTradeSpecs()
	return newRecorder(newTradePairSummaryWriter(tradeSpecs))
}
