package controllers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/responses"
	"github.com/filipio/athletics-backend/utils"
	. "github.com/filipio/athletics-backend/utils"
	"github.com/google/go-cmp/cmp"
)

func TestGetPokemonByName(t *testing.T) {
	t.Run("GET pokemon by name", testCase(func(t *testing.T) {
		pokemons := []*models.Pokemon{
			{PokemonName: "Pikachu", Age: 20, Email: "pikachu@gmail.com", Attack: "Ember"},
			{PokemonName: "Bulbasaur", Age: 20, Email: "bulba@gmail.com", Attack: "Vine_Whip"},
			{PokemonName: "Pikachu", Age: 20, Email: "pika2@gmail.com", Attack: "Thunderbolt"},
		}
		dbInstance.Save(&pokemons)

		var expectedPokemons *[]models.Pokemon
		dbInstance.Where("pokemon_name = ?", "Pikachu").Find(&expectedPokemons)

		response, paginatedResponse, err := Get[responses.PaginatedResponse](fmt.Sprintf("/api/v1/pokemons?name=%s", "Pikachu"))
		fetchedPokemons := []models.Pokemon{}
		for _, pokemonMap := range paginatedResponse.Data.([]interface{}) {
			fetchedPokemons = append(fetchedPokemons, ToStruct[models.Pokemon](pokemonMap))
		}

		if err != nil {
			t.Errorf("Error executing request: %s", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", response.StatusCode)
		}

		if !cmp.Equal(*expectedPokemons, fetchedPokemons) {
			t.Error("Difference between fetched and actual pokemons : ", cmp.Diff(expectedPokemons, fetchedPokemons))
		}
	}))
}

func TestGetPokemon(t *testing.T) {

	t.Run("GET pokemon", testCase(func(t *testing.T) {
		pokemon := &models.Pokemon{
			PokemonName: "Pikachu",
			Age:         20,
			Email:       "pikachu@gmail.com",
			Attack:      "Ember",
		}

		dbInstance.Save(pokemon)
		dbInstance.Last(pokemon) // needed to format dates properly (issue with gorm)

		response, fetchedPokemon, err := Get[models.Pokemon](fmt.Sprintf("/api/v1/pokemons/%d", pokemon.ID))

		if err != nil {
			t.Errorf("Error executing request: %s", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", response.StatusCode)
		}

		if !cmp.Equal(pokemon, fetchedPokemon) {
			t.Error("Difference between fetched and actual pokemon : ", cmp.Diff(pokemon, fetchedPokemon))
		}
	}))
}

func TestPokemonValidation(t *testing.T) {
	testData := []struct {
		name               string
		body               any
		expectedStatusCode int
	}{
		{
			name:               "only pokemon name",
			body:               AnyMap{"pokemon_name": "Pikachu"},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "valid",
			body:               AnyMap{"pokemon_name": "Pikachu", "age": 20, "email": "abc@gmail.com", "attack": "Ember"},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, data := range testData {
		t.Run(data.name, testCase(func(t *testing.T) {
			t.Parallel()
			response, _, err := Post[models.Pokemon]("/api/v1/pokemons", data.body)

			if err != nil {
				t.Errorf("Error making request: %s", err.Error())
			}

			if response.StatusCode != data.expectedStatusCode {
				t.Errorf("Expected status code %d, got %d", data.expectedStatusCode, response.StatusCode)
			}
		}))
	}
}

func TestPostPokemon(t *testing.T) {
	body := AnyMap{"pokemon_name": "Pikachu", "age": 20, "email": "abc@gmail.com", "attack": "Ember"}

	t.Run("POST pokemon", testCase(func(t *testing.T) {

		response, createdPokemon, err := Post[models.Pokemon]("/api/v1/pokemons", body)

		if err != nil {
			t.Errorf("Error making request: %s", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, response.StatusCode)
		}

		pokemon := &models.Pokemon{}
		dbInstance.Last(pokemon)

		if !cmp.Equal(pokemon, createdPokemon) {
			t.Error("Difference between fetched and actual pokemon : ", cmp.Diff(pokemon, createdPokemon))
		}

	}))
}

func TestUpdatePokemon(t *testing.T) {
	t.Run("PUT pokemon", testCase(func(t *testing.T) {
		pokemon := &models.Pokemon{
			PokemonName: "Pikachu",
			Age:         20,
			Email:       "foo.bar@pikachu.com",
			Attack:      "Ember",
		}
		dbInstance.Save(pokemon)
		body := AnyMap{"pokemon_name": "Pikachu", "age": 30, "email": "foo@pikachu.com", "attack": "Thunderbolt"}
		response, updatedPokemon, err := Put[models.Pokemon](fmt.Sprintf("/api/v1/pokemons/%d", pokemon.ID), body)

		if err != nil {
			t.Errorf("Error executing request: %s", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", response.StatusCode)
		}

		expectedPokemon := &models.Pokemon{}
		dbInstance.Last(expectedPokemon)

		if !cmp.Equal(expectedPokemon, updatedPokemon) {
			t.Error("Difference between fetched and actual pokemon : ", cmp.Diff(expectedPokemon, updatedPokemon))
		}
	}))
}

func TestDeletePokemon(t *testing.T) {
	t.Run("DELETE pokemon", testCase(func(t *testing.T) {
		pokemon := &models.Pokemon{
			PokemonName: "Pikachu",
			Age:         20,
			Email:       "pikachu@gmail.com",
			Attack:      "Ember",
		}
		dbInstance.Save(pokemon)

		response, parsedResponse, err := Delete[utils.AnyMap](fmt.Sprintf("/api/v1/pokemons/%d", pokemon.ID))

		if err != nil {
			t.Errorf("Error executing request: %s", err.Error())
		}

		if response.StatusCode != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", response.StatusCode)
		}

		if len(*parsedResponse) != 0 {
			t.Error("Expected empty response")
		}
	}))

}
