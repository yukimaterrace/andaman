package recorder

import "yukimaterrace/andaman/trader"

// Factory is a struct for recorder factory
type Factory struct {
	traderBuilder *trader.Builder
}

// Create is a factory method to create recorder
func (factory *Factory) Create() *Recorder {
	return newRecorder(factory.traderBuilder.BuildTradeRun())
}

// NewFactory is a constructor for factory
func NewFactory(traderBuilder *trader.Builder) *Factory {
	return &Factory{traderBuilder: traderBuilder}
}
