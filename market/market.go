package market

import "yukimaterrace/andaman/config"

// Market is an interface for market
type Market interface {
	Start()
	Pricer
	Orderer
}

// Pricer is an interface for pricer
type Pricer interface {
	Prices(instrument Instrument, granularity Granularity, count int) <-chan *PriceSequenceStatus

	TradePrice(instrument Instrument) <-chan *TradePriceStatus
}

// Orderer is an interface for orderer
type Orderer interface {
	Orders(instrument Instrument) <-chan *OrdersStatus

	Asset() <-chan *AssetStatus

	MakeOrder(instrument Instrument, orderType OrderType, unit float64) <-chan *MadeOrderStatus

	CloseOrder(orderID string) <-chan *ClosedOrderStatus
}

type adaptor interface {
	prices(instrument Instrument, granularity Granularity, count int) *PriceSequenceStatus

	tradePrice(instrument Instrument) *TradePriceStatus

	orders(instrument Instrument) *OrdersStatus

	asset() *AssetStatus

	makeOrder(instrument Instrument, orderType OrderType, unit float64) *MadeOrderStatus

	closeOrder(orderID string) *ClosedOrderStatus
}

type routine struct {
	request chan request
	adaptor
}

func newRoutine(adaptor adaptor) *routine {
	return &routine{
		make(chan request, config.MarketChanCapacity),
		adaptor,
	}
}

// Start is a method to start routine
func (routine *routine) Start() {
	go routine.run()
}

// Prices is a method for prices request
func (routine *routine) Prices(instrument Instrument, granularity Granularity, count int) <-chan *PriceSequenceStatus {
	replyTo := make(chan *PriceSequenceStatus, 1)

	routine.request <- &pricesRequest{
		instrument:  instrument,
		granularity: granularity,
		count:       count,
		replyTo:     replyTo,
	}
	return replyTo
}

// TradePrice is a method for latest price request
func (routine *routine) TredePrice(instrument Instrument) <-chan *TradePriceStatus {
	replyTo := make(chan *TradePriceStatus, 1)

	routine.request <- &tradePriceRequest{
		instrument: instrument,
		replyTo:    replyTo,
	}
	return replyTo
}

// Orders is a method for order request
func (routine *routine) Orders(instrument Instrument) <-chan *OrdersStatus {
	replyTo := make(chan *OrdersStatus, 1)

	routine.request <- &ordersRequest{
		instrument: instrument,
		replyTo:    replyTo,
	}
	return replyTo
}

func (routine *routine) Asset() <-chan *AssetStatus {
	replyTo := make(chan *AssetStatus, 1)

	routine.request <- &assetRequest{
		replyTo: replyTo,
	}
	return replyTo
}

func (routine *routine) MakeOrder(instrument Instrument, orderType OrderType, unit float64) <-chan *MadeOrderStatus {
	replyTo := make(chan *MadeOrderStatus, 1)

	routine.request <- &makeOrderRequest{
		instrument: instrument,
		orderType:  orderType,
		unit:       unit,
		replyTo:    replyTo,
	}
	return replyTo
}

func (routine *routine) CloseOrder(orderID string) <-chan *ClosedOrderStatus {
	replyTo := make(chan *ClosedOrderStatus, 1)

	routine.request <- &closerOrderRequest{
		orderID: orderID,
		replyTo: replyTo,
	}
	return replyTo
}

func (routine *routine) run() {
	for {
		req := <-routine.request

		switch req := req.(type) {
		case *pricesRequest:
			req.replyTo <- routine.prices(req.instrument, req.granularity, req.count)

		case *tradePriceRequest:
			req.replyTo <- routine.tradePrice(req.instrument)

		case *ordersRequest:
			req.replyTo <- routine.orders(req.instrument)

		case *assetRequest:
			req.replyTo <- routine.asset()

		case *closerOrderRequest:
			req.replyTo <- routine.closeOrder(req.orderID)

		default:
			panic("request type not allowed")
		}
	}
}
