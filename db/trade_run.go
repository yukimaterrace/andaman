package db

import "yukimaterrace/andaman/model"

// GetLastTradeRun is a method to get last trade run
func GetLastTradeRun() (*model.TradeRun, error) {
	q := "select * from trade_run order by trade_run_id desc"

	tradeRun := model.TradeRun{}
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

// AddTradeRun is a method to add trade run
func AddTradeRun(tradeSetID int, state model.TradeRunState, createdAt int, updatedAt int) error {
	q := "insert into trade_run (trade_set_id, state, created_at, updated_at) values (?, ?, ?, ?)"

	if _, err := db.Exec(q, tradeSetID, state, createdAt, updatedAt); err != nil {
		return err
	}
	return nil
}

// UpdateTradeRun is a method to update trade run
func UpdateTradeRun(tradeRunID int, state model.TradeRunState, updatedAt int) error {
	q := "update trade_run set state = ?, updated_at = ? where trade_run_id = ?"

	if _, err := db.Exec(q, state, updatedAt, tradeRunID); err != nil {
		return err
	}
	return nil
}
