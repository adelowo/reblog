package models

type UserStore interface {
	FindByEmail(email string) (*User, error)
}

type User struct {
	ID        int
	Moniker   string `db:"moniker"`
	Type      int `db:"type"`
	Name      string `db:"full_name"`
	About     string `db:"about"`
	Email     string `db:email`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

func (db *DB ) FindByEmail(email string) (*User, error)  {


}