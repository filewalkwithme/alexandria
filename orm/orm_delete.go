package orm

import (
	"fmt"
	"reflect"
	"strconv"
)

func (deleter Deleter) deleteByID(table interface{}, id int) int {
	typeOfTable := reflect.TypeOf(table)
	tableName := typeOfTable.Name()

	sqlInstruction := "delete from " + tableName
	sqlInstruction = sqlInstruction + " where id = " + strconv.Itoa(id) + ";"

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

func (deleter Deleter) deleteAll(table interface{}) int {
	typeOfTable := reflect.TypeOf(table)

	tableName := typeOfTable.Name()

	sqlInstruction := "delete from " + tableName + ";"

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
