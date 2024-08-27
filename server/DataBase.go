package server

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

// var goToSQLTypes = map[string]string{
// 	"int":     "integer",
// 	"int64":   "bigint",
// 	"string":  "varchar",
// 	"float32": "numeric",
// 	"float64": "decimal",
// }

var (
	errIsntStruct = errors.New("item value must be struct")
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DataBase struct {
	ConnectionString string
	DBName           string
	connection       *sql.DB
	user             string
	password         string
}

func NewDataBase(login string, password string, dbname string, sslmode string) *DataBase {
	var db *DataBase = new(DataBase)
	db.DBName = dbname
	db.user = login
	db.password = password
	db.ConnectionString = fmt.Sprintf("user=%v password=%v dbname=%v sslmode=%v", login, password, dbname, sslmode)

	var err error
	db.connection, err = sql.Open("postgres", db.ConnectionString)

	// log
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return db
}

// func (db *DataBase) CreateTable(item interface{}, tableName string) error {
// 	itemType := reflect.TypeOf(item)
// 	itemValue := reflect.ValueOf(item)

// 	if itemType.Kind() != reflect.Struct {
// 		return errIsntStruct
// 	}

// 	return nil
// }

func (db *DataBase) Add(item interface{}, tableName string) error {
	itemType := reflect.TypeOf(item)
	itemValue := reflect.ValueOf(item)

	if itemType.Kind() != reflect.Struct {
		return errIsntStruct
	}

	fields := ""

	for i := 0; i < itemValue.NumField(); i++ {
		if itemValue.Field(i).Kind() == reflect.String {
			fields += fmt.Sprintf("'%v',", itemValue.Field(i).Interface())
		} else if itemValue.Field(i).Kind() == reflect.Slice {
			array := fmt.Sprintf("'%v',", itemValue.Field(i).Interface())
			array = strings.Replace(array, "[", "{", -1)
			array = strings.Replace(array, " ", ",", -1)
			array = strings.Replace(array, "]", "}", -1)
			fields += array
		} else {
			fields += fmt.Sprintf("%v,", itemValue.Field(i).Interface())
		}
	}

	fields = fields[:len(fields)-1]

	query := fmt.Sprintf("Insert into \"%v\" values (%v);", tableName, fields)
	result, err := db.connection.Exec(query)
	log.Println("Kek")
	if err != nil {
		return err
	}

	affectedRows, _ := result.RowsAffected()
	log.Printf("Rows affected - %v\n", affectedRows)
	return nil
}

func (db *DataBase) Select(id int, tableName string, outItem interface{}) error {
	query := fmt.Sprintf("Select * From \"%v\" Where id = %v", tableName, id)
	rows, err := db.connection.Query(query)

	if err != nil {
		return err
	}

	for rows.Next() {
		val := reflect.ValueOf(outItem)

		if val.Kind() != reflect.Ptr {
			return errors.New("dest must be a pointer to a struct")
		}
		val = val.Elem() // get the underlying struct value
		if val.Kind() != reflect.Struct {
			return errors.New("dest must be a pointer to a struct")
		}

		numCols := val.NumField() // now this should work
		columns := make([]interface{}, numCols)
		for i := 0; i < numCols; i++ {
			columns[i] = val.Field(i).Addr().Interface()
		}
		rows.Scan(columns...)
	}

	return nil
}

func (db *DataBase) SelectAll(tableName string, outItem interface{}) ([]interface{}, error) {
	// out item is item that will be last and used like a typeOf out fields
	query := fmt.Sprintf("Select * From \"%v\"", tableName)
	rows, err := db.connection.Query(query)

	var result []interface{}

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		val := reflect.ValueOf(outItem)

		if val.Kind() != reflect.Ptr {
			return nil, errors.New("dest must be a pointer to a struct")
		}
		val = val.Elem() // get the underlying struct value
		if val.Kind() != reflect.Struct {
			return nil, errors.New("dest must be a pointer to a struct")
		}

		numCols := val.NumField() // now this should work
		columns := make([]interface{}, numCols)

		for i := 0; i < numCols; i++ {
			columns[i] = val.Field(i).Addr().Interface()
		}

		rows.Scan(columns...)
		result = append(result, outItem)
	}

	return result, nil
}

func (db *DataBase) Delete(id int, tableName string) error {
	query := fmt.Sprintf("Delete From \"%v\" Where id = %v", tableName, id)
	result, err := db.connection.Exec(query)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	log.Printf("Deleted %d row(s)\n", rowsAffected)
	return nil
}

func (db *DataBase) Update(id int, tableName string, item interface{}) error {
	itemType := reflect.TypeOf(item)
	itemValue := reflect.ValueOf(item)

	if itemType.Kind() != reflect.Struct {
		return errIsntStruct
	}

	fields := ""

	for i := 0; i < itemValue.NumField(); i++ {
		if itemValue.Field(i).Kind() == reflect.String {
			fields += fmt.Sprintf("%v = '%v',", strings.ToLower(itemType.Field(i).Name), itemValue.Field(i).Interface())
		} else {
			fields += fmt.Sprintf("%v = %v,", strings.ToLower(itemType.Field(i).Name), itemValue.Field(i).Interface())
		}
	}

	fields = fields[:len(fields)-1]
	query := fmt.Sprintf("Update \"%v\" Set %v Where id = %v", tableName, fields, id)
	result, err := db.connection.Exec(query)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	log.Printf("Updated %d row(s)\n", rowsAffected)
	return nil
}

func (db *DataBase) Close() {
	db.connection.Close()
}
