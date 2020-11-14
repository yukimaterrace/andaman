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
		_type  tradeRunTypeParam `query:"type" validate:"required"`
		count  intParam          `query:"count" validate:"gte=0,lte=100"`
		offset intParam          `query:"offset" validate:"gte=0"`
	}

	createTradeParams struct {
		tradeSetName string            `query:"trade_set_name" validate:"required"`
		_type        tradeRunTypeParam `query:"type" validate:"required"`
		start        intParam          `query:"start" validate:"gte=0"`
		end          intParam          `query:"end" validate:"gte=0"`
	}

	changeTradeModeParams struct {
		tradeMode tradeModeParam `query:"trade_mode" validate:"required"`
	}
)

func getTradeRuns(c echo.Context) error {
	p := getTradeRunsParams{
		count:  20,
		offset: 0,
	}
	if err := c.Bind(&p); err != nil {
		return err
	}
	if err := c.Validate(&p); err != nil {
		return err
	}

	resp, err := service.GetTradeRunDetails(model.TradeRunType(p._type), int(p.count), int(p.offset))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, resp)
}

func createTrade(c echo.Context) error {
	p := createTradeParams{}
	if err := c.Bind(&p); err != nil {
		return err
	}
	if err := c.Validate(&p); err != nil {
		return err
	}

	if err := _createTrade(p.tradeSetName, model.TradeRunType(p._type), int(p.start), int(p.end)); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, success)
}

func changeTradeMode(c echo.Context) error {
	p := changeTradeModeParams{}
	if err := c.Bind(&p); err != nil {
		return err
	}
	if err := c.Validate(&p); err != nil {
		return err
	}

	_changeTradeMode(flow.TradeMode(p.tradeMode))
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
