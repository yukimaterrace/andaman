package controller

import (
	"net/http"
	"strconv"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/model"

	"github.com/go-playground/validator/v10"
)

func paramError(err error) error {
	return &APIError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	}
}

type customValidator struct {
	validate *validator.Validate
}

func (v *customValidator) Validate(i interface{}) error {
	return v.validate.Struct(i)
}

type (
	intParam                int
	tradeSetTypeParam       model.TradeSetType
	tradeRunTypeParam       model.TradeRunType
	tradeModeParam          flow.TradeMode
	tradePairParam          model.TradePair
	timezoneParam           model.Timezone
	tradeDirectionParam     model.TradeDirection
	tradeAlgorithmTypeParam model.TradeAlgorithmType
)

func validateIota(param string, converter func(int64) model.IotaValidator) (int, error) {
	i, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, paramError(err)
	}

	if err := converter(i).IsValid(); err != nil {
		return 0, paramError(err)
	}

	return int(i), nil
}

func (p *intParam) UnmarshalParam(param string) error {
	i, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return paramError(err)
	}

	*p = intParam(i)
	return nil
}

func (p *tradeSetTypeParam) UnmarshalParam(param string) error {
	converter := func(i int64) model.IotaValidator {
		return model.TradeSetType(i)
	}

	i, err := validateIota(param, converter)
	if err != nil {
		return err
	}

	*p = tradeSetTypeParam(i)
	return nil
}

func (p *tradeRunTypeParam) UnmarshalParam(param string) error {
	converter := func(i int64) model.IotaValidator {
		return model.TradeRunType(i)
	}

	i, err := validateIota(param, converter)
	if err != nil {
		return err
	}

	*p = tradeRunTypeParam(i)
	return nil
}

func (p *tradeModeParam) UnmarshalParam(param string) error {
	converter := func(i int64) model.IotaValidator {
		return flow.TradeMode(i)
	}

	i, err := validateIota(param, converter)
	if err != nil {
		return err
	}

	*p = tradeModeParam(i)
	return nil
}

func (p *tradePairParam) UnmarshalParam(param string) error {
	converter := func(i int64) model.IotaValidator {
		return model.TradePair(i)
	}

	i, err := validateIota(param, converter)
	if err != nil {
		return err
	}

	*p = tradePairParam(i)
	return nil
}

func (p *timezoneParam) UnmarshalParam(param string) error {
	converter := func(i int64) model.IotaValidator {
		return model.Timezone(i)
	}

	i, err := validateIota(param, converter)
	if err != nil {
		return err
	}

	*p = timezoneParam(i)
	return nil
}

func (p *tradeDirectionParam) UnmarshalParam(param string) error {
	converter := func(i int64) model.IotaValidator {
		return model.TradeDirection(i)
	}

	i, err := validateIota(param, converter)
	if err != nil {
		return err
	}

	*p = tradeDirectionParam(i)
	return nil
}

func (p *tradeAlgorithmTypeParam) UnmarshalParam(param string) error {
	converter := func(i int64) model.IotaValidator {
		return model.TradeAlgorithmType(i)
	}

	i, err := validateIota(param, converter)
	if err != nil {
		return err
	}

	*p = tradeAlgorithmTypeParam(i)
	return nil
}
