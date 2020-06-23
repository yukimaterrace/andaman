package recorder

import (
	"yukimaterrace/andaman/indicator"
	"yukimaterrace/andaman/market"
)

type (
	// Material is a definition for record material
	Material struct {
		Instrument  market.Instrument
		Prices      *market.PriceSequence
		Indicators  []*indicator.Value
		MadeOrder   *market.MadeOrder
		ClosedOrder []*market.ClosedOrder
	}
)
