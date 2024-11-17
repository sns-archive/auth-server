package repository

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id       uuid.UUID
	Name     string
	Email    string
	Password string
}

type Repositorier interface {
	InsertUsers(db *sqlx.DB, user User) (sql.Result, error)
}

type Repo struct{}

func NewRepo() Repositorier {
	return &Repo{}
}

func (r *Repo) InsertUsers(db *sqlx.DB, user User) (sql.Result, error) {
	sql := `INSERT INTO users (id, username, email, password)
					VALUES
					(UUID_TO_BIN(?, 1), ?, ?, ?)
					ON DUPLICATE KEY UPDATE
					username = VALUES(username),
					email = VALUES(email),
					password = VALUES(password);`
	return db.Exec(sql, user.Id, user.Name, user.Email, user.Password)
}
