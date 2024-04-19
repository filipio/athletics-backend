package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/filipio/athletics-backend/controllers"
	"github.com/filipio/athletics-backend/controllers/pokemons"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") != "prod" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file in local environment")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	// mux.Handle("/favicon.ico", http.NotFoundHandler())
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
	mux.Handle("GET /hello", controllers.HandleSomething())
	mux.Handle("PUT /hello/{id}", controllers.HandleSomePut())
	mux.Handle("PUT /hello/{id}/something/{otherId}", controllers.HandleDeepPut())
	mux.Handle("GET /pokemon", pokemons.HandlePokemon())
	mux.HandleFunc("POST /pokemon", pokemons.CreatePokemon)

	addr := fmt.Sprintf(":%s", port)

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	fmt.Println("Server is running on port", port)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
