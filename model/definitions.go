package model

import (
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
		TradeSets []*TradeSet `json:"trade_sets"`
	}

	// definitions for service

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
		Type                TradeSetType               `json:"type"`
		ConfigurationParams []*TradeConfigurationParam `json:"configuration_params"`
	}
)

// Key is a method to calculate key of trade configuration detail
func (detail *TradeConfigurationDetail) Key() TradeConfigurationKey {
	tradePair := strconv.FormatInt(int64(detail.TradePair), 10)
	timezone := strconv.FormatInt(int64(detail.TradePair), 10)
	algorithmID := strconv.FormatInt(int64(detail.Algorithm.TradeAlgorithmID), 10)

	return TradeConfigurationKey(fmt.Sprintf("%s-%s-%s", tradePair, timezone, algorithmID))
}
