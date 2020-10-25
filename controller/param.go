package controller

import (
	"net/http"
	"strconv"
	"yukimaterrace/andaman/model"
)

func paramError(err error) error {
	return &APIError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
}

type param string

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

func (p param) tradeSetType() (model.TradeSetType, error) {
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
