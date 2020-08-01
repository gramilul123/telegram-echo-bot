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

// CreateTable creates table
func CreateTable(model interface{}) {
	query := getTableCreationRequest(model)
	GetDBConnect().DB.MustExec(query)
}

// GetSelectRequest return select request
func GetSelectRequest(object interface{}, selectWhere map[string]interface{}) string {
	var model string
	var selectMap []string

	reflectValue := reflect.ValueOf(object)
	varFullName := reflectValue.Type().String()
	varSlice := strings.Split(varFullName, ".")
	model = varSlice[len(varSlice)-1]

	selectMap = append(selectMap, fmt.Sprintf("%d", 1))
	for field, value := range selectWhere {
		switch value.(type) {
		case int, int64:
			selectMap = append(selectMap, fmt.Sprintf("%s = %d", field, value))
		case string:
			selectMap = append(selectMap, fmt.Sprintf("%s = '%s'", field, value))
		}
	}

	return fmt.Sprintf("SELECT * FROM %s WHERE %s", model, strings.Join(selectMap, " AND "))
}

// Insert inserts row
func Insert(object interface{}) {
	query := getInsertRequest(object)
	GetDBConnect().DB.NamedExec(query, object)
}

// Delete delete row by where
func Delete(object interface{}, deleteWhere map[string]interface{}) {
	var model, query string
	var deletetMap []string

	reflectValue := reflect.ValueOf(object)
	varFullName := reflectValue.Type().String()
	varSlice := strings.Split(varFullName, ".")
	model = varSlice[len(varSlice)-1]

	deletetMap = append(deletetMap, fmt.Sprintf("%d", 1))
	for field, _ := range deleteWhere {
		deletetMap = append(deletetMap, fmt.Sprintf("%s = :%s", field, field))
	}

	query = fmt.Sprintf("DELETE FROM %s WHERE %s", model, strings.Join(deletetMap, " AND "))

	GetDBConnect().DB.NamedExec(query, deleteWhere)
}

// UpdateRow func updates row from object by field
func UpdateRow(object interface{}, field string) {
	query := getUpdatetRequest(object, field)
	GetDBConnect().DB.NamedExec(query, object)
}

// getUpdatetRequest func returns update request string
func getUpdatetRequest(object interface{}, rowField string) string {
	var model, whereValue string
	var updateMap []string

	reflectValue := reflect.ValueOf(object)
	varFullName := reflectValue.Type().String()
	varSlice := strings.Split(varFullName, ".")
	model = varSlice[len(varSlice)-1]

	for i := 0; i < reflectValue.NumField(); i++ {
		field := reflectValue.Type().Field(i)

		if len(field.Tag.Get("db")) > 0 && len(field.Tag.Get("extra")) == 0 {
			updateMap = append(updateMap, fmt.Sprintf("%s=:%s", field.Tag.Get("db"), field.Tag.Get("db")))

			if field.Tag.Get("db") == rowField {
				switch field.Type.String() {
				case "string":
					whereValue = fmt.Sprintf("'%s'", reflectValue.Field(i).Interface())
				case "int", "int64":
					whereValue = fmt.Sprintf("%d", reflectValue.Field(i).Interface())
				}
			}
		}
	}

	return fmt.Sprintf("UPDATE %s SET %s WHERE %s=%s;", model, strings.Join(updateMap, ", "), rowField, whereValue)
}

// getInsertRequest returns insert row request
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

// getTableCreationRequest returns create table request
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

			switch field.Tag.Get("type") {
			case "text":
				tableField = fmt.Sprintf("%s TEXT", field.Tag.Get("db"))
			default:
				tableField = fmt.Sprintf("%s VARCHAR(%s) NOT NULL DEFAULT ''", field.Tag.Get("db"), field.Tag.Get("len"))
			}

		case "int", "int64":

			switch field.Tag.Get("extra") {
			case "AUTO_INCREMENT":
				tableField = fmt.Sprintf("%s INT %s", field.Tag.Get("db"), field.Tag.Get("extra"))
			default:
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
