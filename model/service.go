package model

import (
	"database/sql"
	"encoding/json"
	"time"
)

// GetTradeSets is a method to get trade sets
func GetTradeSets(_type TradeSetType, count int, offset int) (*TradeSetsResponse, error) {
	tradeSets, err := getTradeSetsByType(_type, count, offset)
	if err != nil {
		return nil, handleError(err)
	}

	return &TradeSetsResponse{TradeSets: tradeSets}, nil
}

// GetTradeSetDetail is a method to get trade set detail
func GetTradeSetDetail(name string, paramObjectCreator TradeParamObjectCreator) (*TradeSetDetail, error) {
	tradeSet, err := getTradeSetByName(name)
	if err != nil {
		return nil, err
	}

	tradeConfigurationDetails, err := getTradeConfigurationDetailsByTradeSetID(tradeSet.TradeSetID)
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

	tradeSetDetail := &TradeSetDetail{
		TradeSet:       tradeSet,
		Configurations: tradeConfigurationDetails,
	}

	return tradeSetDetail, nil
}

// AddTradeRun is a method to add trade run
func AddTradeRun(tradeSetName string) (*TradeRun, error) {
	tradeSet, err := getTradeSetByName(tradeSetName)
	if err != nil {
		return nil, err
	}

	now := int(time.Now().Unix())
	if err := addTradeRun(tradeSet.TradeSetID, Running, now, now); err != nil {
		return nil, err
	}

	tradeRun, err := getLastTradeRun()
	if err != nil {
		return nil, err
	}
	return tradeRun, nil
}

// AddCreatedOrder is a method to add created order
func AddCreatedOrder(
	tradeRunID int, brokerOrderID int, tradeConfigurationID int,
	units float64, tradeDirection TradeDirection, timeAtOpen int, priceAtOpen float64) error {

	err := addOrder(
		tradeRunID, brokerOrderID, tradeConfigurationID,
		units, tradeDirection, Open, 0, timeAtOpen, priceAtOpen, 0, 0,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrderForClose is a method to update order for close
func UpdateOrderForClose(tradeRunID int, brokerOrderID int, profit float64, timeAtClose int, priceAtClose float64) error {
	if err := updateOrderForClose(tradeRunID, brokerOrderID, Closed, profit, timeAtClose, priceAtClose); err != nil {
		return err
	}
	return nil
}

// UpdateOrderForProfit is a method to update order for profit
func UpdateOrderForProfit(tradeRunID int, brokerOrderID int, profit float64) error {
	order, err := getOrderByTradeRunAndBrokerOrder(tradeRunID, brokerOrderID)
	if err != nil {
		return err
	}

	if order.State != Open {
		return nil
	}

	if err := updateOrderForProfit(order.OrderID, profit); err != nil {
		return err
	}
	return nil
}

// AddTradeSet is a method to add trade set
func AddTradeSet(tradeSetParam *TradeSetParam) error {
	tradeSet, err := getTradeSetByName(tradeSetParam.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			if err := addTradeSet(tradeSetParam.Name, tradeSetParam.Type); err != nil {
				return err
			}

			tradeSet, err = getTradeSetByName(tradeSetParam.Name)
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
		tradeAlgorithm, err := getTradeAlgorithmByTypeAndParam(algorithmParam.Type, param)
		if err != nil {
			if err == sql.ErrNoRows {
				if err := addTradeAlgorithm(algorithmParam.Type, param, algorithmParam.TradeDirection); err != nil {
					return err
				}

				tradeAlgorithm, err = getTradeAlgorithmByTypeAndParam(algorithmParam.Type, param)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		tradeConfiguration, err := getTradeConfigurationByFields(configurationParam.TradePair, configurationParam.Timezone, tradeAlgorithm.TradeAlgorithmID)
		if err != nil {
			if err == sql.ErrNoRows {
				if err := addTradeConfiguration(configurationParam.TradePair, configurationParam.Timezone, tradeAlgorithm.TradeAlgorithmID); err != nil {
					return err
				}

				tradeConfiguration, err = getTradeConfigurationByFields(configurationParam.TradePair, configurationParam.Timezone, tradeAlgorithm.TradeAlgorithmID)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}

		_, err = getTradeSetConfigurationRelByTradeSetIDAndTradeConfigurationID(tradeSet.TradeSetID, tradeConfiguration.TradeConfigurationID)
		if err != nil {
			if err == sql.ErrNoRows {
				if err := addTradeSetConfigurationRel(tradeSet.TradeSetID, tradeConfiguration.TradeConfigurationID); err != nil {
					return err
				}
			} else {
				return err
			}
		}
	}

	return nil
}
