package db

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DBConnect struct {
	DB    *sqlx.DB
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

		if connect.DB, connect.Error = sqlx.Connect("mysql", dbURI); connect.Error != nil {
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

func CreateTable(model interface{}) {
	query := getTableCreationRequest(model)
	GetDBConnect().DB.MustExec(query)
}

func Insert(object interface{}) {
	query := getInsertRequest(object)
	tx := GetDBConnect().DB.MustBegin()
	tx.NamedExec(query, object)
	tx.Commit()

}

func getInsertRequest(object interface{}) string {
	var model string
	var insertMap, valueMap []string

	reflectValue := reflect.ValueOf(object)
	varFullName := reflectValue.Type().String()
	varSlice := strings.Split(varFullName, ".")
	model = varSlice[len(varSlice)-1]

	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Type().Field(i)

		if len(field.Tag.Get("db")) > 0 {
			insertMap = append(insertMap, field.Tag.Get("db"))
		}
	}

	for _, value := range insertMap {
		valueMap = append(valueMap, ":"+value)
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);", model, strings.Join(insertMap, ", "), strings.Join(valueMap, ", "))
}

func getTableCreationRequest(model interface{}) string {
	var queryFields, queryPrimary []string
	reflectValue := reflect.ValueOf(model)
	varFullName := reflectValue.Type().String()
	varSlice := strings.Split(varFullName, ".")
	query := fmt.Sprintf("CREATE TABLE %s", varSlice[len(varSlice)-1])
	val := reflectValue.Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		tableField := ""
		switch field.Type.String() {
		case "string":
			tableField = fmt.Sprintf("%s VARCHAR(%s) NOT NULL DEFAULT ''", field.Tag.Get("db"), field.Tag.Get("len"))
		case "int":
			if field.Tag.Get("extra") == "AUTO_INCREMENT" {
				tableField = fmt.Sprintf("%s INT %s", field.Tag.Get("db"), field.Tag.Get("extra"))
			} else {
				tableField = fmt.Sprintf("%s INT NOT NULL default 0", field.Tag.Get("db"))
			}
		}
		queryFields = append(queryFields, tableField)

		if field.Tag.Get("key") == "primary" {
			queryPrimary = append(queryPrimary, field.Tag.Get("db"))
		}

	}

	if len(queryPrimary) > 0 {
		tablePrimary := fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(queryPrimary, ", "))
		queryFields = append(queryFields, tablePrimary)
	}

	return fmt.Sprintf("%s (%s) ENGINE=InnoDB DEFAULT CHARSET=utf8;", query, strings.Join(queryFields, ", "))
}
