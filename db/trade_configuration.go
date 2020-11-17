package db

import "yukimaterrace/andaman/model"

// GetTradeConfigurationByFields is a method to get trade configuration by fields
func GetTradeConfigurationByFields(tradePair model.TradePair, timezone model.Timezone, tradeAlgorithmID int) (*model.TradeConfiguration, error) {
	q := "select * from trade_configuration where trade_pair = ? and timezone = ? and trade_algorithm_id = ?"

	tc := model.TradeConfiguration{}
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

// AddTradeConfiguration is a method to add trade configuration
func AddTradeConfiguration(tradePair model.TradePair, timezone model.Timezone, tradeAlgorithmID int) error {
	q := "insert into trade_configuration (trade_pair, timezone, trade_algorithm_id) values (?, ?, ?)"

	if _, err := db.Exec(q, tradePair, timezone, tradeAlgorithmID); err != nil {
		return err
	}
	return nil
}

// DeleteTradeConfiguration is a method to delete trade configuration
func DeleteTradeConfiguration(tradeConfigurationID int) error {
	q := "delete from trade_configuration where trade_configuration_id = ?"

	if _, err := db.Exec(q, tradeConfigurationID); err != nil {
		return err
	}
	return nil
}

// GetTradeConfigurationDetailsByTradeSetID is a method to get trade configuration details by trade set id
func GetTradeConfigurationDetailsByTradeSetID(tradeSetID int) ([]*model.TradeConfigurationDetail, error) {
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
			trade_configuration
		inner join
			trade_algorithm
		on
			trade_algorithm.trade_algorithm_id = trade_configuration.trade_algorithm_id
		inner join
			trade_set_configuration_rel
		on
			trade_set_configuration_rel.trade_configuration_id = trade_configuration.trade_configuration_id and
			trade_set_configuration_rel.trade_set_id = ? 
	`

	rows, err := db.Query(q, tradeSetID)
	if err != nil {
		return nil, err
	}

	var details []*model.TradeConfigurationDetail
	for rows.Next() {
		d := model.TradeConfigurationDetail{}
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
