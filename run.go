package athleticsbackend

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/controllers"
	"github.com/filipio/athletics-backend/models"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

const shutdownTimeout = 10 * time.Second

func addRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	mux.Handle("GET /api/v1/pokemons", controllers.GetAll[models.Pokemon](db))
	mux.Handle("GET /api/v1/pokemons/{id}", controllers.Get[models.Pokemon](db))
	mux.Handle("POST /api/v1/pokemons", controllers.CreateOrUpdate[models.Pokemon](db))
	mux.Handle("PUT /api/v1/pokemons/{id}", controllers.CreateOrUpdate[models.Pokemon](db))
	mux.Handle("DELETE /api/v1/pokemons/{id}", controllers.Delete[models.Pokemon](db))

	mux.Handle("GET /api/v1/humans", controllers.GetAll[models.Human](db))
	mux.Handle("GET /api/v1/humans/{id}", controllers.Get[models.Human](db))
	mux.Handle("POST /api/v1/humans", controllers.CreateOrUpdate[models.Human](db))
	mux.Handle("PUT /api/v1/humans/{id}", controllers.CreateOrUpdate[models.Human](db))
	mux.Handle("DELETE /api/v1/humans/{id}", controllers.Delete[models.Human](db))

}

func Run(ctx context.Context, envPath string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	if os.Getenv("APP_ENV") != "prod" {
		log.Print("loading .env file in local environment")
		err := godotenv.Load(envPath)
		if err != nil {
			log.Fatal("error loading .env file in local environment")
			return err
		}
	}

	db := config.DatabaseConnection()
	log.Print("established connection to database")
	db.AutoMigrate(&models.Pokemon{}, &models.Other{}, &models.Human{})
	log.Print("migrated database")

	pokemon := models.Pokemon{
		PokemonName: "Pikachu",
		Age:         5,
		Email:       "abc@gmail.com",
		Attack:      "Ember",
	}

	log.Println(pokemon)

	db.Save(&pokemon)
	db.Delete(&pokemon)

	pokemon.Age = 6
	pokemon.Email = "qwe@gmail.com"
	db.Save(&pokemon)
	db.Save(&pokemon)
	db.Save(&pokemon)

	// print id of the created pokemon
	log.Printf("created pokemon with id: %d\n", pokemon.ID)

	// var fetched models.Pokemon
	// db.First(&fetched, 1)
	// log.Printf("fetched pokemon: %+v\n", fetched)

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	handler := newServerHandler(db)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
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
		// make a new context for the Shutdown
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	wg.Wait()
	return nil
}

func newServerHandler(db *gorm.DB) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, db)

	var handler http.Handler = mux
	handler = loggingMiddleware(handler)

	return handler
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(startTime)
		log.Printf("%s %s %s\n", r.Method, r.RequestURI, duration)
	})
}
