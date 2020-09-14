package flow

import (
	"yukimaterrace/andaman/broker"
)

type trader interface {
	trade(material tradeMaterial, mode TradeMode) (recordMaterial, bool)
}

// TraderFactory provides factory method of trader
type TraderFactory interface {
	create(broker broker.Broker, ordererFactory broker.OrdererFactory) trader
}

type tradeMaterial interface{}

type (
	tradeRequest struct {
		material tradeMaterial
	}

	changeTradeModeRequest struct {
		mode TradeMode
	}
)

type tradeWorker struct {
	trader
	*recordWorker
	mode TradeMode
	ch   chan interface{}
}

func (tradeWorker *tradeWorker) tradeRequest(material tradeMaterial) {
	tradeWorker.ch <- &tradeRequest{material: material}
}

func (tradeWorker *tradeWorker) changeTradeModeRequest(mode TradeMode) {
	tradeWorker.ch <- &changeTradeModeRequest{mode}
}

func (tradeWorker *tradeWorker) shutdown() {
	req := newShutdownRequest()
	tradeWorker.ch <- req
	<-req.done
}

func (tradeWorker *tradeWorker) work(exit chan<- bool) {
	request := <-tradeWorker.ch

	switch req := request.(type) {
	case *tradeRequest:
		recordMaterial, ok := tradeWorker.trade(req.material, tradeWorker.mode)

		if ok {
			tradeWorker.recordRequest(recordMaterial)
		}

	case *changeTradeModeRequest:
		tradeWorker.mode = req.mode

	case *shutdownRequest:
		tradeWorker.recordWorker.shutdown()
		req.done <- true
		exit <- true
	}
}

// TradeMode is iota value for trade mode
type TradeMode int

const (
	// Watch is watch mode
	Watch TradeMode = iota
	// Trade is trade mode
	Trade
	// Terminate is terminate mode
	Terminate
)
