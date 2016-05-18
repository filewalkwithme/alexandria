package orm

import (
	"database/sql"
	//"fmt"

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
	db    *sql.DB
}

//All returns an array containing all results of a SELECT
func (f Finder) All() []interface{} {
	return f.findAll(f.table)
}

//Find perform a SELECT operation
func (orm Orm) Find(table interface{}) Finder {
	return Finder{db: orm.db, table: table}
}
