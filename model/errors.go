package model

import (
	"database/sql"
	"errors"
	"net/http"
)

// Error is an erro for model
type Error struct {
	Code    int
	Message string
}

func (err Error) Error() string {
	return err.Message
}

// HandleError is a method to handle error
func HandleError(err error) *Error {
	code := http.StatusInternalServerError
	if err == sql.ErrNoRows {
		code = http.StatusNotFound
	}
	return &Error{
		Code:    code,
		Message: err.Error(),
	}
}

// ErrNotFound is an error for not founc
var ErrNotFound = errors.New("not found")
