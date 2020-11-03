package controller

import (
	"net/http"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/service"

	"github.com/labstack/echo/v4"
)

type ordersQueryParam struct {
	tradeRunID     int
	tradePair      model.TradePair
	timezone       model.Timezone
	tradeDirection model.TradeDirection
	algorithmType  model.TradeAlgorithmType
	count          int
	offset         int
}

func _getOrdersQueryParam(c echo.Context) (*ordersQueryParam, error) {
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

	count, err := param(c.QueryParam("count")).int(false, 20)
	if err != nil {
		return nil, err
	}

	offset, err := param(c.QueryParam("offset")).int(false, 0)
	if err != nil {
		return nil, err
	}

	return &ordersQueryParam{
		tradeRunID:     tradeRunID,
		tradePair:      tradePair,
		timezone:       timezone,
		tradeDirection: tradeDirection,
		algorithmType:  algorithmType,
		count:          count,
		offset:         offset,
	}, nil
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
