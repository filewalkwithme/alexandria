package orm

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSave(t *testing.T) {
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

	//Check if we get an error when we try to insert a diferent object for this handler
	dslTestWithoutID := DSLTestWithoutID{FieldString: "teststring", FieldInt: 123, FieldBool: false}
	err = ormTest.Save(&dslTestWithoutID)
	if err.Error() != "Object table name (DSLTestWithoutID) is diferent from handler table name (DSLTest)" {
		t.Fatalf("Want: `Object table name (DSLTestWithoutID) is diferent from handler table name (DSLTest)`, got: `%v`", err)
	}

	//save a new object
	dslTest := DSLTest{FieldString: "teststring", FieldInt: 123, FieldBool: true}
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
		t.Fatalf("want: a valid object, got nil")
	}

	if obj.FieldInt != 222 {
		t.Fatalf("want: 222, got: %v", obj.FieldInt)
	}

	//check if we got an error when trying to insert an object with negative ID
	dslTestNegativeID := DSLTest{FieldString: "teststring", FieldInt: 123, FieldBool: true, FieldFloat: 1.23, ID: -1}
	err = ormTest.Save(&dslTestNegativeID)
	if err.Error() != "Negative ID not allowed: {-1 teststring 123 true 1.23}" {
		t.Fatalf("want: `Negative ID not allowed: {-1 teststring 123 true 1.23}`, got `%v`", err)
	}
}

func TestInsertAndUpdate(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
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
	dslTest := DSLTest{FieldString: "teststring", FieldInt: 123, FieldBool: false, FieldFloat: 1.23}
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

	if obj.FieldBool != false {
		t.Fatalf("want: false, got: %v", obj.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj.FieldFloat) != "1.23" {
		t.Fatalf("want: 1.23, got: %v", obj.FieldFloat)
	}

	//change FieldInt atribute
	obj.FieldInt = 222
	obj.FieldBool = true
	obj.FieldFloat = 2.22

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

	if obj.FieldBool != true {
		t.Fatalf("want: true, got: %v", obj.FieldBool)
	}

	if fmt.Sprintf("%.2f", obj.FieldFloat) != "2.22" {
		t.Fatalf("want: 2.22, got: %v", obj.FieldFloat)
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
	orm, scream := ConnectToPostgres(dbURL)
	if scream != nil {
		panic(scream)
	}
	//create a new handler for DSLTest structure
	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	argurmentsMap := []saveField{{name: "FieldString", fieldType: "string"}, {name: "FieldInt", fieldType: "int"}, {name: "FieldBool", fieldType: "bool"}, {name: "FieldFloat", fieldType: "float64"}, {name: "ID", fieldType: "int"}}
	objectPtr := &DSLTest{ID: 1, FieldString: "teststring", FieldInt: 123, FieldBool: true, FieldFloat: 1.23}
	object := reflect.ValueOf(objectPtr).Elem()

	res := ormTest.assembleValuesArray(argurmentsMap, object)
	if len(res) != 5 {
		t.Fatalf("want: 5, got: %v", len(res))
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
			t.Fatalf("want: 123, got: %v", res[1].(int))
		}
	}

	switch v := res[2].(type) {
	default:
		t.Fatalf("want: bool, got: %v", v)
	case bool:
		if res[2].(bool) != true {
			t.Fatalf("want: true, got: %v", res[2].(bool))
		}
	}

	switch v := res[3].(type) {
	default:
		t.Fatalf("want: float64, got: %v", v)
	case float64:
		if fmt.Sprintf("%.2f", res[3].(float64)) != "1.23" {
			t.Fatalf("want: 1.23, got: %v", res[3].(float64))
		}
	}

	switch v := res[4].(type) {
	default:
		t.Fatalf("want: int, got: %v", v)
	case int:
		if res[4].(int) != 1 {
			t.Fatalf("want: teststring, got: %v", res[4].(int))
		}
	}
}

func TestAssembleSQLInsertStatement(t *testing.T) {
	var handler Handler
	handler.table = DSLTest{}

	//check if the sqlInsert string is assembled correctly
	handler.assembleSQLInsert()
	expected := `insert into DSLTest(FieldString, FieldInt, FieldBool, FieldFloat) values ($1, $2, $3, $4) RETURNING id;`
	got := handler.insertSQL
	if got != expected {
		t.Fatalf("\nExpected:\t %v\nGot:\t\t %v\n", expected, got)
	}

	//expectedMap := []saveField{}
	mapInsert := handler.insertMap
	if len(mapInsert) != 4 {
		t.Fatalf("\nmapInsert lengh:\t %v\nWant:\t\t\t 4\n", len(mapInsert))
	}

	if mapInsert[0].name != "FieldString" || mapInsert[0].fieldType != "string" {
		t.Fatalf("\nmapInsert[0] got:\t %v\nWant:\t\t\t {FieldString string}\n", mapInsert[0])
	}

	if mapInsert[1].name != "FieldInt" || mapInsert[1].fieldType != "int" {
		t.Fatalf("\nmapInsert[1] got:\t %v\nWant:\t\t\t {FieldInt int}\n", mapInsert[1])
	}

	if mapInsert[2].name != "FieldBool" || mapInsert[2].fieldType != "bool" {
		t.Fatalf("\nmapInsert[2] got:\t %v\nWant:\t\t\t {FieldBool bool}\n", mapInsert[2])
	}

	if mapInsert[3].name != "FieldFloat" || mapInsert[3].fieldType != "float64" {
		t.Fatalf("\nmapInsert[3] got:\t %v\nWant:\t\t\t {FieldBool bool}\n", mapInsert[3])
	}
}

func TestAssembleSQLUpdate(t *testing.T) {
	var handler Handler
	handler.table = DSLTest{}

	handler.assembleSQLUpdate()
	expected := `update DSLTest set FieldString = $1, FieldInt = $2, FieldBool = $3, FieldFloat = $4 where id = $5`
	got := handler.updateSQL
	if got != expected {
		t.Fatalf("\nExpected:\t %v\nGot:\t\t %v\n", expected, got)
	}

	updateMap := handler.updateMap
	if len(updateMap) != 5 {
		t.Fatalf("updateMap lenght: %v, want: 5", len(updateMap))
	}

	if updateMap[0].name != "FieldString" || updateMap[0].fieldType != "string" {
		t.Fatalf("updateMap[0] got:\t %v\nWant:\t\t\t {FieldString string}\n", updateMap[0])
	}

	if updateMap[1].name != "FieldInt" || updateMap[1].fieldType != "int" {
		t.Fatalf("updateMap[1] got:\t %v\nWant:\t\t\t {FieldInt int}\n", updateMap[1])
	}

	if updateMap[2].name != "FieldBool" || updateMap[2].fieldType != "bool" {
		t.Fatalf("updateMap[2] got: %v, want: {FieldBool bool}\n", updateMap[2])
	}

	if updateMap[3].name != "FieldFloat" || updateMap[3].fieldType != "float64" {
		t.Fatalf("updateMap[3] got: %v, want: {FieldBool float64}\n", updateMap[3])
	}

	if updateMap[4].name != "ID" || updateMap[4].fieldType != "int" {
		t.Fatalf("updateMap[4] got: %v, want: {ID int}\n", updateMap[3])
	}
}
