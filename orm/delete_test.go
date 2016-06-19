package orm

import (
	"testing"
)

func TestDeleteByID(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
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
	dslTest := DSLTest{FieldString: "teststring", FieldInt: 123, FieldBool: false}
	err = ormTest.Save(&dslTest)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	id := dslTest.ID
	if id != 1 {
		t.Fatalf("\ndslTest.ID got:\t %v\nWant:\t\t\t 1\n", id)
	}

	n, err := ormTest.Delete().ByID(id)

	if n != 1 {
		t.Fatalf("\ndeleted got:\t %v\nWant:\t\t\t 1\n", n)
	}

	if err != nil {
		t.Fatalf("\nerr got:\t %v\nWant:\t\t\t nil\n", err)
	}

	dslTestFind, err := ormTest.Select().ByID(id)
	if dslTestFind != nil {
		t.Fatalf("\ndslTest got:\t %v\nWant:\t\t\t nil\n", dslTest)
	}

	if err.Error() != "sql: no rows in result set" {
		t.Fatalf("want: `sql: no rows in result set`, got: `%v`", err)
	}

	oldDeleteSQL := ormTest.deleteSQL
	ormTest.deleteSQL = "wrong-sql"
	n, err = ormTest.Delete().ByID(id)
	if err.Error() != "pq: syntax error at or near \"wrong\"" {
		t.Fatalf("want: `pq: syntax error at or near \"wrong\"`, got: `%v`", err)
	}
	ormTest.deleteSQL = oldDeleteSQL
}

func TestDeleteWhere(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
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
	dslTest1 := DSLTest{FieldString: "teststring1", FieldInt: 111, FieldBool: true}
	err = ormTest.Save(&dslTest1)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	dslTest2 := DSLTest{FieldString: "teststring2", FieldInt: 222, FieldBool: false}
	err = ormTest.Save(&dslTest2)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	dslTestResults, err := ormTest.Select().Where("FieldString like 'teststring%'")
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	if len(dslTestResults) != 2 {
		t.Fatalf("want: 2, got: %v", len(dslTestResults))
	}

	n, err := ormTest.Delete().Where("FieldString like 'teststring%'")
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	if n != 2 {
		t.Fatalf("want: 2, got: %v", n)
	}

	oldDeleteSQL := ormTest.deleteSQL
	ormTest.deleteSQL = "wrong-sql"
	n, err = ormTest.Delete().Where("FieldString like 'teststring%'")
	if err.Error() != "pq: syntax error at or near \"wrong\"" {
		t.Fatalf("want: `pq: syntax error at or near \"wrong\"`, got: `%v`", err)
	}
	ormTest.deleteSQL = oldDeleteSQL
}

func TestDeleteAll(t *testing.T) {
	//connect to Postgres
	orm, scream := ConnectToPostgres(dbURL)
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
	dslTest1 := DSLTest{FieldString: "teststring1", FieldInt: 111, FieldBool: true}
	err = ormTest.Save(&dslTest1)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	dslTest2 := DSLTest{FieldString: "teststring2", FieldInt: 222, FieldBool: false}
	err = ormTest.Save(&dslTest2)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	dslTest3 := DSLTest{FieldString: "teststring3", FieldInt: 333, FieldBool: true}
	err = ormTest.Save(&dslTest3)
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	dslTestResults, err := ormTest.Select().All()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	if len(dslTestResults) != 3 {
		t.Fatalf("want: 3, got: %v", len(dslTestResults))
	}

	n, err := ormTest.Delete().All()
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	if n != 3 {
		t.Fatalf("want: 3, got: %v", n)
	}

	oldDeleteSQL := ormTest.deleteSQL
	ormTest.deleteSQL = "wrong-sql"
	n, err = ormTest.Delete().All()
	if err.Error() != "pq: syntax error at or near \"wrong\"" {
		t.Fatalf("want: `pq: syntax error at or near \"wrong\"`, got: `%v`", err)
	}
	ormTest.deleteSQL = oldDeleteSQL
}

func TestAssembleSQLDelete(t *testing.T) {
	orm, scream := ConnectToPostgres(dbURL)
	if scream != nil {
		panic(scream)
	}

	ormTest, err := orm.NewHandler(DSLTest{})
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if ormTest.deleteSQL != "delete from DSLTest" {
		t.Fatalf("want: 'delete from DSLTest', got: '%v'", ormTest.deleteSQL)
	}
}
