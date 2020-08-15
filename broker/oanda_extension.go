package broker

import "math"

// Bid is a method to get bid
func (oandaClientPrice *OandaClientPrice) Bid() float64 {
	return oandaClientPrice.Bids[0].Price
}

// Ask is a method to get ask
func (oandaClientPrice *OandaClientPrice) Ask() float64 {
	return oandaClientPrice.Asks[0].Price
}

// OrderID is a method to get order ID
func (oandaOrderCreated *OandaOrderCreated) OrderID() int {
	return oandaOrderCreated.OrderFillTransaction.TradeOpened.TradeID
}

// TradePair is a method to get trade pair
func (oandaOrderCreated *OandaOrderCreated) TradePair() TradePair {
	return TradePair(oandaOrderCreated.OrderFillTransaction.Instrument)
}

// TimeAtOpen is a method to get time at open
func (oandaOrderCreated *OandaOrderCreated) TimeAtOpen() int {
	return int(oandaOrderCreated.OrderFillTransaction.Time)
}

// PriceAtOpen is a method to get price at open
func (oandaOrderCreated *OandaOrderCreated) PriceAtOpen() float64 {
	return oandaOrderCreated.OrderFillTransaction.TradeOpened.Price
}

// Units is a method to get units
func (oandaOrderCreated *OandaOrderCreated) Units() float64 {
	return oandaOrderCreated.OrderFillTransaction.TradeOpened.Units
}

// IsLong is a method to get is long
func (oandaOrderCreated *OandaOrderCreated) IsLong() bool {
	return oandaOrderCreated.OrderFillTransaction.TradeOpened.Units > 0
}

// OrderID is a methdo to get id
func (oandaTrade *OandaTrade) OrderID() int {
	return oandaTrade.ID
}

// TradePair is a method to get trade pair
func (oandaTrade *OandaTrade) TradePair() TradePair {
	return TradePair(oandaTrade.Instrument)
}

// TimeAtOpen is a method to get open time
func (oandaTrade *OandaTrade) TimeAtOpen() int {
	return int(oandaTrade.OpenTime)
}

// PriceAtOpen is a method to get open price
func (oandaTrade *OandaTrade) PriceAtOpen() float64 {
	return oandaTrade.Price
}

// Units is a method to get units
func (oandaTrade *OandaTrade) Units() float64 {
	return math.Abs(oandaTrade.CurrentUnits)
}

// IsLong is a method to get is long value
func (oandaTrade *OandaTrade) IsLong() bool {
	return oandaTrade.CurrentUnits > 0
}

// UnrealizedProfit is a method to get unrealized profit
func (oandaTrade *OandaTrade) UnrealizedProfit() float64 {
	return oandaTrade.UnrealizedPL
}

// OrderID is a method to get order id
func (oandaTradeClosed *OandaTradeClosed) OrderID() int {
	return oandaTradeClosed.OrderFillTransaction.TradesClosed[0].TradeID
}

// TradePair is a method to get trade pair
func (oandaTradeClosed *OandaTradeClosed) TradePair() TradePair {
	return TradePair(oandaTradeClosed.OrderFillTransaction.Instrument)
}

// TimeAtClose is a method to get time to close
func (oandaTradeClosed *OandaTradeClosed) TimeAtClose() int {
	return int(oandaTradeClosed.OrderFillTransaction.Time)
}

// PriceAtClose is a method to get price at close
func (oandaTradeClosed *OandaTradeClosed) PriceAtClose() float64 {
	return oandaTradeClosed.OrderFillTransaction.TradesClosed[0].Price
}

// RealizedProfit is a method to get realized profit
func (oandaTradeClosed *OandaTradeClosed) RealizedProfit() float64 {
	return oandaTradeClosed.OrderFillTransaction.TradesClosed[0].RealizedPL
}
