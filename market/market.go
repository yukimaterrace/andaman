package market

// Market is an interface for market
type Market interface {
	Prices(instrument Instrument, granularity Granularity, count int, from int64) chan PriceSequence

	LatestPrice(instrument Instrument) chan PriceDetail

	Orders(instrument Instrument) chan Orders

	Asset() chan AssetStatus

	MakeOrder(instrument Instrument, orderType OrderType, unit float64) chan MakeOrderStatus

	CloseOrder(orderID string) chan CloseOrderStatus
}

// Adaptor is an interface for market adaptor
type Adaptor interface {
	Prices(instrument Instrument, granularity Granularity, count int, from int64) PriceSequence

	LatestPrice(instrument Instrument) PriceDetail

	Orders(instrument Instrument) Orders

	Asset() AssetStatus

	MakeOrder(instrument Instrument, orderType OrderType, unit float64) MakeOrderStatus

	CloseOrder(orderID string) CloseOrderStatus
}

// Routine is an implementation for market routine
type Routine struct {
	Adaptor
}
