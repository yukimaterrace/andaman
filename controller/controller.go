package controller

import (
	"net/http"
	"yukimaterrace/andaman/config"
	"yukimaterrace/andaman/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// CreateController is a factory method to create controller
func CreateController() *echo.Echo {
	e := echo.New()

	e.HTTPErrorHandler = httpErrorHandler

	e.Use(authMiddleware)
	e.Use(middleware.Logger())

	e.GET("/trade_sets", getTradeSets)

	return e
}

var authMiddleware = middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
	KeyLookup: "header:api-key",
	Validator: func(key string, c echo.Context) (bool, error) {
		return key == config.APIKey, nil
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

	case *model.Error:
		apiErr = &APIError{
			Code:    _err.Code,
			Message: _err.Message,
		}

	default:
		apiErr = &APIError{
			Code:    http.StatusInternalServerError,
			Message: _err.Error(),
		}
	}

	c.JSON(apiErr.Code, apiErr)
}
