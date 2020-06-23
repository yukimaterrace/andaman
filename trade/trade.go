package trade

import (
	"yukimaterrace/andaman/market"
	"yukimaterrace/andaman/recorder"
)

// Trader is an interface for trader
type Trader interface {
	Start()
}

type routine struct {
	instrument market.Instrument
	market     market.Market
	recorder   recorder.Recorder
}

type orderer interface {
	order()

	close()
}

func (routine *routine) Start() {
	go routine.run()
}

func (routine *routine) run() {

}
