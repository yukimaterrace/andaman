package indicator

import "yukimaterrace/andaman/market"

// Indicator is an interface for indicator
type Indicator interface {
	Indicate(priceSequence *market.PriceSequence) *Value
}

type calculator interface {
	Type() Type
	calculate(priceSequence *market.PriceSequence) *Value
}

type routine struct {
	cache   *Value
	archive *market.PriceSequence
	calculator
}

func (routine *routine) compare(new *market.PriceSequence) bool {
	old := routine.archive

	condInstrument := old.Instrument == new.Instrument
	condGranularity := old.Granularity == new.Granularity
	condType := old.Type == new.Type
	condPrices := len(old.Prices) > 0 && len(old.Prices) == len(new.Prices) && old.Prices[0].Time == new.Prices[0].Time

	return condInstrument && condGranularity && condType && condPrices
}

func (routine *routine) Indicate(priceSequence *market.PriceSequence) *Value {
	if routine.compare(priceSequence) {
		return routine.cache
	}

	routine.archive = priceSequence
	value := routine.calculate(priceSequence)
	routine.cache = value

	return value
}
