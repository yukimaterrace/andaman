package service

import (
	"time"
	"yukimaterrace/andaman/db"
	"yukimaterrace/andaman/model"
)

// AddTradeRun is a method to add trade run
func AddTradeRun(tradeSetName string) (*model.TradeRun, error) {
	tradeSet, err := db.GetTradeSetByName(tradeSetName)
	if err != nil {
		return nil, err
	}

	now := int(time.Now().Unix())
	if err := db.AddTradeRun(tradeSet.TradeSetID, model.Running, now, now); err != nil {
		return nil, err
	}

	tradeRun, err := db.GetLastTradeRun()
	if err != nil {
		return nil, err
	}
	return tradeRun, nil
}
