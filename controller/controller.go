package controller

import (
	"database/sql"
	"net/http"
	"strings"
	"yukimaterrace/andaman/model"
	"yukimaterrace/andaman/trader"
	"yukimaterrace/andaman/util"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CreateController is a factory method to create controller
func CreateController() *echo.Echo {
	apiKey := util.GetEnv("API_KEY")
	origins := util.GetEnv("ORIGINS")

	e := echo.New()

	e.HTTPErrorHandler = httpErrorHandler
	e.Binder = newCustomBinder()

	e.Use(authMiddleware(apiKey))
	e.Use(corsMiddleware(origins))
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

	model.TradeParamObjectCreator = trader.TradeParamObjectCreator

	return e
}

func authMiddleware(apiKey string) echo.MiddlewareFunc {
	return middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:api-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == apiKey, nil
		},
	})
}

func corsMiddleware(origins string) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: strings.Split(origins, ","),
	})
}

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
