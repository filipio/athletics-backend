package athleticsbackend

import (
	"net/http"

	"github.com/filipio/athletics-backend/controllers"
	m "github.com/filipio/athletics-backend/middlewares"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

// Contains all the routes of the application. This is the only place where routes are defined.
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

	mux.Handle("GET /api/v1/pokemons", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll[models.Pokemon]())))
	mux.Handle("GET /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Get[models.Pokemon]()))
	mux.Handle("POST /api/v1/pokemons", m.ErrorsMiddleware(m.UserOnly(controllers.Create[models.Pokemon]())))
	mux.Handle("PUT /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Update[models.Pokemon]()))
	mux.Handle("DELETE /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Delete[models.Pokemon]()))

	mux.Handle("POST /api/v1/register", m.ErrorsMiddleware(controllers.Register()))
	mux.Handle("POST /api/v1/login", m.ErrorsMiddleware(controllers.Login()))

	mux.Handle("GET /api/v1/users", m.ErrorsMiddleware(m.AdminOnly(controllers.GetAll[models.User]())))
	mux.Handle("GET /api/v1/users/{id}", m.ErrorsMiddleware(m.AdminOnly(controllers.Get[models.User]())))
	mux.Handle("POST /api/v1/users", m.ErrorsMiddleware(m.AdminOnly(controllers.Create[models.User]())))
	mux.Handle("PUT /api/v1/users/{id}", m.ErrorsMiddleware(m.AdminOnly(controllers.Update[models.User]())))
	mux.Handle("DELETE /api/v1/users/{id}", m.ErrorsMiddleware(m.AdminOnly(controllers.Delete[models.User]())))

	mux.Handle("GET /api/v1/athletes", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll[models.Athlete]())))
	mux.Handle("GET /api/v1/athletes/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get[models.Athlete]())))

	mux.Handle("GET /api/v1/disciplines", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll[models.Discipline]())))
	mux.Handle("GET /api/v1/disciplines/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get[models.Discipline]())))

	mux.Handle("GET /api/v1/events", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll[models.Event]())))
	mux.Handle("GET /api/v1/events/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get[models.Event]())))
	mux.Handle("POST /api/v1/events", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Create[models.Event]())))
	mux.Handle("PUT /api/v1/events/{id}", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Update[models.Event]())))
	mux.Handle("DELETE /api/v1/events/{id}", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Delete[models.Event]())))

	mux.Handle("GET /api/v1/questions", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll[models.Question]())))
	mux.Handle("GET /api/v1/questions/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get[models.Question]())))
	mux.Handle("POST /api/v1/questions", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Create[models.Question]())))
	mux.Handle("PUT /api/v1/questions/{id}", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Update[models.Question]())))
	mux.Handle("DELETE /api/v1/questions/{id}", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Delete[models.Question]())))

	mux.Handle("GET /api/v1/users/me/answers", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll[models.Answer]())))
	mux.Handle("GET /api/v1/users/me/answers/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get[models.Answer]())))
	mux.Handle("POST /api/v1/users/me/answers", m.ErrorsMiddleware(m.UserOnly(controllers.Create[models.Answer]())))
	mux.Handle("PUT /api/v1/users/me/answers/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Update[models.Answer]())))
	mux.Handle("DELETE /api/v1/users/me/answers/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Delete[models.Answer]())))

	mux.Handle("GET /api/v1/answers", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll[models.Answer]())))
	mux.Handle("GET /api/v1/answers/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get[models.Answer]())))

	mux.Handle("GET /api/v1/ranking", m.ErrorsMiddleware(m.UserOnly(controllers.GetRanking())))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := utils.ErrorsResponse{
			ErrorType: "not_found_error",
			Details:   "path not found",
		}
		utils.Encode(w, r, http.StatusNotFound, response)
	})
}
