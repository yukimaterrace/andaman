package broker

type (
	// Broker is an interface for broker
	Broker interface{}

	// SimulationBroker is an interface for simulation broker
	SimulationBroker interface {
		Broker
		Update(priceExtractor PriceExtractor)
	}

	// Orderer is an interface for orderer
	Orderer interface {
		CreateOrder(tradePair TradePair, units float64, isLong bool) <-chan *CreateOrderResult
		OpenOrders() <-chan *OpenOrdersResult
		CloseOrder(orderID OrderID) <-chan *CloseOrderResult
	}

	// OrdererFactory is a factory for orderer
	OrdererFactory interface {
		Create(broker Broker) Orderer
	}

	// OrderID is a type definition for order ID
	OrderID int

	// TimeExtractor is an interface for time extractor
	TimeExtractor interface {
		Time() int64
	}

	// PriceExtractor is an interface for price extractor
	PriceExtractor interface {
		TradePairs() []TradePair
		Price(tradePair TradePair) Price
		TimeExtractor
	}

	// Price is an interface for price
	Price interface {
		Bid() float64
		Ask() float64
	}

	// CreateOrderResult is a result for create order
	CreateOrderResult struct {
		CreatedOrder CreatedOrder
		Err          error
	}

	// OpenOrdersResult is a result for open orders
	OpenOrdersResult struct {
		OpenOrders []OpenOrder
		Err        error
	}

	// CloseOrderResult is a result for close order
	CloseOrderResult struct {
		ClosedOrder ClosedOrder
		Err         error
	}

	// CreatedOrder is an interface for created order
	CreatedOrder interface {
		OrderID() OrderID
		TradePair() TradePair
		TimeAtOpen() int64
		PriceAtOpen() float64
		Units() float64
		IsLong() bool
	}

	// OpenOrder is an interface for open order
	OpenOrder interface {
		CreatedOrder
		UnrealizedProfit() float64
	}

	// ClosedOrder is an inteface for closed order
	ClosedOrder interface {
		OrderID() OrderID
		TradePair() TradePair
		TimeAtClose() int64
		PriceAtClose() float64
		RealizedProfit() float64
	}
)
