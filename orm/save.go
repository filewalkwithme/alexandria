package orm

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type saveField struct {
	name      string
	fieldType string
}

// save takes a pointer to an object and performs an INSERT if the ID is zero.
// If the ID is greather than zero, then an UPDATE will be performed instead.
// Otherwise, if the ID is negative, an error will be returned
func (handler Handler) save(objectPtr interface{}) error {
	object := reflect.ValueOf(objectPtr).Elem()
	tableName := reflect.TypeOf(objectPtr).Elem().Name()
	if tableName != handler.tableName {
		return fmt.Errorf("Object table name (%v) is diferent from handler table name (%v)", tableName, handler.tableName)
	}

	id := object.FieldByName("ID").Int()

	if id == 0 {
		//insert
		return handler.insert(object)
	} else if id > 0 {
		//update
		return handler.update(object)
	} else {
		//error
		return fmt.Errorf("Negative ID not allowed: %v", object)
	}
}

// insert takes an object, grab its fields and performs an INSERT operation
func (handler Handler) insert(object reflect.Value) error {
	values := handler.assembleValuesArray(handler.insertMap, object)
	var id int

	//run INSERT and grab the last insert id
	err := handler.db.QueryRow(handler.insertSQL, values...).Scan(&id)
	if err != nil {
		return err
	}
	object.FieldByName("ID").SetInt(int64(id))

	return nil
}

// insert takes an object, grab its fields and performs an UPDATE operation
func (handler Handler) update(object reflect.Value) error {
	values := handler.assembleValuesArray(handler.updateMap, object)

	//run Update
	_, err := handler.db.Exec(handler.updateSQL, values...)
	if err != nil {
		return err
	}

	return nil
}

// assembleArguments takes an object and return an array of its values populated
// in the exact order required by the given arguments map. The resulting array
// of values is intended to be consumed by insert() or update()
func (handler Handler) assembleValuesArray(argurmentsMap []saveField, object reflect.Value) []interface{} {
	//build the values array
	var values []interface{}
	for _, argument := range argurmentsMap {
		if argument.fieldType == "int" {
			values = append(values, int(object.FieldByName(argument.name).Int()))
		}
		if argument.fieldType == "float64" {
			values = append(values, float64(object.FieldByName(argument.name).Float()))
		}
		if argument.fieldType == "string" {
			values = append(values, string(object.FieldByName(argument.name).String()))
		}
		if argument.fieldType == "bool" {
			values = append(values, bool(object.FieldByName(argument.name).Bool()))
		}
		if argument.fieldType == "Time" {
			values = append(values, time.Time(object.FieldByName(argument.name).Interface().(time.Time)))
		}
	}
	return values
}

// assembleSQLInsert creates a SQL Insert string in the current handler.
// This string will be used by the save() method
func (handler *Handler) assembleSQLInsert() {
	typeOfTable := reflect.TypeOf(handler.table)
	tableName := typeOfTable.Name()

	sqlInstruction := "insert into " + tableName + "("

	sqlFields := ""
	sqlValues := ""
	j := 1
	var fieldMap []saveField

	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name

		if fieldName == "ID" {
			continue
		}

		fieldMap = append(fieldMap, saveField{name: typeOfTable.Field(i).Name, fieldType: typeOfTable.Field(i).Type.Name()})
		sqlFields = sqlFields + fieldName + ", "
		sqlValues = sqlValues + "$" + strconv.Itoa(j) + ", "
		j = j + 1
	}

	sqlFields = sqlFields[:len(sqlFields)-2]
	sqlValues = sqlValues[:len(sqlValues)-2]
	sqlInstruction = sqlInstruction + sqlFields + ") values (" + sqlValues + ") RETURNING id;"

	handler.insertSQL = sqlInstruction
	handler.insertMap = fieldMap
}

// assembleSQLUpdate creates a SQL Update string in the current handler.
// This string will be used by the save() method
func (handler *Handler) assembleSQLUpdate() {
	typeOfTable := reflect.TypeOf(handler.table)
	tableName := typeOfTable.Name()

	j := 1
	sqlInstruction := "update " + tableName + " set "
	var fieldMap []saveField
	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name

		if fieldName == "ID" {
			continue
		}

		fieldMap = append(fieldMap, saveField{name: typeOfTable.Field(i).Name, fieldType: typeOfTable.Field(i).Type.Name()})
		sqlInstruction = sqlInstruction + fieldName + " = $" + strconv.Itoa(j) + ", "
		j = j + 1
	}
	fieldMap = append(fieldMap, saveField{name: "ID", fieldType: "int"})

	sqlInstruction = sqlInstruction[:len(sqlInstruction)-2]
	sqlInstruction = sqlInstruction + " where id = $" + strconv.Itoa(j)

	handler.updateSQL = sqlInstruction
	handler.updateMap = fieldMap
}
