package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type job struct {
	index int
}

func main() {
	dsn := env("MYSQL_DSN", "root:123456@tcp(127.0.0.1:3306)/micro_admin?charset=utf8mb4&parseTime=true&loc=Local")
	total := envInt("DEVICE_SEED_TOTAL", 10000)
	workers := envInt("DEVICE_SEED_WORKERS", 16)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(workers + 4)

	typeIDs, err := loadTypeIDs(db)
	if err != nil {
		log.Fatal(err)
	}
	if len(typeIDs) == 0 {
		log.Fatal("device_type is empty, please run init.sql first")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	jobs := make(chan job, workers*2)
	results := make(chan error, total)
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for item := range jobs {
				results <- insertDevice(ctx, db, typeIDs, item.index, workerID)
			}
		}(i + 1)
	}

	go func() {
		for i := 1; i <= total; i++ {
			jobs <- job{index: i}
		}
		close(jobs)
		wg.Wait()
		close(results)
	}()

	success := 0
	failed := 0
	for err := range results {
		if err != nil {
			failed++
			log.Printf("insert failed: %v", err)
			continue
		}
		success++
	}
	log.Printf("seed finished, success=%d failed=%d", success, failed)
}

func loadTypeIDs(db *sql.DB) ([]int64, error) {
	rows, err := db.Query(`SELECT id FROM device_type ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func insertDevice(ctx context.Context, db *sql.DB, typeIDs []int64, index int, workerID int) error {
	statuses := []string{"online", "offline", "repair"}
	typeID := typeIDs[rand.Intn(len(typeIDs))]
	status := statuses[rand.Intn(len(statuses))]
	code := fmt.Sprintf("BATCH-%05d", index)
	_, err := db.ExecContext(ctx, `
		INSERT INTO device(name, code, type_id, status, location, remark)
		VALUES(?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE status = VALUES(status), updated_at = NOW()`,
		fmt.Sprintf("批量设备-%05d", index), code, typeID, status, fmt.Sprintf("实验楼-%02d", index%20+1), fmt.Sprintf("worker-%d", workerID))
	return err
}

func env(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}
