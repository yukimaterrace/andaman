package app

import (
	"yukimaterrace/andaman/market"
	"yukimaterrace/andaman/recorder"
	"yukimaterrace/andaman/trade"
)

// CreateApp is a factory method for app
func CreateApp() *App {
	market := market.CreateOandaMarket()

	recorders := make([]recorder.Recorder, 0)
	traders := make([]trade.Trader, 0)
	for _, instrument := range tradeInstruments {
		r := recorder.CreateSimpleFileRecorder(instrument)
		t := trade.CreateZoneTrader(instrument, market, r)

		recorders = append(recorders, r)
		traders = append(traders, t)
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
