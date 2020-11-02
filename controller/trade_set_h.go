package controller

import (
	"net/http"
	"yukimaterrace/andaman/factory"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/service"

	"github.com/labstack/echo/v4"
)

func getTradeSets(c echo.Context) error {
	_type, err := param(c.QueryParam("type")).tradeSetType(true)
	if err != nil {
		return err
	}

	count, err := param(c.QueryParam("count")).int(false, 20)
	if err != nil {
		return err
	}

	offset, err := param(c.QueryParam("offset")).int(false, 0)
	if err != nil {
		return err
	}

	resp, err := service.GetTradeSets(_type, count, offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func addTradeSetByPreset(c echo.Context) error {
	name, err := param(c.FormValue("name")).string(true)
	if err != nil {
		return err
	}

	switch name {
	case factory.SimulationTradeSetName:
		factory.AddSimulationTradeSet()

	case factory.GridSearchTradeSetName:
		factory.AddGridSearchTradeSet()

	default:
		return paramError(model.ErrNotFound)
	}

	return c.JSON(http.StatusOK, success)
}

func addTradeSetByParam(c echo.Context) error {
	param := &model.TradeSetParam{}
	if err := c.Bind(param); err != nil {
		return err
	}

	if err := service.AddTradeSet(param); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, success)
}
