package models

import "github.com/jmoiron/sqlx"

type DataStore interface {
	UserStore
	PostStore
}

type DB struct {
	*sqlx.DB
}
