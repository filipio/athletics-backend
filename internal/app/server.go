package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/filipio/athletics-backend/pkg/config"
	"github.com/filipio/athletics-backend/internal/email"
	m "github.com/filipio/athletics-backend/internal/middleware"
	"github.com/filipio/athletics-backend/internal/models"
	"github.com/filipio/athletics-backend/pkg/httpio"
	"github.com/filipio/athletics-backend/internal/workers"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

const shutdownTimeout = 10 * time.Second

func seed(db *gorm.DB) {
	adminRole := models.Role{Name: httpio.AdminRole}
	userRole := models.Role{Name: httpio.UserRole}
	organizerRole := models.Role{Name: httpio.OrganizerRole}

	db.FirstOrCreate(&adminRole, adminRole)
	db.FirstOrCreate(&userRole, userRole)
	db.FirstOrCreate(&organizerRole, organizerRole)

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	adminUsername := os.Getenv("ADMIN_USERNAME")

	if adminEmail != "" && adminPassword != "" && adminUsername != "" {
		adminUser := models.User{
			Email:    adminEmail,
			Password: adminPassword,
			Username: adminUsername,
		}

		db.FirstOrCreate(&adminUser, models.User{Email: adminEmail})
		db.Model(&adminUser).Association("Roles").Append(&adminRole)
	}
}

// the only function to be used in order to add new workers
func appWorkers(deps *config.Dependencies) *river.Workers {
	riverWorkers := river.NewWorkers()

	river.AddWorker(riverWorkers, workers.NewSortWorker(deps))
	river.AddWorker(riverWorkers, workers.NewPokemonWorker(deps))
	river.AddWorker(riverWorkers, workers.NewPointsGranterWorker(deps))

	return riverWorkers
}

func Run(ctx context.Context, envPath string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	httpio.SetupLogger()

	if os.Getenv("APP_ENV") != "prod" {
		slog.Info("loading .env file in local environment")

		err := godotenv.Load(envPath)
		if err != nil {
			slog.Error("error loading .env file in local environment")
			return err
		}
	}

	db := config.DatabaseConnection()
	slog.Info("established connection to database")

	seed(db)
	slog.Info("seeded database")

	httpio.RegisterValidations(db)

	// Initialize email sender
	emailSender := email.GetDefaultEmailSender()

	// Create dependencies container (workers set to nil temporarily)
	deps := config.NewDependencies(db, nil, emailSender)

	// Create workers with dependencies
	workersClient := config.SetupWorkersClient(ctx, db, appWorkers(deps))
	slog.Info("started workers client")

	// Update Workers in deps
	deps.Workers = workersClient.InsertClient

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	handler := newServerHandler(deps)
	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		slog.Info("listening", "address", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("error listening and serving", "error", err)
			os.Exit(1)
		}
	}()

	// below is used for graceful shutdown
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		<-ctx.Done()
		// make a new context for the Shutdown
		shutdownCtx := context.Background()
		httpShutdownCtx, cancelHttp := context.WithTimeout(shutdownCtx, shutdownTimeout)
		workersShutdownCtx, cancelWorkers := context.WithTimeout(shutdownCtx, shutdownTimeout)

		defer cancelWorkers()
		defer cancelHttp()

		if err := workersClient.Shutdown(workersShutdownCtx); err != nil {
			slog.Error("error shutting down workers client", "error", err)
			os.Exit(1)
		}

		if err := httpServer.Shutdown(httpShutdownCtx); err != nil {
			slog.Error("error shutting down http server", "error", err)
			os.Exit(1)
		}

	}()

	waitGroup.Wait()
	return nil
}

func newServerHandler(deps *config.Dependencies) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, deps)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	var handler http.Handler = c.Handler(mux)
	handler = m.OnlyCurrentUserMiddleware(handler)
	handler = m.LoggingMiddleware(handler)

	return handler
}
