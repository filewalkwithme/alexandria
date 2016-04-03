package main

import (
	"fmt"
	"reflect"

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
}
