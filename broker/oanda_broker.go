package broker

import (
	"log"
)

// OandaBroker is a struct for oanda broker
type OandaBroker struct {
	*OandaClient
}

// NewOandaBroker is a constructor for oanda broker
func NewOandaBroker() *OandaBroker {
	return &OandaBroker{NewOandaClient()}
}

// CreateOrder is a method to creates order
func (broker *OandaBroker) CreateOrder(accountID AccountID, tradePair TradePair, units float64, isLong bool) <-chan CreatedOrder {
	done := make(chan CreatedOrder, 1)

	go func() {
		var u float64
		if isLong {
			u = units
		} else {
			u = -units
		}

		createdOrder, err := broker.OandaClient.CreateOrder(string(accountID), "MARKET", string(tradePair), u)
		if err != nil {
			log.Printf(err.Error())
			done <- nil
			return
		}

		done <- createdOrder
		return
	}()

	return done
}

// OpenOrders is a method to get open orders
func (broker *OandaBroker) OpenOrders(accountID AccountID) <-chan []OpenOrder {
	done := make(chan []OpenOrder, 1)

	go func() {
		trades, err := broker.OandaClient.OpenTrades(string(accountID))
		if err != nil {
			log.Println(err.Error())
			done <- nil
			return
		}

		openOrders := make([]OpenOrder, len(trades.Trades))
		for i, trade := range trades.Trades {
			openOrders[i] = &trade
		}

		done <- openOrders
		return
	}()

	return done
}

// CloseOrder is a method to close order
func (broker *OandaBroker) CloseOrder(accountID AccountID, orderID OrderID) <-chan ClosedOrder {
	done := make(chan ClosedOrder, 1)

	go func() {
		tradeClosed, err := broker.OandaClient.CloseTrade(string(accountID), int(orderID))
		if err != nil {
			log.Println(err.Error())
			done <- nil
		}

		done <- tradeClosed
		return
	}()

	return done
}
