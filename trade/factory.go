package trade

import (
	"yukimaterrace/andaman/market"
	"yukimaterrace/andaman/record"
)

// CreateZoneTrader is a factory method for zone trade strategy
func CreateZoneTrader(instrument market.Instrument, market market.Market, recorder record.Recorder) Trader {
	return newRoutine(instrument, market, recorder, &zoneTradeStrategy{})
}
