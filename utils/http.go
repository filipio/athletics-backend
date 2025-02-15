package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type HandlerWithError func(http.ResponseWriter, *http.Request) error

func (f HandlerWithError) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	return f(w, r)
}

type ErrorsResponse struct {
	ErrorType string `json:"error_type"`
	Details   any    `json:"details"` // should be of type: string, utils.AnyMap
}

type PaginationParams struct {
	PageNo         int
	PerPage        int
	OrderBy        string
	OrderDirection string
}

func IntPathValue(r *http.Request, key string) int {
	value := r.PathValue(key)
	if value == "" {
		return 0
	}
	result, _ := strconv.Atoi(value)
	return result
}

func IntQueryValue(r *http.Request, key string) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0
	}
	result, _ := strconv.Atoi(value)
	return result
}

func BuildPaginationParams(r *http.Request) *PaginationParams {
	pageNo := IntQueryValue(r, "page_no")
	if pageNo == 0 {
		pageNo = DefaultPageNumber
	}

	perPage := IntQueryValue(r, "per_page")
	if perPage == 0 {
		perPage = DefaultPageSize
	}

	orderBy := r.URL.Query().Get("order_by")
	if orderBy == "" {
		orderBy = DefaultOrderBy
	}

	orderDirection := r.URL.Query().Get("order_dir")
	if orderDirection == "" {
		orderDirection = "asc"
	}

	return &PaginationParams{
		PageNo:         pageNo,
		PerPage:        perPage,
		OrderBy:        orderBy,
		OrderDirection: orderDirection,
	}
}

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(&v); err != nil {
		return err
	}
	return nil
}

func DecodeAndValidate[T Validatable](r *http.Request) (T, error) {
	var record T

	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		if err, ok := err.(*json.UnmarshalTypeError); ok {
			return record, AppValidationError{
				FieldPath: err.Field,
				AppError: AppError{
					Message: fmt.Sprintf("must be of type %s", err.Type.String()),
				},
			}
		}

		return record, fmt.Errorf("decode json: %w", err)
	}

	if err := validate.Struct(&record); err != nil {
		return record, err
	}

	if err := record.Validate(r); err != nil {
		return record, err
	}

	return record, nil
}

// validates the instance of a struct using defined tags on its fields
func Validate(instance any) error {
	if err := validate.Struct(instance); err != nil {
		return err
	}

	return nil
}
