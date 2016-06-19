package orm

import (
	"testing"
)

//DSLTest is used to test SQL creation
type DSLTest struct {
	ID          int
	FieldString string
	FieldInt    int
	FieldBool   bool
	FieldFloat  float64
}

//DSLTestWithoutID is used to test SQL creation
type DSLTestWithoutID struct {
	FieldString string
	FieldInt    int
	FieldBool   bool
	FieldFloat  float64
}

func TestCreateTable(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
	if scream != nil {
		panic(scream)
	}

	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ormTest.DropTable()
	err = ormTest.CreateTable()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ormTest.sqlCreateTable = "super wrong sql"
	err = ormTest.CreateTable()

	if err.Error() != "pq: syntax error at or near \"super\"" {
		t.Fatalf("%v", err.Error())
	}

	ormTest, err = orm.NewHandler(DSLTestWithoutID{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ormTest.sqlCreateTable = ""
	ormTest.DropTable()
	err = ormTest.CreateTable()
	if err.Error() != "ID field not found on struct DSLTestWithoutID" {
		t.Fatalf("err: %v", err.Error())
	}
}

func TestDropTable(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
	if scream != nil {
		panic(scream)
	}

	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ormTest.CreateTable()
	err = ormTest.DropTable()
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	ormTest.CreateTable()
	ormTest.sqlDropTable = "super wrong sql"
	err = ormTest.DropTable()
	if err.Error() != "pq: syntax error at or near \"super\"" {
		t.Fatalf("%v", err.Error())
	}
}

func TestAssembleSQLCreateTable(t *testing.T) {
	var handler Handler
	handler.table = DSLTest{}

	expected := `create table DSLTest (ID serial NOT NULL, FieldString character varying, FieldInt integer, FieldBool boolean, FieldFloat real, constraint DSLTest_pkey primary key (id));`
	got, _ := handler.assembleSQLCreateTable()
	if got != expected {
		t.Fatalf("\nExpected:\t %v\nGot:\t\t %v\n", expected, got)
	}

	handler.table = DSLTestWithoutID{}
	got, err := handler.assembleSQLCreateTable()
	if got != "" {
		t.Fatalf("\nExpected:\t \"\"\nGot:\t\t %v\n", got)
	}

	if err.Error() != "ID field not found on struct DSLTestWithoutID" {
		t.Fatalf("\nExpected:\t %v\nGot:\t\t %v\n", "ID field not found on struct DSLTestWithoutID", err.Error())
	}
}

func TestAssembleSQLDropTable(t *testing.T) {
	var handler Handler
	handler.table = DSLTest{}

	expected := `drop table DSLTest;`
	got := handler.assembleSQLDropTable()
	if got != expected {
		t.Fatalf("\nExpected:\t %v\nGot:\t\t %v\n", expected, got)
	}
}
