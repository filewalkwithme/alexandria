package orm

import (
	//	"fmt"
	"testing"
)

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

func TestInsert(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres()
	if scream != nil {
		panic(scream)
	}

	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	//check if the object is stored in the table and if the ID is populated after insert
	ormTest.DropTable()
	ormTest.CreateTable()
	dslTest := DSLTest{FieldString: "teststring", FieldInt: 123}
	err = ormTest.Save(&dslTest)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	if dslTest.ID != 1 {
		t.Fatalf("\ndslTest.ID got:\t %v\nWant:\t\t\t 1\n", dslTest.ID)
	}

	//force an error by passing an object with a wrong type
	dslTestWithoutID := DSLTestWithoutID{FieldString: "teststring", FieldInt: 123}
	err = ormTest.Save(&dslTestWithoutID)
	if err.Error() != "Object table name (DSLTestWithoutID) is diferent from handler table name (DSLTest)" {
		t.Fatalf("Error expected: 'Object table name (DSLTestWithoutID) is diferent from handler table name (DSLTest)'\nError got: %v\n", err)
	}

	//force an error by droping the table before the insert
	ormTest.DropTable()
	err = ormTest.Save(&dslTest)
	if err.Error() != `pq: relation "dsltest" does not exist` {
		t.Fatalf("Error expected: 'pq: relation \"dsltest\" does not exist'\nError got: %v\n", err)
	}
}
