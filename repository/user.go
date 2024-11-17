package repository

import (
	"database/sql"
	"fmt"

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
	GetAllUsers(db *sqlx.DB) ([]User, error)
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

func (r *Repo) GetAllUsers(db *sqlx.DB) ([]User, error) {
	// TODO: users.usernameではなく、users.nameカラムにカラム名を変える
	query := `SELECT id, username AS name, email, password FROM users;`
	var users []User
	err := db.Select(&users, query)
	if err != nil {
		fmt.Printf("ユーザーを取得できませんでした: %v\n", err)
		return nil, err
	}
	return users, nil
}
