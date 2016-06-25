package orm

import (
	"fmt"
	"reflect"
)

// createTable() must execute the sql CREATE TABLE instruction
func (handler Handler) createTable() (err error) {
	sqlInstruction := ""

	if handler.sqlCreateTable != "" {
		sqlInstruction = handler.sqlCreateTable
	} else {
		sqlInstruction, err = handler.assembleSQLCreateTable()
		if err != nil {
			return err
		}
	}

	_, err = handler.db.Exec(sqlInstruction)
	return err
}

// createTable() must execute the sql DROP TABLE instruction
func (handler Handler) dropTable() error {
	sqlInstruction := ""

	if handler.sqlDropTable != "" {
		sqlInstruction = handler.sqlDropTable
	} else {
		sqlInstruction = handler.assembleSQLDropTable()
	}

	_, err := handler.db.Exec(sqlInstruction)

	return err
}

// createTableSQL() must traverse the table structure, colect its fields and
// assemble the sql CREATE TABLE instruction
func (handler Handler) assembleSQLCreateTable() (string, error) {
	var idExists = false

	typeOfTable := reflect.TypeOf(handler.table)
	tableName := typeOfTable.Name()

	fieldsList := ""
	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name
		fieldType := ""
		switch typeOfTable.Field(i).Type.Name() {
		case "int":
			fieldType = "integer"
		case "float64":
			fieldType = "real"
		case "string":
			fieldType = "character varying"
		case "bool":
			fieldType = "boolean"
		case "Time":
			fieldType = "timestamp without time zone"
		default:
			continue
		}

		if fieldName == "ID" {
			idExists = true
			fieldsList = fieldsList + fieldName + " serial NOT NULL, "
		} else {
			fieldsList = fieldsList + fieldName + " " + fieldType + ", "
		}
	}

	if idExists == false {
		return "", fmt.Errorf("ID field not found on struct %v", tableName)
	}

	primaryKey := "constraint " + tableName + "_pkey primary key (id)"
	sqlInstruction := "create table " + tableName + " (" + fieldsList + primaryKey + ");"

	handler.sqlCreateTable = sqlInstruction

	return sqlInstruction, nil
}

// assembleSQLDropTable() assemble the SQL Drop Table instruction
func (handler Handler) assembleSQLDropTable() string {
	tableName := reflect.TypeOf(handler.table).Name()
	sqlInstruction := "drop table " + tableName + ";"

	handler.sqlDropTable = sqlInstruction

	return sqlInstruction
}
