package main

import (
	"fmt"
	"reflect"
	"strconv"

	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var db *sql.DB

//Book is on the table
type Book struct {
	ID    int
	Name  string
	Pages int
}

func initDB() {
	tmpDB, err := sql.Open("postgres", "user=docker password=docker dbname=docker sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	err = tmpDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db = tmpDB
}

func save(table interface{}) {

	typeOfTable := reflect.TypeOf(table)
	valueOfTable := reflect.ValueOf(table)

	tableName := typeOfTable.Name()

	fields := make(map[string]string)
	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name

		if fieldName == "ID" && valueOfTable.Field(i).Int() <= 0 {
			continue
		}

		if typeOfTable.Field(i).Type.Name() == "int" {
			fields[fieldName] = strconv.Itoa(int(valueOfTable.Field(i).Int()))
		}
		if typeOfTable.Field(i).Type.Name() == "string" {
			fields[fieldName] = "'" + valueOfTable.Field(i).String() + "'"
		}
	}

	sqlInstruction := "insert into " + tableName + "("
	sqlFields := ""
	sqlValues := ""
	for fieldName, value := range fields {
		sqlFields = sqlFields + fieldName + ", "
		sqlValues = sqlValues + value + ", "
	}
	sqlFields = sqlFields[:len(sqlFields)-2]
	sqlValues = sqlValues[:len(sqlValues)-2]
	sqlInstruction = sqlInstruction + sqlFields + ") values (" + sqlValues + ");"

	result, err := db.Exec(sqlInstruction)
	fmt.Printf("sqlInstruction: %v\n", sqlInstruction)
	fmt.Printf("result: %v\n", result)
	fmt.Printf("err: %v\n", err)
}

func createTable(table interface{}) {

	typeOfTable := reflect.TypeOf(table)

	tableName := typeOfTable.Name()

	fieldsList := ""
	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name
		fieldType := ""
		if typeOfTable.Field(i).Type.Name() == "int" {
			fieldType = "integer"
		}
		if typeOfTable.Field(i).Type.Name() == "string" {
			fieldType = "character varying"
		}

		if fieldName == "ID" {
			fieldsList = fieldsList + fieldName + " " + fieldType + " NOT NULL, "
		} else {
			fieldsList = fieldsList + fieldName + " " + fieldType + ", "
		}
	}

	primaryKey := "constraint " + tableName + "_pkey primary key (id)"
	sqlInstruction := "create table " + tableName + " (" + fieldsList + " " + primaryKey + ");\n"

	result, err := db.Exec(sqlInstruction)
	fmt.Printf("result: %v\n", result)
	fmt.Printf("err: %v\n", err)
}

func main() {
	initDB()
	createTable(Book{})
	save(Book{Name: "moby dick", Pages: 199})
}
