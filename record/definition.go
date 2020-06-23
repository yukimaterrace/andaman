package record

import (
	"yukimaterrace/andaman/indicate"
	"yukimaterrace/andaman/market"
)

type (
	// Material is a definition for record material
	Material struct {
		Instrument  market.Instrument
		Prices      *market.PriceSequence
		Indicators  []*indicate.Value
		MadeOrder   *market.MadeOrder
		ClosedOrder []*market.ClosedOrder
	}
)
