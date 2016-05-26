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

//Handler manipulates the table (create/destroy/save/finde/delete)
type Handler struct {
	table interface{}
	db    *sql.DB
}

//Finder represents the result of a find operation
type Finder struct {
	table interface{}
	db    *sql.DB
}

//Handle returns a Handler object to manipulate the given table
func (orm Orm) Handle(table interface{}) Handler {
	return Handler{db: orm.db, table: table}
}

//Save perform an INSERT operation
func (handler Handler) Save(table interface{}) (interface{}, error) {
	return handler.save(table)
}

//Find returns a Finder object
func (handler Handler) Find() Finder {
	return Finder{db: handler.db, table: handler.table}
}

//Where returns an array containing all results of a SELECT
func (f Finder) Where(where string) []interface{} {
	return f.findWhere(f.table, where)
}

//----------------------------

//Deleter represents a delete operation
type Deleter struct {
	table interface{}
	db    *sql.DB
}

//All returns an array containing all results of a SELECT
func (f Finder) All() []interface{} {
	return f.findAll(f.table)
}

//One returns a single object which matches a given ID
func (f Finder) One() interface{} {
	return f.find(f.table)
}

//Find perform a SELECT operation
func (orm Orm) Find(table interface{}) Finder {
	return Finder{db: orm.db, table: table}
}

//Save perform an INSERT operation
func (orm Orm) Save(table interface{}) (interface{}, error) {
	//return orm.save(table)
	return table, nil
}

//Delete perform a DELETE operation
func (orm Orm) Delete(table interface{}) {
	orm.delete(table)
}
