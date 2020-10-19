package model

import (
	"database/sql"
	"fmt"
	"log"
	"os"
)

var db *sql.DB

func panicForEnv(name string) {
	log.Panicf("environment %s has not been set", name)
}

func init() {
	name := "DB_USER"
	user := os.Getenv(name)
	if user == "" {
		panicForEnv(name)
	}

	name = "DB_PASSWORD"
	password := os.Getenv(name)
	if password == "" {
		panicForEnv(name)
	}

	name = "DB_HOST"
	host := os.Getenv(name)
	if host == "" {
		panicForEnv(name)
	}

	name = "DB_PORT"
	port := os.Getenv(name)
	if port == "" {
		panicForEnv(name)
	}

	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/andaman", user, password, host, port)

	var err error
	db, err = sql.Open("mysql", sourceName)
	if err != nil {
		panic(err)
	}

	if db.Ping() != nil {
		panic(err)
	}
}
