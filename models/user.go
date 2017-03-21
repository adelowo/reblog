package models

import (
	"github.com/adelowo/gotils/hasher"
	"github.com/adelowo/reblog/utils"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserStore interface {
	FindByEmail(email string) (User, error)
	DoesUserExist(email, moniker string) bool
	FindByMoniker(moniker string) (User, error)
	CreateUser(u *User) error
	CreateCollaborator(email string) error
	FindCollaboratorByToken(token string) (Collaborator, error)
	DeleteCollaborator(Collaborator) error
}

type User struct {
	ID        int       `db:"id"`
	Moniker   string    `db:"moniker"`
	Type      int       `db:"type"`
	Name      string    `db:"full_name"`
	About     string    `db:"about"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Password  string    `db:"password"`
}

type Collaborator struct {
	ID        int       `db:"id"`
	Token     string    `db:"token"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
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
	//_, err1 := db.FindByEmail(email)
	//_, err2 := db.FindByMoniker(moniker)
	//
	//return err1 == nil && err2 == nil

	var u User
	stmt, err := db.Preparex("SELECT * FROM users WHERE email=? OR moniker=?")

	//Just silence the error
	//All we want is a bool
	if err != nil {
		return false
	}

	rows := stmt.QueryRowx(email, moniker)

	err = rows.Scan(&u)

	if err != nil {
		return false
	}

	return true
}

func (db *DB) CreateUser(u *User) error {

	hashed, err := hasher.NewBcryptHasher(bcrypt.DefaultCost).Hash(u.Password)

	if err != nil {
		return nil
	}

	now := time.Now().UTC()

	stmt, err := db.Preparex("INSERT INTO users(moniker,full_name,password,email,created_at,updated_at) VALUES(?,?,?,?,?,?)")

	if err != nil {
		return errors.Wrap(err, "Could not prepare the insert statement")
	}

	count, err := stmt.MustExec(u.Moniker, u.Name, hashed, u.Email, now, now).
		RowsAffected()

	if err == nil && count == 1 {
		return nil
	}

	return errors.Wrap(err, "Could not create user")
}

func (db *DB) CreateCollaborator(email string) error {

	token, err := utils.NewTokenGenerator().Generate()

	if err != nil {
		return errors.Wrap(err, "Could not generate token for collaborator")
	}

	var u Collaborator

	stmt, err := db.Preparex("SELECT * FROM collaborator_tokens WHERE email=?")

	if err != nil {
		return errors.Wrap(err, "Could not prepare statement")
	}

	err = stmt.QueryRowx(email).StructScan(&u)

	createdAt := time.Now()

	if err != nil {
		//The user does not exist, we can add the collaborator

		stmt, err = db.Preparex("INSERT INTO collaborator_tokens(email,token,created_at) VALUES(?,?,?)")
		if err == nil {

			if id, err := stmt.MustExec(email, token, createdAt).LastInsertId(); err == nil && id == 1 {
				return nil
			}
		}

		return errors.Wrap(err, "An error occured while preparing the insert statement")
	}

	//THe user def exists, so we update here
	stmt, err = db.Preparex("UPDATE collaborator_tokens SET token=?,created_at=? WHERE email=?")

	if err != nil {
		return errors.Wrap(err, "An error occured while preparing the update statement")
	}

	if id, err := stmt.MustExec(token, createdAt, email).LastInsertId(); err == nil && id == 1 {
		return nil
	}

	return errors.Wrap(err, "An error occured while trying to update the collaborator's row")

}

func (db *DB) FindCollaboratorByToken(token string) (Collaborator, error) {

	var c Collaborator

	stmt, err := db.Preparex("SELECT * FROM collaborator_tokens WHERE token=?")

	if err != nil {
		return Collaborator{}, errors.Wrap(err, "Failed to prepare statement")
	}

	err = stmt.QueryRowx(token).
		StructScan(&c)

	if err != nil {
		return Collaborator{}, errors.Wrap(err, "Collaborator not found")
	}

	return c, nil
}

func (db *DB) DeleteCollaborator(c Collaborator) error {
	return nil
}
