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

type ValidationResponse struct {
	Errors map[string]interface{} `json:"errors"`
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

func Decode[T any](r *http.Request) (T, *ValidationResponse, error) {
	var v T
	// if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
	// 	return v, nil, fmt.Errorf("decode json: %w", err)
	// }

	bytes, err := io.ReadAll(r.Body)
	if err != nil {
		return v, nil, fmt.Errorf("read bytes body: %w", err)
	}

	var toAddBytes []byte

	if structId := r.PathValue("id"); structId != "" {
		toAddBytes = []byte(fmt.Sprintf(`,"id":%s}`, structId))
	}

	var finalBytes []byte

	if len(toAddBytes) > 0 {
		bytes := bytes[:len(bytes)-1]
		finalBytes = append(bytes, toAddBytes...)
	} else {
		finalBytes = bytes
	}

	if err := json.Unmarshal(finalBytes, &v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}

	// var objmap1 map[string]*json.RawMessage

	// if err := json.NewDecoder(r.Body).Decode(&objmap1); err != nil {
	// 	return v, nil, fmt.Errorf("decode json: %w", err)
	// }

	// var objmap2 map[string]*json.RawMessage = make(map[string]*json.RawMessage)

	// if structId := r.PathValue("id"); structId != "" {
	// 	idBytes := []byte(structId)

	// 	objmap2["id"] = (*json.RawMessage)(&idBytes)
	// }

	// for key, value := range objmap2 {
	// 	objmap1[key] = value
	// }

	// bytes, err := json.Marshal(objmap1)
	// if err != nil {
	// 	return v, nil, fmt.Errorf("marshal json: %w", err)
	// }

	// if err := json.Unmarshal(bytes, &v); err != nil {
	// 	return v, nil, fmt.Errorf("decode json: %w", err)
	// }

	fmt.Println("v", v)

	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.Struct(v); err != nil {
		response := ValidationResponse{}
		response.Errors = make(map[string]interface{})

		for _, err := range err.(validator.ValidationErrors) {
			namespace := err.Namespace()
			namespaceParts := strings.Split(namespace, ".")
			nestedElements := namespaceParts[1:]
			snakeCased := make([]string, len(nestedElements))
			for i, element := range nestedElements {
				snakeCased[i] = strcase.ToSnake(element)
			}

			var errorMessage string

			switch err.ActualTag() {
			case "required":
				errorMessage = "must be present"
			case "email":
				errorMessage = "must be a valid email address"
			case "oneof":
				errorMessage = "must be one of the following values: " + err.Param()
			case "lte":
				errorMessage = "must be less than or equal to " + err.Param()
			case "gte":
				errorMessage = "must be greater than or equal to " + err.Param()
			default:
				errorMessage = "unknown error"
			}

			errorType := err.ActualTag()

			currentErrorMap := response.Errors
			for i, element := range snakeCased {
				if i == len(snakeCased)-1 {
					currentErrorMap[element] = ValidationResponseItem{
						Type:    errorType,
						Message: errorMessage,
					}

				} else {
					if _, ok := currentErrorMap[element]; !ok {
						currentErrorMap[element] = make(map[string]interface{})
					}
					currentErrorMap = currentErrorMap[element].(map[string]interface{})
				}
			}
			log.Print("validation error: ", err.Error())
		}

		return v, &response, nil
	}

	return v, nil, nil
}

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}
