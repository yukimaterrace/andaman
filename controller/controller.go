package controller

import (
	"database/sql"
	"net/http"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/util"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var apiKey string

// CreateController is a factory method to create controller
func CreateController() *echo.Echo {
	apiKey = util.GetEnv("API_KEY")

	e := echo.New()

	e.HTTPErrorHandler = httpErrorHandler
	e.Binder = newCustomBinder()

	e.Use(authMiddleware)
	e.Use(middleware.Logger())

	e.GET("/api/trade_sets", getTradeSets)
	e.GET("/api/trade_set", getTradeSet)

	e.POST("/api/add_trade_set_by_preset", addTradeSetByPreset)
	e.POST("/api/add_trade_set_by_param", addTradeSetByParam)

	e.GET("/api/trade_runs", getTradeRuns)
	e.POST("/api/create_trade", createTrade)
	e.POST("/api/change_trade_mode", changeTradeMode)

	e.GET("/api/orders", getOrders)

	e.GET("/api/trade_summaries_a", getTradeSummariesA)
	e.GET("/api/trade_summaries_b", getTradeSummariesB)

	e.GET("/api/trade_count_profits", getTradeCountProfits)
	e.GET("/api/trade_configuration_group_summaries", getTradeConfigurationGroupSummaries)

	return e
}

var authMiddleware = middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
	KeyLookup: "header:api-key",
	Validator: func(key string, c echo.Context) (bool, error) {
		return key == apiKey, nil
	},
})

// APIError is an error for api
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (err APIError) Error() string {
	return err.Message
}

func httpErrorHandler(err error, c echo.Context) {
	var apiErr *APIError

	switch _err := err.(type) {
	case *APIError:
		apiErr = _err

	case *echo.HTTPError:
		apiErr = &APIError{
			Code:    _err.Code,
			Message: _err.Error(),
		}

	case validator.ValidationErrors:
		apiErr = &APIError{
			Code:    http.StatusBadRequest,
			Message: _err.Error(),
		}

	default:
		switch _err {
		case sql.ErrNoRows:
			apiErr = &APIError{
				Code:    http.StatusNotFound,
				Message: "Not Found",
			}

		default:
			apiErr = &APIError{
				Code:    http.StatusInternalServerError,
				Message: _err.Error(),
			}
		}
	}

	c.JSON(apiErr.Code, apiErr)
}

var success = &model.SuccessResponse{Message: "success"}
