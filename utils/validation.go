package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func RegisterValidations(db *gorm.DB) {

	validate.RegisterValidation("id_of", func(fl validator.FieldLevel) bool {
		// TODO: extract question type and validate content json based on it
		// questionType := fl.Parent().FieldByName("Type")
		// fmt.Printf("questionType: %v\n", questionType)

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
	case "id_of":
		return fmt.Sprintf("must be an existing id (integer) of '%s' resource", err.Param())
	case "oneof":
		param := strings.ReplaceAll(err.Param(), " ", ",")
		return fmt.Sprintf("must be one of values:%s", param)
	default:
		return "is invalid"
	}
}
