package dbHelper

import (
	"database/sql"
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
	db, err := sql.Open(globalCtx.driverName, globalCtx.dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

type RowScanner interface {
	Scan(...interface{}) error
}

//run on first time
func CreateDb() error {
	return nil
}
