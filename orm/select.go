package orm

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
)

func (s Selecter) selectByID(id int) (*interface{}, error) {
	sqlInstruction := s.handler.selectSQL + " where id = $1;"

	scanFields := s.handler.selectScanMap
	err := s.handler.db.QueryRow(sqlInstruction, id).Scan(scanFields...)

	if err == nil {
		return s.buildObject(scanFields)
	}

	return nil, err
}

func (s Selecter) selectAll() ([]*interface{}, error) {
	sqlInstruction := s.handler.selectSQL

	rows, err := s.handler.db.Query(sqlInstruction)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.buildArrayOfObjects(rows)
}

func (s Selecter) selectWhere(where string, arguments ...interface{}) ([]*interface{}, error) {
	sqlInstruction := s.handler.selectSQL + " where " + where + ";"

	rows, err := s.handler.db.Query(sqlInstruction, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return s.buildArrayOfObjects(rows)
}

func (s Selecter) buildArrayOfObjects(rows *sql.Rows) ([]*interface{}, error) {
	var result []*interface{}
	for rows.Next() {
		scanFields := s.handler.selectScanMap
		err := rows.Scan(scanFields...)
		if err == nil {
			obj, err := s.buildObject(scanFields)
			if err != nil {
				return nil, err
			}
			result = append(result, obj)
		} else {
			return nil, err
		}
	}
	return result, nil
}

func (s Selecter) buildObject(fields []interface{}) (*interface{}, error) {
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
				//fmt.Printf("handler.childHandlers: %v\n", s.handler.childHandlers["Chapters"])
				//fmt.Printf("handler.childHandlers: %v\n", *fields[0])
				childValues, err := s.handler.childHandlers[typeOfTable.Field(i).Name].Select().Where(fmt.Sprintf("%v = %v", parentFieldName, parentID))
				if err != nil {
					return nil, err
				}

				if len(childValues) > 1 {
					return nil, fmt.Errorf("%v:%v found multiple childs. Expected zero or one", typeOfTable, typeOfTable.Field(i).Name)
				}

				//Super awesome thanks to smith.wi...@gmail.com
				//https://groups.google.com/forum/#!topic/golang-nuts/KB3_Yj3Ny4c
				//https://play.golang.org/p/hYnsAijyCE
				if len(childValues) == 1 {
					childPointer := reflect.New(reflect.TypeOf(*childValues[0]))
					childPointer.Elem().Set(reflect.ValueOf(*childValues[0]))

					vPtr.Elem().Field(i).Set(childPointer)
				}
			}

			if typeOfTable.Field(i).Type.Kind().String() == "slice" {
				childValues, err := s.handler.childHandlers[typeOfTable.Field(i).Name].Select().Where(fmt.Sprintf("%v = %v", parentFieldName, parentID))
				if err != nil {
					return nil, err
				}

				sliceType := reflect.SliceOf(reflect.New(reflect.TypeOf(s.handler.childHandlers[typeOfTable.Field(i).Name].table)).Type())
				sliceLocal := reflect.MakeSlice(sliceType, 0, 0)

				for j := 0; j < len(childValues); j++ {
					childPointer := reflect.New(reflect.TypeOf(s.handler.childHandlers[typeOfTable.Field(i).Name].table))
					childPointer.Elem().Set(reflect.ValueOf(*childValues[j]))
					sliceLocal = reflect.Append(sliceLocal, childPointer)
				}

				vPtr.Elem().Field(i).Set(sliceLocal)
			}
			continue
		}
	}

	ptrInterface := vPtr.Elem().Interface()
	return &ptrInterface, nil
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
