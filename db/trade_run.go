package db

import (
	"database/sql"
	"yukimaterrace/andaman/model"
)

func getTradeRun(row *sql.Row) (*model.TradeRun, error) {
	tradeRun := model.TradeRun{}

	err := row.Scan(
		&tradeRun.TradeRunID,
		&tradeRun.TradeSetID,
		&tradeRun.Type,
		&tradeRun.State,
		&tradeRun.StartAt,
		&tradeRun.FinishAt,
	)
	if err != nil {
		return nil, err
	}

	return &tradeRun, nil
}

// GetTradeRun is a method to get trade run
func GetTradeRun(tradeRunID int) (*model.TradeRun, error) {
	q := "select * from trade_run where trade_run_id = ?"
	row := db.QueryRow(q, tradeRunID)
	return getTradeRun(row)
}

// GetLastTradeRun is a method to get last trade run
func GetLastTradeRun() (*model.TradeRun, error) {
	q := "select * from trade_run order by trade_run_id desc"
	row := db.QueryRow(q)
	return getTradeRun(row)
}

// CountTradeRun is a method to get count of trade run
func CountTradeRun(_type model.TradeRunType) (int, error) {
	q := "select count(1) from trade_run where type = ?"
	row := db.QueryRow(q, _type)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// GetTradeRunDetails is a method to get trade run detail
func GetTradeRunDetails(_type model.TradeRunType, count int, offset int) ([]*model.TradeRunDetail, error) {
	q := `
		select
			trade_run.type,
			trade_run.state,
			trade_run.created_at,
			trade_run.start_at,
			trade_run.finish_at,
			trade_set.name,
			trade_set.type
		from 
			trade_run
		inner join
			trade_set
		on
			trade_set.trade_set_id = trade_run.trade_set_id
		where
			trade_run.type = ? 
		order by
			trade_run.created_at desc 
		limit ?
		offset ?
	`

	rows, err := db.Query(q, _type, count, offset)
	if err != nil {
		return nil, err
	}

	details := []*model.TradeRunDetail{}
	for rows.Next() {
		var d model.TradeRunDetail
		err := rows.Scan(
			&d.Type,
			&d.State,
			&d.CreatedAt,
			&d.StartAt,
			&d.FinishAt,
			&d.TradeSet.Name,
			&d.TradeSet.Type,
		)
		if err != nil {
			return nil, err
		}

		details = append(details, &d)
	}

	return details, nil
}

// AddTradeRun is a method to add trade run
func AddTradeRun(tradeSetID int, _type model.TradeRunType, state model.TradeRunState, createdAt int, startAt int, finishAt int) error {
	q := "insert into trade_run (trade_set_id, type, state, created_at, start_at, finish_at) values (?, ?, ?, ?, ?, ?)"

	if _, err := db.Exec(q, tradeSetID, _type, state, createdAt, startAt, finishAt); err != nil {
		return err
	}
	return nil
}

// UpdateTradeRunForStart is a method to update trade run for start
func UpdateTradeRunForStart(tradeRunID int, startAt int) error {
	q := "update trade_run set state = ?, start_at = ? where trade_run_id = ?"

	if _, err := db.Exec(q, model.Running, startAt, tradeRunID); err != nil {
		return err
	}
	return nil
}

// UpdateTradeRunForFinish is a method to update trade run for finish
func UpdateTradeRunForFinish(tradeRunID int, finishAt int) error {
	q := "update trade_run set state = ?, finish_at = ? where trade_run_id = ?"

	if _, err := db.Exec(q, model.Finished, finishAt, tradeRunID); err != nil {
		return err
	}
	return nil
}
