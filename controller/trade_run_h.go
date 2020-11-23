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

type (
	getTradeRunsParams struct {
		Type   model.TradeRunType `query:"type" validate:"required"`
		Count  int                `query:"count" validate:"min=0,max=100"`
		Offset int                `query:"offset" validate:"min=0"`
	}

	createTradeParams struct {
		TradeSetName string             `form:"trade_set_name" validate:"required"`
		Version      int                `form:"version" validate:"required"`
		Type         model.TradeRunType `form:"type" validate:"required"`
		Start        int                `form:"start" validate:"min=0"`
		End          int                `form:"end" validate:"min=0"`
		Parallel     int                `form:"parallel" validate:"min=0,max=8"`
	}

	changeTradeModeParams struct {
		TradeMode flow.TradeMode `query:"trade_mode" validate:"required"`
	}
)

func getTradeRuns(c echo.Context) error {
	p := getTradeRunsParams{
		Count:  20,
		Offset: 0,
	}
	if err := c.Bind(&p); err != nil {
		return err
	}

	resp, err := service.GetTradeRunDetailsResponse(p.Type, p.Count, p.Offset)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func createTrade(c echo.Context) error {
	p := createTradeParams{
		Parallel: 2,
	}
	if err := c.Bind(&p); err != nil {
		return err
	}

	if err := _createTrade(p.TradeSetName, p.Version, p.Type, p.Start, p.End, p.Parallel); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, success)
}

func changeTradeMode(c echo.Context) error {
	p := changeTradeModeParams{}
	if err := c.Bind(&p); err != nil {
		return err
	}

	_changeTradeMode(flow.TradeMode(p.TradeMode))
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

func _createTrade(
	tradeSetName string, tradeSetVersion int, tradeRunType model.TradeRunType, start int, end int, parallel int) error {
	if _, err := service.GetTradeSet(tradeSetName, tradeSetVersion); err != nil {
		return err
	}

	var _flow *flow.Flow
	switch tradeRunType {
	case model.OandaSimulation:
		_flow = factory.CreateSimulationFlow(tradeSetName, tradeSetVersion, time.Second, start, end, parallel)
	default:
		panic("unknown type")
	}

	if err := _checkAndSetFlow(_flow); err != nil {
		return err
	}

	go func(_flow *flow.Flow) {
		log.Printf("Trade to start by trade set %s [version %d]\n", tradeSetName, tradeSetVersion)

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
