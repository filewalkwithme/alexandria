package orm

import (
	"fmt"
	"reflect"
	"strconv"
)

func (orm Orm) delete(table interface{}) int {
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

	result, err := orm.db.Exec(sqlInstruction)

	n, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return -1
	}
	return int(n)
}
