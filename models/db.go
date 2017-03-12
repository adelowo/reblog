package models

import (
	"github.com/jmoiron/sqlx"
	_"github.com/mattn/go-sqlite3"
)

func MustNewDB(databaseName string) *DB {

	db := sqlx.MustConnect("sqlite3", databaseName)

	return &DB{db}
}
