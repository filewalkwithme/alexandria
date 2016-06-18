package orm

import (
	"reflect"
)

func (deleter Deleter) deleteByID(id int) (int, error) {
	sqlInstruction := deleter.handler.deleteSQL + " where id = $1"

	result, err := deleter.handler.db.Exec(sqlInstruction, id)
	if err != nil {
		return 0, err
	}
	n, err := result.RowsAffected()

	return int(n), err
}

func (deleter Deleter) deleteWhere(where string, arguments ...interface{}) (int, error) {
	sqlInstruction := deleter.handler.deleteSQL + " where " + where

	result, err := deleter.handler.db.Exec(sqlInstruction, arguments...)
	if err != nil {
		return 0, err
	}
	n, err := result.RowsAffected()

	return int(n), err
}

func (deleter Deleter) deleteAll() (int, error) {
	sqlInstruction := deleter.handler.deleteSQL

	result, err := deleter.handler.db.Exec(sqlInstruction)
	if err != nil {
		return 0, err
	}
	n, err := result.RowsAffected()

	return int(n), err
}

func (handler *Handler) assembleSQLDelete() {
	typeOfTable := reflect.TypeOf(handler.table)
	tableName := typeOfTable.Name()

	sqlInstruction := "delete from " + tableName

	handler.deleteSQL = sqlInstruction
}
