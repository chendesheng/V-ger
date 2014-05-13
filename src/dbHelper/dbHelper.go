package dbHelper

import (
	"database/sql"
	"filelock"
	"log"
	"time"
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
	b := time.Now()
	filelock.Lock()
	dur := time.Since(b)
	if dur > 10*time.Millisecond {
		println("filelock", dur.String())
	}

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
