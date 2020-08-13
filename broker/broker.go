package broker

import (
	"strconv"
)

// Broker is an interface for broker
type Broker interface {
	CreateOrder(accountID AccountID, tradePair TradePair, units float64, isLong bool) <-chan CreatedOrder
	OpenOrders(accountID AccountID) <-chan []OpenOrder
	CloseOrder(accountID AccountID, orderID OrderID) <-chan ClosedOrder
}

// SimulationBroker is an interface for simulation broker
type SimulationBroker interface {
	Broker
	Update(priceExtractor PriceExtractor)
}

// AccountID is a type definition for account ID
type AccountID string

// OrderID is a type definition for order ID
type OrderID int

// PriceExtractor is an interface for price extractor
type PriceExtractor interface {
	TradePairs() []TradePair
	Price(tradePair TradePair) *Price
	Time() int
}

// Price is a definition for price
type Price struct {
	Bid float64
	Ask float64
}

// CreatedOrder is an interface for created order
type CreatedOrder interface {
	OrderID() int
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
	OrderID() int
	TradePair() TradePair
	TimeAtClose() int
	PriceAtClose() float64
	RealizedProfit() float64
}

// MakeOrderCsv makes a csv row of order
func MakeOrderCsv(created CreatedOrder, closed ClosedOrder) []string {
	csv := []string{
		strconv.FormatInt(int64(created.OrderID()), 10),
		string(created.TradePair()),
		strconv.FormatInt(int64(created.TimeAtOpen()), 10),
		strconv.FormatFloat(created.PriceAtOpen(), 'f', 6, 64),
		strconv.FormatFloat(created.Units(), 'f', 8, 64),
		strconv.FormatBool(created.IsLong()),
	}

	if closed != nil {
		return append(csv, "not closed", "not closed", "0")
	}

	return append(csv,
		strconv.FormatInt(int64(closed.TimeAtClose()), 10),
		strconv.FormatFloat(closed.PriceAtClose(), 'f', 6, 64),
		strconv.FormatFloat(closed.RealizedProfit(), 'f', 6, 64),
	)
}
