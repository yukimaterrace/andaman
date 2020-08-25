package broker

// Broker is an interface for broker
type Broker interface{}

// SimulationBroker is an interface for simulation broker
type SimulationBroker interface {
	Broker
	Update(priceExtractor PriceExtractor)
}

// Orderer is an interface for orderer
type Orderer interface {
	CreateOrder(tradePair TradePair, units float64, isLong bool) <-chan *CreateOrderResult
	OpenOrders() <-chan *OpenOrdersResult
	CloseOrder(orderID OrderID) <-chan *CloseOrderResult
}

// OrdererFactory is a factory for orderer
type OrdererFactory interface {
	Create(broker Broker) Orderer
}

// OrderID is a type definition for order ID
type OrderID int

// PriceExtractor is an interface for price extractor
type PriceExtractor interface {
	TradePairs() []TradePair
	Price(tradePair TradePair) Price
	Time() int
}

// Price is an interface for price
type Price interface {
	Bid() float64
	Ask() float64
}

// CreateOrderResult is a result for create order
type CreateOrderResult struct {
	CreatedOrder CreatedOrder
	Err          error
}

// OpenOrdersResult is a result for open orders
type OpenOrdersResult struct {
	OpenOrders []OpenOrder
	Err        error
}

// CloseOrderResult is a result for close order
type CloseOrderResult struct {
	ClosedOrder ClosedOrder
	Err         error
}

// CreatedOrder is an interface for created order
type CreatedOrder interface {
	OrderID() OrderID
	TradePair() TradePair
	TimeAtOpen() int
	PriceAtOpen() float64
	Units() float64
	IsLong() bool
}

// OpenOrder is an interface for open order
type OpenOrder interface {
	CreatedOrder
	UnrealizedProfit() float64
}

// ClosedOrder is an inteface for closed order
type ClosedOrder interface {
	OrderID() OrderID
	TradePair() TradePair
	TimeAtClose() int
	PriceAtClose() float64
	RealizedProfit() float64
}
