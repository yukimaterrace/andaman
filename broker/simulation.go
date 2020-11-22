package broker

import (
	"log"
	"yukimaterrace/andaman/config"
	"yukimaterrace/andaman/model"
)

// SimpleSimulationBroker is a broker for simple simulation
type SimpleSimulationBroker struct {
	currentPriceMap map[model.TradePair]Price
	currentTime     int64
}

// NewSimpleSimulationBroker is a constructor for simple simulation broker
func NewSimpleSimulationBroker() *SimpleSimulationBroker {
	return &SimpleSimulationBroker{
		currentPriceMap: map[model.TradePair]Price{},
		currentTime:     0,
	}
}

// Update is a method to update orderer
func (broker *SimpleSimulationBroker) Update(priceExtractor PriceExtractor) {
	for _, tradePair := range priceExtractor.TradePairs() {
		broker.currentPriceMap[tradePair] = priceExtractor.Price(tradePair)
	}

	broker.currentTime = priceExtractor.Time()
}

// SimpleSimulationOrderer is a struct for simple simulation broker
type SimpleSimulationOrderer struct {
	*SimpleSimulationBroker
	currentOrderID OrderID
	currentOrders  []*order
	ch             chan interface{}
}

func (orderer *SimpleSimulationOrderer) price(tradePair model.TradePair) Price {
	price, ok := orderer.currentPriceMap[tradePair]
	if !ok {
		log.Panicf("%s cannot handle in this orderer\n", tradePair.String())
	}
	return price
}

type (
	createOrderRequest struct {
		tradePair model.TradePair
		units     float64
		isLong    bool
		done      chan<- *CreateOrderResult
	}

	openOrdersRequest struct {
		done chan<- *OpenOrdersResult
	}

	closeOrderRequest struct {
		orderID OrderID
		done    chan<- *CloseOrderResult
	}
)

// CreateOrder is a method to create order
func (orderer *SimpleSimulationOrderer) CreateOrder(tradePair model.TradePair, units float64, isLong bool) <-chan *CreateOrderResult {
	done := make(chan *CreateOrderResult, 1)

	orderer.ch <- &createOrderRequest{
		tradePair: tradePair,
		units:     units,
		isLong:    isLong,
		done:      done,
	}

	return done
}

// OpenOrders is a method to open orders
func (orderer *SimpleSimulationOrderer) OpenOrders() <-chan *OpenOrdersResult {
	done := make(chan *OpenOrdersResult, 1)

	orderer.ch <- &openOrdersRequest{
		done: done,
	}

	return done
}

// CloseOrder is a method to close order
func (orderer *SimpleSimulationOrderer) CloseOrder(orderID OrderID) <-chan *CloseOrderResult {
	done := make(chan *CloseOrderResult, 1)

	orderer.ch <- &closeOrderRequest{
		orderID: orderID,
		done:    done,
	}

	return done
}

func (orderer *SimpleSimulationOrderer) createOrder(tradePair model.TradePair, units float64, isLong bool) CreatedOrder {
	price := orderer.price(tradePair)

	order := &order{
		orderID:    orderer.currentOrderID,
		tradePair:  tradePair,
		units:      units,
		isLong:     isLong,
		timeAtOpen: orderer.currentTime,
	}

	if isLong {
		order.priceAtOpen = price.Ask()
	} else {
		order.priceAtOpen = price.Bid()
	}

	order.unrealizedProfit = profitPips(price, order)

	orderer.currentOrderID++
	orderer.currentOrders = append(orderer.currentOrders, order)

	return order
}

func (orderer *SimpleSimulationOrderer) openOrders() []OpenOrder {
	openOrders := make([]OpenOrder, len(orderer.currentOrders))

	for i, order := range orderer.currentOrders {
		price := orderer.price(order.tradePair)
		order.unrealizedProfit = profitPips(price, order)
		openOrders[i] = order
	}

	return openOrders
}

