package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func Connect() (*sql.DB, error) {
	var err error
	Db, err = sql.Open("sqlite3", "user.db?cache=shared&mode=rwc&_journal_mode=WAL&busy_timeout=10000")
	Db.SetMaxOpenConns(1)
	if err != nil {
		return nil, err
	}
	return Db, nil
}

func ConnectTest() (*sql.DB, error) {
	var err error
	Db, err = sql.Open("sqlite3", "file::memory:?cache=shared&mode=rwc&_journal_mode=WAL&busy_timeout=10000")
	Db.SetMaxOpenConns(1)
	if err != nil {
		return nil, err
	}
	return Db, nil
}
