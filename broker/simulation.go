package broker

import (
	"fmt"
	"log"
)

// SimpleSimulationBroker is a broker for simple simulation
type SimpleSimulationBroker struct {
	currentPriceMap map[TradePair]Price
	currentTime     int
}

// SimpleSimulationOrderer is a struct for simple simulation broker
type SimpleSimulationOrderer struct {
	*SimpleSimulationBroker
	currentOrderIDMap map[AccountID]OrderID
	currentOrdersMap  map[AccountID][]*order
}

// SimpleSimulationOrdererFactory is a factory for simple simulation orderer
type SimpleSimulationOrdererFactory struct{}

// Create is a factory method for simple simulation orderer
func (factory *SimpleSimulationOrdererFactory) Create(broker Broker) Orderer {
	simpleSimulationBroker, ok := broker.(*SimpleSimulationBroker)
	if !ok {
		panic("wrong type object has been passed")
	}

	return &SimpleSimulationOrderer{
		SimpleSimulationBroker: simpleSimulationBroker,
		currentOrderIDMap:      map[AccountID]OrderID{},
		currentOrdersMap:       map[AccountID][]*order{},
	}
}

// Update is a method to update orderer
func (orderer *SimpleSimulationOrderer) Update(priceExtractor PriceExtractor) {
	for _, tradePair := range priceExtractor.TradePairs() {
		orderer.currentPriceMap[tradePair] = priceExtractor.Price(tradePair)
	}

	orderer.currentTime = priceExtractor.Time()
}

func (orderer *SimpleSimulationOrderer) price(tradePair TradePair) Price {
	price, ok := orderer.currentPriceMap[tradePair]
	if !ok {
		log.Panicf("%s cannot handle in this orderer\n", string(tradePair))
	}
	return price
}

// CreateOrder is a method to create order
func (orderer *SimpleSimulationOrderer) CreateOrder(accountID AccountID, tradePair TradePair, units float64, isLong bool) <-chan *CreateOrderResult {
	price := orderer.price(tradePair)

	if _, ok := orderer.currentOrderIDMap[accountID]; !ok {
		orderer.currentOrderIDMap[accountID] = 0
	}

	if _, ok := orderer.currentOrdersMap[accountID]; !ok {
		orderer.currentOrdersMap[accountID] = make([]*order, 0)
	}

	order := &order{
		orderID:    orderer.currentOrderIDMap[accountID],
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

	currentOrderID := orderer.currentOrderIDMap[accountID]
	orderer.currentOrderIDMap[accountID] = currentOrderID + 1

	currentOrders := orderer.currentOrdersMap[accountID]
	orderer.currentOrdersMap[accountID] = append(currentOrders, order)

	done := make(chan *CreateOrderResult, 1)
	done <- &CreateOrderResult{order, nil}

	return done
}

// OpenOrders is a method to open orders
func (orderer *SimpleSimulationOrderer) OpenOrders(accountID AccountID) <-chan *OpenOrdersResult {
	if _, ok := orderer.currentOrdersMap[accountID]; !ok {
		orderer.currentOrdersMap[accountID] = make([]*order, 0)
	}

	orders := orderer.currentOrdersMap[accountID]
	openOrders := make([]OpenOrder, len(orders))

	for i, order := range orders {
		price := orderer.price(order.tradePair)
		order.unrealizedProfit = profitPips(price, order)
		openOrders[i] = order
	}

	done := make(chan *OpenOrdersResult, 1)
	done <- &OpenOrdersResult{openOrders, nil}

	return done
}

// CloseOrder is a method to close order
func (orderer *SimpleSimulationOrderer) CloseOrder(accountID AccountID, orderID OrderID) <-chan *CloseOrderResult {
	if _, ok := orderer.currentOrdersMap[accountID]; !ok {
		orderer.currentOrdersMap[accountID] = make([]*order, 0)
	}

	orders := orderer.currentOrdersMap[accountID]

	var order *order
	var pos int
	for i, o := range orders {
		if o.orderID == orderID {
			order = o
			pos = i
			break
		}
	}

	if order == nil {
		panic(fmt.Sprintf("no open order exists: orderID %d", int(orderID)))
	}

	price := orderer.price(order.tradePair)

	order.timeAtClose = orderer.currentTime

	if order.isLong {
		order.priceAtClose = price.Bid()
	} else {
		order.priceAtClose = price.Ask()
	}

	order.realizedProfit = profitPips(price, order)

	orderer.currentOrdersMap[accountID] = append(orders[:pos], orders[pos+1:]...)

	done := make(chan *CloseOrderResult, 1)
	done <- &CloseOrderResult{order, nil}

	return done
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

type order struct {
	orderID          OrderID
	tradePair        TradePair
	timeAtOpen       int
	priceAtOpen      float64
	units            float64
	isLong           bool
	timeAtClose      int
	priceAtClose     float64
	realizedProfit   float64
	unrealizedProfit float64
}

func (order *order) OrderID() OrderID {
	return order.orderID
}

func (order *order) TradePair() TradePair {
	return order.tradePair
}

func (order *order) TimeAtOpen() int {
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

func (order *order) TimeAtClose() int {
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
