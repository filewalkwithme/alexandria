package orm

import (
	"reflect"
	"testing"
)

func TestSave(t *testing.T) {
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

	//Check if we get an error when we try to insert a diferent object for this handler
	dslTestWithoutID := DSLTestWithoutID{FieldString: "teststring", FieldInt: 123}
	err = ormTest.Save(&dslTestWithoutID)
	if err.Error() != "Object table name (DSLTestWithoutID) is diferent from handler table name (DSLTest)" {
		t.Fatalf("Want: `Object table name (DSLTestWithoutID) is diferent from handler table name (DSLTest)`, got: `%v`", err)
	}

	//save a new object
	dslTest := DSLTest{FieldString: "teststring", FieldInt: 123}
	err = ormTest.Save(&dslTest)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//check if the id is 1
	id := dslTest.ID
	if id != 1 {
		t.Fatalf("want: 1; got: %v", id)
	}

	//check if the object was persisted
	dslTestFind, err := ormTest.Select().ByID(id)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	obj := dslTestFind.(DSLTest)

	if dslTestFind == nil {
		t.Fatalf("want: a valida object, got nil")
	}

	//update the FieldInt atribute
	obj.FieldInt = 222
	err = ormTest.Save(&obj)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//check if the changes were persisted
	dslTestFindUptated, err := ormTest.Select().ByID(id)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	obj = dslTestFindUptated.(DSLTest)

	if dslTestFindUptated == nil {
		t.Fatalf("want: a valida object, got nil")
	}

	if obj.FieldInt != 222 {
		t.Fatalf("want: 222, got: %v", obj.FieldInt)
	}

	//check if we got an error when trying to insert an object with negative ID
	dslTestNegativeID := DSLTest{FieldString: "teststring", FieldInt: 123, ID: -1}
	err = ormTest.Save(&dslTestNegativeID)
	if err.Error() != "Negative ID not allowed: {-1 teststring 123}" {
		t.Fatalf("want: `Negative ID not allowed: {-1 teststring 123}`, got `%v`", err)
	}
}

func TestInsertAndUpdate(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres()
	if scream != nil {
		panic(scream)
	}

	//create a new handler for DSLTest structure
	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//DropTable & CreateTable
	ormTest.DropTable()
	ormTest.CreateTable()

	//create a test object
	dslTest := DSLTest{FieldString: "teststring", FieldInt: 123}
	object := reflect.ValueOf(&dslTest).Elem()

	//insert on db
	err = ormTest.insert(object)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	if dslTest.ID != 1 {
		t.Fatalf("\ndslTest.ID got:\t %v\nWant:\t\t\t 1\n", dslTest.ID)
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

	if obj.FieldInt != 123 {
		t.Fatalf("want: 123, got: %v", obj.FieldInt)
	}

	//change FieldInt atribute
	obj.FieldInt = 222

	//update on db
	object = reflect.ValueOf(&obj).Elem()
	err = ormTest.update(object)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	//check if the object was correctly updated
	dslTestFind, err = ormTest.Select().ByID(1)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}
	obj = dslTestFind.(DSLTest)

	if dslTestFind == nil {
		t.Fatalf("want: a valida object, got nil")
	}

	if obj.FieldInt != 222 {
		t.Fatalf("want: 222, got: %v", obj.FieldInt)
	}

	//force a error on insert
	oldInsertSQL := ormTest.insertSQL
	ormTest.insertSQL = "wrong-sql"
	err = ormTest.insert(object)
	if err.Error() != "pq: syntax error at or near \"wrong\"" {
		t.Fatalf("want: `pq: syntax error at or near \"wrong\"`, got: `%v`", err)
	}
	ormTest.insertSQL = oldInsertSQL

	//force a error on insert
	oldUpdateSQL := ormTest.updateSQL
	ormTest.updateSQL = "wrong-sql"
	err = ormTest.update(object)
	if err.Error() != "pq: syntax error at or near \"wrong\"" {
		t.Fatalf("want: `pq: syntax error at or near \"wrong\"`, got: `%v`", err)
	}
	ormTest.updateSQL = oldUpdateSQL
}

