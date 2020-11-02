package db

import (
	"crypto/sha256"
	"encoding/hex"
	"yukimaterrace/andaman/model"
)

func getHash(data string) string {
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

// GetTradeAlgorithmByTypeAndParam is a method to get trade algorithm by type and param
func GetTradeAlgorithmByTypeAndParam(_type model.TradeAlgorithmType, param string) (*model.TradeAlgorithm, error) {
	paramHash := getHash(param)

	q := "select * from trade_algorithm where type = ? and param_hash = ?"
	row := db.QueryRow(q, _type, paramHash)

	ta := model.TradeAlgorithm{}
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

// AddTradeAlgorithm is a method to add trade algorithm
func AddTradeAlgorithm(_type model.TradeAlgorithmType, param string, tradeDirection model.TradeDirection) error {
	paramHash := getHash(param)

	q := "insert into trade_algorithm (type, param_hash, param, trade_direction) values (?, ?, ?, ?)"
	if _, err := db.Exec(q, _type, paramHash, param, tradeDirection); err != nil {
		return err
	}
	return nil
}

// DeleteTradeAlgorithm is a method to delete trade algorithm
func DeleteTradeAlgorithm(tradeAlgorithmID int) error {
	q := "delete from trade_algorithm where trade_algorithm_id = ?"
	if _, err := db.Exec(q, tradeAlgorithmID); err != nil {
		return err
	}
	return nil
}