func (orderer *SimpleSimulationOrderer) closeOrder(orderID OrderID) ClosedOrder {
	var order *order
	var pos int
	for i, o := range orderer.currentOrders {
		if o.orderID == orderID {
			order = o
			pos = i
			break
		}
	}

	if order == nil {
		log.Panicf("no open order exists: orderID %d", int(orderID))
	}

	price := orderer.price(order.tradePair)

	order.timeAtClose = orderer.currentTime

	if order.isLong {
		order.priceAtClose = price.Bid()
	} else {
		order.priceAtClose = price.Ask()
	}

	order.realizedProfit = profitPips(price, order)

	currentOrders := orderer.currentOrders
	orderer.currentOrders = append(currentOrders[:pos], currentOrders[pos+1:]...)

	return order
}

func (orderer *SimpleSimulationOrderer) run() {
	go func() {
		for {
			orderer.work()
		}
	}()
}

func (orderer *SimpleSimulationOrderer) work() {
	request := <-orderer.ch

	switch req := request.(type) {
	case *createOrderRequest:
		createdOrder := orderer.createOrder(req.tradePair, req.units, req.isLong)
		req.done <- &CreateOrderResult{
			CreatedOrder: createdOrder,
			Err:          nil,
		}

	case *openOrdersRequest:
		openOrders := orderer.openOrders()
		req.done <- &OpenOrdersResult{
			OpenOrders: openOrders,
			Err:        nil,
		}

	case *closeOrderRequest:
		closedOrder := orderer.closeOrder(req.orderID)
		req.done <- &CloseOrderResult{
			ClosedOrder: closedOrder,
			Err:         nil,
		}
	}
}

func profitPips(price Price, order *order) float64 {
	var diff float64
	if order.isLong {
		diff = price.Bid() - order.priceAtOpen
	} else {
		diff = order.priceAtOpen - price.Ask()
	}
	return diff / order.tradePair.PricePerPip()
}

// SimpleSimulationOrdererFactory is a factory for simple simulation orderer
type SimpleSimulationOrdererFactory struct{}

// NewSimpleSimulationOrdererFactory is a constructor for simple simulation orderer
func NewSimpleSimulationOrdererFactory() *SimpleSimulationOrdererFactory {
	return &SimpleSimulationOrdererFactory{}
}

// Create is a factory method for simple simulation orderer
func (factory *SimpleSimulationOrdererFactory) Create(broker Broker) Orderer {
	simpleSimulationBroker, ok := broker.(*SimpleSimulationBroker)
	if !ok {
		panic(model.ErrWrongType)
	}

	orderer := &SimpleSimulationOrderer{
		SimpleSimulationBroker: simpleSimulationBroker,
		currentOrderID:         0,
		currentOrders:          []*order{},
		ch:                     make(chan interface{}, config.SimulationOrdererChanCap),
	}

	orderer.run()
	return orderer
}

type order struct {
	orderID          OrderID
	tradePair        model.TradePair
	timeAtOpen       int64
	priceAtOpen      float64
	units            float64
	isLong           bool
	timeAtClose      int64
	priceAtClose     float64
	realizedProfit   float64
	unrealizedProfit float64
}

func (order *order) OrderID() OrderID {
	return order.orderID
}

func (order *order) TradePair() model.TradePair {
	return order.tradePair
}

func (order *order) TimeAtOpen() int64 {
	return order.timeAtOpen
}

func (order *order) PriceAtOpen() float64 {
	return order.priceAtOpen
}

func (order *order) Units() float64 {
	return order.units
}

func (order *order) IsLong() bool {
	return order.isLong
}

func (order *order) TimeAtClose() int64 {
	return order.timeAtClose
}

func (order *order) PriceAtClose() float64 {
	return order.priceAtClose
}

func (order *order) RealizedProfit() float64 {
	return order.realizedProfit
}

func (order *order) UnrealizedProfit() float64 {
	return order.unrealizedProfit
}
