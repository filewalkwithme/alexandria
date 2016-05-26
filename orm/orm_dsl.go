package orm

import (
	"fmt"
	"reflect"
)

type dbField struct {
	name   string
	dbType string
}

type dbTable struct {
	name   string
	fields []dbField
}

var tables = make(map[string]dbTable)

func loadToMemory(table interface{}) {
	tmpTable := dbTable{}

	typeOfTable := reflect.TypeOf(table)
	tableName := typeOfTable.Name()
	if _, exists := tables[tableName]; exists == false {
		tmpTable.name = tableName

		for i := 0; i < typeOfTable.NumField(); i++ {
			fieldName := typeOfTable.Field(i).Name
			fieldType := ""
			if typeOfTable.Field(i).Type.Name() == "int" {
				fieldType = "integer"
			}
			if typeOfTable.Field(i).Type.Name() == "string" {
				fieldType = "character varying"
			}

			tmpField := dbField{}
			tmpField.name = fieldName
			tmpField.dbType = fieldType
			tmpTable.fields = append(tmpTable.fields, tmpField)
		}

		tables[tableName] = tmpTable
	}
}

func (handler Handler) createTable() error {
	loadToMemory(handler.table)

	tableName := reflect.TypeOf(handler.table).Name()
	dbTable := tables[tableName]

	fieldsList := ""
	for _, field := range dbTable.fields {
		fieldName := field.name
		fieldType := field.dbType

		if fieldName == "ID" {
			fieldsList = fieldsList + fieldName + " serial NOT NULL, "
		} else {
			fieldsList = fieldsList + fieldName + " " + fieldType + ", "
		}
	}

	primaryKey := "constraint " + tableName + "_pkey primary key (id)"
	sqlInstruction := "create table " + tableName + " (" + fieldsList + " " + primaryKey + ");\n"

	result, err := handler.db.Exec(sqlInstruction)

	fmt.Printf("result: %v\n", result)
	fmt.Printf("err: %v\n", err)
	return err
}
