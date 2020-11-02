package controller

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
	"yukimaterrace/andaman/factory"
	"yukimaterrace/andaman/flow"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/service"

	"github.com/labstack/echo/v4"
)

func getTradeRuns(c echo.Context) error {
	_type, err := param(c.QueryParam("type")).tradeRunType(true)
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

func createTrade(c echo.Context) error {
	tradeSetName, err := param(c.FormValue("trade_set_name")).string(true)
	if err != nil {
		return err
	}

	_type, err := param(c.QueryParam("type")).tradeRunType(true)
	if err != nil {
		return err
	}

	start, err := param(c.QueryParam("start")).int(false, 0)
	if err != nil {
		return err
	}

	end, err := param(c.QueryParam(("end"))).int(false, 0)
	if err != nil {
		return err
	}

	if err := _createTrade(tradeSetName, _type, start, end); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, success)
}

func changeTradeMode(c echo.Context) error {
	tradeMode, err := param(c.FormValue("trade_mode")).tradeMode(true)
	if err != nil {
		return err
	}

	_changeTradeMode(tradeMode)
	return c.JSON(http.StatusOK, success)
}

var currentFlow = struct {
	sync.Mutex
	flow *flow.Flow
}{}

func _checkAndSetFlow(flow *flow.Flow) error {
	currentFlow.Lock()
	defer currentFlow.Unlock()

	if currentFlow.flow != nil {
		return errors.New("another trade is runnning")
	}

	currentFlow.flow = flow
	return nil
}

func _unsetFlow() {
	currentFlow.Lock()
	defer currentFlow.Unlock()

	currentFlow.flow = nil
}

func _createTrade(tradeSetName string, tradeRunType model.TradeRunType, start int, end int) error {
	if _, err := service.GetTradeSetByName(tradeSetName); err != nil {
		return err
	}

	var _flow *flow.Flow
	switch tradeRunType {
	case model.OandaSimulation:
		_flow = factory.CreateSimulationFlow(tradeSetName, time.Minute, start, end)
	default:
		panic("unknown type")
	}

	if err := _checkAndSetFlow(_flow); err != nil {
		return err
	}

	go func(_flow *flow.Flow) {
		log.Printf("Trade to start by TradeSet %s\n", tradeSetName)

		_flow.Start()
		_flow.WaitForCompletion()

		_unsetFlow()
	}(_flow)

	return nil
}

func _changeTradeMode(mode flow.TradeMode) {
	currentFlow.Lock()
	defer currentFlow.Unlock()

	currentFlow.flow.ChangeTradeMode(mode)
}
