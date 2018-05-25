package dbtile

// This file init connection (should be used just by main)
// And it share a globale "db" variable that is the connection pool to be used by all dbmap package

import (
	"fmt"

	sqlx "github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Connection pool to db remote
// this variable is global to dbremote package
var dbConn *sqlx.DB

// InitDb connection to remote database
func InitDb(sqliteFile string) error {
	// Prepare database connection
	var dbConnErr error
	dbConn, dbConnErr = sqlx.Open("sqlite3", sqliteFile)
	if nil != dbConnErr {
		fmt.Println("Database initializing connection error : ", dbConnErr.Error())
		return dbConnErr
	}

	// Check database connection ready
	dbConnErr = dbConn.Ping()
	if nil != dbConnErr {
		fmt.Println("Database first connection error : ", dbConnErr.Error())
	}

	return dbConnErr
}

// CloseDb close connection to remote database
// Should be called in main with defer keyword
func CloseDb() {
	dbConn.Close()
}
