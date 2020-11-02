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

func (p param) tradeSetType(required bool) (model.TradeSetType, error) {
	if required && p == "" {
		return 0, errParamRequired
	}

	i, err := strconv.ParseInt(string(p), 10, 64)
	if err != nil {
		return 0, paramError(err)
	}

	_type := model.TradeSetType(i)
	if err := _type.IsValid(); err != nil {
		return 0, paramError(err)
	}

	return _type, nil
}

func (p param) tradeRunType(required bool) (model.TradeRunType, error) {
	if required && p == "" {
		return 0, errParamRequired
	}

	i, err := strconv.ParseInt(string(p), 10, 64)
	if err != nil {
		return 0, paramError(err)
	}

	_type := model.TradeRunType(i)
	if err := _type.IsValid(); err != nil {
		return 0, paramError(err)
	}

	return _type, nil
}

func (p param) tradeMode(required bool) (flow.TradeMode, error) {
	if required && p == "" {
		return 0, errParamRequired
	}

	i, err := strconv.ParseInt(string(p), 10, 64)
	if err != nil {
		return 0, paramError(err)
	}

	_type := flow.TradeMode(i)
	if err := _type.IsValid(); err != nil {
		return 0, paramError(err)
	}

	return _type, nil
}
