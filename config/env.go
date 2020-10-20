package config

import (
	"log"
	"os"
)

var (
	// OandaHost is an env for oanda host
	OandaHost string
	// OandaToken is an env for oanda token
	OandaToken string
	// Units is an env for units
	Units string
	// DBUser is an env for db user
	DBUser string
	// DBPassword is an env for db password
	DBPassword string
	// DBHost is an env for DBHost
	DBHost string
	// DBPort is an env for DBPort
	DBPort string
	// APIKey is an env for APIKey
	APIKey string
)

func panicForEnv(name string) {
	log.Panicf("environment variable %s has not been set", name)
}

func init() {
	var name string

	name = "OANDA_HOST"
	OandaHost = os.Getenv(name)
	if OandaHost == "" {
		panicForEnv(name)
	}

	name = "OANDA_TOKEN"
	OandaToken = os.Getenv(name)
	if OandaToken == "" {
		panicForEnv(name)
	}

	name = "UNITS"
	Units = os.Getenv(name)
	if Units == "" {
		panicForEnv(name)
	}

	name = "DB_USER"
	DBUser = os.Getenv(name)
	if DBUser == "" {
		panicForEnv(name)
	}

	name = "DB_PASSWORD"
	DBPassword = os.Getenv(name)
	if DBPassword == "" {
		panicForEnv(name)
	}

	name = "DB_HOST"
	DBHost = os.Getenv(name)
	if DBHost == "" {
		panicForEnv(name)
	}

	name = "DB_PORT"
	DBPort = os.Getenv(name)
	if DBPort == "" {
		panicForEnv(name)
	}

	name = "API_KEY"
	APIKey = os.Getenv(name)
	if APIKey == "" {
		panicForEnv(name)
	}
}
