package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	athleticsbackend "github.com/filipio/athletics-backend"
	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

// following values should not be changed outside of this file - they are read-only
var host string
var dbInstance *gorm.DB
var ctx = context.Background()
var adminToken string

const envPath = "../../.env.test"
const hostFormula = "http://localhost:%s"

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := athleticsbackend.Run(ctx, envPath); err != nil {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	if err := waitForReady(ctx); err != nil {
		slog.Error("failed to wait for server to be ready", "error", err)
		os.Exit(1)
	}

	dbInstance = config.DatabaseConnection()

	payload := utils.AnyMap{"email": os.Getenv("ADMIN_EMAIL"), "password": os.Getenv("ADMIN_PASSWORD")}
	_, tokenResponse, err := Post[map[string]string]("/api/v1/login", payload)
	if err != nil {
		panic(err)
	}

	adminToken = (*tokenResponse)["token"]

	code := m.Run()
	os.Exit(code)
}

func waitForReady(ctx context.Context) error {
	startTime := time.Now()
	for {
		host = fmt.Sprintf(hostFormula, os.Getenv("PORT"))
		endpoint := fmt.Sprintf("%s/api/readyz", host)
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
			slog.Info("Error making request", "error", err)
			continue
		}
		if resp.StatusCode == http.StatusOK {
			slog.Info("Endpoint is ready!")
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
