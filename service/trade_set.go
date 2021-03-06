package service

import (
	"database/sql"
	"encoding/json"
	"yukimaterrace/andaman/db"
	"yukimaterrace/andaman/model"
)

// GetTradeSet is a method to get trade set
func GetTradeSet(name string, version int) (*model.TradeSet, error) {
	return db.GetTradeSet(name, version)
}

// GetTradeSetsResponse is a method to get trade sets
func GetTradeSetsResponse(_type model.TradeSetType, count int, offset int) (*model.TradeSetsResponse, error) {
	tradeSets, err := db.GetTradeSetsByType(_type, count, offset)
	if err != nil {
		return nil, err
	}

	all, err := db.CountTradeSet(_type)
	if err != nil {
		return nil, err
	}

	paging := model.OffsetPaging{All: all, Count: len(tradeSets), Offset: offset}
	return &model.TradeSetsResponse{TradeSets: tradeSets, Paging: &paging}, nil
}

// GetTradeSetDetailResponse is a method to get trade set detail
func GetTradeSetDetailResponse(name string, version int) (*model.TradeSetDetailResponse, error) {
	detail, err := GetTradeSetDetail(name, version)
	if err != nil {
		return nil, err
	}

	resp := &model.TradeSetDetailResponse{
		TradeSet: detail,
	}
	return resp, nil
}

// GetTradeSetDetail is a method to get trade set detail
func GetTradeSetDetail(name string, version int) (*model.TradeSetDetail, error) {
	tradeSet, err := db.GetTradeSet(name, version)
	if err != nil {
		return nil, err
	}

	tradeConfigurationDetails, err := db.GetTradeConfigurationDetailsByTradeSetID(tradeSet.TradeSetID)
	if err != nil {
		return nil, err
	}

	tradeSetDetail := &model.TradeSetDetail{
		TradeSet:       tradeSet,
		Configurations: tradeConfigurationDetails,
	}

	return tradeSetDetail, nil
}

// AddTradeSet is a method to add trade set
func AddTradeSet(tradeSetParam *model.TradeSetParam) error {
	version, err := db.GetTradeSetLastVersionByName(tradeSetParam.Name)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	version++
	if err := db.AddTradeSet(tradeSetParam.Name, version, tradeSetParam.Type); err != nil {
		return err
	}

	tradeSet, err := db.GetTradeSet(tradeSetParam.Name, version)
	if err != nil {
		return err
	}

	for _, configurationParam := range tradeSetParam.ConfigurationParams {
		algorithmParam := configurationParam.AlgorithmParam
		_param, err := json.Marshal(algorithmParam.Param)
		if err != nil {
			return err
		}

		param := string(_param)
		tradeAlgorithm, err := db.GetTradeAlgorithmByTypeAndParam(algorithmParam.Type, param)
		if err != nil {
			if err == sql.ErrNoRows {
				if err := db.AddTradeAlgorithm(algorithmParam.Type, param, algorithmParam.TradeDirection); err != nil {
					return err
				}

				tradeAlgorithm, err = db.GetTradeAlgorithmByTypeAndParam(algorithmParam.Type, param)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		tradeConfiguration, err := db.GetTradeConfigurationByFields(configurationParam.TradePair, configurationParam.Timezone, tradeAlgorithm.TradeAlgorithmID)
		if err != nil {
			if err == sql.ErrNoRows {
				if err := db.AddTradeConfiguration(configurationParam.TradePair, configurationParam.Timezone, tradeAlgorithm.TradeAlgorithmID); err != nil {
					return err
				}

				tradeConfiguration, err = db.GetTradeConfigurationByFields(configurationParam.TradePair, configurationParam.Timezone, tradeAlgorithm.TradeAlgorithmID)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		_, err = db.GetTradeSetConfigurationRelByTradeSetIDAndTradeConfigurationID(tradeSet.TradeSetID, tradeConfiguration.TradeConfigurationID)
		if err != nil {
			if err == sql.ErrNoRows {
				if err := db.AddTradeSetConfigurationRel(tradeSet.TradeSetID, tradeConfiguration.TradeConfigurationID); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
