package controller

import (
	"net/http"
	"yukimaterrace/andaman/service"

	"github.com/labstack/echo/v4"
)

func getTradeRuns(c echo.Context) error {
	_type, err := param(c.QueryParam("type")).tradeRunType()
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

	resp, err := service.GetTradeRunDetails(_type, count, offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}
