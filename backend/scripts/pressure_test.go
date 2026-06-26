package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"
)

type loginResp struct {
	Code int `json:"code"`
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

func TestPressure1000Users(t *testing.T) {
	baseURL := testEnv("API_BASE", "http://127.0.0.1:8080")
	total := testEnvInt("PRESSURE_TOTAL", 1000)
	concurrency := testEnvInt("PRESSURE_CONCURRENCY", 100)

	start := time.Now()
	jobs := make(chan int, total)
	results := make(chan error, total)
	var wg sync.WaitGroup

	client := &http.Client{Timeout: 3 * time.Second}
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range jobs {
				results <- doOneRequest(client, baseURL, id)
			}
		}()
	}

	go func() {
		for i := 0; i < total; i++ {
			jobs <- i
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
			t.Log(err)
			continue
		}
		success++
	}
	elapsed := time.Since(start)
	t.Logf("pressure finished total=%d success=%d failed=%d elapsed=%s avg=%s", total, success, failed, elapsed, elapsed/time.Duration(total))
	if failed > 0 {
		t.Fatalf("pressure test has %d failed requests", failed)
	}
	if elapsed > 3*time.Second {
		t.Logf("warning: total elapsed > 3s, please check machine resources and service logs")
	}
}

func doOneRequest(client *http.Client, baseURL string, id int) error {
	body := []byte(`{"username":"admin","password":"admin123"}`)
	resp, err := client.Post(baseURL+"/api/auth/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("login %d http=%d body=%s", id, resp.StatusCode, string(data))
	}
	var lr loginResp
	if err := json.Unmarshal(data, &lr); err != nil {
		return err
	}
	if lr.Code != 0 || lr.Data.Token == "" {
		return fmt.Errorf("login %d invalid response: %s", id, string(data))
	}
	req, err := http.NewRequest(http.MethodGet, baseURL+"/api/auth/profile", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+lr.Data.Token)
	profileResp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer profileResp.Body.Close()
	if profileResp.StatusCode != 200 {
		data, _ := io.ReadAll(profileResp.Body)
		return fmt.Errorf("profile %d http=%d body=%s", id, profileResp.StatusCode, string(data))
	}
	return nil
}

func testEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func testEnvInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}
