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
		Type   model.TradeSetType `query:"type" validate:"required"`
		Count  int                `query:"count" validate:"min=0,max=100"`
		Offset int                `query:"offset" validate:"min=0"`
	}

	addTradeSetByPresetParams struct {
		Name string `form:"name" validate:"required"`
	}
)

func getTradeSets(c echo.Context) error {
	p := getTradeSetsParams{
		Count:  20,
		Offset: 0,
	}
	if err := c.Bind(&p); err != nil {
		return err
	}

	resp, err := service.GetTradeSets(p.Type, p.Count, p.Offset)
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

	switch p.Name {
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
