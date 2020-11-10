package service

import (
	"sort"
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

type tradePairTradeSummaries []*model.TradePairTradeSummary

func (s tradePairTradeSummaries) Len() int {
	return len(s)
}

func (s tradePairTradeSummaries) Less(i, j int) bool {
	return s[i].TradePair < s[j].TradePair
}

func (s tradePairTradeSummaries) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type tradeConfigurationTradeSummariesOrderedByKey []*model.TradeConfigurationTradeSummary

func (s tradeConfigurationTradeSummariesOrderedByKey) Len() int {
	return len(s)
}

func (s tradeConfigurationTradeSummariesOrderedByKey) Less(i, j int) bool {
	ci := s[i].TradeConfiguration
	cj := s[j].TradeConfiguration

	if ci.TradePair != cj.TradePair {
		return ci.TradePair < cj.TradePair
	}

	if ci.Timezone != cj.Timezone {
		return ci.Timezone < cj.Timezone
	}

	if ci.Algorithm.TradeDirection != cj.Algorithm.TradeDirection {
		return ci.Algorithm.TradeDirection < cj.Algorithm.TradeDirection
	}

	return ci.Algorithm.Type < cj.Algorithm.Type
}

func (s tradeConfigurationTradeSummariesOrderedByKey) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type tradeConfigurationTradeSummariesOrderedByProfit []*model.TradeConfigurationTradeSummary

func (s tradeConfigurationTradeSummariesOrderedByProfit) Len() int {
	return len(s)
}

func (s tradeConfigurationTradeSummariesOrderedByProfit) Less(i, j int) bool {
	if s[i].Closed.Profit != s[j].Closed.Profit {
		return s[i].Closed.Profit < s[j].Closed.Profit
	}
	return s[i].Open.Profit < s[j].Open.Profit
}

func (s tradeConfigurationTradeSummariesOrderedByProfit) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
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

	sort.Sort(tradePairTradeSummaries(tradeSummaries))
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

	open, err := db.GetTradeCountProfitByFilter2(tradeRunID, model.Open, tradePair, timezone, start, end)
	if err != nil {
		return nil, err
	}

	closed, err := db.GetTradeCountProfitByFilter2(tradeRunID, model.Closed, tradePair, timezone, start, end)
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

	tradeSummaries := []*model.TradeConfigurationTradeSummary{}
	for key, cps := range oc {
		var ts model.TradeSummary
		if cps[0] != nil {
			ts.Open = *cps[0]
		}
		if cps[1] != nil {
			ts.Closed = *cps[1]
		}

		var err error
		key.Algorithm.ParamObject, err = paramObjectCreator(key.Algorithm.Type, key.Algorithm.Param)
		if err != nil {
			return nil, err
		}

		tcts := model.TradeConfigurationTradeSummary{
			TradeConfiguration: key,
			TradeSummary:       ts,
		}

		tradeSummaries = append(tradeSummaries, &tcts)
	}

	sort.Sort(tradeConfigurationTradeSummariesOrderedByKey(tradeSummaries))
	resp := &model.TradeSummariesResponseB{
		UnrealizedProfit: unrealizedProfit,
		RealizedProfit:   realizedProfit,
		TradeSummaries:   tradeSummaries,
	}

	return resp, nil
}

func getTradeCount(totalCount int, positiveCount int) model.TradeCount {
	return model.TradeCount{
		PositiveProfitCount: positiveCount,
		NegativeProfitCount: totalCount - positiveCount,
	}
}

// GetTradeCountProfits is a method to get trade count profits
func GetTradeCountProfits(
	tradeRunID int, tradePair model.TradePair, timezone model.Timezone, tradeDirection model.TradeDirection, algorithmType model.TradeAlgorithmType,
	start int, end int, count int, paramObjectCreator model.TradeParamObjectCreator) (*model.TradeCountProfitsResponse, error) {

	unrealizedCount, err := db.GetCountByPeriod(tradeRunID, model.Open, tradePair, timezone, tradeDirection, algorithmType, start, end)
	if err != nil {
		return nil, err
	}

	unrealizedPositiveCount, err := db.GetPositiveProfitCountByPeriod(tradeRunID, model.Open, tradePair, timezone, tradeDirection, algorithmType, start, end)
	if err != nil {
		return nil, err
	}

	realizedCount, err := db.GetCountByPeriod(tradeRunID, model.Closed, tradePair, timezone, tradeDirection, algorithmType, start, end)
	if err != nil {
		return nil, err
	}

	realizedPositiveCount, err := db.GetPositiveProfitCountByPeriod(tradeRunID, model.Closed, tradePair, timezone, tradeDirection, algorithmType, start, end)
	if err != nil {
		return nil, err
	}

	cps, err := db.GetTradeConfigurationTradeCountProfits(tradeRunID, tradePair, timezone, tradeDirection, algorithmType, start, end, count, 0)
	if err != nil {
		return nil, err
	}

	resp := &model.TradeCountProfitsResponse{
		UnrealizedTradeCount: getTradeCount(unrealizedCount, unrealizedPositiveCount),
		RealizedTradeCount:   getTradeCount(realizedCount, realizedPositiveCount),
		TradeCountProfit:     cps,
	}

	return resp, nil
}

// GetTradeConfigurationGroupSummaries is a method to get trade configuration group summaries
func GetTradeConfigurationGroupSummaries(tradeRunID int, count int, offset int) (*model.TradeConfigurationGroupSummariesResponse, error) {
	return nil, nil
}
