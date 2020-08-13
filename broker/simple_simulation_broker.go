package broker

import "fmt"

// SimpleSimulationBroker is a struct for simple simulation broker
type SimpleSimulationBroker struct {
	currentPriceMap   map[TradePair]*Price
	currentTime       int
	currentOrderIDMap map[AccountID]int
	currentOrdersMap  map[AccountID][]*order
}

// Update is a method to update broker
func (broker *SimpleSimulationBroker) Update(priceExtractor PriceExtractor) {
	for _, tradePair := range priceExtractor.TradePairs() {
		broker.currentPriceMap[tradePair] = priceExtractor.Price(tradePair)
	}

	broker.currentTime = priceExtractor.Time()
}

func (broker *SimpleSimulationBroker) price(tradePair TradePair) *Price {
	price, ok := broker.currentPriceMap[tradePair]
	if !ok {
		panic(fmt.Sprintf("%s cannot handle in this broker", string(tradePair)))
	}
	return price
}

// CreateOrder is a method to create order
func (broker *SimpleSimulationBroker) CreateOrder(accountID AccountID, tradePair TradePair, units float64, isLong bool) <-chan CreatedOrder {
	price := broker.price(tradePair)

	if _, ok := broker.currentOrderIDMap[accountID]; !ok {
		broker.currentOrderIDMap[accountID] = 0
	}

	if _, ok := broker.currentOrdersMap[accountID]; !ok {
		broker.currentOrdersMap[accountID] = make([]*order, 0)
	}

	order := &order{
		orderID:    broker.currentOrderIDMap[accountID],
		tradePair:  tradePair,
		units:      units,
		isLong:     isLong,
		timeAtOpen: broker.currentTime,
	}

	if isLong {
		order.priceAtOpen = price.Ask
	} else {
		order.priceAtOpen = price.Bid
	}

	order.unrealizedProfit = profitPips(price, order)

	currentOrderID := broker.currentOrderIDMap[accountID]
	broker.currentOrderIDMap[accountID] = currentOrderID + 1

	currentOrders := broker.currentOrdersMap[accountID]
	broker.currentOrdersMap[accountID] = append(currentOrders, order)

	done := make(chan CreatedOrder, 1)
	done <- order

	return done
}

// OpenOrders is a method to open orders
func (broker *SimpleSimulationBroker) OpenOrders(accountID AccountID) <-chan []OpenOrder {
	if _, ok := broker.currentOrdersMap[accountID]; !ok {
		broker.currentOrdersMap[accountID] = make([]*order, 0)
	}

	orders := broker.currentOrdersMap[accountID]
	openOrders := make([]OpenOrder, len(orders))

	for i, order := range orders {
		price := broker.price(order.tradePair)
		order.unrealizedProfit = profitPips(price, order)
		openOrders[i] = order
	}

	done := make(chan []OpenOrder, 1)
	done <- openOrders

	return done
}

// CloseOrder is a method to close order
func (broker *SimpleSimulationBroker) CloseOrder(accountID AccountID, orderID OrderID) <-chan ClosedOrder {
	if _, ok := broker.currentOrdersMap[accountID]; !ok {
		broker.currentOrdersMap[accountID] = make([]*order, 0)
	}

	orders := broker.currentOrdersMap[accountID]

	var order *order
	var pos int
	for i, o := range orders {
		if o.orderID == int(orderID) {
			order = o
			pos = i
			break
		}
	}

	if order == nil {
		panic(fmt.Sprintf("no open order exists: orderID %d", int(orderID)))
	}

	price := broker.price(order.tradePair)

	order.timeAtClose = broker.currentTime

	if order.isLong {
		order.priceAtClose = price.Bid
	} else {
		order.priceAtClose = price.Ask
	}

	order.realizedProfit = profitPips(price, order)

	broker.currentOrdersMap[accountID] = append(orders[:pos], orders[pos+1:]...)

	done := make(chan ClosedOrder, 1)
	done <- order

	return done
}

func profitPips(price *Price, order *order) float64 {
	var diff float64
	if order.isLong {
		diff = price.Bid - order.priceAtOpen
	} else {
		diff = order.priceAtOpen - price.Ask
	}
	return diff / priceGap2Pips(order.tradePair)
}

type order struct {
	orderID          int
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

func (order *order) OrderID() int {
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
