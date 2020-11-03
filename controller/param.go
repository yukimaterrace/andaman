package controller

import (
	"errors"
	"net/http"
	"strconv"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/model"
)

func paramError(err error) error {
	return &APIError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
}

var errParamRequired = errors.New("param required")

type param string

func (p param) string(required bool) (string, error) {
	if required && p == "" {
		return "", paramError(errParamRequired)
	}

	return string(p), nil
}

func (p param) int(required bool, _default int) (int, error) {
	if !required && p == "" {
		return _default, nil
	}

	i, err := strconv.ParseInt(string(p), 10, 64)
	if err != nil {
		return 0, paramError(err)
	}
	return int(i), nil
}

func (p param) validateIota(required bool, converter func(int64) model.IotaValidator) (int, error) {
	if required && p == "" {
		return 0, paramError(errParamRequired)
	}

	i, err := strconv.ParseInt(string(p), 10, 64)
	if err != nil {
		return 0, paramError(err)
	}

	if err := converter(i).IsValid(); err != nil {
		return 0, paramError(err)
	}

	return int(i), nil
}

func (p param) tradeSetType(required bool) (model.TradeSetType, error) {
	converter := func(i int64) model.IotaValidator {
		return model.TradeSetType(i)
	}

	i, err := p.validateIota(required, converter)
	if err != nil {
		return 0, err
	}

	return model.TradeSetType(i), nil
}

func (p param) tradeRunType(required bool) (model.TradeRunType, error) {
	converter := func(i int64) model.IotaValidator {
		return model.TradeRunType(i)
	}

	i, err := p.validateIota(required, converter)
	if err != nil {
		return 0, err
	}

	return model.TradeRunType(i), nil
}

func (p param) tradeMode(required bool) (flow.TradeMode, error) {
	converter := func(i int64) model.IotaValidator {
		return flow.TradeMode(i)
	}

	i, err := p.validateIota(required, converter)
	if err != nil {
		return 0, err
	}

	return flow.TradeMode(i), nil
}

func (p param) tradePair(required bool) (model.TradePair, error) {
	converter := func(i int64) model.IotaValidator {
		return model.TradePair(i)
	}

	i, err := p.validateIota(required, converter)
	if err != nil {
		return 0, err
	}

	return model.TradePair(i), nil
}

func (p param) timezone(required bool) (model.Timezone, error) {
	converter := func(i int64) model.IotaValidator {
		return model.Timezone(i)
	}

	i, err := p.validateIota(required, converter)
	if err != nil {
		return 0, err
	}

	return model.Timezone(i), nil
}

func (p param) tradeDirection(required bool) (model.TradeDirection, error) {
	converter := func(i int64) model.IotaValidator {
		return model.TradeDirection(i)
	}

	i, err := p.validateIota(required, converter)
	if err != nil {
		return 0, err
	}

	return model.TradeDirection(i), nil
}

func (p param) tradeAlgorithmType(required bool) (model.TradeAlgorithmType, error) {
	converter := func(i int64) model.IotaValidator {
		return model.TradeAlgorithmType(i)
	}

	i, err := p.validateIota(required, converter)
	if err != nil {
		return 0, err
	}

	return model.TradeAlgorithmType(i), nil
}
