package main

import (
	"context"
	"fmt"
	"os"

	athleticsbackend "github.com/filipio/athletics-backend"
)

func main() {
	ctx := context.Background()
	if err := athleticsbackend.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
