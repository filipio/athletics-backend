package athleticsbackend

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"time"

	"github.com/filipio/athletics-backend/controllers"
	"github.com/filipio/athletics-backend/controllers/pokemons"
)

func addRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
	mux.Handle("GET /hello", controllers.HandleSomething())
	mux.Handle("PUT /hello/{id}", controllers.HandleSomePut())
	mux.Handle("PUT /hello/{id}/something/{otherId}", controllers.HandleDeepPut())
	mux.Handle("GET /pokemon", pokemons.HandlePokemon())
	mux.HandleFunc("POST /pokemon", pokemons.CreatePokemon)
	mux.HandleFunc("POST /persons", pokemons.Create[pokemons.Person])
	mux.Handle("POST /personsAdvaced", pokemons.CreateAdvaced[pokemons.Person](nil))
	// add routes for health checks
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
}

func Run(ctx context.Context) error {

	if os.Getenv("APP_ENV") != "prod" {
		log.Print("loading .env file in local environment")
		// get absolute location of current file
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			log.Fatal(err)
		}

		log.Printf(dir)

		// load .env file
		// TODO: this fails
		// err = godotenv.Load(filepath.Join(dir, ".env"))
		// if err != nil {
		// 	log.Fatal("Error loading .env file in local environment")
		// 	return err
		// }

		// err := godotenv.Load()
		// if err != nil {
		// 	log.Fatal("Error loading .env file in local environment")
		// }
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	srv := newServer()
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: srv,
	}

	go func() {
		log.Printf("listening on %s\n", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "error listening and serving: %s\n", err)
		}
	}()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		// make a new context for the Shutdown (thanks Alessandro Rosetti)
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}

func newServer() http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux)
	return mux
}
