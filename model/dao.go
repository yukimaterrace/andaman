package model

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"
	"yukimaterrace/andaman/util"

	_ "github.com/go-sql-driver/mysql" // import driver
)

var db *sql.DB

func init() {
	dbUser := util.GetEnv("DB_USER")
	dbPassword := util.GetEnv("DB_PASSWORD")
	dbHost := util.GetEnv("DB_HOST")
	dbPort := util.GetEnv("DB_PORT")

	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/andaman", dbUser, dbPassword, dbHost, dbPort)

	var err error
	db, err = sql.Open("mysql", sourceName)
	if err != nil {
		panic(err)
	}

	if db.Ping() != nil {
		panic(err)
	}
}

func getTradeSetsByType(_type TradeSetType, count int, offset int) ([]*TradeSet, error) {
	q := "select * from trade_set where type = ? order by created_at desc limit ? offset ?"

	rows, err := db.Query(q, _type, count, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tradeSets := []*TradeSet{}
	for rows.Next() {
		ts := TradeSet{}

		err := rows.Scan(
			&ts.TradeSetID,
			&ts.Name,
			&ts.Type,
			&ts.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		tradeSets = append(tradeSets, &ts)
	}

	return tradeSets, nil
}

func getTradeSetByName(name string) (*TradeSet, error) {
	q := "select * from trade_set where name = ?"

	row := db.QueryRow(q, name)

	ts := TradeSet{}
	err := row.Scan(
		&ts.TradeSetID,
		&ts.Name,
		&ts.Type,
		&ts.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &ts, nil
}

func addTradeSet(name string, _type TradeSetType) error {
	q := "insert into trade_set (name, type, created_at) values (?, ?, ?)"

	now := time.Now().Unix()
	if _, err := db.Exec(q, name, _type, now); err != nil {
		return err
	}
	return nil
}

func deleteTradeSet(tradeSetID int) error {
	q := "delete from trade_set where trade_set_id = ?"

	if _, err := db.Exec(q, tradeSetID); err != nil {
		return err
	}
	return nil
}

func getHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func getTradeAlgorithmByTypeAndParam(_type TradeAlgorithmType, param string) (*TradeAlgorithm, error) {
	paramHash := getHash(param)

	q := "select * from trade_algorithm where type = ? and param_hash = ?"
	row := db.QueryRow(q, _type, paramHash)

	ta := TradeAlgorithm{}
	err := row.Scan(
		&ta.TradeAlgorithmID,
		&ta.Type,
		&ta.ParamHash,
		&ta.Param,
		&ta.TradeDirection,
	)
	if err != nil {
		return nil, err
	}

	return &ta, nil
}

func addTradeAlgorithm(_type TradeAlgorithmType, param string, tradeDirection TradeDirection) error {
	paramHash := getHash(param)

	q := "insert into trade_algorithm (type, param_hash, param, trade_direction) values (?, ?, ?, ?)"
	if _, err := db.Exec(q, _type, paramHash, param, tradeDirection); err != nil {
		return err
	}
	return nil
}

func deleteTradeAlgorithm(tradeAlgorithmID int) error {
	q := "delete from trade_algorithm where trade_algorithm_id = ?"
	if _, err := db.Exec(q, tradeAlgorithmID); err != nil {
		return err
	}
	return nil
}

func getTradeConfigurationByFields(tradePair TradePair, timezone Timezone, tradeAlgorithmID int) (*TradeConfiguration, error) {
	q := "select * from trade_configuration where trade_pair = ? and timezone = ? and trade_algorithm_id = ?"

	tc := TradeConfiguration{}
	row := db.QueryRow(q, tradePair, timezone, tradeAlgorithmID)
	err := row.Scan(
		&tc.TradeConfigurationID,
		&tc.TradePair,
		&tc.Timezone,
		&tc.TradeAlgorithmID,
	)
	if err != nil {
		return nil, err
	}

	return &tc, nil
}

func addTradeConfiguration(tradePair TradePair, timezone Timezone, tradeAlgorithmID int) error {
	q := "insert into trade_configuration (trade_pair, timezone, trade_algorithm_id) values (?, ?, ?)"

	if _, err := db.Exec(q, tradePair, timezone, tradeAlgorithmID); err != nil {
		return err
	}
	return nil
}

func deleteTradeConfiguration(tradeConfigurationID int) error {
	q := "delete from trade_configuration where trade_configuration_id = ?"

	if _, err := db.Exec(q, tradeConfigurationID); err != nil {
		return err
	}
	return nil
}

func getTradeSetConfigurationRelByTradeSetIDAndTradeConfigurationID(tradeSetID int, tradeConfigurationID int) (*TradeSetConfigurationRel, error) {
	q := "select * from trade_set_configuration_rel where trade_set_id = ?, trade_configuration_id"

	row := db.QueryRow(q, tradeSetID, tradeConfigurationID)

	rel := TradeSetConfigurationRel{}
	if err := row.Scan(&rel.TradeSetID, &rel.TradeConfigurationID); err != nil {
		return nil, err
	}

	return &rel, nil
}

func addTradeSetConfigurationRel(tradeSetID int, tradeConfigurationID int) error {
	q := "insert into trade_set_configuration_rel (trade_set_id, trade_configuration_id) values (?, ?)"

	if _, err := db.Exec(q, tradeSetID, tradeConfigurationID); err != nil {
		return err
	}
	return nil
}

func deleteTradeSetConfigurationRelByTradeSetID(tradeSetID int) error {
	q := "delete from trade_set_configuration where trade_set_id = ?"

	if _, err := db.Exec(q, tradeSetID); err != nil {
		return err
	}
	return nil
}

func getTradeConfigurationDetailsByTradeSetID(tradeSetID int) ([]*TradeConfigurationDetail, error) {
	q := `
		select
			trade_configuration.trade_configuration_id,
			trade_configuration.trade_pair,
			trade_configuration.timezone,
			trade_algorithm.trade_algorithm_id,
			trade_algorithm.type,
			trade_algorithm.param,
			trade_algorithm.trade_direction
		from 
			trade_configuration,
			trade_set_configuration_rel,
			trade_algorithm
		where
			trade_set_configuration_rel.trade_set_id = ? and
			trade_set_configuration_rel.trade_configuration_id = trade_configuration.trade_configuration_id and
			trade_configuratgion.trade_algorithm_id = trade_algorithm.trade_algorithm_id
	`

	rows, err := db.Query(q, tradeSetID)
	if err != nil {
		return nil, err
	}

	var details []*TradeConfigurationDetail
	for rows.Next() {
		d := TradeConfigurationDetail{}
		err := rows.Scan(
			&d.TradeConfigurationID,
			&d.TradePair,
			&d.Timezone,
			&d.Algorithm.TradeAlgorithmID,
			&d.Algorithm.Type,
			&d.Algorithm.Param,
			&d.Algorithm.TradeDirection,
		)
		if err != nil {
			return nil, err
		}

		details = append(details, &d)
	}

	return details, nil
}

func getOrderByTradeRunAndBrokerOrder(tradeRunID int, brokerOrderID int) (*Order, error) {
	q := "select * from order where trade_run_id = ? and broker_order_id = ?"

	order := Order{}
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

func getOrdersByTradeRunAndState(tradeRunID int, state OrderState) ([]*Order, error) {
	q := "select * from order where trade_run_id = ? and state = ?"

	rows, err := db.Query(q, tradeRunID, state)
	if err != nil {
		return nil, err
	}

	var orders []*Order
	for rows.Next() {
		order := Order{}
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

func addOrder(
	tradeRunID int, brokerOrderID int, tradeConfigurationID int, units float64, tradeDirection TradeDirection,
	state OrderState, profit float64, timeAtOpen int, priceAtOpen float64, timeAtClose int, priceAtClose float64) error {
	q := `
		insert into order (
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

func updateOrderForProfit(orderID int, profit float64) error {
	q := "update order set profit = ? order_id = ?"

	if _, err := db.Exec(q, profit, orderID); err != nil {
		return err
	}
	return nil
}

func updateOrderForClose(tradeRunID int, brokerOrderID int, state OrderState, profit float64, timeAtClose int, priceAtClose float64) error {
	q := `
		update
			order
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

func getLastTradeRun() (*TradeRun, error) {
	q := "select * from trade_run order by trade_run_id desc"

	tradeRun := TradeRun{}
	row := db.QueryRow(q)

	err := row.Scan(
		&tradeRun.TradeRunID,
		&tradeRun.TradeSetID,
		&tradeRun.State,
		&tradeRun.CreatedAt,
		&tradeRun.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &tradeRun, nil
}

func addTradeRun(tradeSetID int, state TradeRunState, createdAt int, updatedAt int) error {
	q := "insert into trade_run (trade_set_id, state, created_at, updated_at) values (?, ?, ?, ?)"

	if _, err := db.Exec(q, tradeSetID, state, createdAt, updatedAt); err != nil {
		return err
	}
	return nil
}

func updateTradeRun(tradeRunID int, state TradeRunState, updatedAt int) error {
	q := "update trade_run set state = ?, updated_at = ? where trade_run_id = ?"

	if _, err := db.Exec(q, state, updatedAt, tradeRunID); err != nil {
		return err
	}
	return nil
}
