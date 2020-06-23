package app

import (
	"yukimaterrace/andaman/market"
	"yukimaterrace/andaman/record"
	"yukimaterrace/andaman/trade"
)

// CreateApp is a factory method for app
func CreateApp() *App {
	market := market.CreateOandaMarket()

	recorders := make([]record.Recorder, 0)
	traders := make([]trade.Trader, 0)
	for _, instrument := range tradeInstruments {
		recorder := record.CreateSimpleFileRecorder(instrument)
		trader := trade.CreateZoneTrader(instrument, market, r)

		recorders = append(recorders, recorder)
		traders = append(traders, trader)
	}

	return &App{
		market:    market,
		traders:   traders,
		recorders: recorders,
	}
}

// CreatePracticeApp is a factory method for practice app
func CreatePracticeApp() *App {
	return nil
}
