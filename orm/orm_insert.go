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
func (handler *Handler) assembleSQLInsertStatement() error {
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
	sqlInstruction = sqlInstruction + sqlFields + ") values (" + sqlValues + ");"

	handler.sqlInsert = sqlInstruction
	handler.mapInsert = fieldMap

	return nil
}

// assembleSQLInsertStatement traverse the the object
// returns a SQL insert instruction and a string array containing the exact
// parameters order
func (handler Handler) assembleSQLUpdateStatement(object interface{}) (string, []saveField, error) {
	typeOfTable := reflect.TypeOf(object)
	tableName := typeOfTable.Name()
	if tableName != handler.tableName {
		return "", nil, fmt.Errorf("Object table name (%v) is diferent from handler table name (%v)", tableName, handler.tableName)
	}

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
	sqlInstruction = sqlInstruction + " where id = ?;"

	return sqlInstruction, fieldMap, nil
}

func (handler Handler) insert(object interface{}) error {
	typeOfTable := reflect.TypeOf(object)
	valueOfTable := reflect.ValueOf(object)
	tableName := typeOfTable.Name()

	if tableName != handler.tableName {
		return fmt.Errorf("Object table name (%v) is diferent from handler table name (%v)", tableName, handler.tableName)
	}

	var args []interface{}
	for _, field := range handler.mapInsert {
		if field.fieldType == "int" {
			args = append(args, int(valueOfTable.FieldByName(field.name).Int()))
		}
		if field.fieldType == "string" {
			args = append(args, string(valueOfTable.FieldByName(field.name).String()))
		}
	}
	fmt.Printf("%v\n", handler.sqlInsert)
	_, err := handler.db.Exec(handler.sqlInsert, args...)

	return err
}
