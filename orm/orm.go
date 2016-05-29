package orm

import (
	"database/sql"
	//"fmt"
	"reflect"

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
	table          interface{}
	tableName      string
	db             *sql.DB
	sqlCreateTable string
	sqlDropTable   string
	sqlInsert      string
	mapInsert      []saveField
}

//Finder represents the result of a find operation
type Finder struct {
	table interface{}
	db    *sql.DB
}

//Deleter represents a delete operation
type Deleter struct {
	table interface{}
	db    *sql.DB
}

//NewHandler returns a Handler object to manipulate a given table
func (orm Orm) NewHandler(table interface{}) (Handler, error) {
	typeOfTable := reflect.TypeOf(table)
	tableName := typeOfTable.Name()

	handler := Handler{db: orm.db, table: table, tableName: tableName}

	//insert
	err := handler.assembleSQLInsertStatement()
	if err != nil {
		return Handler{}, err
	}

	return handler, err
}

//CreateTable is just a wrapper for the internal method createTable
func (handler Handler) CreateTable() error {
	return handler.createTable()
}

//DropTable is just a wrapper for the internal method dropTable
func (handler Handler) DropTable() error {
	return handler.dropTable()
}

//Insert is just a wrapper for the internal method insert
func (handler Handler) Insert(object interface{}) error {
	err := handler.insert(object)
	return err
}

//Find returns a Finder object
func (handler Handler) Find() Finder {
	return Finder{db: handler.db, table: handler.table}
}

//Where returns an array containing all results of a SELECT
func (f Finder) Where(where string) []interface{} {
	return f.findWhere(f.table, where)
}

//ByID returns an array containing all results of a SELECT
func (f Finder) ByID(id int) interface{} {
	return f.findByID(f.table, id)
}

//All returns an array containing all results of a SELECT
func (f Finder) All() []interface{} {
	return f.findAll(f.table)
}

//Delete returns a Finder object
func (handler Handler) Delete() Deleter {
	return Deleter{db: handler.db, table: handler.table}
}

//Where perform a DELETE operation
func (d Deleter) Where(where string) int {
	return d.deleteWhere(d.table, where)
}

//ByID perform a DELETE operation
func (d Deleter) ByID(id int) int {
	return d.deleteByID(d.table, id)
}

//All perform a DELETE operation
func (d Deleter) All() int {
	return d.deleteAll(d.table)
}

//----------------------------
