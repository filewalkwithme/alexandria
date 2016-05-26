package orm

import (
	"fmt"
	"reflect"
	"strconv"
)

func (deleter Deleter) delete(table interface{}) int {
	typeOfTable := reflect.TypeOf(table)
	valueOfTable := reflect.ValueOf(table)

	tableName := typeOfTable.Name()
	id := ""

	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name
		if fieldName == "ID" {
			id = strconv.Itoa(int(valueOfTable.Field(i).Int()))
			break
		}
	}

	sqlInstruction := "delete from " + tableName
	sqlInstruction = sqlInstruction + " where id = " + id + ";"

	result, err := deleter.db.Exec(sqlInstruction)

	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return -1
	}
	return int(n)
}

func (deleter Deleter) deleteWhere(table interface{}, where string) int {
	typeOfTable := reflect.TypeOf(table)

	tableName := typeOfTable.Name()

	sqlInstruction := "delete from " + tableName + " where " + where + ";"

	result, err := deleter.db.Exec(sqlInstruction)

	//TODO: make save return or populate the ID field
	fmt.Printf("sqlInstruction: %v\n", sqlInstruction)
	fmt.Printf("result: %v\n", result)
	fmt.Printf("err: %v\n", err)

	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return -1
	}
	return int(n)
}
