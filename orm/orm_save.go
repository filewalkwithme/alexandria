package orm

import (
	"reflect"
	"strconv"
)

func (orm Orm) save(table interface{}) (interface{}, error) {

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

	_, err := orm.db.Exec(sqlInstruction)

	//TODO: make save return the new object
	return nil, err
}
