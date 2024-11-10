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

	return riverWorkers
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

	seed(db)
	log.Print("seeded database")

	utils.RegisterValidations(db)

	workersClient := config.SetupWorkersClient(ctx, db, appWorkers())
	log.Print("started workers client")

	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	handler := newServerHandler(db, workersClient.InsertClient)
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
		httpShutdownCtx, cancelHttp := context.WithTimeout(shutdownCtx, shutdownTimeout)
		workersShutdownCtx, cancelWorkers := context.WithTimeout(shutdownCtx, shutdownTimeout)

		defer cancelWorkers()
		defer cancelHttp()

		if err := workersClient.Shutdown(workersShutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down workers client: %s\n", err)
		}

		if err := httpServer.Shutdown(httpShutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "error shutting down http server: %s\n", err)
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
	handler = m.LoggingMiddleware(handler)

	return handler
}
