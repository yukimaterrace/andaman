package db

import (
	"database/sql"
	"yukimaterrace/andaman/model"
)

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

func getOrders(
	tradeRunID int, state model.OrderState,
	tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType,
	count int, offset int) ([]*model.Order, error) {
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
		limit ?
		offset ?
	`

	rows, err := db.Query(q, tradeRunID, state, tradePair, timezone, algorithmType, tradeDirection, count, offset)
	if err != nil {
		return nil, err
	}

	orders := []*model.Order{}
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

func getTotalProfitForOrders(
	tradeRunID int, state model.OrderState,
	tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) (float64, error) {
	q := `
		select
			sum(order_.profit)
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

	row := db.QueryRow(q, tradeRunID, state, tradePair, timezone, algorithmType, tradeDirection)

	var totalProfit float64
	if err := row.Scan(&totalProfit); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return totalProfit, nil
}

func getCountForOrders(
	tradeRunID int, state model.OrderState,
	tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) (int, error) {
	q := `
		select
			count(1)
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

	row := db.QueryRow(q, tradeRunID, state, tradePair, timezone, algorithmType, tradeDirection)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// GetOpenOrders is a method to get open orders
func GetOpenOrders(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType,
	count int, offset int) ([]*model.Order, error) {

	return getOrders(tradeRunID, model.Open, tradePair, timezone, tradeDirection, algorithmType, count, offset)
}

// GetTotalProfitForOpenOrders is a method to get total profit for open orders
func GetTotalProfitForOpenOrders(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) (float64, error) {

	return getTotalProfitForOrders(tradeRunID, model.Open, tradePair, timezone, tradeDirection, algorithmType)
}

// GetCountForOpenOrders is a method to get count for open orders
func GetCountForOpenOrders(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) (int, error) {

	return getCountForOrders(tradeRunID, model.Open, tradePair, timezone, tradeDirection, algorithmType)
}

// GetClosedOrders is a method to get closed orders
func GetClosedOrders(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType,
	count int, offset int) ([]*model.Order, error) {

	return getOrders(tradeRunID, model.Closed, tradePair, timezone, tradeDirection, algorithmType, count, offset)
}

// GetTotalProfitForClosedOrders is a method to get total profit for closed orders
func GetTotalProfitForClosedOrders(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) (float64, error) {

	return getTotalProfitForOrders(tradeRunID, model.Closed, tradePair, timezone, tradeDirection, algorithmType)
}

// GetCountForClosedOrders is a method to get count for closed orders
func GetCountForClosedOrders(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) (int, error) {

	return getCountForOrders(tradeRunID, model.Closed, tradePair, timezone, tradeDirection, algorithmType)
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
