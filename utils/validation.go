package utils

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type ValidationResponseItem struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func ValidationErrorMap(errors validator.ValidationErrors) AnyMap {
	errorsMap := make(AnyMap)

	for _, err := range errors {
		log.Print("validation error: ", err.Error())

		errorPath := buildErrorPath(err.Namespace())
		errorMessage := buildErrorMessage(err)
		errorType := err.ActualTag()

		addErrorToMap(errorsMap, errorPath, errorType, errorMessage)
	}

	return errorsMap
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

func addErrorToMap(errorsMap AnyMap, errorPath []string, errorType string, errorMessage string) {
	currentErrorMap := errorsMap
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
