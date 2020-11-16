package controller

import (
	"net/http"
	"yukimaterrace/andaman/factory"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/service"

	"github.com/labstack/echo/v4"
)

type (
	getTradeSetsParams struct {
		_type  model.TradeSetType `query:"type" validate:"required"`
		count  int                `query:"count" validate:"min=0,max=100"`
		offset int                `query:"offset" validate:"min=0"`
	}

	addTradeSetByPresetParams struct {
		name string `form:"name" validate:"required"`
	}
)

func getTradeSets(c echo.Context) error {
	p := getTradeSetsParams{
		count:  20,
		offset: 0,
	}
	if err := c.Bind(&p); err != nil {
		return err
	}

	resp, err := service.GetTradeSets(model.TradeSetType(p._type), int(p.count), int(p.offset))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func addTradeSetByPreset(c echo.Context) error {
	p := addTradeSetByPresetParams{}
	if err := c.Bind(&p); err != nil {
		return err
	}

	switch p.name {
	case factory.SimulationTradeSetName:
		factory.AddSimulationTradeSet()

	case factory.GridSearchTradeSetName:
		factory.AddGridSearchTradeSet()

	default:
		return paramError("not found trade set")
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
