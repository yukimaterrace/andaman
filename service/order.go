package service

import (
	"yukimaterrace/andaman/db"
	"yukimaterrace/andaman/model"
)

// GetOpenOrders is a method to get open orders
func GetOpenOrders(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone,
	tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType, count int, offset int) (*model.OrdersResponse, error) {

	orders, err := db.GetOpenOrders(tradeRunID, tradePair, timezone, tradeDirection, algorithmType, count, offset)
	if err != nil {
		return nil, err
	}

	totalProfit, err := db.GetTotalProfitForOpenOrders(tradeRunID, tradePair, timezone, tradeDirection, algorithmType)
	if err != nil {
		return nil, err
	}

	all, err := db.GetCountForOpenOrders(tradeRunID, tradePair, timezone, tradeDirection, algorithmType)
	if err != nil {
		return nil, err
	}

	paging := &model.OffsetPaging{All: all, Count: len(orders), Offset: offset}
	return &model.OrdersResponse{Orders: orders, TotalProfit: totalProfit, Paging: paging}, nil
}

// GetClosedOrders is a method to get closed orders
func GetClosedOrders(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone,
	tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType, count int, offset int) (*model.OrdersResponse, error) {

	orders, err := db.GetClosedOrders(tradeRunID, tradePair, timezone, tradeDirection, algorithmType, count, offset)
	if err != nil {
		return nil, err
	}

	totalProfit, err := db.GetTotalProfitForClosedOrders(tradeRunID, tradePair, timezone, tradeDirection, algorithmType)
	if err != nil {
		return nil, err
	}

	all, err := db.GetCountForClosedOrders(tradeRunID, tradePair, timezone, tradeDirection, algorithmType)
	if err != nil {
		return nil, err
	}

	paging := &model.OffsetPaging{All: all, Count: len(orders), Offset: offset}
	return &model.OrdersResponse{Orders: orders, TotalProfit: totalProfit, Paging: paging}, nil
}

// AddCreatedOrder is a method to add created order
func AddCreatedOrder(
	tradeRunID int, brokerOrderID int, tradeConfigurationID int,
	units float64, tradeDirection model.TradeDirection, timeAtOpen int, priceAtOpen float64) error {

	err := db.AddOrder(
		tradeRunID, brokerOrderID, tradeConfigurationID,
		units, tradeDirection, model.Open, 0, timeAtOpen, priceAtOpen, 0, 0,
	)
	if err != nil {
		return err
	}
	return nil
}

// UpdateOrderForClose is a method to update order for close
func UpdateOrderForClose(tradeRunID int, brokerOrderID int, profit float64, timeAtClose int, priceAtClose float64) error {
	if err := db.UpdateOrderForClose(tradeRunID, brokerOrderID, model.Closed, profit, timeAtClose, priceAtClose); err != nil {
		return err
	}
	return nil
}

// UpdateOrderForProfit is a method to update order for profit
func UpdateOrderForProfit(tradeRunID int, brokerOrderID int, profit float64) error {
	order, err := db.GetOrderByTradeRunAndBrokerOrder(tradeRunID, brokerOrderID)
	if err != nil {
		return err
	}

	if order.State != model.Open {
		return nil
	}

	if err := db.UpdateOrderForProfit(order.OrderID, profit); err != nil {
		return err
	}
	return nil
}
