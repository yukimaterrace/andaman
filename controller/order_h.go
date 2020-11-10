package controller

import (
	"net/http"
	"time"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/service"
	"yukimaterrace/andaman/trader"

	"github.com/labstack/echo/v4"
)

type ordersQueryParam struct {
	*tradeGrainQueryParam
	*countOffsetQueryParam
}

type tradeGrainQueryParam struct {
	tradeRunID     int
	tradePair      model.TradePair
	timezone       model.Timezone
	tradeDirection model.TradeDirection
	algorithmType  model.TradeAlgorithmType
}

type countOffsetQueryParam struct {
	count  int
	offset int
}

func _getTradeGrainQueryParam(c echo.Context) (*tradeGrainQueryParam, error) {
	tradeRunID, err := param(c.QueryParam("trade_run_id")).int(true, 0)
	if err != nil {
		return nil, err
	}

	tradePair, err := param(c.QueryParam("trade_pair")).tradePair(true)
	if err != nil {
		return nil, err
	}

	timezone, err := param(c.QueryParam("timezone")).timezone(true)
	if err != nil {
		return nil, err
	}

	tradeDirection, err := param(c.QueryParam("trade_direction")).tradeDirection(true)
	if err != nil {
		return nil, err
	}

	algorithmType, err := param(c.QueryParam("trade_algorithm_type")).tradeAlgorithmType(true)
	if err != nil {
		return nil, err
	}

	return &tradeGrainQueryParam{
		tradeRunID:     tradeRunID,
		tradePair:      tradePair,
		timezone:       timezone,
		tradeDirection: tradeDirection,
		algorithmType:  algorithmType,
	}, nil
}

func _getCountOffsetQueryParam(c echo.Context) (*countOffsetQueryParam, error) {
	count, err := param(c.QueryParam("count")).int(false, 20)
	if err != nil {
		return nil, err
	}

	offset, err := param(c.QueryParam("offset")).int(false, 0)
	if err != nil {
		return nil, err
	}

	return &countOffsetQueryParam{
		count:  count,
		offset: offset,
	}, nil
}

func _getOrdersQueryParam(c echo.Context) (*ordersQueryParam, error) {
	tradeGrainQueryParam, err := _getTradeGrainQueryParam(c)
	if err != nil {
		return nil, err
	}

	countOffsetQueryParam, err := _getCountOffsetQueryParam(c)
	if err != nil {
		return nil, err
	}

	return &ordersQueryParam{tradeGrainQueryParam, countOffsetQueryParam}, nil
}

func _getOrders(c echo.Context, orderState model.OrderState) error {
	p, err := _getOrdersQueryParam(c)
	if err != nil {
		return err
	}

	resp, err := service.GetOrders(p.tradeRunID, orderState, p.tradePair, p.timezone, p.tradeDirection, p.algorithmType, p.count, p.offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func getOpenOrders(c echo.Context) error {
	return _getOrders(c, model.Open)
}

func getClosedOrders(c echo.Context) error {
	return _getOrders(c, model.Closed)
}

func getTradeSummariesA(c echo.Context) error {
	tradeRunID, err := param(c.QueryParam("trade_run_id")).int(true, 0)
	if err != nil {
		return err
	}

	start, err := param(c.QueryParam("start")).int(true, 0)
	if err != nil {
		return err
	}

	resp, err := service.GetTradeSummariesA(tradeRunID, start, int(time.Now().Unix()))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func getTradeSummariesB(c echo.Context) error {
	tradeRunID, err := param(c.QueryParam("trade_run_id")).int(true, 0)
	if err != nil {
		return err
	}

	tradePair, err := param(c.QueryParam("trade_pair")).tradePair(true)
	if err != nil {
		return err
	}

	timezone, err := param(c.QueryParam("timezone")).timezone(true)
	if err != nil {
		return err
	}

	start, err := param(c.QueryParam("start")).int(true, 0)
	if err != nil {
		return err
	}

	resp, err := service.GetTradeSummariesB(tradeRunID, tradePair, timezone, start, int(time.Now().Unix()), trader.TradeParamObjectCreator)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func getTradeCountProfits(c echo.Context) error {
	tradeRunID, err := param(c.QueryParam("trade_run_id")).int(true, 0)
	if err != nil {
		return err
	}

	tradePair, err := param(c.QueryParam("trade_pair")).tradePair(true)
	if err != nil {
		return err
	}

	timezone, err := param(c.QueryParam("timezone")).timezone(true)
	if err != nil {
		return err
	}

	tradeDirection, err := param(c.QueryParam("trade_direction")).tradeDirection(true)
	if err != nil {
		return err
	}

	algorithmType, err := param(c.QueryParam("trade_algorithm_type")).tradeAlgorithmType(true)
	if err != nil {
		return err
	}

	start, err := param(c.QueryParam("start")).int(true, 0)
	if err != nil {
		return err
	}

	resp, err := service.GetTradeCountProfits(
		tradeRunID, tradePair, timezone, tradeDirection, algorithmType, start, int(time.Now().Unix()), 100, trader.TradeParamObjectCreator,
	)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
