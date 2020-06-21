package market

import "yukimaterrace/andaman/config"

// Market is an interface for market
type Market interface {
	Start()

	Prices(instrument Instrument, granularity Granularity, count int, from int64) <-chan *PriceSequence

	LatestPrice(instrument Instrument) <-chan *PriceDetail

	Orders(instrument Instrument) <-chan *Orders

	Asset() <-chan *AssetStatus

	MakeOrder(instrument Instrument, orderType OrderType, unit float64) <-chan *MakeOrderStatus

	CloseOrder(orderID string) <-chan *CloseOrderStatus
}

type adaptor interface {
	prices(instrument Instrument, granularity Granularity, count int, from int64) *PriceSequence

	latestPrice(instrument Instrument) *PriceDetail

	orders(instrument Instrument) *Orders

	asset() *AssetStatus

	makeOrder(instrument Instrument, orderType OrderType, unit float64) *MakeOrderStatus

	closeOrder(orderID string) *CloseOrderStatus
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
func (routine *routine) Prices(instrument Instrument, granularity Granularity, count int, from int64) <-chan *PriceSequence {
	replyTo := make(chan *PriceSequence, 1)

	routine.request <- &pricesRequest{
		instrument:  instrument,
		granularity: granularity,
		count:       count,
		from:        from,
		replyTo:     replyTo,
	}
	return replyTo
}

// LatestPrice is a method for latest price request
func (routine *routine) LatestPrice(instrument Instrument) <-chan *PriceDetail {
	replyTo := make(chan *PriceDetail, 1)

	routine.request <- &latestPriceRequest{
		instrument: instrument,
		replyTo:    replyTo,
	}
	return replyTo
}

// Orders is a method for order request
func (routine *routine) Orders(instrument Instrument) <-chan *Orders {
	replyTo := make(chan *Orders, 1)

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

func (routine *routine) MakeOrder(instrument Instrument, orderType OrderType, unit float64) <-chan *MakeOrderStatus {
	replyTo := make(chan *MakeOrderStatus, 1)

	routine.request <- &makeOrderRequest{
		instrument: instrument,
		orderType:  orderType,
		unit:       unit,
		replyTo:    replyTo,
	}
	return replyTo
}

func (routine *routine) CloseOrder(orderID string) <-chan *CloseOrderStatus {
	replyTo := make(chan *CloseOrderStatus, 1)

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
			req.replyTo <- routine.prices(req.instrument, req.granularity, req.count, req.from)

		case *latestPriceRequest:
			req.replyTo <- routine.latestPrice(req.instrument)

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
