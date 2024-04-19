package controllers

import (
	"fmt"
	"net/http"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	fmt.Println("TestHealthcheck")
	// execute real http request to the server
	// and check the response

	code, err := http.Get("http://localhost:8080/healthz")

	if err != nil {
		t.Errorf("Error making request: %s", err.Error())
	}

	if code.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", code.StatusCode)
	}
}
