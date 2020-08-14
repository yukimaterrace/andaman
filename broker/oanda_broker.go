package broker

// OandaBroker is a struct for oanda broker
type OandaBroker struct {
	*OandaClient
}

// NewOandaBroker is a constructor for oanda broker
func NewOandaBroker() *OandaBroker {
	return &OandaBroker{NewOandaClient()}
}

// CreateOrder is a method to creates order
func (broker *OandaBroker) CreateOrder(accountID AccountID, tradePair TradePair, units float64, isLong bool) <-chan *CreateOrderResult {
	done := make(chan *CreateOrderResult, 1)

	go func() {
		var u float64
		if isLong {
			u = units
		} else {
			u = -units
		}

		createdOrder, err := broker.OandaClient.CreateOrder(string(accountID), "MARKET", string(tradePair), u)

		done <- &CreateOrderResult{createdOrder, err}
		return
	}()

	return done
}

// OpenOrders is a method to get open orders
func (broker *OandaBroker) OpenOrders(accountID AccountID) <-chan *OpenOrdersResult {
	done := make(chan *OpenOrdersResult, 1)

	go func() {
		trades, err := broker.OandaClient.OpenTrades(string(accountID))
		if err != nil {
			done <- &OpenOrdersResult{nil, err}
			return
		}

		openOrders := make([]OpenOrder, len(trades.Trades))
		for i, trade := range trades.Trades {
			openOrders[i] = &trade
		}

		done <- &OpenOrdersResult{openOrders, nil}
		return
	}()

	return done
}

// CloseOrder is a method to close order
func (broker *OandaBroker) CloseOrder(accountID AccountID, orderID OrderID) <-chan *CloseOrderResult {
	done := make(chan *CloseOrderResult, 1)

	go func() {
		tradeClosed, err := broker.OandaClient.CloseTrade(string(accountID), int(orderID))

		done <- &CloseOrderResult{tradeClosed, err}
		return
	}()

	return done
}
