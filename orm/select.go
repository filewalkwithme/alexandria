package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

func (s Selecter) selectByID(id int) error {
	object := reflect.ValueOf(s.handler.object).Elem()
	sqlInstruction := s.handler.selectSQL + " where id = $1;"
	scanFields := s.handler.selectScanMap
	err := s.handler.db.QueryRow(sqlInstruction, id).Scan(scanFields...)

	if err == nil {
		obj, err := s.buildObject(scanFields)
		if err == nil {
			object.Set(reflect.ValueOf(obj).Elem())
		}
		return err
	}

	return err
}

func (s Selecter) selectAll() (interface{}, error) {
	sqlInstruction := s.handler.selectSQL

	rows, err := s.handler.db.Query(sqlInstruction)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.buildArrayOfObjects(rows)
}

func (s Selecter) selectWhere(where string, arguments ...interface{}) error {
	object := reflect.ValueOf(s.handler.object).Elem()
	sqlInstruction := s.handler.selectSQL + " where " + where + ";"

	rows, err := s.handler.db.Query(sqlInstruction, arguments...)
	if err != nil {
		return err
	}
	defer rows.Close()

	obj, err := s.buildArrayOfObjects(rows)
	if err == nil {
		object.Set(reflect.ValueOf(obj))
		return err
	}
	return err
}

func (s Selecter) buildArrayOfObjects(rows *sql.Rows) (interface{}, error) {
	v := reflect.MakeSlice(reflect.ValueOf(s.handler.object).Elem().Type(), 0, 0)

	for rows.Next() {
		scanFields := s.handler.selectScanMap
		err := rows.Scan(scanFields...)
		if err == nil {
			obj, err := s.buildObject(scanFields)
			if err != nil {
				return nil, err
			}

			v = reflect.Append(v, reflect.ValueOf(obj))
		} else {
			return nil, err
		}
	}
	return v.Interface(), nil
}

func (s Selecter) buildObject(fields []interface{}) (interface{}, error) {
	typeOfTable := reflect.TypeOf(s.handler.table)
	parentFieldName := typeOfTable.Name() + "ID"
	parentID := int64(-1)
	for i := 0; i < typeOfTable.NumField(); i++ {
		if typeOfTable.Field(i).Name == "ID" && typeOfTable.Field(i).Type.Name() == "int" {
			parentID = int64(*(fields[i].(*int)))
		}
	}

	v := reflect.ValueOf(s.handler.table)
	vPtr := reflect.New(v.Type())
	for i := 0; i < typeOfTable.NumField(); i++ {
		switch typeOfTable.Field(i).Type.Name() {
		case "int":
			vPtr.Elem().FieldByName(s.handler.selectFieldNamesMap[i]).SetInt(int64(*(fields[i].(*int))))
		case "float64":
			vPtr.Elem().FieldByName(s.handler.selectFieldNamesMap[i]).SetFloat(float64(*(fields[i].(*float64))))
		case "string":
			vPtr.Elem().FieldByName(s.handler.selectFieldNamesMap[i]).SetString(*(fields[i].(*string)))
		case "bool":
			vPtr.Elem().FieldByName(s.handler.selectFieldNamesMap[i]).SetBool(*(fields[i].(*bool)))
		case "Time":
			vPtr.Elem().FieldByName(s.handler.selectFieldNamesMap[i]).Set(reflect.ValueOf(*(fields[i].(*time.Time))))
		default:
			if typeOfTable.Field(i).Type.Kind().String() == "ptr" {
				tableObject := s.handler.childHandlers[typeOfTable.Field(i).Name].table
				pointerToTableObject := reflect.New(reflect.ValueOf(tableObject).Type())
				sliceOfObjects := reflect.MakeSlice(reflect.SliceOf(pointerToTableObject.Type()), 0, 0).Interface()
				pointerToSliceOfObjects := reflect.New(reflect.TypeOf(sliceOfObjects)).Interface()

				err := s.handler.childHandlers[typeOfTable.Field(i).Name].Select(pointerToSliceOfObjects).Where(fmt.Sprintf("%v = %v", parentFieldName, parentID))
				if err != nil {
					return nil, err
				}

				sliceSize := reflect.ValueOf(pointerToSliceOfObjects).Elem().Len()

				if sliceSize > 1 {
					return nil, fmt.Errorf("%v:%v found multiple childs. Expected zero or one", typeOfTable, typeOfTable.Field(i).Name)
				}

				//Super awesome thanks to smith.wi...@gmail.com
				//https://groups.google.com/forum/#!topic/golang-nuts/KB3_Yj3Ny4c
				//https://play.golang.org/p/hYnsAijyCE
				if sliceSize == 1 {
					vPtr.Elem().Field(i).Set(reflect.ValueOf(pointerToSliceOfObjects).Elem().Index(0))
				}
			}

			if typeOfTable.Field(i).Type.Kind().String() == "slice" {
				tableObject := s.handler.childHandlers[typeOfTable.Field(i).Name].table
				pointerToTableObject := reflect.New(reflect.ValueOf(tableObject).Type())
				sliceOfObjects := reflect.MakeSlice(reflect.SliceOf(pointerToTableObject.Type()), 0, 0).Interface()
				pointerToSliceOfObjects := reflect.New(reflect.TypeOf(sliceOfObjects)).Interface()

				err := s.handler.childHandlers[typeOfTable.Field(i).Name].Select(pointerToSliceOfObjects).Where(fmt.Sprintf("%v = %v", parentFieldName, parentID))
				if err != nil {
					return nil, err
				}

				vPtr.Elem().Field(i).Set(reflect.ValueOf(pointerToSliceOfObjects).Elem())
			}
			continue
		}
	}

	//ptrInterface := vPtr.Elem().Interface()
	//return &ptrInterface, nil
	return vPtr.Interface(), nil
}

func (handler *Handler) assembleSQLSelect() {
	typeOfTable := reflect.TypeOf(handler.table)
	tableName := typeOfTable.Name()

	var fieldNames = make([]string, 0)
	var scanFields = make([]interface{}, 0)

	sqlInstruction := "select "
	for i := 0; i < typeOfTable.NumField(); i++ {
		switch typeOfTable.Field(i).Type.Name() {
		case "int":
			scanFields = append(scanFields, new(int))
		case "float64":
			scanFields = append(scanFields, new(float64))
		case "string":
			scanFields = append(scanFields, new(string))
		case "bool":
			scanFields = append(scanFields, new(bool))
		case "Time":
			scanFields = append(scanFields, new(time.Time))
		default:
			continue
		}

		fieldName := typeOfTable.Field(i).Name
		fieldNames = append(fieldNames, fieldName)
		sqlInstruction = sqlInstruction + fieldName + ", "
	}

	sqlInstruction = sqlInstruction[:len(sqlInstruction)-2]
	sqlInstruction = sqlInstruction + " from " + tableName

	handler.selectSQL = sqlInstruction
	handler.selectScanMap = scanFields
	handler.selectFieldNamesMap = fieldNames
}
