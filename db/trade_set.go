package db

import (
	"time"
	"yukimaterrace/andaman/model"
)

// GetTradeSetsByType is a method to get trade sets by type
func GetTradeSetsByType(_type model.TradeSetType, count int, offset int) ([]*model.TradeSet, error) {
	q := "select * from trade_set where type = ? order by created_at desc limit ? offset ?"

	rows, err := db.Query(q, _type, count, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tradeSets := []*model.TradeSet{}
	for rows.Next() {
		ts := model.TradeSet{}

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

// GetTradeSetByName is a method to get trade sets by name
func GetTradeSetByName(name string) (*model.TradeSet, error) {
	q := "select * from trade_set where name = ?"

	row := db.QueryRow(q, name)

	ts := model.TradeSet{}
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

// CountTradeSet is a method to get trade count
func CountTradeSet(_type model.TradeSetType) (int, error) {
	q := "select count(1) from trade_set where type = ?"
	row := db.QueryRow(q, _type)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, nil
	}
	return count, nil
}

// AddTradeSet is a method to add trade set
func AddTradeSet(name string, _type model.TradeSetType) error {
	q := "insert into trade_set (name, type, created_at) values (?, ?, ?)"

	now := time.Now().Unix()
	if _, err := db.Exec(q, name, _type, now); err != nil {
		return err
	}
	return nil
}

// DeleteTradeSet is a method to delete trade set
func DeleteTradeSet(tradeSetID int) error {
	q := "delete from trade_set where trade_set_id = ?"

	if _, err := db.Exec(q, tradeSetID); err != nil {
		return err
	}
	return nil
}

// GetTradeSetConfigurationRelByTradeSetIDAndTradeConfigurationID is a method to get trade set configuration rel
func GetTradeSetConfigurationRelByTradeSetIDAndTradeConfigurationID(tradeSetID int, tradeConfigurationID int) (*model.TradeSetConfigurationRel, error) {
	q := "select * from trade_set_configuration_rel where trade_set_id = ? and trade_configuration_id = ?"

	row := db.QueryRow(q, tradeSetID, tradeConfigurationID)

	rel := model.TradeSetConfigurationRel{}
	if err := row.Scan(&rel.TradeSetID, &rel.TradeConfigurationID); err != nil {
		return nil, err
	}

	return &rel, nil
}

// AddTradeSetConfigurationRel is a method to add trade set configuration rel
func AddTradeSetConfigurationRel(tradeSetID int, tradeConfigurationID int) error {
	q := "insert into trade_set_configuration_rel (trade_set_id, trade_configuration_id) values (?, ?)"

	if _, err := db.Exec(q, tradeSetID, tradeConfigurationID); err != nil {
		return err
	}
	return nil
}

// DeleteTradeSetConfigurationRel is a method to delete trade set configuration rel
func DeleteTradeSetConfigurationRel(tradeSetID int, tradeConfigurationID int) error {
	q := "delete from trade_set_configuration where trade_set_id = ? and trade_configuration_id = ?"

	if _, err := db.Exec(q, tradeSetID, tradeConfigurationID); err != nil {
		return err
	}
	return nil
}
