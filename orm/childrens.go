package orm

import (
	"fmt"
	"reflect"
)

func (handler *Handler) detectChildHandlers(orm *Orm) error {
	typeOfTable := reflect.TypeOf(handler.table)
	handler.childHandlers = make(map[string]Handler)

	parentFieldName := typeOfTable.Name() + "ID"
	for i := 0; i < typeOfTable.NumField(); i++ {
		if typeOfTable.Field(i).Type.Kind().String() == "slice" {
			elemType := typeOfTable.Field(i).Type.Elem().Elem()
			newElemInterface := reflect.New(elemType).Elem().Interface()
			typeOfSubTable := reflect.TypeOf(newElemInterface)

			if typeOfTable.Field(i).Type.Elem().Kind().String() != "ptr" {
				return fmt.Errorf("Struct %v got non-pointer type for slice field: %v, expected: pointer", typeOfSubTable.String(), typeOfTable.Field(i))
			}

			childHandler, err := orm.NewHandler(newElemInterface)
			if err != nil {
				return err
			}

			// check if type contains parent id
			parentFieldFound := false

			// list type fields
			for j := 0; j < typeOfSubTable.NumField(); j++ {
				if parentFieldName == typeOfSubTable.Field(j).Name {
					if typeOfSubTable.Field(j).Type.String() != "int" {
						return fmt.Errorf("Struct %v got wrong type for %v field: %v, expected: int", typeOfSubTable.String(), parentFieldName, typeOfSubTable.Field(j).Type.String())
					}
					parentFieldFound = true
				}
			}

			fieldName := typeOfTable.Field(i).Name
			if !parentFieldFound {
				return fmt.Errorf("%v field not found on struct %v", parentFieldName, typeOfSubTable.String())
			}

			handler.childHandlers[fieldName] = childHandler
		}

		if typeOfTable.Field(i).Type.Kind().String() == "ptr" {
			structType := typeOfTable.Field(i).Type.Elem()
			newStructInterface := reflect.New(structType).Elem().Interface()

			childHandler, err := orm.NewHandler(newStructInterface)
			if err != nil {
				return err
			}

			// check if type contains parent id
			parentFieldFound := false

			// list type fields
			typeOfSubTable := reflect.TypeOf(newStructInterface)
			for j := 0; j < typeOfSubTable.NumField(); j++ {
				if parentFieldName == typeOfSubTable.Field(j).Name {
					if typeOfSubTable.Field(j).Type.String() != "int" {
						return fmt.Errorf("Struct %v got wrong type for %v field: %v, expected: int", typeOfSubTable.String(), parentFieldName, typeOfSubTable.Field(j).Type.String())
					}
					parentFieldFound = true
				}
			}

			fieldName := typeOfTable.Field(i).Name
			if !parentFieldFound {
				return fmt.Errorf("%v field not found on struct %v", parentFieldName, typeOfSubTable.String())
			}

			handler.childHandlers[fieldName] = childHandler
		}

	}
	return nil
}
