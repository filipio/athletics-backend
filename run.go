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
	m "github.com/filipio/athletics-backend/middlewares"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

const shutdownTimeout = 10 * time.Second

func addRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })

	mux.Handle("GET /api/v1/pokemons", m.UserOnly(controllers.GetAll[models.Pokemon](db), db))
	mux.Handle("GET /api/v1/pokemons/{id}", controllers.Get[models.Pokemon](db))
	mux.Handle("POST /api/v1/pokemons", m.AdminOnly(controllers.Create[models.Pokemon](db), db))
	mux.Handle("PUT /api/v1/pokemons/{id}", controllers.Update[models.Pokemon](db))
	mux.Handle("DELETE /api/v1/pokemons/{id}", controllers.Delete[models.Pokemon](db))

	mux.Handle("POST /api/v1/register", controllers.Register(db))
	mux.Handle("POST /api/v1/login", controllers.Login(db))
}

func executeMigrations(db *gorm.DB) {
	db.AutoMigrate(&models.Pokemon{}, &models.User{}, &models.Role{})
}

func seed(db *gorm.DB) {
	if db.Where("name = ?", utils.UserRole).First(&models.Role{}).RowsAffected > 0 {
		return
	}

	db.Create(&models.Role{Name: utils.UserRole})

	if db.Where("name = ?", utils.AdminRole).First(&models.Role{}).RowsAffected > 0 {
		return
	}
	db.Create(&models.Role{Name: utils.AdminRole})
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

	executeMigrations(db)
	log.Print("migrated database")

	seed(db)
	log.Print("seeded database")

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

	// belowe is used for graceful shutdown
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		<-ctx.Done()
		// make a new context for the Shutdown
		shutdownCtx := context.Background()
		shutdownCtx, cancel := context.WithTimeout(ctx, shutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
		}
	}()

	waitGroup.Wait()
	return nil
}

func newServerHandler(db *gorm.DB) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, db)

	var handler http.Handler = mux
	handler = m.LoggingMiddleware(handler)

	return handler
}