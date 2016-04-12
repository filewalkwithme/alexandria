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

func findAll(table interface{}) []interface{} {
	typeOfTable := reflect.TypeOf(table)

	tableName := typeOfTable.Name()

	var destFieds = make([]interface{}, typeOfTable.NumField())
	var fields = make([]string, typeOfTable.NumField())
	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name

		if typeOfTable.Field(i).Type.Name() == "int" {
			destFieds[i] = new(int)
		}
		if typeOfTable.Field(i).Type.Name() == "string" {
			destFieds[i] = new(string)
		}

		fields[i] = fieldName
	}

	sqlInstruction := "select "
	for _, fieldName := range fields {
		sqlInstruction = sqlInstruction + fieldName + ", "
	}
	sqlInstruction = sqlInstruction[:len(sqlInstruction)-2]
	sqlInstruction = sqlInstruction + " from " + tableName + ";"

	fmt.Printf("sqlInstruction: %v\n", sqlInstruction)

	var res []interface{}

	rows, err := db.Query(sqlInstruction)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(destFieds...)
		var book Book

		if err == nil {
			p := reflect.ValueOf(&book)
			for i := 0; i < typeOfTable.NumField(); i++ {
				if typeOfTable.Field(i).Type.Name() == "int" {
					p.Elem().FieldByName(fields[i]).SetInt(int64(*(destFieds[i].(*int))))
				}
				if typeOfTable.Field(i).Type.Name() == "string" {
					p.Elem().FieldByName(fields[i]).SetString(*(destFieds[i].(*string)))
				}
			}
			res = append(res, book)
		} else {
			fmt.Printf("err: %v\n", err)
		}
	}

	return res
}

func find(table interface{}) interface{} {
	typeOfTable := reflect.TypeOf(table)
	valueOfTable := reflect.ValueOf(table)

	tableName := typeOfTable.Name()
	id := ""

	var destFieds = make([]interface{}, typeOfTable.NumField())
	var fields = make([]string, typeOfTable.NumField())
	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name
		if fieldName == "ID" {
			id = strconv.Itoa(int(valueOfTable.Field(i).Int()))
		}

		if typeOfTable.Field(i).Type.Name() == "int" {
			destFieds[i] = new(int)
		}
		if typeOfTable.Field(i).Type.Name() == "string" {
			destFieds[i] = new(string)
		}

		fields[i] = fieldName
	}

	sqlInstruction := "select "
	for _, fieldName := range fields {
		sqlInstruction = sqlInstruction + fieldName + ", "
	}
	sqlInstruction = sqlInstruction[:len(sqlInstruction)-2]
	sqlInstruction = sqlInstruction + " from " + tableName
	sqlInstruction = sqlInstruction + " where id = " + id + ";"

	fmt.Printf("sqlInstruction: %v\n", sqlInstruction)

	err := db.QueryRow(sqlInstruction).Scan(destFieds...)
	res := Book{}
	if err == nil {
		p := reflect.ValueOf(&res)
		for i := 0; i < typeOfTable.NumField(); i++ {
			if typeOfTable.Field(i).Type.Name() == "int" {
				p.Elem().FieldByName(fields[i]).SetInt(int64(*(destFieds[i].(*int))))
			}
			if typeOfTable.Field(i).Type.Name() == "string" {
				p.Elem().FieldByName(fields[i]).SetString(*(destFieds[i].(*string)))
			}
		}

	} else {
		fmt.Printf("err: %v\n", err)
	}

	return res
}

func save(table interface{}) {

	typeOfTable := reflect.TypeOf(table)
	valueOfTable := reflect.ValueOf(table)

	tableName := typeOfTable.Name()

	fields := make(map[string]string)
	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name

		if typeOfTable.Field(i).Type.Name() == "int" {
			fields[fieldName] = strconv.Itoa(int(valueOfTable.Field(i).Int()))
		}
		if typeOfTable.Field(i).Type.Name() == "string" {
			fields[fieldName] = "'" + valueOfTable.Field(i).String() + "'"
		}
	}
	sqlInstruction := "select 1"
	id, _ := strconv.Atoi(fields["ID"])
	if id > 0 {
		//update
		sqlInstruction = "update " + tableName + " set "
		for fieldName, value := range fields {
			if fieldName == "ID" {
				continue
			}
			sqlInstruction = sqlInstruction + fieldName + " = " + value + ", "
		}
		sqlInstruction = sqlInstruction[:len(sqlInstruction)-2]
		sqlInstruction = sqlInstruction + " where id = " + fields["ID"] + ";"
	} else {
		//insert
		sqlInstruction = "insert into " + tableName + "("
		sqlFields := ""
		sqlValues := ""
		for fieldName, value := range fields {
			if fieldName == "ID" {
				continue
			}
			sqlFields = sqlFields + fieldName + ", "
			sqlValues = sqlValues + value + ", "
		}
		sqlFields = sqlFields[:len(sqlFields)-2]
		sqlValues = sqlValues[:len(sqlValues)-2]
		sqlInstruction = sqlInstruction + sqlFields + ") values (" + sqlValues + ");"
	}

	result, err := db.Exec(sqlInstruction)

	//TODO: make save return or populate the ID field
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
			fieldsList = fieldsList + fieldName + " serial NOT NULL, "
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
	//createTable(Book{})
	//save(Book{Name: "moby dick", Pages: 199})
	//save(Book{ID: 1, Name: "moby dick2", Pages: 299})
	//book := find(Book{ID: 1})
	books := findAll(Book{})
	fmt.Printf("book: %v\n", books)
}
