package controller

import (
	"net/http"
	"time"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/service"
	"yukimaterrace/andaman/trader"

	"github.com/labstack/echo/v4"
)

type (
	getOrdersQueryParams struct {
		TradeRunID     int                      `query:"trade_run_id" validate:"required"`
		OrderState     model.OrderState         `query:"order_state" validate:"required"`
		TradePair      model.TradePair          `query:"trade_pair" validate:"required"`
		Timezone       model.Timezone           `query:"timezone" validate:"required"`
		TradeDirection model.TradeDirection     `query:"trade_direction" validate:"required"`
		AlgorithmType  model.TradeAlgorithmType `query:"algorithm_type" validate:"required"`
		Count          int                      `query:"count" validate:"min=0,max=100"`
		Offset         int                      `query:"offset" validate:"min=0"`
	}

	getTradeSummariesAParams struct {
		TradeRunID int `query:"trade_run_id" validate:"required"`
		Start      int `query:"start" validate:"required"`
	}

	getTradeSummariesBParams struct {
		TradeRunID int             `query:"trade_run_id" validate:"required"`
		TradePair  model.TradePair `query:"trade_pair" validate:"required"`
		Timezone   model.Timezone  `query:"timezone" validate:"required"`
		Start      int             `query:"start" validate:"required"`
	}

	getTradeCountProfitsParams struct {
		TradeRunID     int                      `query:"trade_run_id" validate:"required"`
		TradePair      model.TradePair          `query:"trade_pair" validate:"required"`
		Timezone       model.Timezone           `query:"timezone" validate:"required"`
		TradeDirection model.TradeDirection     `query:"trade_direction" validate:"required"`
		AlgorithmType  model.TradeAlgorithmType `query:"algorithm_type" validate:"required"`
	}

	getTradeConfigurationGroupSummariesParams struct {
		TradeRunID int `query:"trade_run_id" validate:"required"`
		Count      int `query:"count" validate:"min=0,max=100"`
		Offset     int `query:"offset" validate:"min=0"`
	}
)

func getOrders(c echo.Context) error {
	p := getOrdersQueryParams{
		Count:  20,
		Offset: 0,
	}
	if err := c.Bind(&p); err != nil {
		return err
	}

	resp, err := service.GetOrdersResponse(p.TradeRunID, p.OrderState, p.TradePair, p.Timezone, p.TradeDirection, p.AlgorithmType, p.Count, p.Offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func getTradeSummariesA(c echo.Context) error {
	p := getTradeSummariesAParams{}
	if err := c.Bind(&p); err != nil {
		return err
	}

	resp, err := service.GetTradeSummariesAResponse(p.TradeRunID, p.Start, int(time.Now().Unix()))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func getTradeSummariesB(c echo.Context) error {
	p := getTradeSummariesBParams{}
	if err := c.Bind(&p); err != nil {
		return err
	}

	resp, err := service.GetTradeSummariesBResposne(
		p.TradeRunID, p.TradePair, p.Timezone, p.Start, int(time.Now().Unix()), trader.TradeParamObjectCreator,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func getTradeCountProfits(c echo.Context) error {
	p := getTradeCountProfitsParams{}
	if err := c.Bind(&p); err != nil {
		return err
	}

	resp, err := service.GetTradeCountProfitsResponse(
		p.TradeRunID, p.TradePair, p.Timezone, p.TradeDirection, p.AlgorithmType, 100, trader.TradeParamObjectCreator,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func getTradeConfigurationGroupSummaries(c echo.Context) error {
	p := getTradeConfigurationGroupSummariesParams{}
	if err := c.Bind(&p); err != nil {
		return err
	}

	resp, err := service.GetTradeConfigurationGroupSummariesResponse(p.TradeRunID, p.Count, p.Offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
