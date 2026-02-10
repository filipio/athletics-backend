package athleticsbackend

import (
	"net/http"

	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/controllers"
	m "github.com/filipio/athletics-backend/middlewares"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
)

// Contains all the routes of the application. This is the only place where routes are defined.
func addRoutes(mux *http.ServeMux, deps *config.Dependencies) {
	db := deps.DB
	auth := m.NewAuthMiddleware(deps)

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

	mux.Handle("GET /api/v1/pokemons", m.ErrorsMiddleware(auth.UserOnly(controllers.GetAll[models.Pokemon](deps))))
	mux.Handle("GET /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Get[models.Pokemon](deps)))
	mux.Handle("POST /api/v1/pokemons", m.ErrorsMiddleware(auth.UserOnly(controllers.Create[models.Pokemon](deps))))
	mux.Handle("PUT /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Update[models.Pokemon](deps)))
	mux.Handle("DELETE /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Delete[models.Pokemon](deps)))

	mux.Handle("POST /api/v1/auth/register/request-verification", m.ErrorsMiddleware(controllers.RequestVerification(deps)))
	mux.Handle("POST /api/v1/auth/verify-email", m.ErrorsMiddleware(controllers.VerifyEmail(deps)))
	mux.Handle("POST /api/v1/login", m.ErrorsMiddleware(controllers.Login(deps)))
	mux.Handle("POST /api/v1/auth/refresh", m.ErrorsMiddleware(controllers.RefreshToken(deps)))
	mux.Handle("POST /api/v1/auth/logout", m.ErrorsMiddleware(auth.UserOnly(controllers.Logout(deps))))

	mux.Handle("GET /api/v1/users", m.ErrorsMiddleware(auth.AdminOnly(controllers.GetAll[models.User](deps))))
	mux.Handle("GET /api/v1/users/{id}", m.ErrorsMiddleware(auth.AdminOnly(controllers.Get[models.User](deps))))
	mux.Handle("POST /api/v1/users", m.ErrorsMiddleware(auth.AdminOnly(controllers.CreateUser(deps))))
	mux.Handle("PUT /api/v1/users/{id}", m.ErrorsMiddleware(auth.AdminOnly(controllers.Update[models.User](deps))))
	mux.Handle("DELETE /api/v1/users/{id}", m.ErrorsMiddleware(auth.AdminOnly(controllers.Delete[models.User](deps))))

	mux.Handle("GET /api/v1/athletes", m.ErrorsMiddleware(auth.UserOnly(controllers.GetAll[models.Athlete](deps))))
	mux.Handle("GET /api/v1/athletes/{id}", m.ErrorsMiddleware(auth.UserOnly(controllers.Get[models.Athlete](deps))))

	mux.Handle("GET /api/v1/disciplines", m.ErrorsMiddleware(auth.UserOnly(controllers.GetAll[models.Discipline](deps))))
	mux.Handle("GET /api/v1/disciplines/{id}", m.ErrorsMiddleware(auth.UserOnly(controllers.Get[models.Discipline](deps))))

	mux.Handle("GET /api/v1/events", m.ErrorsMiddleware(auth.UserOnly(controllers.GetAll[models.Event](deps))))
	mux.Handle("GET /api/v1/events/{id}", m.ErrorsMiddleware(auth.UserOnly(controllers.Get[models.Event](deps))))
	mux.Handle("POST /api/v1/events", m.ErrorsMiddleware(auth.OrganizerOnly(controllers.Create[models.Event](deps))))
	mux.Handle("PUT /api/v1/events/{id}", m.ErrorsMiddleware(auth.OrganizerOnly(controllers.Update[models.Event](deps))))
	mux.Handle("DELETE /api/v1/events/{id}", m.ErrorsMiddleware(auth.OrganizerOnly(controllers.Delete[models.Event](deps))))
	mux.Handle("POST /api/v1/events/{id}/publish", m.ErrorsMiddleware(auth.OrganizerOnly(controllers.PublishEvent(deps))))
	mux.Handle("POST /api/v1/events/{id}/unpublish", m.ErrorsMiddleware(auth.OrganizerOnly(controllers.UnpublishEvent(deps))))

	mux.Handle("GET /api/v1/questions", m.ErrorsMiddleware(auth.UserOnly(controllers.GetAll[models.Question](deps))))
	mux.Handle("GET /api/v1/questions/{id}", m.ErrorsMiddleware(auth.UserOnly(controllers.Get[models.Question](deps))))
	mux.Handle("POST /api/v1/questions", m.ErrorsMiddleware(auth.OrganizerOnly(controllers.CreateQuestion(deps))))
	mux.Handle("PUT /api/v1/questions/{id}", m.ErrorsMiddleware(auth.OrganizerOnly(controllers.UpdateQuestion(deps))))
	mux.Handle("DELETE /api/v1/questions/{id}", m.ErrorsMiddleware(auth.OrganizerOnly(controllers.Delete[models.Question](deps))))

	mux.Handle("GET /api/v1/users/me/answers", m.ErrorsMiddleware(auth.UserOnly(controllers.GetAll[models.Answer](deps))))
	mux.Handle("GET /api/v1/users/me/answers/{id}", m.ErrorsMiddleware(auth.UserOnly(controllers.Get[models.Answer](deps))))
	mux.Handle("POST /api/v1/users/me/answers", m.ErrorsMiddleware(auth.UserOnly(controllers.CreateAnswer(deps))))
	mux.Handle("PUT /api/v1/users/me/answers/{id}", m.ErrorsMiddleware(auth.UserOnly(controllers.UpdateAnswer(deps))))
	mux.Handle("DELETE /api/v1/users/me/answers/{id}", m.ErrorsMiddleware(auth.UserOnly(controllers.Delete[models.Answer](deps))))

	mux.Handle("GET /api/v1/users/me", m.ErrorsMiddleware(auth.UserOnly(controllers.Get[models.User](deps))))
	mux.Handle("GET /api/v1/users/me/ranking", m.ErrorsMiddleware(auth.UserOnly(controllers.GetMyRanking(deps))))

	mux.Handle("GET /api/v1/answers", m.ErrorsMiddleware(auth.UserOnly(controllers.GetAll[models.Answer](deps))))
	mux.Handle("GET /api/v1/answers/{id}", m.ErrorsMiddleware(auth.UserOnly(controllers.Get[models.Answer](deps))))

	mux.Handle("GET /api/v1/ranking", m.ErrorsMiddleware(auth.UserOnly(controllers.GetRanking(deps))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := utils.ErrorsResponse{
			ErrorType: "not_found_error",
			Details:   "path not found",
		}
		utils.Encode(w, r, http.StatusNotFound, response)
	})
}
