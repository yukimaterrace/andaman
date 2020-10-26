package model

import (
	"fmt"
	"strconv"
)

type (
	// definitions for API

	// TradeSetsResponse is a response for trade sets
	TradeSetsResponse struct {
		TradeSets []*TradeSet `json:"trade_sets"`
	}

	// definitions for Service

	// TradeParamObjectCreator is a definition for param object creator
	TradeParamObjectCreator func(_type TradeAlgorithmType, param string) (interface{}, error)

	// TradeAlgorithmDetail is a struct for trade algorithm detail
	TradeAlgorithmDetail struct {
		TradeAlgorithm
		ParamObject interface{}
	}

	// TradeConfigurationDetail is a struct for trade configurtation detail
	TradeConfigurationDetail struct {
		TradeConfiguration
		Algorithm TradeAlgorithmDetail
	}

	// TradeConfigurationKey is a key for trade configuration
	TradeConfigurationKey string

	// TradeSetDetail is a struct for trade set detail
	TradeSetDetail struct {
		*TradeSet
		Configurations []*TradeConfigurationDetail
	}
)

// Key is a method to calculate key of trade configuration detail
func (detail *TradeConfigurationDetail) Key() TradeConfigurationKey {
	tradePair := strconv.FormatInt(int64(detail.TradePair), 10)
	timezone := strconv.FormatInt(int64(detail.TradePair), 10)
	algorithmID := strconv.FormatInt(int64(detail.Algorithm.TradeAlgorithmID), 10)

	return TradeConfigurationKey(fmt.Sprintf("%s-%s-%s", tradePair, timezone, algorithmID))
}
