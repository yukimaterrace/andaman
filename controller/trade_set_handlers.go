package controller

import (
	"net/http"
	"yukimaterrace/andaman/service"

	"github.com/labstack/echo/v4"
)

func getTradeSets(c echo.Context) error {
	_type, err := param(c.QueryParam("type")).tradeSetType()
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
