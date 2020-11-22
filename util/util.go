package util

import (
	"fmt"
	"os"
)

// GetEnv gets environment variable
func GetEnv(name string) string {
	env := os.Getenv(name)
	if env == "" {
		panic(fmt.Sprintf("%s has not been set", name))
	}
	return env
}
