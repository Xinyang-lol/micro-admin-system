package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func OpenMySQL(dsn string) (*sql.DB, error) {
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	database.SetMaxOpenConns(40)
	database.SetMaxIdleConns(10)
	database.SetConnMaxLifetime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	var lastErr error
	for i := 0; i < 15; i++ {
		if err := database.PingContext(ctx); err == nil {
			return database, nil
		} else {
			lastErr = err
			time.Sleep(time.Second)
		}
	}
	return nil, fmt.Errorf("connect mysql failed: %w", lastErr)
}
