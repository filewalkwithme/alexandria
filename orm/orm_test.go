package orm

import (
	//"log"
	"testing"
)

var dbURL = "user=docker password=docker dbname=docker sslmode=disable"

func TestConnectToPostgres(t *testing.T) {
	dbDriver = "not-found"
	orm, scream := ConnectToPostgres(dbURL)
	if orm.db != nil {
		t.Fatalf("Want: nil, got: %v", orm)
	}

	if scream.Error() != `sql: unknown driver "not-found" (forgotten import?)` {
		t.Fatalf("Want: `sql: unknown driver \"not-found\" (forgotten import?)`, got: `%v`", scream)
	}

	dbDriver = "postgres"
	orm, scream = ConnectToPostgres("user=WRONG password=docker dbname=docker sslmode=disable")
	if orm.db != nil {
		t.Fatalf("Want: nil, got: %v", orm)
	}

	if scream.Error() != `pq: password authentication failed for user "WRONG"` {
		t.Fatalf("Want: `pq: password authentication failed for user \"WRONG\"`, got: `%v`", scream)
	}
}
