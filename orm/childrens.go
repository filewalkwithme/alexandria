package orm

import (
	"fmt"
	"reflect"
)

func (handler *Handler) detectChildHandlers(orm *Orm) {
	fmt.Printf("detectChildHandlers\n")
	typeOfTable := reflect.TypeOf(handler.table)
	//v := reflect.ValueOf(handler.table)

	for i := 0; i < typeOfTable.NumField(); i++ {
		if typeOfTable.Field(i).Type.Kind().String() == "slice" {
			x := typeOfTable.Field(i).Type.Elem()
			y := reflect.New(x).Elem().Interface()
			fmt.Printf("->%v\n", reflect.TypeOf(y))

			ormChild, _ := orm.NewHandler(y)
			fmt.Printf("ormChild-> %v\n", ormChild)
		}
	}
}
