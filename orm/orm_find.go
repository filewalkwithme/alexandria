package orm

import (
	"fmt"
	"reflect"
	"strconv"
)

func (f Finder) findByID(table interface{}, id int) interface{} {
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
	sqlInstruction = sqlInstruction + " from " + tableName
	sqlInstruction = sqlInstruction + " where id = " + strconv.Itoa(id) + ";"

	fmt.Printf("sqlInstruction: %v\n", sqlInstruction)

	err := f.db.QueryRow(sqlInstruction).Scan(destFieds...)
	v := reflect.ValueOf(table)

	if err == nil {
		vPtr := reflect.New(v.Type())
		for i := 0; i < typeOfTable.NumField(); i++ {
			if typeOfTable.Field(i).Type.Name() == "int" {
				vPtr.Elem().FieldByName(fields[i]).SetInt(int64(*(destFieds[i].(*int))))
			}
			if typeOfTable.Field(i).Type.Name() == "string" {
				vPtr.Elem().FieldByName(fields[i]).SetString(*(destFieds[i].(*string)))
			}
		}
		return vPtr.Elem().Interface()
	}

	fmt.Printf("err: %v\n", err)
	return nil
}

func (f Finder) findAll(table interface{}) []interface{} {
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

	//fmt.Printf("sqlInstruction: %v\n", sqlInstruction)

	var res []interface{}

	rows, err := f.db.Query(sqlInstruction)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(destFieds...)
		v := reflect.ValueOf(table)

		if err == nil {
			vPtr := reflect.New(v.Type())
			for i := 0; i < typeOfTable.NumField(); i++ {
				if typeOfTable.Field(i).Type.Name() == "int" {
					vPtr.Elem().FieldByName(fields[i]).SetInt(int64(*(destFieds[i].(*int))))
				}
				if typeOfTable.Field(i).Type.Name() == "string" {
					vPtr.Elem().FieldByName(fields[i]).SetString(*(destFieds[i].(*string)))
				}
			}
			res = append(res, vPtr.Elem().Interface())
		} else {
			fmt.Printf("err: %v\n", err)
		}
	}

	return res
}

func (f Finder) findWhere(table interface{}, where string) []interface{} {
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
	sqlInstruction = sqlInstruction + " from " + tableName
	sqlInstruction = sqlInstruction + " where " + where + ";"

	fmt.Printf("sqlInstruction: %v\n", sqlInstruction)

	var res []interface{}

	rows, err := f.db.Query(sqlInstruction)
	defer rows.Close()

	if err == nil {
		for rows.Next() {
			err := rows.Scan(destFieds...)
			v := reflect.ValueOf(table)

			if err == nil {
				vPtr := reflect.New(v.Type())
				for i := 0; i < typeOfTable.NumField(); i++ {
					if typeOfTable.Field(i).Type.Name() == "int" {
						vPtr.Elem().FieldByName(fields[i]).SetInt(int64(*(destFieds[i].(*int))))
					}
					if typeOfTable.Field(i).Type.Name() == "string" {
						vPtr.Elem().FieldByName(fields[i]).SetString(*(destFieds[i].(*string)))
					}
				}
				res = append(res, vPtr.Elem().Interface())
			} else {
				fmt.Printf("err: %v\n", err)
			}
		}
	}

	return res
}
