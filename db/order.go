package db

import "yukimaterrace/andaman/model"

// GetOrderByTradeRunAndBrokerOrder is a method to get order by trade run and broker order
func GetOrderByTradeRunAndBrokerOrder(tradeRunID int, brokerOrderID int) (*model.Order, error) {
	q := "select * from order_ where trade_run_id = ? and broker_order_id = ?"

	order := model.Order{}
	row := db.QueryRow(q, tradeRunID, brokerOrderID)

	err := row.Scan(
		&order.OrderID,
		&order.TradeRunID,
		&order.BrokerOrderID,
		&order.TradeConfigurationID,
		&order.Units,
		&order.TradeDirection,
		&order.State,
		&order.Profit,
		&order.TimeAtOpen,
		&order.PriceAtOpen,
		&order.TimeAtClose,
		&order.PriceAtClose,
	)

	if err != nil {
		return nil, err
	}
	return &order, nil
}

// GetOrdersByTradeRunAndState is a method to get orders by trade run and state
func GetOrdersByTradeRunAndState(tradeRunID int, state model.OrderState) ([]*model.Order, error) {
	q := "select * from order_ where trade_run_id = ? and state = ?"

	rows, err := db.Query(q, tradeRunID, state)
	if err != nil {
		return nil, err
	}

	var orders []*model.Order
	for rows.Next() {
		order := model.Order{}
		err := rows.Scan(
			&order.OrderID,
			&order.TradeRunID,
			&order.BrokerOrderID,
			&order.TradeConfigurationID,
			&order.Units,
			&order.TradeDirection,
			&order.Profit,
			&order.TimeAtOpen,
			&order.PriceAtOpen,
			&order.TimeAtClose,
			&order.PriceAtClose,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

// GetOpenOrders is a method to get open orders
func GetOpenOrders(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) ([]*model.Order, error) {
	q := `
		select
			order_.trade_run_id,
			order_.broker_order_id,
			order_.units,
			order_.trade_direction,
			order_.state,
			order_.profit,
			order_.time_at_open,
			order_.price_at_open,
			order_.time_at_close,
			order_.price_at_close
		from
			order_,
			trade_configuration,
			trade_algorithm
		where
			order_.trade_run_id = ? and
			order_.state = ? and
			order_.trade_configuration_id = trade_configuration.trade_configuration_id and
			trade_configuration.trade_algorithm_id = trade_algorithm.trade_algorithm_id and
			trade_configuration.trade_pair = ? and
			trade_configuration.timezone = ? and
			trade_algorithm.type = ? and
			trade_algorithm.trade_direction = ?
		`

	orders := []*model.Order{}

	rows, err := db.Query(q, tradeRunID, model.Open, tradePair, timezone, algorithmType, tradeDirection)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order model.Order
		err := rows.Scan(
			&order.TradeRunID,
			&order.BrokerOrderID,
			&order.Units,
			&order.TradeDirection,
			&order.State,
			&order.Profit,
			&order.TimeAtOpen,
			&order.PriceAtOpen,
			&order.TimeAtClose,
			&order.PriceAtClose,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders)
	}

	return orders, nil
}

// AddOrder is a method to add order
func AddOrder(
	tradeRunID int, brokerOrderID int, tradeConfigurationID int, units float64, tradeDirection model.TradeDirection,
	state model.OrderState, profit float64, timeAtOpen int, priceAtOpen float64, timeAtClose int, priceAtClose float64) error {
	q := `
		insert into order_ (
			trade_run_id,
			broker_order_id,
			trade_configuration_id,
			units,
			trade_direction,
			state,
			profit,
			time_at_open,
			price_at_open,
			time_at_close,
			price_at_close
		) values (
			?, ?, ?, ?, ?, ?, ?, ?, ? ,? ,?
		)
		`

	_, err := db.Exec(
		q,
		brokerOrderID,
		tradeConfigurationID,
		units,
		tradeDirection,
		state,
		profit,
		timeAtOpen,
		priceAtOpen,
		timeAtClose,
		priceAtClose,
	)

	if err != nil {
		return err
	}
	return nil
}

// UpdateOrderForProfit is a method to update order for profit
func UpdateOrderForProfit(orderID int, profit float64) error {
	q := "update order_ set profit = ? order_id = ?"

	if _, err := db.Exec(q, profit, orderID); err != nil {
		return err
	}
	return nil
}

// UpdateOrderForClose is a method to update order for close
func UpdateOrderForClose(tradeRunID int, brokerOrderID int, state model.OrderState, profit float64, timeAtClose int, priceAtClose float64) error {
	q := `
		update
			order_
		set
			state = ?,
			profit = ?,
			time_at_close = ?,
			price_at_close = ?
		where
			trade_run_id = ? and
			broker_order_id = ?
		`

	if _, err := db.Exec(q, state, profit, timeAtClose, priceAtClose, tradeRunID, brokerOrderID); err != nil {
		return err
	}
	return nil
}
