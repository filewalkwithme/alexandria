package orm

import (
	"fmt"
	"testing"
)

func TestSelectByID(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
	if scream != nil {
		panic(scream)
	}

	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//DropTable & CreateTable
	ormTest.DropTable()
	ormTest.CreateTable()

	//save a new object
	dslTest := DSLTest{FieldString: "teststring", FieldInt: 123, FieldBool: true, FieldFloat: 1.23}
	err = ormTest.Save(&dslTest)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//check if the object was persisted
	dslTestFind, err := ormTest.Select().ByID(1)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	obj := dslTestFind.(DSLTest)

	if dslTestFind == nil {
		t.Fatalf("want: a valida object, got nil")
	}

	if obj.ID != 1 {
		t.Fatalf("want: 1 got `%v`", obj.ID)
	}

	if obj.FieldString != "teststring" {
		t.Fatalf("want: `teststring`, got `%v`", obj.FieldString)
	}

	if obj.FieldInt != 123 {
		t.Fatalf("want: 123 got `%v`", obj.FieldInt)
	}

	if obj.FieldBool != true {
		t.Fatalf("want: true got `%v`", obj.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj.FieldFloat) != "1.23" {
		t.Fatalf("want: 1.23 got `%v`", obj.FieldFloat)
	}

	//force a error on select
	oldSelectSQL := ormTest.selectSQL
	ormTest.selectSQL = "wrong-sql"
	dslTestFind, err = ormTest.Select().ByID(1)
	if err.Error() != "pq: syntax error at or near \"wrong\"" {
		t.Fatalf("want: `pq: syntax error at or near \"wrong\"`, got: `%v`", err)
	}
	ormTest.selectSQL = oldSelectSQL
}

func TestSelectAll(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
	if scream != nil {
		panic(scream)
	}

	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//DropTable & CreateTable
	ormTest.DropTable()
	ormTest.CreateTable()

	//save a new object - 1
	dslTest1 := DSLTest{FieldString: "teststring1", FieldInt: 111, FieldBool: true, FieldFloat: 1.11}
	err = ormTest.Save(&dslTest1)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//save a new object - 2
	dslTest2 := DSLTest{FieldString: "teststring2", FieldInt: 222, FieldBool: false, FieldFloat: 2.22}
	err = ormTest.Save(&dslTest2)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//check if the object was persisted
	dslTestFindAll, err := ormTest.Select().All()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	if dslTestFindAll == nil {
		t.Fatalf("want: a valida object, got nil")
	}

	if len(dslTestFindAll) != 2 {
		t.Fatalf("want: 2, got: %v", len(dslTestFindAll))
	}

	obj1 := dslTestFindAll[0].(DSLTest)

	if obj1.ID != 1 {
		t.Fatalf("want: 1 got `%v`", obj1.ID)
	}

	if obj1.FieldString != "teststring1" {
		t.Fatalf("want: `teststring1`, got `%v`", obj1.FieldString)
	}

	if obj1.FieldInt != 111 {
		t.Fatalf("want: 111 got `%v`", obj1.FieldInt)
	}

	if obj1.FieldBool != true {
		t.Fatalf("want: true got `%v`", obj1.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj1.FieldFloat) != "1.11" {
		t.Fatalf("want: 1.11 got `%v`", obj1.FieldFloat)
	}

	obj2 := dslTestFindAll[1].(DSLTest)

	if obj2.ID != 2 {
		t.Fatalf("want: 2 got `%v`", obj2.ID)
	}

	if obj2.FieldString != "teststring2" {
		t.Fatalf("want: `teststring2`, got `%v`", obj2.FieldString)
	}

	if obj2.FieldInt != 222 {
		t.Fatalf("want: 222 got `%v`", obj2.FieldInt)
	}

	if obj2.FieldBool != false {
		t.Fatalf("want: false got `%v`", obj2.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj2.FieldFloat) != "2.22" {
		t.Fatalf("want: 2.22 got `%v`", obj2.FieldFloat)
	}

	//force a error on select
	oldSelectSQL := ormTest.selectSQL
	ormTest.selectSQL = "wrong-sql"
	dslTestFindAll, err = ormTest.Select().All()
	if err.Error() != "pq: syntax error at or near \"wrong\"" {
		t.Fatalf("want: `pq: syntax error at or near \"wrong\"`, got: `%v`", err)
	}
	ormTest.selectSQL = oldSelectSQL
}

func TestSelectWhere(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
	if scream != nil {
		panic(scream)
	}

	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//DropTable & CreateTable
	ormTest.DropTable()
	ormTest.CreateTable()

	//save a new object - 1
	dslTest1 := DSLTest{FieldString: "teststring1", FieldInt: 111, FieldBool: true, FieldFloat: 1.11}
	err = ormTest.Save(&dslTest1)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//save a new object - 2
	dslTest2 := DSLTest{FieldString: "teststring2", FieldInt: 222, FieldBool: false, FieldFloat: 2.22}
	err = ormTest.Save(&dslTest2)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//save a new object - 3
	dslTest3 := DSLTest{FieldString: "teststring3", FieldInt: 333, FieldBool: true, FieldFloat: 3.33}
	err = ormTest.Save(&dslTest3)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//save a new object - 4
	dslTest4 := DSLTest{FieldString: "teststring4", FieldInt: 444, FieldBool: false, FieldFloat: 4.44}
	err = ormTest.Save(&dslTest4)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//check if the object was persisted
	dslTestFindAll, err := ormTest.Select().Where("id > 1")
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	if dslTestFindAll == nil {
		t.Fatalf("want: a valida object, got nil")
	}

	if len(dslTestFindAll) != 3 {
		t.Fatalf("want: 3, got: %v", len(dslTestFindAll))
	}

	obj1 := dslTestFindAll[0].(DSLTest)

	if obj1.ID != 2 {
		t.Fatalf("want: 2 got `%v`", obj1.ID)
	}

	if obj1.FieldString != "teststring2" {
		t.Fatalf("want: `teststring2`, got `%v`", obj1.FieldString)
	}

	if obj1.FieldInt != 222 {
		t.Fatalf("want: 222 got `%v`", obj1.FieldInt)
	}

	if obj1.FieldBool != false {
		t.Fatalf("want: false got `%v`", obj1.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj1.FieldFloat) != "2.22" {
		t.Fatalf("want: 2.22 got `%v`", obj1.FieldFloat)
	}

	obj2 := dslTestFindAll[1].(DSLTest)

	if obj2.ID != 3 {
		t.Fatalf("want: 3 got `%v`", obj2.ID)
	}

	if obj2.FieldString != "teststring3" {
		t.Fatalf("want: `teststring3`, got `%v`", obj2.FieldString)
	}

	if obj2.FieldInt != 333 {
		t.Fatalf("want: 333 got `%v`", obj2.FieldInt)
	}

	if obj2.FieldBool != true {
		t.Fatalf("want: true got `%v`", obj2.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj2.FieldFloat) != "3.33" {
		t.Fatalf("want: 3.33 got `%v`", obj2.FieldFloat)
	}

	obj3 := dslTestFindAll[2].(DSLTest)

	if obj3.ID != 4 {
		t.Fatalf("want: 4 got `%v`", obj3.ID)
	}

	if obj3.FieldString != "teststring4" {
		t.Fatalf("want: `teststring4`, got `%v`", obj3.FieldString)
	}

	if obj3.FieldInt != 444 {
		t.Fatalf("want: 444 got `%v`", obj3.FieldInt)
	}

	if obj3.FieldBool != false {
		t.Fatalf("want: false got `%v`", obj3.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj3.FieldFloat) != "4.44" {
		t.Fatalf("want: 4.44 got `%v`", obj3.FieldFloat)
	}

	//force a error on select
	oldSelectSQL := ormTest.selectSQL
	ormTest.selectSQL = "wrong-sql"
	dslTestFindAll, err = ormTest.Select().Where("id > 1")
	if err.Error() != "pq: syntax error at or near \"wrong\"" {
		t.Fatalf("want: `pq: syntax error at or near \"wrong\"`, got: `%v`", err)
	}
	ormTest.selectSQL = oldSelectSQL
}

func TestBuildArrayOfObjects(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
	if scream != nil {
		panic(scream)
	}

	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//DropTable & CreateTable
	ormTest.DropTable()
	ormTest.CreateTable()

	//save a new object - 1
	dslTest1 := DSLTest{FieldString: "teststring1", FieldInt: 111, FieldBool: false, FieldFloat: 1.11}
	err = ormTest.Save(&dslTest1)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//save a new object - 2
	dslTest2 := DSLTest{FieldString: "teststring2", FieldInt: 222, FieldBool: true, FieldFloat: 2.22}
	err = ormTest.Save(&dslTest2)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	rows, err := ormTest.db.Query(ormTest.selectSQL)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	defer rows.Close()

	objects, err := ormTest.Select().buildArrayOfObjects(rows)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	if objects == nil {
		t.Fatalf("want: a valida object, got nil")
	}

	if len(objects) != 2 {
		t.Fatalf("want: 2, got: %v", len(objects))
	}

	obj1 := objects[0].(DSLTest)

	if obj1.ID != 1 {
		t.Fatalf("want: 1 got `%v`", obj1.ID)
	}

	if obj1.FieldString != "teststring1" {
		t.Fatalf("want: `teststring1`, got `%v`", obj1.FieldString)
	}

	if obj1.FieldInt != 111 {
		t.Fatalf("want: 111 got `%v`", obj1.FieldInt)
	}

	if obj1.FieldBool != false {
		t.Fatalf("want: false got `%v`", obj1.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj1.FieldFloat) != "1.11" {
		t.Fatalf("want: 1.11 got `%v`", obj1.FieldFloat)
	}

	obj2 := objects[1].(DSLTest)

	if obj2.ID != 2 {
		t.Fatalf("want: 2 got `%v`", obj2.ID)
	}

	if obj2.FieldString != "teststring2" {
		t.Fatalf("want: `teststring2`, got `%v`", obj2.FieldString)
	}

	if obj2.FieldInt != 222 {
		t.Fatalf("want: 222 got `%v`", obj2.FieldInt)
	}

	if obj2.FieldBool != true {
		t.Fatalf("want: true got `%v`", obj2.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj2.FieldFloat) != "2.22" {
		t.Fatalf("want: 2.22 got `%v`", obj2.FieldFloat)
	}

	oldSelectScanMap := ormTest.selectScanMap
	ormTest.selectScanMap = make([]interface{}, 0)
	rows, err = ormTest.db.Query(ormTest.selectSQL)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	defer rows.Close()

	objects, err = ormTest.Select().buildArrayOfObjects(rows)
	if err.Error() != "sql: expected 5 destination arguments in Scan, not 0" {
		t.Fatalf("Expected: `sql: expected 5 destination arguments in Scan, not 0`, got: %v", err)
	}
	ormTest.selectScanMap = oldSelectScanMap
}

func TestBuildObject(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
	if scream != nil {
		panic(scream)
	}

	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//DropTable & CreateTable
	ormTest.DropTable()
	ormTest.CreateTable()

	scan := make([]interface{}, 5)
	id := 1
	fieldString := "teststring1"
	fieldInt := 111
	fieldBool := true
	fieldFloat := 1.23
	scan[0] = &id
	scan[1] = &fieldString
	scan[2] = &fieldInt
	scan[3] = &fieldBool
	scan[4] = &fieldFloat

	obj := (ormTest.Select().buildObject(scan).(DSLTest))

	if obj.ID != 1 {
		t.Fatalf("want: 1 got `%v`", obj.ID)
	}

	if obj.FieldString != "teststring1" {
		t.Fatalf("want: `teststring1`, got `%v`", obj.FieldString)
	}

	if obj.FieldInt != 111 {
		t.Fatalf("want: 111 got `%v`", obj.FieldInt)
	}

	if obj.FieldBool != true {
		t.Fatalf("want: true got `%v`", obj.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj.FieldFloat) != "1.23" {
		t.Fatalf("want: 1.23 got `%v`", obj.FieldFloat)
	}
}

func TestAssembleSQLSelect(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
	if scream != nil {
		panic(scream)
	}

	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if ormTest.selectSQL != "select ID, FieldString, FieldInt, FieldBool, FieldFloat from DSLTest" {
		t.Fatalf("Want: `select ID, FieldString, FieldInt, FieldBool, FieldFloat from DSLTest`, Got: `%v`", ormTest.selectSQL)
	}

	if len(ormTest.selectScanMap) != 5 {
		t.Fatalf("Want: 5, Got: %v", len(ormTest.selectScanMap))
	}

	obj := (ormTest.Select().buildObject(ormTest.selectScanMap).(DSLTest))

	if obj.ID != 0 {
		t.Fatalf("want: 0 got: %v", obj.ID)
	}

	if obj.FieldString != "" {
		t.Fatalf("want: ``, got: `%v`", obj.FieldString)
	}

	if obj.FieldInt != 0 {
		t.Fatalf("want: 0 got: %v", obj.FieldInt)
	}

	if obj.FieldBool != false {
		t.Fatalf("want: false got: %v", obj.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj.FieldFloat) != "0.00" {
		t.Fatalf("want: 0.00 got `%v`", obj.FieldFloat)
	}

	if len(ormTest.selectFieldNamesMap) != 5 {
		t.Fatalf("Want: 5, Got: %v", len(ormTest.selectFieldNamesMap))
	}

	if ormTest.selectFieldNamesMap[0] != "ID" {
		t.Fatalf("want: `ID` got `%v`", ormTest.selectFieldNamesMap[0])
	}

	if ormTest.selectFieldNamesMap[1] != "FieldString" {
		t.Fatalf("want: `FieldString` got `%v`", ormTest.selectFieldNamesMap[1])
	}

	if ormTest.selectFieldNamesMap[2] != "FieldInt" {
		t.Fatalf("want: `FieldInt` got `%v`", ormTest.selectFieldNamesMap[2])
	}

	if ormTest.selectFieldNamesMap[3] != "FieldBool" {
		t.Fatalf("want: `FieldBool` got `%v`", ormTest.selectFieldNamesMap[3])
	}

	if ormTest.selectFieldNamesMap[4] != "FieldFloat" {
		t.Fatalf("want: `FieldFloat` got `%v`", ormTest.selectFieldNamesMap[4])
	}
}
