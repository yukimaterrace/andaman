package controller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type customBinder struct {
	binder   *echo.DefaultBinder
	validate *validator.Validate
}

func newCustomBinder() *customBinder {
	return &customBinder{
		binder:   &echo.DefaultBinder{},
		validate: validator.New(),
	}
}

func (b *customBinder) Bind(i interface{}, c echo.Context) error {
	if err := b.binder.Bind(i, c); err != nil {
		return err
	}
	return b.validate.Struct(i)
}

func paramError(message string) *APIError {
	return &APIError{
		Code:    http.StatusBadRequest,
		Message: message,
	}
}
