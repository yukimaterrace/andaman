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

// Type is indicator type
type Type int

const (
	// SMA is simple moving average type
	SMA Type = iota
)

func (t Type) String() string {
	switch t {
	case SMA:
		return "SMA"
	default:
		return "Unknown"
	}
}
