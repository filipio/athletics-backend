package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type AnyMap map[string]any

type ErrorsResponse struct {
	Errors AnyMap `json:"errors"`
}

type ValidationResponseItem struct {
	Type    string `json:"type"`
	Message string `json:"message"`
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
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func Decode[T any](r *http.Request) (T, *ErrorsResponse, error) {
	var v T

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return v, nil, err
	}

	if err := json.Unmarshal(bytes, &v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}

	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(v); err != nil {
		return v, validationErrorResponse(err), nil
	}

	return v, nil, nil
}

func validationErrorResponse(err error) *ErrorsResponse {
	response := ErrorsResponse{}
	response.Errors = make(AnyMap)

	for _, err := range err.(validator.ValidationErrors) {
		log.Print("validation error: ", err.Error())

		errorPath := buildErrorPath(err.Namespace())
		errorMessage := buildErrorMessage(err)
		errorType := err.ActualTag()

		addErrorToResponse(&response, errorPath, errorType, errorMessage)
	}

	return &response
}

func buildErrorPath(errorNamespace string) []string {
	namespaceParts := strings.Split(errorNamespace, ".")
	pathElements := namespaceParts[1:]

	pathElementsSnakeCased := make([]string, len(pathElements))
	for i, element := range pathElements {
		pathElementsSnakeCased[i] = strcase.ToSnake(element)
	}

	return pathElementsSnakeCased
}

func addErrorToResponse(response *ErrorsResponse, errorPath []string, errorType string, errorMessage string) {
	currentErrorMap := response.Errors
	for i, element := range errorPath {
		if i == len(errorPath)-1 {
			currentErrorMap[element] = ValidationResponseItem{
				Type:    errorType,
				Message: errorMessage,
			}

		} else {
			if _, ok := currentErrorMap[element]; !ok {
				currentErrorMap[element] = make(AnyMap)
			}
			currentErrorMap = currentErrorMap[element].(AnyMap)
		}
	}
}

// defines the mapping between the error tag (used in models 'validate' tag)
func buildErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "This field must be an email"
	case "gte":
		return fmt.Sprintf("This field must be greater than or equal to %s", err.Param())
	case "lte":
		return fmt.Sprintf("This field must be less than or equal to %s", err.Param())
	case "len":
		return fmt.Sprintf("This field must have exactly %s characters", err.Param())
	case "max":
		return fmt.Sprintf("This field must be less than %s", err.Param())
	case "min":
		return fmt.Sprintf("This field must be greater than %s", err.Param())
	default:
		return "This field is invalid"
	}
}
