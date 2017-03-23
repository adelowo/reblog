package models

import (
	"github.com/pkg/errors"
	"time"
)

type PostStore interface {
	CreatePost(p Post, userType int) error
	FindPostBySlug(slug string) (Post, error)
	FindPostByTitle(title string) (Post, error)
}

type Post struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Slug      string    `db:"slug"`
	Content   string    `db:"content"`
	Status    int       `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	User      User      `db:"user_id"`
}

func (db *DB) CreatePost(p Post, userType int) error {

	now := time.Now()

	p.CreatedAt = now
	p.UpdatedAt = now

	stmt, err := db.Preparex("INSERT INTO posts(title, slug, content, status, created_at, updated_at, user_id")

	if err != nil {
		return errors.Wrap(err, "An error occurred while we tried preparing the statement")
	}

	r, err := stmt.MustExec(p.Title, p.Slug, p.Content, p.Status, p.CreatedAt, p.UpdatedAt, userType).
		RowsAffected()

	if err == nil && r == 1 {
		return nil
	}

	return errors.Wrap(err, "An error occurred while we tried creating the post")
}

func (db *DB) FindPostBySlug(slug string) (Post, error) {
	var p Post

	stmt, err := db.Preparex("SELECT * FROM posts WHERE slug=?")

	if err != nil {
		return p, errors.Wrap(err, "COuld not prepare statement")
	}

	rows, err := stmt.Queryx(slug)

	if err != nil {
		return p, errors.Wrap(err, "An error occurred while trying to replace the prepared statement placeholder")
	}

	err = rows.StructScan(&p)

	if err != nil {
		return p, errors.Wrap(err, "Post does not exists")
	}

	return p, nil
}

func (db *DB) FindPostByTitle(title string) (Post, error) {
	var p Post

	stmt, err := db.Preparex("SELECT * FROM posts WHERE title=?")

	if err != nil {
		return p, errors.Wrap(err, "COuld not prepare statement")
	}

	rows, err := stmt.Queryx(title)

	if err != nil {
		return p, errors.Wrap(err, "An error occurred while trying to replace the prepared statement placeholder")
	}

	err = rows.StructScan(&p)

	if err != nil {
		return p, errors.Wrap(err, "Post does not exists")
	}

	return p, nil
}
