package flow

import (
	"errors"
	"log"
	"yukimaterrace/andaman/broker"
)

// Pricer is an interface for pricer
type Pricer interface {
	CreatePrice(done chan<- *CreatePriceResult)
}

// PricerFactory is a factory of pricer
type PricerFactory interface {
	Create(broker broker.Broker, tradePairs []broker.TradePair) Pricer
}

// ErrNoMorePrice is an error for no more price
var ErrNoMorePrice = errors.New("no more price")

// CreatePriceResult is a result for create price
type CreatePriceResult struct {
	TradeMaterial TradeMaterial
	Err           error
}

type priceWorker struct {
	Pricer
	*tradeWorker
	createPriceResult chan *CreatePriceResult
	ch                chan interface{}
	init              bool
}

func (priceWorker *priceWorker) shutdown() {
	req := newShutdownRequest()
	priceWorker.ch <- req
	<-req.done
}

func (priceWorker *priceWorker) work(exit chan<- bool) {
	if !priceWorker.init {
		go func() {
			for {
				priceWorker.CreatePrice(priceWorker.createPriceResult)
			}
		}()

		priceWorker.init = true
	}

	select {
	case result := <-priceWorker.createPriceResult:
		if result.Err == nil {
			priceWorker.tradeRequest(result.TradeMaterial)
		} else {
			log.Println(result.Err.Error())

			if result.Err == ErrNoMorePrice {
				priceWorker.tradeWorker.shutdown()
				exit <- true
			}
		}

	case request := <-priceWorker.ch:
		switch req := request.(type) {
		case *shutdownRequest:
			priceWorker.tradeWorker.shutdown()
			req.done <- true
			exit <- true
		}
	}
}
