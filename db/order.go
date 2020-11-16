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

// GetOrders is a method to get orders
func GetOrders(
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
			order_
		inner join
			trade_configuration
		on
			trade_configuration.trade_configuration_id = order_.trade_configuration_id and
			trade_configuration.trade_pair = ? and
			trade_configuration.timezone = ?
		inner join
			trade_algorithm
		on
			trade_configuration.trade_algorithm_id = trade_algorithm.trade_algorithm_id and
			trade_algorithm.type = ? and
			trade_algorithm.trade_direction = ?
		where
			order_.trade_run_id = ? and
			order_.state = ?
		limit ?
		offset ?
	`

	rows, err := db.Query(q, tradePair, timezone, algorithmType, tradeDirection, tradeRunID, state, count, offset)
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

// GetTotalProfitForOrders is a method to get total progit for orders
func GetTotalProfitForOrders(
	tradeRunID int, state model.OrderState,
	tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) (float64, error) {
	q := `
		select
			sum(order_.profit)
		from
			order_
		inner join
			trade_configuration
		on
			trade_configuration.trade_configuration_id = order_.trade_configuration_id and
			trade_configuration.trade_pair = ? and
			trade_configuration.timezone = ?
		inner join
			trade_algorithm
		on
			trade_algorithm.trade_algorithm_id = trade_configuration.trade_algorithm_id and
			trade_algorithm.type = ? and
			trade_algorithm.trade_direction = ?
		where
			order_.trade_run_id = ? and
			order_.state = ?
	`

	row := db.QueryRow(q, tradePair, timezone, algorithmType, tradeDirection, tradeRunID, state)

	var totalProfit float64
	if err := row.Scan(&totalProfit); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}

	return totalProfit, nil
}

// GetCountForOrders is a method to get count for orders
func GetCountForOrders(
	tradeRunID int, state model.OrderState,
	tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) (int, error) {
	q := `
		select
			count(1)
		from
			order_
		inner join
			trade_configuration
		on
			trade_configuration.trade_configuration_id = order_.trade_configuration_id and
			trade_configuration.trade_pair = ? and
			trade_configuration.timezone = ? and
		inner join
			trade_algorithm
		on
			trade_algorithm.trade_algorithm_id = trade_configuration.trade_algorithm_id and
			trade_algorithm.type = ? and
			trade_algorithm.trade_direction = ?
		where
			order_.trade_run_id = ? and
			order_.state = ?
	`

	row := db.QueryRow(q, tradePair, timezone, algorithmType, tradeDirection, tradeRunID, state)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// GetCountByPeriod is a method to get count by period
func GetCountByPeriod(
	tradeRunID int, state model.OrderState, tradePair model.TradePair, timezone model.Timezone,
	tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) (int, error) {
	q := `
		select
			count(1)
		from
			order_
		inner join
			trade_configuration
		on
			order_.trade_configuration_id = trade_configuration.trade_configuration_id and
			trade_configuration.trade_pair = ? and
			trade_configuration.timezone = ?
		inner join
			trade_algorithm
		on
			trade_algorithm.trade_algorithm_id = trade_configuration.trade_algorithm_id and
			trade_algorithm.type = ? and
			trade_algorithm.trade_direction = ?
		where
			order_.trade_run_id = ? and
			order_.state = ?
	`

	row := db.QueryRow(q, tradePair, timezone, algorithmType, tradeDirection, tradeRunID, state)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

// GetPositiveProfitCountByPeriod is a method to get positive profit count by period
func GetPositiveProfitCountByPeriod(
	tradeRunID int, state model.OrderState, tradePair model.TradePair, timezone model.Timezone,
	tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) (int, error) {
	q := `
		select
			count(1)
		from
			order_
		inner join
			trade_configuration
		on
			trade_configuration.trade_configuration_id = order_.trade_configuration_id and
			trade_configuration.trade_pair = ? and
			trade_configuration.timezone = ?
		inner join
			trade_algorithm
		on
			trade_algorithm.trade_algorithm_id = trade_configuration.trade_algorithm_id and
			trade_algorithm.type = ? and
			trade_algorithm.trade_direction = ?
		where
			order_.trade_run_id = ? and
			order_.state = ? and
			order_.profit >= 0
	`

	row := db.QueryRow(q, tradePair, timezone, algorithmType, tradeDirection, tradeRunID, state)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
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

// GetTradeCountProfitByFilter1 is a method to get trade count profit 1
func GetTradeCountProfitByFilter1(tradeRunID int, state model.OrderState, start int, end int) (map[model.TradePair]*model.TradeCountProfit, error) {
	q := `
		select
			trade_configuration.trade_pair,
			count(order_.order_id),
			sum(order_.profit)
		from
			order_
		inner join
			trade_configuration
		on
			trade_configuration.trade_configuration_id = order_.trade_configuration_id
		where
			order_.trade_run_id = ? and
			order_.state = ? and
			order_.time_at_open > ? and
			order_.time_at_open < ?
		group by
			trade_configuration.trade_pair
	`

	rows, err := db.Query(q, tradeRunID, state, start, end)
	if err != nil {
		return nil, err
	}

	var m map[model.TradePair]*model.TradeCountProfit
	for rows.Next() {
		var tradePair model.TradePair
		var cp model.TradeCountProfit

		err := rows.Scan(
			&tradePair,
			&cp.Count,
			&cp.Profit,
		)
		if err != nil {
			return nil, err
		}

		m[tradePair] = &cp
	}
	return m, nil
}

// GetTradeCountProfitByFilter2 is a method to get trade count profit by filter 2
func GetTradeCountProfitByFilter2(
	tradeRunID int, state model.OrderState, tradePair model.TradePair, timezone model.Timezone,
	start int, end int) (map[model.TradeConfigurationDetail]*model.TradeCountProfit, error) {

	q := `
		select
			trade_algorithm.type,
			trade_algorithm.trade_direction,
			trade_algorithm.param_hash,
			trade_algorithm.param,
			count(order_.order_id),
			sum(order_.profit)
		from
			order_
		inner join
			trade_configuration
		on
			order_.trade_configuration_id = trade_configuration.trade_configuration_id and
			trade_configuration.trade_pair = ? and
			trade_configuration.timezone = ?
		inner join
			trade_algorithm
		on
			trade_algorithm.trade_algorithm_id = trade_configuration.trade_algorithm_id
		where
			order_.trade_run_id = ? and
			order_.state = ? and
			order_.time_at_open > ? and
			order_.time_at_open < ?
		group by
			trade_algorithm.type,
			trade_algorithm.trade_direction,
			trade_algorithm.param_hash
	`

	rows, err := db.Query(q, tradePair, timezone, tradeRunID, state, start, end)
	if err != nil {
		return nil, err
	}

	var m map[model.TradeConfigurationDetail]*model.TradeCountProfit
	for rows.Next() {
		var key model.TradeConfigurationDetail
		var cp model.TradeCountProfit

		err := rows.Scan(
			&key.Algorithm.Type,
			&key.Algorithm.TradeDirection,
			&key.Algorithm.ParamHash,
			&key.Algorithm.Param,
			&cp.Count,
			&cp.Profit,
		)
		if err != nil {
			return nil, err
		}

		key.TradePair = tradePair
		key.Timezone = timezone

		m[key] = &cp
	}
	return m, nil
}

// GetTradeCountProfitByFilter3 is a method to get trade count profit by filter 3
func GetTradeCountProfitByFilter3(
	tradeRunID int, state model.OrderState, tradePair model.TradePair, timezone model.Timezone,
	tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType, start int, end int) (map[model.TradeConfigurationDetail]*model.TradeCountProfit, error) {

	q := `
		select
			trade_configuration.trade_pair,
			trade_configuration.timezone,
			trade_algorithm.type,
			trade_algorithm.trade_direction,
			trade_algorithm.param_hash,
			trade_algorithm.param,
			count(order_.order_id),
			sum(order_.profit)
		from
			order_
		inner join
			trade_configuration
		on
			trade_configuration.trade_configuration_id = order_.trade_configuration_id and
			trade_configuration.trade_pair = ? and
			trade_configuration.timezone = ?
		inner join
			trade_algorithm
		on
			trade_algorithm.trade_algorithm_id = trade_configuration.trade_algorithm_id and
			trade_algorithm.type = ? and
			trade_algorithm.trade_direction = ?
		where
			order_.trade_run_id = ? and
			order_.state = ? and
			order_.time_at_open > ? and
			order_.time_at_open < ?
		group by
			trade_configuration.trade_pair,
			trade_configuration.timezone,
			trade_algorithm.type,
			trade_algorithm.trade_direction,
			trade_algorithm.param_hash
	`

	rows, err := db.Query(q, tradePair, timezone, algorithmType, tradeDirection, tradeRunID, state, start, end)
	if err != nil {
		return nil, err
	}

	var m map[model.TradeConfigurationDetail]*model.TradeCountProfit
	for rows.Next() {
		var key model.TradeConfigurationDetail
		var cp model.TradeCountProfit

		err := rows.Scan(
			&key.TradePair,
			&key.Timezone,
			&key.Algorithm.Type,
			&key.Algorithm.TradeDirection,
			&key.Algorithm.ParamHash,
			&key.Algorithm.Param,
			&cp.Count,
			&cp.Profit,
		)
		if err != nil {
			return nil, err
		}

		m[key] = &cp
	}
	return m, nil
}

// GetTradeConfigurationTradeCountProfits is a method to get trade count profit by filter 4
func GetTradeConfigurationTradeCountProfits(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType,
	count int, offset int) ([]*model.TradeConfigurationTradeCountProfit, error) {

	q := `
		select
			trade_configuration.trade_pair,
			trade_configuration.timezone,
			trade_algorithm.type,
			trade_algorithm.trade_direction,
			trade_algorithm.param_hash,
			trade_algorithm.param,
			count(order_.order_id),
			sum(order_.profit) as profit
		from
			order_
		inner join
			trade_configuration
		on
			trade_configuration.trade_configuration_id = order_.trade_configuration_id and
			trade_configuration.trade_pair = ? and
			trade_configuration.timezone = ?
		inner join
			trade_algorithm
		on
			trade_algorithm.trade_algorithm_id = trade_configuration.trade_algorithm_id and
			trade_algorithm.type = ? and
			trade_algorithm.trade_direction = ?
		where
			order_.trade_run_id = ?
		group by
			trade_configuration.trade_pair,
			trade_configuration.timezone,
			trade_algorithm.type,
			trade_algorithm.trade_direction,
			trade_algorithm.param_hash
		order by
			profit
			desc
		limit
			?
		offset
			?
	`

	rows, err := db.Query(q, tradePair, timezone, algorithmType, tradeDirection, tradeRunID, count, offset)
	if err != nil {
		return nil, err
	}

	cps := []*model.TradeConfigurationTradeCountProfit{}

	for rows.Next() {
		var cp model.TradeConfigurationTradeCountProfit
		err := rows.Scan(
			&cp.TradeConfiguration.TradePair,
			&cp.TradeConfiguration.Timezone,
			&cp.TradeConfiguration.Algorithm.Type,
			&cp.TradeConfiguration.Algorithm.TradeDirection,
			&cp.TradeConfiguration.Algorithm.ParamHash,
			&cp.TradeConfiguration.Algorithm.Param,
			&cp.TradeCountProfit.Count,
			&cp.TradeCountProfit.Profit,
		)
		if err != nil {
			return nil, err
		}

		cps = append(cps, &cp)
	}

	return cps, nil
}

// GetFirstTradeConfigurationTradeCountProfit is a method to get trade count profit by filter 5
func GetFirstTradeConfigurationTradeCountProfit(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone,
	tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType) (*model.TradeConfigurationTradeCountProfit, error) {

	q := `
		select
			trade_configuration.trade_pair,
			trade_configuration.timezone,
			trade_algorithm.type,
			trade_algorithm.trade_direction,
			trade_algorithm.param_hash,
			trade_algorithm.param,
			count(order_.order_id),
			sum(order_.profit) as profit
		from
			order_
		inner join
			trade_configuration
		on
			trade_configuration.trade_configuration_id = order_.trade_configuration_id and
			trade_configuration.trade_pair = ? and
			trade_configuration.timezone = ?
		inner join
			trade_algorithm
		on
			trade_algorithm.trade_algorithm_id = trade_configuration.trade_algorithm_id and
			trade_algorithm.type = ? and
			trade_algorithm.trade_direction = ?
		where
			order_.trade_run_id = ?
		group by
			trade_configuration.trade_pair,
			trade_configuration.timezone,
			trade_algorithm.type,
			trade_algorithm.trade_direction,
			trade_algorithm.param_hash
		order by
			profit
			desc
		limit
			1
	`

	row := db.QueryRow(q, tradePair, timezone, algorithmType, tradeDirection, tradeRunID)

	var tctp model.TradeConfigurationTradeCountProfit

	err := row.Scan(
		&tctp.TradeConfiguration.TradePair,
		&tctp.TradeConfiguration.Timezone,
		&tctp.TradeConfiguration.Algorithm.Type,
		&tctp.TradeConfiguration.Algorithm.TradeDirection,
		&tctp.TradeConfiguration.Algorithm.ParamHash,
		&tctp.TradeConfiguration.Algorithm.Param,
		&tctp.TradeCountProfit.Count,
		&tctp.TradeCountProfit.Profit,
	)
	if err != nil {
		return nil, err
	}

	return &tctp, nil
}

// GetTotalProfitByFilter1 is a method to get total profit 1
func GetTotalProfitByFilter1(tradeRunID int, state model.OrderState, start int, end int) (float64, error) {
	q := `
		select
			sum(order_.profit)
		from
			order_
		where
			order_.trade_run_id = ? and
			order_.state = ? and
			order_.time_at_open > ? and
			order_.time_at_open < ?
	`

	row := db.QueryRow(q, tradeRunID, state, start, end)

	var profit float64
	if err := row.Scan(&profit); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return profit, nil
}

// GetTotalProfitByFilter2 is a method to get total profit 2
func GetTotalProfitByFilter2(tradeRunID int, state model.OrderState, tradePair model.TradePair, timezone model.Timezone, start int, end int) (float64, error) {
	q := `
		select
			sum(order_.profit)
		from
			order_
		inner join
			trade_configuration
		on
			trade_configuration.trade_configuration_id = order_.trade_configuration_id and
			trade_configuration.trade_pair = ? and
			trade_configuration.timezone = ?
		where
			order_.trade_run_id = ? and
			order_.state = ? and
			order_.time_at_open > ? and
			order_.time_at_open < ?
	`

	row := db.QueryRow(q, tradePair, timezone, tradeRunID, state, start, end)

	var profit float64
	if err := row.Scan(&profit); err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return profit, nil
}

// GetTradeConfigurationGroupsForOrder is a method to get trade configuration groups
func GetTradeConfigurationGroupsForOrder(tradeRunID int, count int, offset int) ([]*model.TradeConfigurationGroup, error) {
	q := `
		select
			trade_configuration.trade_pair,
			trade_configuration.timezone,
			trade_algorithm.trade_direction,
			trade_algorithm.type
		from
			order_
		inner join
			trade_configuration
		on
			trade_configuration.trade_configuration_id = order_.trade_configuration_id
		inner join
			trade_algorithm
		on
			trade_algorithm.trade_algorithm_id = trade_configuration.trade_algorithm_id
		where
			order_.trade_run_id = ?
		group by
			trade_configuration.trade_pair,
			trade_configuration.timezone,
			trade_algorithm.trade_direction,
			trade_algorithm.type
		order by
			trade_configuration.trade_pair,
			trade_configuration.timezone,
			trade_algorithm.trade_direction,
			trade_algorithm.type
		limit ?
		offset ?
	`

	rows, err := db.Query(q, tradeRunID, count, offset)
	if err != nil {
		return nil, err
	}

	groups := []*model.TradeConfigurationGroup{}
	for rows.Next() {
		var g model.TradeConfigurationGroup
		err := rows.Scan(
			&g.TradePair,
			&g.Timezone,
			&g.TradeDirection,
			&g.TradeAlgorithmType,
		)
		if err != nil {
			return nil, err
		}

		groups = append(groups, &g)
	}

	return groups, nil
}
