package app

import (
	"yukimaterrace/andaman/market"
	"yukimaterrace/andaman/recorder"
	"yukimaterrace/andaman/trade"
)

// App is an application structure for Andaman
type App struct {
	market    market.Market
	traders   []trade.Trader
	recorders []recorder.Recorder
}

// Start is a method to start Andaman
func (app *App) Start() {
	app.market.Start()

	for _, recorder := range app.recorders {
		recorder.Start()
	}

	for _, trader := range app.traders {
		trader.Start()
	}
}
