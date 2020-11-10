package service

import (
	"yukimaterrace/andaman/db"
	"yukimaterrace/andaman/model"
)

// GetOrders is a method to get orders
func GetOrders(
	tradeRunID int, state model.OrderState, tradePair model.TradePair, timezone model.Timezone,
	tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType, count int, offset int) (*model.OrdersResponse, error) {

	orders, err := db.GetOrders(tradeRunID, state, tradePair, timezone, tradeDirection, algorithmType, count, offset)
	if err != nil {
		return nil, err
	}

	totalProfit, err := db.GetTotalProfitForOrders(tradeRunID, state, tradePair, timezone, tradeDirection, algorithmType)
	if err != nil {
		return nil, err
	}

	all, err := db.GetCountForOrders(tradeRunID, state, tradePair, timezone, tradeDirection, algorithmType)
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

// GetTradeSummariesA is a method to get trade summaries A
func GetTradeSummariesA(tradeRunID int, start int, end int) (*model.TradeSummariesResponseA, error) {
	unrealizedProfit, err := db.GetTotalProfitByFilter1(tradeRunID, model.Open, start, end)
	if err != nil {
		return nil, err
	}

	realizedProfit, err := db.GetTotalProfitByFilter1(tradeRunID, model.Closed, start, end)
	if err != nil {
		return nil, err
	}

	open, err := db.GetTradeCountProfitByFilter1(tradeRunID, model.Open, start, end)
	if err != nil {
		return nil, err
	}

	closed, err := db.GetTradeCountProfitByFilter1(tradeRunID, model.Closed, start, end)
	if err != nil {
		return nil, err
	}

	var oc map[model.TradePair][2]*model.TradeCountProfit

	for key, cp := range open {
		oc[key] = [2]*model.TradeCountProfit{cp, nil}
	}

	for key, cp := range closed {
		if cps, ok := oc[key]; ok {
			oc[key] = [2]*model.TradeCountProfit{cps[0], cp}
		} else {
			oc[key] = [2]*model.TradeCountProfit{nil, cp}
		}
	}

	tradeSummaries := []*model.TradePairTradeSummary{}
	for key, cps := range oc {
		var ts model.TradeSummary
		if cps[0] != nil {
			ts.Open = *cps[0]
		}
		if cps[1] != nil {
			ts.Closed = *cps[1]
		}

		tpts := model.TradePairTradeSummary{
			TradePair:    key,
			TradeSummary: ts,
		}

		tradeSummaries = append(tradeSummaries, &tpts)
	}

	resp := &model.TradeSummariesResponseA{
		UnrealizedProfit: unrealizedProfit,
		RealizedProfit:   realizedProfit,
		TradeSummaries:   tradeSummaries,
	}

	return resp, nil
}

// GetTradeSummariesB is a method to get trade summaries B
func GetTradeSummariesB(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone, start int, end int,
	paramObjectCreator model.TradeParamObjectCreator) (*model.TradeSummariesResponseB, error) {

	unrealizedProfit, err := db.GetTotalProfitByFilter2(tradeRunID, model.Open, tradePair, timezone, start, end)
	if err != nil {
		return nil, err
	}

	realizedProfit, err := db.GetTotalProfitByFilter2(tradeRunID, model.Closed, tradePair, timezone, start, end)
	if err != nil {
		return nil, err
	}

	open, err := db.GetTradeCountProfitByFilter2(tradeRunID, model.Open, start, end)
	if err != nil {
		return nil, err
	}

	closed, err := db.GetTradeCountProfitByFilter2(tradeRunID, model.Closed, start, end)
	if err != nil {
		return nil, err
	}

	var oc map[model.TradeConfigurationDetail][2]*model.TradeCountProfit

	for key, cp := range open {
		oc[key] = [2]*model.TradeCountProfit{cp, nil}
	}

	for key, cp := range closed {
		if cps, ok := oc[key]; ok {
			oc[key] = [2]*model.TradeCountProfit{cps[0], cp}
		} else {
			oc[key] = [2]*model.TradeCountProfit{nil, cp}
		}
	}

	tradeSummaries := []*model.TradeConfigurationDetailTradeSummary{}
	for key, cps := range oc {
		var ts model.TradeSummary
		if cps[0] != nil {
			ts.Open = *cps[0]
		}
		if cps[1] != nil {
			ts.Closed = *cps[1]
		}

		key.Algorithm.ParamObject, err = paramObjectCreator(key.Algorithm.Type, key.Algorithm.Param)
		if err != nil {
			return nil, err
		}

		tcts := model.TradeConfigurationDetailTradeSummary{
			TradeConfiguration: key,
			TradeSummary:       ts,
		}

		tradeSummaries = append(tradeSummaries, &tcts)
	}

	resp := &model.TradeSummariesResponseB{
		UnrealizedProfit: unrealizedProfit,
		RealizedProfit:   realizedProfit,
		TradeSummaries:   tradeSummaries,
	}

	return resp, nil
}
