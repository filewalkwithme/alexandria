package orm

import (
	"database/sql"
	"reflect"
)

func (s Selecter) selectByID(id int) (interface{}, error) {
	sqlInstruction := s.handler.selectSQL + " where id = $1;"

	scanFields := s.handler.selectScanMap
	err := s.handler.db.QueryRow(sqlInstruction, id).Scan(scanFields...)

	if err == nil {
		return s.buildObject(scanFields), nil
	}

	return nil, err
}

func (s Selecter) selectAll() ([]interface{}, error) {
	sqlInstruction := s.handler.selectSQL

	rows, err := s.handler.db.Query(sqlInstruction)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.buildArrayOfObjects(rows)
}

func (s Selecter) selectWhere(where string, arguments ...interface{}) ([]interface{}, error) {
	sqlInstruction := s.handler.selectSQL + " where " + where + ";"

	rows, err := s.handler.db.Query(sqlInstruction, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.buildArrayOfObjects(rows)
}

func (s Selecter) buildArrayOfObjects(rows *sql.Rows) ([]interface{}, error) {
	var result []interface{}
	for rows.Next() {
		scanFields := s.handler.selectScanMap
		err := rows.Scan(scanFields...)
		if err == nil {
			result = append(result, s.buildObject(scanFields))
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (s Selecter) buildObject(fields []interface{}) interface{} {
	typeOfTable := reflect.TypeOf(s.handler.table)
	v := reflect.ValueOf(s.handler.table)
	vPtr := reflect.New(v.Type())
	for i := 0; i < typeOfTable.NumField(); i++ {
		if typeOfTable.Field(i).Type.Name() == "int" {
			vPtr.Elem().FieldByName(s.handler.selectFieldNamesMap[i]).SetInt(int64(*(fields[i].(*int))))
		}
		if typeOfTable.Field(i).Type.Name() == "string" {
			vPtr.Elem().FieldByName(s.handler.selectFieldNamesMap[i]).SetString(*(fields[i].(*string)))
		}
		if typeOfTable.Field(i).Type.Name() == "bool" {
			vPtr.Elem().FieldByName(s.handler.selectFieldNamesMap[i]).SetBool(*(fields[i].(*bool)))
		}
	}
	return vPtr.Elem().Interface()
}

func (handler *Handler) assembleSQLSelect() {
	typeOfTable := reflect.TypeOf(handler.table)
	tableName := typeOfTable.Name()

	var fieldNames = make([]string, typeOfTable.NumField())
	var scanFieds = make([]interface{}, typeOfTable.NumField())

	sqlInstruction := "select "
	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name
		fieldNames[i] = fieldName

		if typeOfTable.Field(i).Type.Name() == "int" {
			scanFieds[i] = new(int)
		}
		if typeOfTable.Field(i).Type.Name() == "string" {
			scanFieds[i] = new(string)
		}
		if typeOfTable.Field(i).Type.Name() == "bool" {
			scanFieds[i] = new(bool)
		}

		sqlInstruction = sqlInstruction + fieldName + ", "
	}

	sqlInstruction = sqlInstruction[:len(sqlInstruction)-2]
	sqlInstruction = sqlInstruction + " from " + tableName

	handler.selectSQL = sqlInstruction
	handler.selectScanMap = scanFieds
	handler.selectFieldNamesMap = fieldNames
}
