package orm

import (
	"testing"
)

func TestSelectByID(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres()
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
	dslTest := DSLTest{FieldString: "teststring", FieldInt: 123}
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
	orm, scream := ConnectToPostgres()
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
	dslTest1 := DSLTest{FieldString: "teststring1", FieldInt: 111}
	err = ormTest.Save(&dslTest1)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//save a new object - 2
	dslTest2 := DSLTest{FieldString: "teststring2", FieldInt: 222}
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
	orm, scream := ConnectToPostgres()
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
	dslTest1 := DSLTest{FieldString: "teststring1", FieldInt: 111}
	err = ormTest.Save(&dslTest1)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//save a new object - 2
	dslTest2 := DSLTest{FieldString: "teststring2", FieldInt: 222}
	err = ormTest.Save(&dslTest2)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//save a new object - 3
	dslTest3 := DSLTest{FieldString: "teststring3", FieldInt: 333}
	err = ormTest.Save(&dslTest3)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//save a new object - 4
	dslTest4 := DSLTest{FieldString: "teststring4", FieldInt: 444}
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
	orm, scream := ConnectToPostgres()
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
	dslTest1 := DSLTest{FieldString: "teststring1", FieldInt: 111}
	err = ormTest.Save(&dslTest1)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//save a new object - 2
	dslTest2 := DSLTest{FieldString: "teststring2", FieldInt: 222}
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

	oldSelectScanMap := ormTest.selectScanMap
	ormTest.selectScanMap = make([]interface{}, 0)
	rows, err = ormTest.db.Query(ormTest.selectSQL)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	defer rows.Close()

	objects, err = ormTest.Select().buildArrayOfObjects(rows)
	if err.Error() != "sql: expected 3 destination arguments in Scan, not 0" {
		t.Fatalf("Expected: `sql: expected 3 destination arguments in Scan, not 0`, got: %v", err)
	}
	ormTest.selectScanMap = oldSelectScanMap
}

func TestBuildObject(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres()
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

	scan := make([]interface{}, 3)
	id := 1
	fieldString := "teststring1"
	fieldInt := 111
	scan[0] = &id
	scan[1] = &fieldString
	scan[2] = &fieldInt

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
}
