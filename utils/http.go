package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

type HandlerWithError func(http.ResponseWriter, *http.Request) error

func (f HandlerWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

type ErrorsResponse struct {
	ErrorType string `json:"error_type"`
	Details   any    `json:"details"` // should be of type: string, utils.AnyMap
}

func IntPathValue(r *http.Request, key string) int {
	value := r.PathValue(key)
	if value == "" {
		return 0
	}
	result, _ := strconv.Atoi(value)
	return result
}

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&v); err != nil {
		return err
	}
	return nil
}

func Decode[T any](r *http.Request) (T, error) {
	var record T

	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		return record, fmt.Errorf("decode json: %w", err)
	}

	if err := validate.Struct(record); err != nil {
		log.Println(err.Error())
		return record, err
	}

	return record, nil
}
