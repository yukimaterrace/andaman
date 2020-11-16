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
		tradeRunID     int                      `query:"trade_run_id" validate:"required"`
		orderState     model.OrderState         `query:"order_state" validate:"required"`
		tradePair      model.TradePair          `query:"trade_pair" validate:"required"`
		timezone       model.Timezone           `query:"timezone" validate:"required"`
		tradeDirection model.TradeDirection     `query:"trade_direction" validate:"required"`
		algorithmType  model.TradeAlgorithmType `query:"algorithm_type" validate:"required"`
		count          int                      `query:"count" validate:"min=0,max=100"`
		offset         int                      `query:"offset" validate:"min=0"`
	}

	getTradeSummariesAParams struct {
		tradeRunID int `query:"trade_run_id" validate:"required"`
		start      int `query:"start" validate:"required"`
	}

	getTradeSummariesBParams struct {
		tradeRunID int             `query:"trade_run_id" validate:"required"`
		tradePair  model.TradePair `query:"trade_pair" validate:"required"`
		timezone   model.Timezone  `query:"timezone" validate:"required"`
		start      int             `query:"start" validate:"required"`
	}

	getTradeCountProfitsParams struct {
		tradeRunID     int                      `query:"trade_run_id" validate:"required"`
		tradePair      model.TradePair          `query:"trade_pair" validate:"required"`
		timezone       model.Timezone           `query:"timezone" validate:"required"`
		tradeDirection model.TradeDirection     `query:"trade_direction" validate:"required"`
		algorithmType  model.TradeAlgorithmType `query:"algorithm_type" validate:"required"`
		start          int                      `query:"start" validate:"required"`
	}
)

func getOrders(c echo.Context) error {
	p := getOrdersQueryParams{
		count:  20,
		offset: 0,
	}
	if err := c.Bind(&p); err != nil {
		return err
	}

	resp, err := service.GetOrders(p.tradeRunID, p.orderState, p.tradePair, p.timezone, p.tradeDirection, p.algorithmType, p.count, p.offset)
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

	resp, err := service.GetTradeSummariesA(p.tradeRunID, p.start, int(time.Now().Unix()))
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

	resp, err := service.GetTradeSummariesB(
		p.tradeRunID, p.tradePair, p.timezone, p.start, int(time.Now().Unix()), trader.TradeParamObjectCreator,
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

	resp, err := service.GetTradeCountProfits(
		p.tradeRunID, p.tradePair, p.timezone, p.tradeDirection, p.algorithmType, p.start,
		int(time.Now().Unix()), 100, trader.TradeParamObjectCreator,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
