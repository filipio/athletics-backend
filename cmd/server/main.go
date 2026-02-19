package main

import (
	"context"
	"fmt"
	"os"

	"github.com/filipio/athletics-backend/internal/app"
)

const envPath = ".env"

func main() {
	ctx := context.Background()
	if err := app.Run(ctx, envPath); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
