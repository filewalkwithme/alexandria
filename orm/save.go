package orm

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type saveField struct {
	name      string
	fieldType string
}

// save takes a pointer to an object and performs an INSERT if the ID is zero.
// If the ID is greather than zero, then an UPDATE will be performed instead.
// Otherwise, if the ID is negative, an error will be returned
func (handler Handler) save(objectPtr interface{}) error {
	object := reflect.ValueOf(objectPtr).Elem()
	tableName := reflect.TypeOf(objectPtr).Elem().Name()
	if tableName != handler.tableName {
		return fmt.Errorf("Object table name (%v) is diferent from handler table name (%v)", tableName, handler.tableName)
	}

	id := object.FieldByName("ID").Int()

	if id >= 0 {
		var err error
		if id == 0 {
			//insert
			err = handler.insert(object)
		} else {
			//update
			err = handler.update(object)
		}

		err = handler.saveChilds(object)
		if err != nil {
			return err
		}
	} else {
		//error
		return fmt.Errorf("Negative ID not allowed: %v", object)
	}

	return nil
}

// insert takes an object, grab its fields and performs an INSERT operation
func (handler Handler) insert(object reflect.Value) error {
	values := handler.assembleValuesArray(handler.insertMap, object)
	var id int

	//run INSERT and grab the last insert id
	err := handler.db.QueryRow(handler.insertSQL, values...).Scan(&id)
	if err != nil {
		return err
	}
	object.FieldByName("ID").SetInt(int64(id))

	return nil
}

// insert takes an object, grab its fields and performs an UPDATE operation
func (handler Handler) update(object reflect.Value) error {
	values := handler.assembleValuesArray(handler.updateMap, object)

	//run Update
	_, err := handler.db.Exec(handler.updateSQL, values...)
	if err != nil {
		return err
	}

	return nil
}

func (handler Handler) saveChilds(object reflect.Value) error {
	typeOfTable := reflect.TypeOf(handler.table)
	parentFieldName := typeOfTable.Name() + "ID"
	parentID := object.FieldByName("ID").Int()

	for fieldName := range handler.childHandlers {
		field, exists := typeOfTable.FieldByName(fieldName)

		if exists && field.Type.Kind().String() == "slice" {
			fieldValue := object.FieldByName(fieldName)
			elements := fieldValue.Len()
			var ids string
			for i := 0; i < elements; i++ {
				obj := fieldValue.Index(i).Elem().Addr().Interface()
				fieldValue.Index(i).Elem().FieldByName(parentFieldName).SetInt(parentID)
				err := handler.childHandlers[fieldName].Save(obj)
				if err != nil {
					return err
				}

				v := reflect.ValueOf(obj)
				ids = ids + strconv.Itoa(int(v.Elem().FieldByName("ID").Int())) + ", "
			}
			if len(ids) > 0 {
				ids = ids[:len(ids)-2]
			}

			//delete old records
			var err error
			if len(ids) > 0 {
				_, err = handler.childHandlers[fieldName].Delete().Where(fmt.Sprintf("%v = %v and id not in (%v)", parentFieldName, parentID, ids))
			} else {
				_, err = handler.childHandlers[fieldName].Delete().Where(fmt.Sprintf("%v = %v", parentFieldName, parentID))
			}
			if err != nil {
				return err
			}
		}

		if exists && field.Type.Kind().String() == "ptr" {
			if object.FieldByName(fieldName).IsNil() {
				_, err := handler.childHandlers[fieldName].Delete().Where(fmt.Sprintf("%v = %v", parentFieldName, parentID))
				if err != nil {
					return err
				}
			} else {
				fieldValue := object.FieldByName(fieldName).Elem()
				obj := fieldValue.Addr().Interface()

				v := reflect.ValueOf(obj)
				if v.Elem().FieldByName("ID").Int() <= 0 {
					_, err := handler.childHandlers[fieldName].Delete().Where(fmt.Sprintf("%v = %v", parentFieldName, parentID))
					if err != nil {
						return err
					}
				}

				fieldValue.FieldByName(parentFieldName).SetInt(parentID)
				err := handler.childHandlers[fieldName].Save(obj)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// assembleArguments takes an object and return an array of its values populated
// in the exact order required by the given arguments map. The resulting array
// of values is intended to be consumed by insert() or update()
func (handler Handler) assembleValuesArray(argurmentsMap []saveField, object reflect.Value) []interface{} {
	//build the values array
	var values []interface{}
	for _, argument := range argurmentsMap {
		switch argument.fieldType {
		case "int":
			values = append(values, int(object.FieldByName(argument.name).Int()))
		case "float64":
			values = append(values, float64(object.FieldByName(argument.name).Float()))
		case "string":
			values = append(values, string(object.FieldByName(argument.name).String()))
		case "bool":
			values = append(values, bool(object.FieldByName(argument.name).Bool()))
		case "Time":
			values = append(values, time.Time(object.FieldByName(argument.name).Interface().(time.Time)))
		default:
			continue
		}
	}
	return values
}

// assembleSQLInsert creates a SQL Insert string in the current handler.
// This string will be used by the save() method
func (handler *Handler) assembleSQLInsert() {
	typeOfTable := reflect.TypeOf(handler.table)
	tableName := typeOfTable.Name()
	validTypes := make(map[string]bool)
	validTypes["int"] = true
	validTypes["float64"] = true
	validTypes["string"] = true
	validTypes["bool"] = true
	validTypes["Time"] = true

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

		if !validTypes[typeOfTable.Field(i).Type.Name()] {
			continue
		}

		fieldMap = append(fieldMap, saveField{name: typeOfTable.Field(i).Name, fieldType: typeOfTable.Field(i).Type.Name()})
		sqlFields = sqlFields + fieldName + ", "
		sqlValues = sqlValues + "$" + strconv.Itoa(j) + ", "
		j = j + 1
	}

	sqlFields = sqlFields[:len(sqlFields)-2]
	sqlValues = sqlValues[:len(sqlValues)-2]
	sqlInstruction = sqlInstruction + sqlFields + ") values (" + sqlValues + ") RETURNING id;"

	handler.insertSQL = sqlInstruction
	handler.insertMap = fieldMap
}

// assembleSQLUpdate creates a SQL Update string in the current handler.
// This string will be used by the save() method
func (handler *Handler) assembleSQLUpdate() {
	typeOfTable := reflect.TypeOf(handler.table)
	tableName := typeOfTable.Name()
	validTypes := make(map[string]bool)
	validTypes["int"] = true
	validTypes["float64"] = true
	validTypes["string"] = true
	validTypes["bool"] = true
	validTypes["Time"] = true

	j := 1
	sqlInstruction := "update " + tableName + " set "
	var fieldMap []saveField
	for i := 0; i < typeOfTable.NumField(); i++ {
		fieldName := typeOfTable.Field(i).Name

		if fieldName == "ID" {
			continue
		}

		if !validTypes[typeOfTable.Field(i).Type.Name()] {
			continue
		}

		fieldMap = append(fieldMap, saveField{name: typeOfTable.Field(i).Name, fieldType: typeOfTable.Field(i).Type.Name()})
		sqlInstruction = sqlInstruction + fieldName + " = $" + strconv.Itoa(j) + ", "
		j = j + 1
	}
	fieldMap = append(fieldMap, saveField{name: "ID", fieldType: "int"})

	sqlInstruction = sqlInstruction[:len(sqlInstruction)-2]
	sqlInstruction = sqlInstruction + " where id = $" + strconv.Itoa(j)

	handler.updateSQL = sqlInstruction
	handler.updateMap = fieldMap
}
