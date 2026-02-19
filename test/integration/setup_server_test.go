package controllers

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/filipio/athletics-backend/internal/app"
	"github.com/filipio/athletics-backend/pkg/config"
	"github.com/filipio/athletics-backend/pkg/httpio"
	"gorm.io/gorm"
)

// following values should not be changed outside of this file - they are read-only
var host string
var dbInstance *gorm.DB
var ctx = context.Background()
var adminToken string
var adminEmail string
var adminPassword string

const envPath = "../../.env.test"
const hostFormula = "http://localhost:%s"

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := app.Run(ctx, envPath); err != nil {
			slog.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	if err := waitForReady(ctx); err != nil {
		slog.Error("failed to wait for server to be ready", "error", err)
		os.Exit(1)
	}

	dbInstance = config.DatabaseConnection()

	// Store admin credentials from environment variables for use in tests
	adminEmail = os.Getenv("ADMIN_EMAIL")
	adminPassword = os.Getenv("ADMIN_PASSWORD")

	payload := httpio.AnyMap{"email": adminEmail, "password": adminPassword}
	_, tokenResponse, err := Post[map[string]interface{}]("/api/v1/login", payload)
	if err != nil {
		panic(err)
	}

	adminToken = (*tokenResponse)["access_token"].(string)

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

// getAdminCredentials returns the admin user's email and password from environment variables
// This is useful for tests that need to login as admin
func getAdminCredentials() (email, password string) {
	return adminEmail, adminPassword
}
