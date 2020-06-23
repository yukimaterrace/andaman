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
	// SMA25 is 25 points simple moving average type
	SMA25 Type = iota
	// SMA75 is 75 points simple moving average type
	SMA75
	// SMA150 is 150 points simple moving average type
	SMA150
	// SMA600 is 600 points simple moving average type
	SMA600
)

func (t Type) String() string {
	switch t {
	case SMA:
		return "SMA"
	default:
		return "Unknown"
	}
}
