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

func TestMain(m *testing.M) {
	log.Print("test main is running")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := athleticsbackend.Run(ctx); err != nil {
			fmt.Printf("failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	if err := waitForReady(ctx, "http://localhost:8080/readyz"); err != nil {
		fmt.Printf("failed to wait for server to be ready: %v\n", err)
		os.Exit(1)
	}

	// Run the tests
	code := m.Run()

	// Exit with the code returned from m.Run
	os.Exit(code)
}

func waitForReady(
	ctx context.Context,
	endpoint string,
) error {
	client := http.Client{}
	startTime := time.Now()
	for {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			endpoint,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error making request: %s\n", err.Error())
			continue
		}
		if resp.StatusCode == http.StatusOK {
			fmt.Println("Endpoint is ready!")
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
