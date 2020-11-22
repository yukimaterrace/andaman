package model

import "errors"

var (
	// ErrWrongType is an error for wrong type
	ErrWrongType = errors.New("wrong type has been passed")
	// ErrInconsistentLogic is an error for inconsistent logic
	ErrInconsistentLogic = errors.New("inconsistent logic detected")

	// ErrTimezoneNotDefined is an error for timezone not defined
	ErrTimezoneNotDefined = errors.New("timezone not defined")
)
