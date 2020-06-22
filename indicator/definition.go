package indicator

import "yukimaterrace/andaman/market"

type (
	// Value is a definition for indicator value
	Value struct {
		Instrument  market.Instrument
		Granularity market.Granularity
		Type        Type
		Elements    []*Element
	}

	// Element is a definition for indicator element
	Element struct {
		ElementType string
		Pieces      []*Piece
	}

	// Piece is a definition for indicator piece
	Piece struct {
		Value float64
		Time  int64
	}
)

type (
	// Type is indicator type
	Type interface {
		String() string
	}

	// SMA is simple moving average type
	SMA Type
)

func (sma SMA) String() string {
	return "SMA"
}
