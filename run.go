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
	queries "github.com/filipio/athletics-backend/queries"
	"github.com/filipio/athletics-backend/responses"
	"github.com/filipio/athletics-backend/utils"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

const shutdownTimeout = 10 * time.Second

func addRoutes(mux *http.ServeMux, db *gorm.DB) {
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
	mux.HandleFunc("GET /api/readyz", func(w http.ResponseWriter, r *http.Request) {
		result := db.Exec("SELECT 1")
		if result.Error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(result.Error.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	mux.Handle("GET /api/v1/pokemons", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll(db, queries.GetPokemonsQuery, responses.BuildDefaultResponse[models.Pokemon]), db)))
	mux.Handle("GET /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Get(db, queries.GetByIdQuery, responses.BuildDefaultResponse[models.Pokemon])))
	mux.Handle("POST /api/v1/pokemons", m.ErrorsMiddleware(m.AdminOnly(controllers.Create[models.Pokemon](db, responses.BuildDefaultResponse), db)))
	mux.Handle("PUT /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Update[models.Pokemon](db, responses.BuildDefaultResponse)))
	mux.Handle("DELETE /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Delete[models.Pokemon](db)))

	mux.Handle("POST /api/v1/register", m.ErrorsMiddleware(controllers.Register(db)))
	mux.Handle("POST /api/v1/login", m.ErrorsMiddleware(controllers.Login(db)))

	mux.Handle("GET /api/v1/users", m.ErrorsMiddleware(m.AdminOnly(controllers.GetAll(db, queries.GetUsersQuery, responses.BuildUserResponse), db)))
	mux.Handle("GET /api/v1/users/{id}", m.ErrorsMiddleware(m.AdminOnly(controllers.Get(db, queries.GetUserQuery, responses.BuildUserResponse), db)))
	mux.Handle("POST /api/v1/users", m.ErrorsMiddleware(m.AdminOnly(controllers.Create(db, responses.BuildUserResponse), db)))
	mux.Handle("PUT /api/v1/users/{id}", m.ErrorsMiddleware(m.AdminOnly(controllers.Update(db, responses.BuildUserResponse), db)))
	mux.Handle("DELETE /api/v1/users/{id}", m.ErrorsMiddleware(m.AdminOnly(controllers.Delete[models.User](db), db)))

	mux.Handle("GET /api/v1/athletes", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll(db, queries.GetAthletesQuery, responses.BuildAthleteResponse), db)))
	mux.Handle("GET /api/v1/athletes/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get(db, queries.GetByIdQuery, responses.BuildAthleteResponse), db)))
	mux.Handle("POST /api/v1/athletes", m.ErrorsMiddleware(m.UserOnly(controllers.Create(db, responses.BuildAthleteResponse), db)))
	mux.Handle("PUT /api/v1/athletes/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Update(db, responses.BuildAthleteResponse), db)))
	mux.Handle("DELETE /api/v1/athletes/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Delete[models.Athlete](db), db)))

	mux.Handle("GET /api/v1/disciplines", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll(db, queries.DefaultQuery, responses.BuildDisciplineResponse), db)))
	mux.Handle("GET /api/v1/disciplines/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get(db, queries.GetByIdQuery, responses.BuildDisciplineResponse), db)))

	mux.Handle("GET /api/v1/events", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll(db, queries.GetEventsQuery, responses.BuildDefaultResponse[models.Event]), db)))
	mux.Handle("GET /api/v1/events/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get(db, queries.GetByIdQuery, responses.BuildDefaultResponse[models.Event]), db)))
	mux.Handle("POST /api/v1/events", m.ErrorsMiddleware(m.UserOnly(controllers.Create[models.Event](db, responses.BuildDefaultResponse), db)))
	mux.Handle("PUT /api/v1/events/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Update[models.Event](db, responses.BuildDefaultResponse), db)))
	mux.Handle("DELETE /api/v1/events/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Delete[models.Event](db), db)))

	mux.Handle("GET /api/v1/questions", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll(db, queries.GetQuestionsQuery, responses.BuildDefaultResponse[models.Question]), db)))
	mux.Handle("GET /api/v1/questions/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get(db, queries.GetByIdQuery, responses.BuildDefaultResponse[models.Question]), db)))
	mux.Handle("POST /api/v1/questions", m.ErrorsMiddleware(m.UserOnly(controllers.Create[models.Question](db, responses.BuildDefaultResponse), db)))
	mux.Handle("PUT /api/v1/questions/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Update[models.Question](db, responses.BuildDefaultResponse), db)))
	mux.Handle("DELETE /api/v1/questions/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Delete[models.Question](db), db)))

	// what is needed:
	// fetch events (all available)
}

func seed(db *gorm.DB) {
	adminRole := models.Role{Name: utils.AdminRole}
	userRole := models.Role{Name: utils.UserRole}

	db.FirstOrCreate(&adminRole, models.Role{Name: utils.AdminRole})
	db.FirstOrCreate(&userRole, models.Role{Name: utils.UserRole})

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
