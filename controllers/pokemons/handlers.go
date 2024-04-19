package pokemons

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/iancoleman/strcase"
)

type Address struct {
	Street string `json:"street" validate:"required"`
	City   string `json:"city" validate:"min=2,max=100"`
}

type Child struct {
	Height       *int     `json:"height" validate:"required"`
	Weight       *float32 `json:"weight" validate:"required"`
	ChildAddress *Address `json:"child_address" validate:"required"`
}

type Pokemon struct {
	PokemonName  *string `json:"pokemon_name" validate:"required,oneof=D Pikachu Bulbasaur Charmander Squirtle,min=2,max=100"`
	Age          *int    `json:"age" validate:"required,gte=0,lte=130"`
	Email        *string `json:"email" validate:"required,email"`
	PokemonChild *Child  `json:"pokemon_child" validate:"required"`
	// PokemonChild Child  `json:"pokemon_child" validate:"required"`
}

type Person struct {
	Name    *string  `json:"name" validate:"required"`
	Age     *int     `json:"age" validate:"required,gte=0,lte=130"`
	Address *Address `json:"address" validate:"required"`
}

type ValidationResponse struct {
	Errors map[string]interface{} `json:"errors"`
}

type ValidationResponseItem struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func decode[T any](r *http.Request) (T, *ValidationResponse, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}

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
			// }

			fmt.Println("printing errors")
			fmt.Println("namespace", err.Namespace())
			fmt.Println("field", err.Field())
			fmt.Println("struct namespace", err.StructNamespace())
			fmt.Println("struct field", err.StructField())
			fmt.Println("tag", err.Tag())
			fmt.Println("actual tag", err.ActualTag())
			fmt.Println("kind", err.Kind())
			fmt.Println("type", err.Type())
			fmt.Println("value", err.Value())
			fmt.Println("param", err.Param())
			fmt.Println("error", err.Error())
		}

		return v, &response, nil
	}

	return v, nil, nil
}

func encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

func HandlePokemon() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {

			pokemonAge := r.URL.Query()["age"]
			fmt.Fprintf(w, "Hello, pokemon with age %s!", pokemonAge)
		},
	)
}

func CreatePokemon(w http.ResponseWriter, r *http.Request) {
	pokemon, validationError, err := decode[Pokemon](r)

	if validationError != nil {
		if err := encode(w, r, http.StatusBadRequest, validationError); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// encode pokemon to response body
	if err := encode(w, r, http.StatusCreated, pokemon); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func Create[T any](w http.ResponseWriter, r *http.Request) {
	decodedStruct, validationError, err := decode[T](r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if validationError != nil {
		if err := encode(w, r, http.StatusBadRequest, validationError); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}

	// encode pokemon to response body
	if err := encode(w, r, http.StatusCreated, decodedStruct); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateAdvaced[T any](afterDecode func(T)) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			decodedStruct, validationError, err := decode[T](r)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if afterDecode != nil {
				afterDecode(decodedStruct)
			}

			if validationError != nil {
				if err := encode(w, r, http.StatusBadRequest, validationError); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			// encode pokemon to response body
			if err := encode(w, r, http.StatusCreated, decodedStruct); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}

func HandleSomePut() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.PathValue("id"))
			if err != nil {
				http.Error(w, "Invalid id", http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "Put request with id %d", id)
		},
	)
}
