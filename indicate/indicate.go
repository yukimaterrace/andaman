package indicate

import "yukimaterrace/andaman/market"

// Indicator is an interface for indicator
type Indicator interface {
	Type() Type
	Indicate(priceSequence *market.PriceSequence) *Value
}

type calculator interface {
	calculate(prices []*market.Price) []*Element
	valueType() Type
}

type routine struct {
	calculator
}

func newRoutine(calculator calculator) *routine {
	return &routine{
		calculator: calculator,
	}
}

func (routine *routine) Type() Type {
	return routine.valueType()
}

func (routine *routine) Indicate(priceSequence *market.PriceSequence) *Value {
	elements := routine.calculate(priceSequence.Prices)

	return &Value{
		Instrument:  priceSequence.Instrument,
		Granularity: priceSequence.Granularity,
		Type:        routine.valueType(),
		Elements:    elements,
	}
}
