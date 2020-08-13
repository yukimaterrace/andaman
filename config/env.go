package config

import (
	"fmt"
	"os"
)

var (
	// OandaHost is config for oanda host
	OandaHost string

	// OandaToken is config for oanda token
	OandaToken string

	// RecordDir is config for record directory
	RecordDir string
)

func panicForEnv(name string) {
	panic(fmt.Sprintf("Environment %s has not been set", name))
}

func init() {
	var name string

	name = "OANDA_HOST"
	if OandaHost = os.Getenv(name); OandaHost == "" {
		panicForEnv(name)
	}

	name = "OANDA_TOKEN"
	if OandaToken = os.Getenv(name); OandaToken == "" {
		panicForEnv(name)
	}

	name = "RECORD_DIR"
	if RecordDir = os.Getenv(name); RecordDir == "" {
		panicForEnv(name)
	}
}
