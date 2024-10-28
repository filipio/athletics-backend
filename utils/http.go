package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

type HandlerWithError func(http.ResponseWriter, *http.Request) error

func RegisterValidations(db *gorm.DB) {

	validate.RegisterValidation("id_of", func(fl validator.FieldLevel) bool {
		tableName := fl.Param() + "s"
		passedId := fl.Field().Uint()
		sqlQuery := fmt.Sprintf("SELECT id FROM %s WHERE id = ?", tableName)

		result := db.Exec(sqlQuery, passedId)
		if result.Error != nil {
			fmt.Println("Error:", result.Error)
			return false
		}
		if result.RowsAffected != 1 {
			fmt.Println("wrong number of rows affected:", result.RowsAffected)
			return false
		}

		return true
	})
}

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

func IntQueryValue(r *http.Request, key string) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return 0
	}
	result, _ := strconv.Atoi(value)
	return result
}

func PaginationParams(r *http.Request) (pageNo int, perPage int, orderBy string, orderDirection string) {
	pageNo = IntQueryValue(r, "page_no")
	if pageNo == 0 {
		pageNo = DefaultPageNumber
	}

	perPage = IntQueryValue(r, "per_page")
	if perPage == 0 {
		perPage = DefaultPageSize
	}

	orderBy = r.URL.Query().Get("order_by")
	if orderBy == "" {
		orderBy = DefaultOrderBy
	}

	orderDirection = r.URL.Query().Get("order_dir")
	if orderDirection == "" {
		orderDirection = "asc"
	}

	return
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
