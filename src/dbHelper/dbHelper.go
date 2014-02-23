package dbHelper

import (
	"database/sql"
	"filelock"
	"log"
)

type connContext struct {
	driverName     string
	dataSourceName string
}

var globalCtx connContext

func Init(driverName, dataSourceName string) {
	globalCtx.driverName = driverName
	globalCtx.dataSourceName = dataSourceName
}

func Open() *sql.DB {
	filelock.Lock()

	db, err := sql.Open(globalCtx.driverName, globalCtx.dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func Close(db *sql.DB) {
	db.Close()

	filelock.Unlock()
}

type RowScanner interface {
	Scan(...interface{}) error
}

//run on first time
func CreateDb() error {
	return nil
}
