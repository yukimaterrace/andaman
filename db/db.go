package db

import (
	"database/sql"
	"fmt"
	"yukimaterrace/andaman/util"

	_ "github.com/go-sql-driver/mysql" // import driver
)

var db *sql.DB

func init() {
	dbUser := util.GetEnv("DB_USER")
	dbPassword := util.GetEnv("DB_PASSWORD")
	dbHost := util.GetEnv("DB_HOST")
	dbPort := util.GetEnv("DB_PORT")

	sourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/andaman", dbUser, dbPassword, dbHost, dbPort)

	var err error
	db, err = sql.Open("mysql", sourceName)
	if err != nil {
		panic(err)
	}

	if db.Ping() != nil {
		panic(err)
	}
}
