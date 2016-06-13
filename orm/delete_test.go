package orm

import (
	"testing"
)

func TestDeleteByID(t *testing.T) {
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
