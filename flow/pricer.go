package flow

import (
	"errors"
	"log"
	"yukimaterrace/andaman/broker"
)

type pricer interface {
	createPrice(done chan<- *createPriceResult)
}

// PricerFactory provide factory method of pricer
type PricerFactory interface {
	create(broker broker.Broker, tradePairs []broker.TradePair) pricer
}

type createPriceResult struct {
	tradeMaterial tradeMaterial
	err           error
}

var errNoMorePrice = errors.New("no more price")

type priceWorker struct {
	pricer
	*tradeWorker
	createPriceResult chan *createPriceResult
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
				priceWorker.createPrice(priceWorker.createPriceResult)
			}
		}()

		priceWorker.init = true
	}

	select {
	case result := <-priceWorker.createPriceResult:
		if result.err == nil {
			priceWorker.tradeRequest(result.tradeMaterial)
		} else {
			log.Println(result.err.Error())

			if result.err == errNoMorePrice {
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
