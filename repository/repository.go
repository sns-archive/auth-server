package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// connectDB はDBに接続する
func ConnectDB(ctx context.Context) (*sqlx.DB, func(), error) {
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
