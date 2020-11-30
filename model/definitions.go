package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type (
	// definitions for API

	// SuccessResponse is a response for success
	SuccessResponse struct {
		Message string `json:"message"`
	}

	// TradeSetsResponse is a response for trade sets
	TradeSetsResponse struct {
		TradeSets []*TradeSet   `json:"trade_sets"`
		Paging    *OffsetPaging `json:"paging"`
	}

	// TradeSetDetailResponse is a response for trade set detail
	TradeSetDetailResponse struct {
		TradeSet *TradeSetDetail `json:"trade_set"`
	}

	// TradeRunDetailsResponse is a response for trade run details
	TradeRunDetailsResponse struct {
		TradeRuns []*TradeRunDetail `json:"trade_runs"`
		Paging    *OffsetPaging     `json:"paging"`
	}

	// OrdersResponse is a response for orders
	OrdersResponse struct {
		Orders      []*Order
		TotalProfit float64
		Paging      *OffsetPaging
	}

	// TradeSummariesResponseA is a A response for trade summaries
	TradeSummariesResponseA struct {
		UnrealizedProfit float64                  `json:"unrealized_profit"`
		RealizedProfit   float64                  `json:"realized_profit"`
		TradeSummaries   []*TradePairTradeSummary `json:"trade_summaries"`
	}

	// TradeSummariesResponseB is a B response for trade summaries
	TradeSummariesResponseB struct {
		UnrealizedProfit float64                           `json:"unrealized_profit"`
		RealizedProfit   float64                           `json:"realized_profit"`
		TradeSummaries   []*TradeConfigurationTradeSummary `json:"trade_summaries"`
	}

	// TradeCountProfitsResponse is a response for trade count profits
	TradeCountProfitsResponse struct {
		TradeCount        *TradeCount                           `json:"trade_count"`
		TradeCountProfits []*TradeConfigurationTradeCountProfit `json:"trade_count_profits"`
	}

	// TradeConfigurationGroupSummariesResponse is a struct for trade configuration grouup summaries response
	TradeConfigurationGroupSummariesResponse struct {
		TradeConfigurationGroupSummaries []*TradeConfigurationGroupSummary `json:"group_summaries"`
		Paging                           *OffsetPaging                     `json:"paging"`
	}

	// definitions for service

	// OffsetPaging is a struct for offset paging
	OffsetPaging struct {
		All    int `json:"all"`
		Count  int `json:"count"`
		Offset int `json:"offset"`
	}
)

// TradeParamObjectCreator is a definition for param object creator
var TradeParamObjectCreator func(_type TradeAlgorithmType, param string) (interface{}, error)

// TradeAlgorithmDetail is a struct for trade algorithm detail
type TradeAlgorithmDetail struct {
	TradeAlgorithm
}

// MarshalJSON is an implementation of Marshaler for TradeAlgorithmDetail
func (detail *TradeAlgorithmDetail) MarshalJSON() ([]byte, error) {
	if TradeParamObjectCreator == nil {
		return nil, ErrInconsistentLogic
	}

	paramObject, err := TradeParamObjectCreator(detail.Type, detail.Param)
	if err != nil {
		return nil, err
	}

	var target = struct {
		TradeAlgorithm
		ParamObject interface{} `json:"param"`
	}{
		detail.TradeAlgorithm,
		paramObject,
	}
	return json.Marshal(&target)
}

type (
	// TradeConfigurationDetail is a struct for trade configurtation detail
	TradeConfigurationDetail struct {
		TradeConfiguration
		Algorithm TradeAlgorithmDetail `json:"algorithm"`
	}

	// TradeConfigurationKey is a key for trade configuration
	TradeConfigurationKey string

	// TradeConfigurationGroup is a struct for trade configuration group
	TradeConfigurationGroup struct {
		TradePair          TradePair          `json:"trade_pair"`
		Timezone           Timezone           `json:"timezone"`
		TradeDirection     TradeDirection     `json:"trade_direction"`
		TradeAlgorithmType TradeAlgorithmType `json:"trade_algorithm_type"`
	}

	// TradeSetDetail is a struct for trade set detail
	TradeSetDetail struct {
		*TradeSet
		Configurations []*TradeConfigurationDetail `json:"configurations"`
	}

	// TradeRunDetail is a struct for trade run detail
	TradeRunDetail struct {
		TradeRun
		TradeSet TradeSet `json:"trade_set"`
	}

	// TradeCountProfit is a struct for trade count profit
	TradeCountProfit struct {
		Count  int     `json:"count"`
		Profit float64 `json:"profit"`
	}

	// TradeCount is a struct for trade count
	TradeCount struct {
		PositiveProfitCount int `json:"positive_profit_count"`
		NegativeProfitCount int `json:"negative_profit_count"`
	}

	// TradeSummary is a struct for trade summary
	TradeSummary struct {
		Open   TradeCountProfit `json:"open"`
		Closed TradeCountProfit `json:"closed"`
	}

	// TradePairTradeSummary is a struct for trade pair trade summary
	TradePairTradeSummary struct {
		TradePair TradePair `json:"trade_pair"`
		TradeSummary
	}

	// TradeConfigurationTradeSummary is a struct for trade configuration detail trade summary
	TradeConfigurationTradeSummary struct {
		TradeConfiguration TradeConfigurationDetail `json:"trade_configuration"`
		TradeSummary
	}

	// TradeConfigurationTradeCountProfit is a struct for trade configuration trade count profit
	TradeConfigurationTradeCountProfit struct {
		TradeConfiguration TradeConfigurationDetail `json:"trade_configuration"`
		TradeCountProfit
	}

	// TradeConfigurationGroupSummary is a struct for trade configuration group
	TradeConfigurationGroupSummary struct {
		TradeConfigurationGroup            *TradeConfigurationGroup            `json:"trade_configuration_group"`
		TradeCount                         *TradeCount                         `json:"trade_count"`
		TradeConfigurationTradeCountProfit *TradeConfigurationTradeCountProfit `json:"first_trade_configuration"`
	}

	// definitions for parameters in service

	// TradeAlgorithmParam is a param for trade algorithm
	TradeAlgorithmParam struct {
		Type           TradeAlgorithmType `json:"type"`
		Param          interface{}        `json:"param"`
		TradeDirection TradeDirection     `json:"trade_direction"`
	}

	// TradeConfigurationParam is a param for trade configuration
	TradeConfigurationParam struct {
		TradePair      TradePair            `json:"trade_pair"`
		Timezone       Timezone             `json:"timezone"`
		AlgorithmParam *TradeAlgorithmParam `json:"algorithm_param"`
	}

	// TradeSetParam is a param for trade set
	TradeSetParam struct {
		Name                string                     `json:"name"`
		Version             int                        `json:"version"`
		Type                TradeSetType               `json:"type"`
		ConfigurationParams []*TradeConfigurationParam `json:"configuration_params"`
	}
)

// Key is a method to calculate key of trade configuration detail
func (detail *TradeConfigurationDetail) Key() TradeConfigurationKey {
	tradePair := strconv.FormatInt(int64(detail.TradePair), 10)
	timezone := strconv.FormatInt(int64(detail.Timezone), 10)
	algorithmID := strconv.FormatInt(int64(detail.Algorithm.TradeAlgorithmID), 10)

	return TradeConfigurationKey(fmt.Sprintf("%s-%s-%s", tradePair, timezone, algorithmID))
}
