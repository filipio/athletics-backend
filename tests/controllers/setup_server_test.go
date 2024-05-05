package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	athleticsbackend "github.com/filipio/athletics-backend"
	"github.com/filipio/athletics-backend/config"
	"gorm.io/gorm"
)

// following values should not be changed outside of this file - they are read-only
var host string
var dbInstance *gorm.DB
var ctx = context.Background()

const envPath = "../../.env.test"
const hostFormula = "http://localhost:%s"

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := athleticsbackend.Run(ctx, envPath); err != nil {
			log.Fatalf("failed to start server: %v\n", err)
		}
	}()

	if err := waitForReady(ctx); err != nil {
		log.Fatalf("failed to wait for server to be ready: %v\n", err)
	}

	dbInstance = config.DatabaseConnection()
	code := m.Run()
	os.Exit(code)
}

func waitForReady(ctx context.Context) error {
	startTime := time.Now()
	for {
		host = fmt.Sprintf(hostFormula, os.Getenv("PORT"))
		endpoint := fmt.Sprintf("%s/readyz", host)
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error making request: %s\n", err.Error())
			continue
		}
		if resp.StatusCode == http.StatusOK {
			log.Println("Endpoint is ready!")
			resp.Body.Close()
			return nil
		}
		resp.Body.Close()

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if time.Since(startTime) >= 5*time.Second {
				return fmt.Errorf("timeout reached while waiting for endpoint")
			}
			// wait a little while between checks
			time.Sleep(250 * time.Millisecond)
		}
	}
}
