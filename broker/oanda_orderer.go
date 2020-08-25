package broker

// OandaOrderer is a struct for oanda orderer
type OandaOrderer struct {
	*OandaBroker
}

// OandaOrdererFactory is a factory for oanda orderer
type OandaOrdererFactory struct{}

// Create is a factory method of oanda orderer
func (factory *OandaOrdererFactory) Create(broker Broker) Orderer {
	oandaBroker, ok := broker.(*OandaBroker)
	if !ok {
		panic("wrong type passed")
	}
	return &OandaOrderer{oandaBroker}
}

// CreateOrder is a method to creates order
func (orderer *OandaOrderer) CreateOrder(tradePair TradePair, units float64, isLong bool) <-chan *CreateOrderResult {
	done := make(chan *CreateOrderResult, 1)

	go func() {
		var u float64
		if isLong {
			u = units
		} else {
			u = -units
		}

		createdOrder, err := orderer.OandaBroker.CreateOrder("MARKET", string(tradePair), u)

		done <- &CreateOrderResult{createdOrder, err}
		return
	}()

	return done
}

// OpenOrders is a method to get open orders
func (orderer *OandaOrderer) OpenOrders() <-chan *OpenOrdersResult {
	done := make(chan *OpenOrdersResult, 1)

	go func() {
		trades, err := orderer.OandaBroker.OpenTrades()
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
func (orderer *OandaOrderer) CloseOrder(orderID OrderID) <-chan *CloseOrderResult {
	done := make(chan *CloseOrderResult, 1)

	go func() {
		tradeClosed, err := orderer.OandaBroker.CloseTrade(int(orderID))

		done <- &CloseOrderResult{tradeClosed, err}
		return
	}()

	return done
}
