package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
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
		log.Println(dbURI)

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

func CreateTable(model interface{}) bool {
	query := getTableCreationRequest(model)
	_, err := GetDBConnect().DB.Exec(query)
	if err != nil {
		log.Fatalf("Could not query db: %v", err)
	}

	return true
}

func getTableCreationRequest(model interface{}) string {
	queryFields := []string{}
	queryPrimary := []string{}
	varFullName := reflect.ValueOf(model).Type().String()
	varSlice := strings.Split(varFullName, ".")
	query := "CREATE TABLE IF NOT EXISTS " + varSlice[len(varSlice)-1]
	val := reflect.ValueOf(model).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tableField := ""
		switch field.Type.String() {
		case "string":
			tableField = fmt.Sprintf("%s VARCHAR(%s) NOT NULL DEFAULT ''", field.Tag.Get("field"), field.Tag.Get("len"))
		case "int":
			tableField = field.Tag.Get("field") + " INT NOT NULL default 0"
		}
		queryFields = append(queryFields, tableField)

		if field.Tag.Get("key") == "primary" {
			queryPrimary = append(queryPrimary, field.Tag.Get("field"))
		}

	}

	if len(queryPrimary) > 0 {
		tablePrimary := fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(queryPrimary, ","))
		queryFields = append(queryFields, tablePrimary)
	}

	return fmt.Sprintf("%s (%s)", query, strings.Join(queryFields, ", "))
}
