package orm

import (
	"fmt"
	"reflect"
	"strconv"
)

type saveField struct {
	name      string
	fieldType string
}

// assembleSQLInsertStatement traverse the the object
// returns a SQL insert instruction and a string array containing the exact
// parameters order
func (handler *Handler) assembleSQLInsertStatement() {
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

	handler.sqlInsert = sqlInstruction
	handler.mapInsert = fieldMap
}

func (handler Handler) insert(objectPtr interface{}) error {
	object := reflect.ValueOf(objectPtr).Elem()
	tableName := reflect.TypeOf(objectPtr).Elem().Name()
	if tableName != handler.tableName {
		return fmt.Errorf("Object table name (%v) is diferent from handler table name (%v)", tableName, handler.tableName)
	}

	//build the arguments array
	var args []interface{}
	for _, field := range handler.mapInsert {
		if field.fieldType == "int" {
			args = append(args, int(object.FieldByName(field.name).Int()))
		}
		if field.fieldType == "string" {
			args = append(args, string(object.FieldByName(field.name).String()))
		}
	}

	//run INSERT and grab the last insert id
	var id int
	err := handler.db.QueryRow(handler.sqlInsert, args...).Scan(&id)
	if err != nil {
		return err
	}
	object.FieldByName("ID").SetInt(int64(id))

	return nil
}
