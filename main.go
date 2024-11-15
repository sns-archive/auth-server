package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type User struct {
	id       string
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

	// TODO: idはUUIDを自動発番できるようにする
	user := User{
		id:       "123e4567-e89b-12d3-a456-426614174001", // 例としてUUIDを使用
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
		return nil, nil, err
	}
	// 接続確認するため
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	// Timeoutタイマーを終了する
	defer cancel()
	// 2秒待つ
	time.Sleep(2 * time.Second)
	if err := xdb.PingContext(ctx); err != nil {
		fmt.Printf("%+v\n", "timeoutしたかもよ")
		fmt.Println("timeout だよ", err.Error())
		return nil, nil, err
	}
	// NOTE: goでは、慣習的に使用しない返り値は「_」に格納する必要がある。
	return xdb, func() { _ = xdb.Close() }, nil
}

func insertUsers(db *sqlx.DB, user User) (sql.Result, error) {
	// NOTE: 一旦ID固定にするため、冪等な処理にしたいのでUpsertにする
	sql := `INSERT INTO users (id, username, email, password)
					VALUES
					(?, ?, ?, ?)
					ON DUPLICATE KEY UPDATE
					username = VALUES(username),
					email = VALUES(email),
					password = VALUES(password);`
	return db.Exec(sql, user.id, user.name, user.email, user.password)
}
