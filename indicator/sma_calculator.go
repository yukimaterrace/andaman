package indicator

import "yukimaterrace/andaman/market"

type smaCalculator struct {
	t       Type
	width   int
	archive []*Element
}

func newSmaCalculator(valueType Type) *smaCalculator {
	var width int
	switch valueType {
	case SMA25:
		width = 25
	case SMA75:
		width = 75
	case SMA150:
		width = 150
	case SMA600:
		width = 600
	default:
		panic("Unknown indicator type")
	}

	return &smaCalculator{
		t:     valueType,
		width: width,
	}
}

func (calculator *smaCalculator) calculate(prices []market.Price) []*Element {
	return nil
}

func (calculator *smaCalculator) valueType() Type {
	return calculator.t
}
