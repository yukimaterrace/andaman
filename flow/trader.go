package flow

import (
	"yukimaterrace/andaman/broker"
	"yukimaterrace/andaman/model"
)

type (
	// Trader is an interface for trader
	Trader interface {
		Trade(material TradeMaterial, mode TradeMode) (RecordMaterial, bool)
	}

	// TraderFactory provides factory method of trader
	TraderFactory interface {
		Create(broker broker.Broker, ordererFactory broker.OrdererFactory) Trader
	}

	// TradeMaterial is an interface for trade material
	TradeMaterial interface{}

	tradeRequest struct {
		material TradeMaterial
	}

	changeTradeModeRequest struct {
		mode TradeMode
	}
)

type tradeWorker struct {
	Trader
	*recordWorker
	mode TradeMode
	ch   chan interface{}
}

func (tradeWorker *tradeWorker) tradeRequest(material TradeMaterial) {
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
		recordMaterial, ok := tradeWorker.Trade(req.material, tradeWorker.mode)

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

const (
	watchString    = "watch"
	tradeString    = "trade"
	terminateStrig = "terminate"
)

func (t *TradeMode) String() string {
	switch *t {
	case Watch:
		return watchString
	case Trade:
		return tradeString
	case Terminate:
		return terminateStrig
	default:
		return model.UnknownString
	}
}

// UnmarshalParam is a method to unmarshal param for trade mode
func (t *TradeMode) UnmarshalParam(param string) error {
	switch param {
	case watchString:
		*t = Watch
	case tradeString:
		*t = Trade
	case terminateStrig:
		*t = Terminate
	default:
		return model.ErrUnknownType
	}
	return nil
}

// MarshalJSON is a method to marshal JSON for trade mode
func (t *TradeMode) MarshalJSON() ([]byte, error) {
	return model.MarshalJSON(t)
}

// UnmarshalJSON is a method to unmarshal JSON for trade mode
func (t *TradeMode) UnmarshalJSON(b []byte) error {
	return t.UnmarshalParam(string(b))
}
