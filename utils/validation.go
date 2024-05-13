package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type ValidationResponseItem struct {
	ErrorsResponse
	Path []string `json:"path"`
}

func ValidationErrorDetails(errors validator.ValidationErrors) []ValidationResponseItem {
	details := make([]ValidationResponseItem, len(errors))

	for i, err := range errors {
		errorPath := buildErrorPath(err.Namespace())
		errorMessage := buildErrorMessage(err)
		errorType := err.ActualTag()

		details[i] = ValidationResponseItem{
			ErrorsResponse: ErrorsResponse{
				ErrorType: errorType,
				Details:   errorMessage,
			},
			Path: errorPath,
		}
	}

	return details
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

// defines the mapping between the error tag (used in models 'validate' tag)
func buildErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "field is required"
	case "email":
		return "must be an email"
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", err.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", err.Param())
	case "len":
		return fmt.Sprintf("must have exactly %s characters", err.Param())
	case "max":
		return fmt.Sprintf("must be less than %s", err.Param())
	case "min":
		return fmt.Sprintf("must be greater than %s", err.Param())
	default:
		return "is invalid"
	}
}
