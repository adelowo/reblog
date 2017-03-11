package models

import "github.com/jmoiron/sqlx"

type DataStore interface {
	UserStore
}

type DB struct {
	*sqlx.DB
}
