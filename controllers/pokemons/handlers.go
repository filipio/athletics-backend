package pokemons

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Child struct {
	Height int `json:"height" validate:"required"`
}

type Pokemon struct {
	PokemonName *string `json:"pokemon_name" validate:"required"`
	Age         *int    `json:"age" validate:"required"`
	// PokemonChild Child  `json:"pokemon_child" validate:"required"`
}

func decode[T any](r *http.Request) (T, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, fmt.Errorf("decode json: %w", err)
	}
	return v, nil
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
	// decode pokemon from request body
	var validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

	pokemon, err := decode[Pokemon](r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validate pokemon

	if err := validate.Struct(pokemon); err != nil {
		for _, err := range err.(validator.ValidationErrors) {

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
	}

	// encode pokemon to response body
	if err := encode(w, r, http.StatusCreated, pokemon); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