func TestAssembleValuesArray(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres()
	if scream != nil {
		panic(scream)
	}
	//create a new handler for DSLTest structure
	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	argurmentsMap := []saveField{{name: "FieldString", fieldType: "string"}, {name: "FieldInt", fieldType: "int"}, {name: "ID", fieldType: "int"}}
	objectPtr := &DSLTest{ID: 1, FieldString: "teststring", FieldInt: 123}
	object := reflect.ValueOf(objectPtr).Elem()

	res := ormTest.assembleValuesArray(argurmentsMap, object)
	if len(res) != 3 {
		t.Fatalf("want: 3, got: %v", len(res))
	}

	switch v := res[0].(type) {
	default:
		t.Fatalf("want: string, got: %v", v)
	case string:
		if res[0].(string) != "teststring" {
			t.Fatalf("want: teststring, got: %v", res[0].(string))
		}
	}

	switch v := res[1].(type) {
	default:
		t.Fatalf("want: int, got: %v", v)
	case int:
		if res[1].(int) != 123 {
			t.Fatalf("want: teststring, got: %v", res[1].(int))
		}
	}

	switch v := res[2].(type) {
	default:
		t.Fatalf("want: int, got: %v", v)
	case int:
		if res[2].(int) != 1 {
			t.Fatalf("want: teststring, got: %v", res[2].(int))
		}
	}
}

func TestAssembleSQLInsertStatement(t *testing.T) {
	var handler Handler
	handler.table = DSLTest{}

	//check if the sqlInsert string is assembled correctly
	handler.assembleSQLInsert()
	expected := `insert into DSLTest(FieldString, FieldInt) values ($1, $2) RETURNING id;`
	got := handler.insertSQL
	if got != expected {
		t.Fatalf("\nExpected:\t %v\nGot:\t\t %v\n", expected, got)
	}

	//expectedMap := []saveField{}
	mapInsert := handler.insertMap
	if len(mapInsert) != 2 {
		t.Fatalf("\nmapInsert lengh:\t %v\nWant:\t\t\t 2\n", len(mapInsert))
	}

	if mapInsert[0].name != "FieldString" || mapInsert[0].fieldType != "string" {
		t.Fatalf("\nmapInsert[0] got:\t %v\nWant:\t\t\t {FieldString string}\n", mapInsert[0])
	}

	if mapInsert[1].name != "FieldInt" || mapInsert[1].fieldType != "int" {
		t.Fatalf("\nmapInsert[0] got:\t %v\nWant:\t\t\t {FieldInt int}\n", mapInsert[1])
	}
}

func TestAssembleSQLUpdate(t *testing.T) {
	var handler Handler
	handler.table = DSLTest{}

	handler.assembleSQLUpdate()
	expected := `update DSLTest set FieldString = $1, FieldInt = $2 where id = $3`
	got := handler.updateSQL
	if got != expected {
		t.Fatalf("\nExpected:\t %v\nGot:\t\t %v\n", expected, got)
	}

	updateMap := handler.updateMap
	if len(updateMap) != 3 {
		t.Fatalf("updateMap lenght: %v, want: 3", len(updateMap))
	}

	if updateMap[0].name != "FieldString" || updateMap[0].fieldType != "string" {
		t.Fatalf("updateMap[0] got:\t %v\nWant:\t\t\t {FieldString string}\n", updateMap[0])
	}

	if updateMap[1].name != "FieldInt" || updateMap[1].fieldType != "int" {
		t.Fatalf("updateMap[1] got:\t %v\nWant:\t\t\t {FieldInt int}\n", updateMap[1])
	}

	if updateMap[2].name != "ID" || updateMap[2].fieldType != "int" {
		t.Fatalf("updateMap[2] got: %v, want: {ID int}\n", updateMap[2])
	}
}
