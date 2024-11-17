package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/sns-archive/jwt-auth-server/repository"
)

func main() {
	ctx := context.Background()
	// データベース接続の確立
	xdb, cleanup, err := repository.ConnectDB(ctx)
	if err != nil {
		fmt.Printf("データベース接続エラー: %v\n", err)
		return
	}
	defer cleanup()

	e := echo.New()
	e.GET("/", handleHello)
	e.POST("/users", func(c echo.Context) error {
		return createUserHandler(c, xdb)
	})
	e.Logger.Fatal(e.Start(":1323"))
}

func handleHello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func createUserHandler(c echo.Context, xdb *sqlx.DB) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	fmt.Printf("生成されたUUIDv7: %s\n", uuid.String())
	email := fmt.Sprintf("example+%v@example.com", rand.Intn(100))
	user := repository.User{
		Id:       uuid, // 例としてUUIDを使用
		Name:     "うんち💩",
		Email:    email,
		Password: "securepassword",
	}

	result, err := repository.NewRepo().InsertUsers(xdb, user)
	fmt.Printf("%+v\n", result)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	return c.String(http.StatusOK, fmt.Sprintf("ユーザーが作成されました。作成した数: %d", rowsAffected))
}
