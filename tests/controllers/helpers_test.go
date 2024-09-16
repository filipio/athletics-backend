package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/filipio/athletics-backend/models"
)

func executeHttp[T any](method string, path string, body any) (*http.Response, *T, error) {
	url := host + path
	var jsonPayload io.Reader = nil
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("cannot marshal into json: %w", err)
		}
		jsonPayload = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		url,
		jsonPayload,
	)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+adminToken)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to do request: %w", err)
	}

	var result T
	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response, &result, nil
}

func Get[T any](path string) (*http.Response, *T, error) {
	return executeHttp[T]("GET", path, nil)
}

func Post[T any](path string, body any) (*http.Response, *T, error) {
	return executeHttp[T]("POST", path, body)
}

func Put[T any](path string, body any) (*http.Response, *T, error) {
	return executeHttp[T]("PUT", path, body)
}

func Delete[T any](path string) (*http.Response, *T, error) {
	return executeHttp[T]("DELETE", path, nil)
}

func beforeEach() {
	dbInstance.Where("1 = 1").Delete(&models.Pokemon{})
}

func afterEach() {
}

func testCase(test func(t *testing.T)) func(*testing.T) {
	return func(t *testing.T) {
		beforeEach()
		defer afterEach()
		test(t)
	}
}
