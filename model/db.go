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

func getTradeSetsByType(_type TradeSetType, count int, offset int) ([]TradeSet, error) {
	q := "select * from trade_set where type = ? order by updated_at desc limit ? offset ?"

	rows, err := db.Query(q, _type, count, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tradeSets := []TradeSet{}
	for rows.Next() {
		ts := TradeSet{}

		err := rows.Scan(
			&ts.TradeSetID,
			&ts.Name,
			&ts.Type,
			&ts.State,
			&ts.CreatedAt,
			&ts.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		tradeSets = append(tradeSets, ts)
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
		&ts.State,
		&ts.CreatedAt,
		&ts.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &ts, nil
}

func addTradeSet(name string, _type TradeSetType, state TradeSetState) error {
	q := "insert into trade_set (name, type, state, created_at, updated_at) values (?, ?, ?, ?, ?)"

	now := time.Now().Unix()
	if _, err := db.Exec(q, name, _type, state, now, now); err != nil {
		return err
	}
	return nil
}

func updateTradeSetByName(name string, state TradeSetState) error {
	q := "update trade_set set state = ?, updated_at = ? where name = ?"

	now := time.Now().Unix()
	if _, err := db.Exec(q, state, now, name); err != nil {
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

func getTradeSetConfigurationRelsByTradeSetID(tradeSetID int) ([]TradeSetConfigurationRel, error) {
	q := "select * from trade_set_configuration_rel where trade_set_id = ?"

	rows, err := db.Query(q, tradeSetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rels []TradeSetConfigurationRel
	for rows.Next() {
		rel := TradeSetConfigurationRel{}
		if err := rows.Scan(&rel.TradeSetID, &rel.TradeConfigurationID); err != nil {
			return nil, err
		}
		rels = append(rels, rel)
	}

	return rels, nil
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
