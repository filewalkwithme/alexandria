package orm

import (
	//	"fmt"
	"testing"
)

func TestAssembleSQLInsertStatement(t *testing.T) {
	var handler Handler
	handler.table = DSLTest{}

	//check if the sqlInsert string is assembled correctly
	handler.assembleSQLInsertStatement()
	expected := `insert into DSLTest(FieldString, FieldInt) values ($1, $2) RETURNING id;`
	got := handler.sqlInsert
	if got != expected {
		t.Fatalf("\nExpected:\t %v\nGot:\t\t %v\n", expected, got)
	}

	//expectedMap := []saveField{}
	mapInsert := handler.mapInsert
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
