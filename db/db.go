package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
)

type DBConnect struct {
	DB    *sql.DB
	Error error
}

var connect *DBConnect
var once sync.Once

func GetDBConnect() *DBConnect {
	once.Do(func() {
		var (
			connectionName = getEnv("CLOUDSQL_CONNECTION_NAME")
			user           = getEnv("CLOUDSQL_USER")
			dbName         = getEnv("CLOUDSQL_DATABASE_NAME")
			password       = getEnv("CLOUDSQL_PASSWORD")
			socket         = getEnv("CLOUDSQL_SOCKET_PREFIX")
		)
		connect = &DBConnect{}

		dbURI := fmt.Sprintf("%s:%s@unix(%s/%s)/%s", user, password, socket, connectionName, dbName)

		if connect.DB, connect.Error = sql.Open("mysql", dbURI); connect.Error != nil {
			panic(fmt.Sprintf("DB: %v", connect.Error))
		}
	})

	return connect
}

func getEnv(key string) string {
	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("%s environment variable not set.", key)
	}

	return value
}

func Query(query string) *sql.Rows {
	rows, err := GetDBConnect().DB.Query(query)

	if err != nil {
		log.Fatalf("Could not query db: %v", err)
	}
	defer rows.Close()

	return rows
}
