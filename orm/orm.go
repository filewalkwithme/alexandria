package orm

import (
	"database/sql"
	"reflect"
	//needed to access postgres
	_ "github.com/lib/pq"
)

//Orm is the main struct on this package
type Orm struct {
	db *sql.DB
}

var dbDriver = "postgres"

//ConnectToPostgres open a connection to Posgres
func ConnectToPostgres(dbURL string) (Orm, error) {
	tmpDB, err := sql.Open(dbDriver, dbURL)
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

	//insert
	insertSQL string
	insertMap []saveField

	//update
	updateSQL string
	updateMap []saveField

	//select
	selectSQL           string
	selectFieldNamesMap []string
	selectScanMap       []interface{}

	//delete
	deleteSQL string

	//child handlers
	childHandlers map[string]Handler
}

//NewHandler returns a Handler object to manipulate a given table
func (orm Orm) NewHandler(table interface{}) (Handler, error) {
	typeOfTable := reflect.TypeOf(table)
	tableName := typeOfTable.Name()

	handler := Handler{db: orm.db, table: table, tableName: tableName}

	//build sql insert
	handler.assembleSQLInsert()

	//build sql update
	handler.assembleSQLUpdate()

	//build sql update
	handler.assembleSQLSelect()

	//build sql update
	handler.assembleSQLDelete()

	//load child handlers
	err := handler.detectChildHandlers(&orm)

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

//Save ....
func (handler Handler) Save(object interface{}) error {
	err := handler.save(object)
	return err
}

//Selecter represents the result of a find operation
type Selecter struct {
	handler Handler
}

//Select returns a Finder object
func (handler Handler) Select() Selecter {
	return Selecter{handler: handler}
}

//Where returns an array containing all results of a SELECT
func (s Selecter) Where(where string, arguments ...interface{}) ([]*interface{}, error) {
	return s.selectWhere(where, arguments...)
}

//ByID returns an array containing all results of a SELECT
func (s Selecter) ByID(id int) (*interface{}, error) {
	return s.selectByID(id)
}

//All returns an array containing all results of a SELECT
func (s Selecter) All() ([]*interface{}, error) {
	return s.selectAll()
}

//Deleter represents a delete operation
type Deleter struct {
	handler Handler
}

//Delete returns a Finder object
func (handler Handler) Delete() Deleter {
	return Deleter{handler: handler}
}

//Where perform a DELETE operation
func (d Deleter) Where(where string, arguments ...interface{}) (int, error) {
	return d.deleteWhere(where, arguments...)
}

//ByID perform a DELETE operation
func (d Deleter) ByID(id int) (int, error) {
	return d.deleteByID(id)
}

//All perform a DELETE operation
func (d Deleter) All() (int, error) {
	return d.deleteAll()
}

//----------------------------
