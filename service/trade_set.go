package service

import (
	"database/sql"
	"encoding/json"
	"yukimaterrace/andaman/db"
	"yukimaterrace/andaman/model"
)

// GetTradeSetByName is a method to get trade set by name
func GetTradeSetByName(name string) (*model.TradeSet, error) {
	return db.GetTradeSetByName(name)
}

// GetTradeSets is a method to get trade sets
func GetTradeSets(_type model.TradeSetType, count int, offset int) (*model.TradeSetsResponse, error) {
	tradeSets, err := db.GetTradeSetsByType(_type, count, offset)
	if err != nil {
		return nil, model.HandleError(err)
	}

	all, err := db.CountTradeSet(_type)
	if err != nil {
		return nil, err
	}

	paging := model.OffsetPaging{All: all, Count: len(tradeSets), Offset: offset}
	return &model.TradeSetsResponse{TradeSets: tradeSets, Paging: &paging}, nil
}

// GetTradeSetDetail is a method to get trade set detail
func GetTradeSetDetail(name string, paramObjectCreator model.TradeParamObjectCreator) (*model.TradeSetDetail, error) {
	tradeSet, err := db.GetTradeSetByName(name)
	if err != nil {
		return nil, err
	}

	tradeConfigurationDetails, err := db.GetTradeConfigurationDetailsByTradeSetID(tradeSet.TradeSetID)
	if err != nil {
		return nil, err
	}

	for _, detail := range tradeConfigurationDetails {
		paramObject, err := paramObjectCreator(detail.Algorithm.Type, detail.Algorithm.Param)
		if err != nil {
			return nil, err
		}
		detail.Algorithm.ParamObject = paramObject
	}

	tradeSetDetail := &model.TradeSetDetail{
		TradeSet:       tradeSet,
		Configurations: tradeConfigurationDetails,
	}

	return tradeSetDetail, nil
}

// AddTradeSet is a method to add trade set
func AddTradeSet(tradeSetParam *model.TradeSetParam) error {
	tradeSet, err := db.GetTradeSetByName(tradeSetParam.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := db.AddTradeSet(tradeSetParam.Name, tradeSetParam.Type); err != nil {
				return err
			}

			tradeSet, err = db.GetTradeSetByName(tradeSetParam.Name)
			if err != nil {
				return err
			}
		} else {
			return err
		}
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
