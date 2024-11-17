package main

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type User struct {
	id       uuid.UUID
	name     string
	email    string
	password string
}

func main() {
	ctx := context.Background()
	// データベース接続の確立
	xdb, cleanup, err := connectDB(ctx)
	if err != nil {
		fmt.Printf("データベース接続エラー: %v\n", err)
		return
	}
	defer cleanup()

	e := echo.New()
	e.GET("/", handleHello)
	// e.POST("/users", func(c echo.Context) error {
	// 	return createUserHandler(c, xdb)
	// })

	// DIをする場合
	repo := NewRepo()
	createUserUsecase := NewCreateUserUseCase(xdb, repo)
	createUserHandler := NewCreateUserHandler(createUserUsecase)
	e.POST("/users", createUserHandler.ServeHTTP)

	e.Logger.Fatal(e.Start(":1323"))
}

func handleHello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

// ********************
// usecase層
// ********************
type CreateUserUsecaseImpl struct {
	DB   *sqlx.DB // interfaceに依存する
	Repo Repositorier
}

func NewCreateUserUseCase(db *sqlx.DB, repo Repositorier) CreateUserUseCaseier {
	return &CreateUserUsecaseImpl{
		DB:   db,
		Repo: repo,
	}
}

func (cu *CreateUserUsecaseImpl) CreateUser(ctx echo.Context) error {
	_, err := cu.Repo.insertUsers(cu.DB, User{
		id:       uuid.New(),
		name:     "うんち💩",
		email:    "hoge@tintin.com",
		password: "securepassword",
	})
	if err != nil {
		return err
	}
	fmt.Println("ユーザーが作成されました")
	return nil
}

// *******************
// handler層
// *******************
type CreateUserUseCaseier interface {
	CreateUser(ctx echo.Context) error
}
type CreateUser struct {
	UseCase CreateUserUseCaseier
}

func NewCreateUserHandler(s CreateUserUseCaseier) *CreateUser {
	return &CreateUser{UseCase: s}
}

func (cu *CreateUser) ServeHTTP(ctx echo.Context) error {
	if err := cu.UseCase.CreateUser(ctx); err != nil {
		return err
	}
	return nil
}

// ******************
// repository層
// ******************
type Repositorier interface {
	insertUsers(db *sqlx.DB, user User) (sql.Result, error)
}

type Repo struct{}

func NewRepo() Repositorier {
	return &Repo{}
}

func (r *Repo) insertUsers(db *sqlx.DB, user User) (sql.Result, error) {
	sql := `INSERT INTO users (id, username, email, password)
					VALUES
					(UUID_TO_BIN(?, 1), ?, ?, ?)
					ON DUPLICATE KEY UPDATE
					username = VALUES(username),
					email = VALUES(email),
					password = VALUES(password);`
	return db.Exec(sql, user.id, user.name, user.email, user.password)
}

func createUserHandler(c echo.Context, xdb *sqlx.DB) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	fmt.Printf("生成されたUUIDv7: %s\n", uuid.String())
	email := fmt.Sprintf("example+%v@example.com", rand.Intn(100))
	user := User{
		id:       uuid, // 例としてUUIDを使用
		name:     "うんち💩",
		email:    email,
		password: "securepassword",
	}

	result, err := insertUsers(xdb, user)
	fmt.Printf("%+v\n", result)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	return c.String(http.StatusOK, fmt.Sprintf("ユーザーが作成されました。作成した数: %d", rowsAffected))
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
