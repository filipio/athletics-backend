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
)

var host string // only for reading outside of this file
var ctx = context.Background()

const envPath = "../../.env"
const hostFormula = "http://localhost:%s"

func Get(path string) (*http.Response, error) {
	url := host + path
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		url,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	return http.DefaultClient.Do(req)
}

func TestMain(m *testing.M) {
	log.Print("test main is running")
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

	// Run the tests
	log.Print("Host is ", host)
	code := m.Run()

	// Exit with the code returned from m.Run
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
