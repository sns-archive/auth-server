package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type User struct {
	id       uuid.UUID
	name     string
	email    string
	password string
}

func main() {
	helloStruct := User{
		name: "うんち💩",
	}
	// NOTE: +vでvalueだけでなく、keyも表示できる
	fmt.Printf("%+v\n", helloStruct)

	ctx := context.Background()
	xdb, cleanup, err := connectDB(ctx)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer cleanup()

	uuid, err := uuid.NewV7()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	fmt.Printf("生成されたUUIDv7: %s\n", uuid.String())
	// TODO: idはUUIDを自動発番できるようにする
	user := User{
		id:       uuid, // 例としてUUIDを使用
		name:     "うんち💩",
		email:    "example + 1@example.com",
		password: "securepassword",
	}

	result, err := insertUsers(xdb, user)
	fmt.Printf("%+v\n", result)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

// connectDB はDBに接続する
func connectDB(ctx context.Context) (*sqlx.DB, func(), error) {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "root",
		Addr:                 fmt.Sprintf("%s:%d", "127.0.0.1", 3307), // 127.0.01:3306
		DBName:               "sns_archive_jwt",
		ParseTime:            true,
		Net:                  "tcp",
		AllowNativePasswords: true,
	}
	xdb, err := sqlx.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open db: %w", err)
	}
	// 接続確認するため
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	// Timeoutタイマーを終了する
	defer cancel()
	if err := xdb.PingContext(ctx); err != nil {
		return nil, nil, fmt.Errorf("cannot ping: %w", err)
	}
	// NOTE: goでは、慣習的に使用しない返り値は「_」に格納する必要がある。
	return xdb, func() { _ = xdb.Close() }, nil
}

func insertUsers(db *sqlx.DB, user User) (sql.Result, error) {
	// NOTE: 一旦ID固定にするため、冪等な処理にしたいのでUpsertにする
	sql := `INSERT INTO users (id, username, email, password)
					VALUES
					(UUID_TO_BIN(?, 1), ?, ?, ?)
					ON DUPLICATE KEY UPDATE
					username = VALUES(username),
					email = VALUES(email),
					password = VALUES(password);`
	return db.Exec(sql, user.id, user.name, user.email, user.password)
}
