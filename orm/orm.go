package orm

import (
	"database/sql"
	"fmt"

	//needed to access postgres
	_ "github.com/lib/pq"
)

//Orm is the main struct on this package
type Orm struct {
	db *sql.DB
}

//ConnectToPostgres open a connection to Posgres
func ConnectToPostgres() (Orm, error) {
	tmpDB, err := sql.Open("postgres", "user=docker password=docker dbname=docker sslmode=disable")
	if err != nil {
		return Orm{}, err
	}

	err = tmpDB.Ping()
	if err != nil {
		return Orm{}, err
	}

	orm := Orm{}
	orm.db = tmpDB
	return orm, nil
}

//Finder represents the result of a find operation
type Finder struct {
	table interface{}
}

//First returns the first position of a collection
func (f Finder) First() interface{} {
	fmt.Printf("First: %v\n", f.table)
	return nil
}

//Find perform a SELECT operation
func (orm Orm) Find(table interface{}) Finder {
	return Finder{}
}
