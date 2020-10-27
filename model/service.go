package model

import (
	"errors"
	"time"
)

// ErrOrderState is an error for order state
var ErrOrderState = errors.New("order state error")

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
		return ErrOrderState
	}

	if err := updateOrderForProfit(order.OrderID, profit); err != nil {
		return err
	}
	return nil
}
