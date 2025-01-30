package athleticsbackend

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/filipio/athletics-backend/config"
	m "github.com/filipio/athletics-backend/middlewares"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"github.com/filipio/athletics-backend/workers"
	"github.com/joho/godotenv"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

const shutdownTimeout = 10 * time.Second

func seed(db *gorm.DB) {
	adminRole := models.Role{Name: utils.AdminRole}
	userRole := models.Role{Name: utils.UserRole}
	organizerRole := models.Role{Name: utils.OrganizerRole}

	db.FirstOrCreate(&adminRole, adminRole)
	db.FirstOrCreate(&userRole, userRole)
	db.FirstOrCreate(&organizerRole, organizerRole)

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	if adminEmail != "" && adminPassword != "" {
		adminUser := models.User{
			Email:    adminEmail,
			Password: adminPassword,
		}

		db.FirstOrCreate(&adminUser, models.User{Email: adminEmail})
		db.Model(&adminUser).Association("Roles").Append(&adminRole)
	}
}

// the only function to be used in order to add new workers
func appWorkers() *river.Workers {
	riverWorkers := river.NewWorkers()

	river.AddWorker(riverWorkers, &workers.SortWorker{})
	river.AddWorker(riverWorkers, &workers.PokemonWorker{})
	river.AddWorker(riverWorkers, &workers.PointsGranterWorker{})

	return riverWorkers
}

func Run(ctx context.Context, envPath string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	utils.SetupLogger()

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

	utils.RegisterValidations(db)

	workersClient := config.SetupWorkersClient(ctx, db, appWorkers())
	slog.Info("started workers client")

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	handler := newServerHandler(db, workersClient.InsertClient)
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

	// belowe is used for graceful shutdown
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

func newServerHandler(db *gorm.DB, insertClient *config.InsertWorkerClient) http.Handler {
	mux := http.NewServeMux()
	addRoutes(mux, db)

	var handler http.Handler = mux
	handler = m.DbMiddleware(handler, db)
	handler = m.WorkersMiddleware(handler, insertClient)
	handler = m.OnlyCurrentUserMiddleware(handler)
	handler = m.LoggingMiddleware(handler)

	return handler
}
