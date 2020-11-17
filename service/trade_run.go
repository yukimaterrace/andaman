package service

import (
	"time"
	"yukimaterrace/andaman/db"
	"yukimaterrace/andaman/model"
)

// GetTradeRunDetailsResponse is a method to get trade run details
func GetTradeRunDetailsResponse(_type model.TradeRunType, count int, offset int) (*model.TradeRunDetailsResponse, error) {
	tradeRuns, err := db.GetTradeRunDetails(_type, count, offset)
	if err != nil {
		return nil, err
	}

	all, err := db.CountTradeRun(_type)
	if err != nil {
		return nil, err
	}

	paging := &model.OffsetPaging{All: all, Count: len(tradeRuns), Offset: offset}
	return &model.TradeRunDetailsResponse{TradeRuns: tradeRuns, Paging: paging}, nil
}

// AddTradeRun is a method to add trade run
func AddTradeRun(tradeSetName string, tradeSetVersion int, tradeRunType model.TradeRunType) (*model.TradeRun, error) {
	tradeSet, err := db.GetTradeSet(tradeSetName, tradeSetVersion)
	if err != nil {
		return nil, err
	}

	now := int(time.Now().Unix())
	if err := db.AddTradeRun(tradeSet.TradeSetID, tradeRunType, model.Pending, now, 0, 0); err != nil {
		return nil, err
	}

	tradeRun, err := db.GetLastTradeRun()
	if err != nil {
		return nil, err
	}
	return tradeRun, nil
}

// UpdateTradeRunForStart is a method to update trade run for start
func UpdateTradeRunForStart(tradeRun *model.TradeRun) (*model.TradeRun, error) {
	if err := db.UpdateTradeRunForStart(tradeRun.TradeRunID, int(time.Now().Unix())); err != nil {
		return nil, err
	}

	return db.GetTradeRun(tradeRun.TradeRunID)
}

// UpdateTradeRunForFinish is a method to update trade run for finish
func UpdateTradeRunForFinish(tradeRun *model.TradeRun) error {
	if err := db.UpdateTradeRunForFinish(tradeRun.TradeRunID, int(time.Now().Unix())); err != nil {
		return err
	}
	return nil
}
