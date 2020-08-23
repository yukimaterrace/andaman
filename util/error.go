package util

import (
	"errors"
	"fmt"
	"os"
)

var (
	// ErrWrongType is an error for wrong type
	ErrWrongType = errors.New("wrong type has been passed")
)

// GetEnv gets environment variable
func GetEnv(name string) string {
	env := os.Getenv(name)
	if env == "" {
		panic(fmt.Sprintf("%s has not been set", env))
	}
	return env
}
