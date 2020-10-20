package model

import (
	"database/sql"
	"fmt"
	"time"
	"yukimaterrace/andaman/config"

	_ "github.com/go-sql-driver/mysql" // import driver
)

var db *sql.DB

func init() {
	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/andaman", config.DBUser, config.DBPassword, config.DBHost, config.DBPort)

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

	var tradeSets []TradeSet
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
	_, err := db.Exec(q, name, _type, state, now, now)
	if err != nil {
		return err
	}
	return nil
}

func updateTradeSetByName(name string, state TradeSetState) error {
	q := "update trade_set set state = ?, updated_at = ? where name = ?"

	now := time.Now().Unix()
	_, err := db.Exec(q, state, now, name)
	if err != nil {
		return err
	}
	return nil
}

func deleteTradeSetByName(name string) error {
	q := "delete from trade_set where name = ?"

	_, err := db.Exec(q, name)
	if err != nil {
		return err
	}
	return nil
}
