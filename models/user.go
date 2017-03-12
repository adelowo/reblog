package models

import (
	"github.com/pkg/errors"
)

type UserStore interface {
	FindByEmail(email string) (User, error)
	DoesUserExist(email, moniker string) bool
	FindByMoniker(moniker string) (User, error)
}

type User struct {
	ID        int    `db:"id"`
	Moniker   string `db:"moniker"`
	Type      int    `db:"type"`
	Name      string `db:"full_name"`
	About     string `db:"about"`
	Email     string `db:"email"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	Password  string `db:"password"`
}

func (db *DB) FindByEmail(email string) (User, error) {

	var u User

	stmt, err := db.Preparex("SELECT * FROM users WHERE email=?")

	if err != nil {
		return User{}, errors.Wrap(err, "An error occurred while we tried preparing this statement")
	}

	row := stmt.QueryRowx(email)

	err = row.StructScan(&u)

	if err != nil {
		return User{}, errors.Wrap(err, "Could not find a user with the specified email address")
	}

	return u, nil
}

func (db *DB) FindByMoniker(moniker string) (User, error) {

	var u User

	stmt, err := db.Preparex("SELECT * FROM users WHERE moniker=?")

	if err != nil {
		return User{}, errors.Wrap(err, "An error occurred while we tried preparing this statement")
	}

	row := stmt.QueryRowx(moniker)

	err = row.StructScan(&u)

	if err != nil {
		return User{}, errors.Wrap(err, "Could not find a user with the specified username")
	}

	return u, nil
}

func (db *DB) DoesUserExist(email, moniker string) bool {
	_, err1 := db.FindByEmail(email)
	_, err2 := db.FindByMoniker(moniker)

	return err1 == nil && err2 == nil
}
