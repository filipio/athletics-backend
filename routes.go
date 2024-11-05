package athleticsbackend

import (
	"net/http"

	"github.com/filipio/athletics-backend/controllers"
	m "github.com/filipio/athletics-backend/middlewares"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

func addRoutes(mux *http.ServeMux, db *gorm.DB) {
	// TODO: remove passing of db where possible, as it can be extracted from the request context

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

	mux.Handle("GET /api/v1/pokemons", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll(models.GetPokemonsQuery, models.BuildDefaultResponse[models.Pokemon]))))
	mux.Handle("GET /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Get(models.GetByIdQuery, models.BuildDefaultResponse[models.Pokemon])))
	mux.Handle("POST /api/v1/pokemons", m.ErrorsMiddleware(m.UserOnly(controllers.Create[models.Pokemon](models.BuildDefaultResponse))))
	mux.Handle("PUT /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Update[models.Pokemon](models.DefaultUpdateQuery, models.BuildDefaultResponse)))
	mux.Handle("DELETE /api/v1/pokemons/{id}", m.ErrorsMiddleware(controllers.Delete[models.Pokemon](models.GetByIdQuery)))

	mux.Handle("POST /api/v1/register", m.ErrorsMiddleware(controllers.Register(db)))
	mux.Handle("POST /api/v1/login", m.ErrorsMiddleware(controllers.Login(db)))

	mux.Handle("GET /api/v1/users", m.ErrorsMiddleware(m.AdminOnly(controllers.GetAll(models.GetUsersQuery, models.BuildUserResponse))))
	mux.Handle("GET /api/v1/users/{id}", m.ErrorsMiddleware(m.AdminOnly(controllers.Get(models.GetUserQuery, models.BuildUserResponse))))
	mux.Handle("POST /api/v1/users", m.ErrorsMiddleware(m.AdminOnly(controllers.Create(models.BuildUserResponse))))
	mux.Handle("PUT /api/v1/users/{id}", m.ErrorsMiddleware(m.AdminOnly(controllers.Update(models.DefaultUpdateQuery, models.BuildUserResponse))))
	mux.Handle("DELETE /api/v1/users/{id}", m.ErrorsMiddleware(m.AdminOnly(controllers.Delete[models.User](models.GetByIdQuery))))

	mux.Handle("GET /api/v1/athletes", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll(models.GetAthletesQuery, models.BuildAthleteResponse))))
	mux.Handle("GET /api/v1/athletes/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get(models.GetByIdQuery, models.BuildAthleteResponse))))

	mux.Handle("GET /api/v1/disciplines", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll(models.DefaultQuery, models.BuildDisciplineResponse))))
	mux.Handle("GET /api/v1/disciplines/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get(models.GetByIdQuery, models.BuildDisciplineResponse))))

	mux.Handle("GET /api/v1/events", m.ErrorsMiddleware(m.OrganizerOnly(controllers.GetAll(models.GetEventsQuery, models.BuildDefaultResponse[models.Event]))))
	mux.Handle("GET /api/v1/events/{id}", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Get(models.GetByIdQuery, models.BuildDefaultResponse[models.Event]))))
	mux.Handle("POST /api/v1/events", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Create[models.Event](models.BuildDefaultResponse))))
	mux.Handle("PUT /api/v1/events/{id}", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Update[models.Event](models.DefaultUpdateQuery, models.BuildDefaultResponse))))
	mux.Handle("DELETE /api/v1/events/{id}", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Delete[models.Event](models.GetByIdQuery))))

	mux.Handle("GET /api/v1/questions", m.ErrorsMiddleware(m.OrganizerOnly(controllers.GetAll(models.GetQuestionsQuery, models.BuildDefaultResponse[models.Question]))))
	mux.Handle("GET /api/v1/questions/{id}", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Get(models.GetByIdQuery, models.BuildDefaultResponse[models.Question]))))
	mux.Handle("POST /api/v1/questions", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Create[models.Question](models.BuildDefaultResponse))))
	mux.Handle("PUT /api/v1/questions/{id}", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Update[models.Question](models.DefaultUpdateQuery, models.BuildDefaultResponse))))
	mux.Handle("DELETE /api/v1/questions/{id}", m.ErrorsMiddleware(m.OrganizerOnly(controllers.Delete[models.Question](models.GetByIdQuery))))

	mux.Handle("GET /api/v1/users/me/answers", m.ErrorsMiddleware(m.UserOnly(controllers.GetAll(models.GetAnswersQuery, models.BuildDefaultResponse[models.Answer]))))
	mux.Handle("GET /api/v1/users/me/answers/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Get(models.GetAnswerQuery, models.BuildDefaultResponse[models.Answer]))))
	mux.Handle("POST /api/v1/users/me/answers", m.ErrorsMiddleware(m.UserOnly(controllers.Create[models.Answer](models.BuildDefaultResponse))))
	mux.Handle("PUT /api/v1/users/me/answers/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Update[models.Answer](models.UpdateAnswerQuery, models.BuildDefaultResponse))))
	mux.Handle("DELETE /api/v1/users/me/answers/{id}", m.ErrorsMiddleware(m.UserOnly(controllers.Delete[models.Answer](models.GetAnswerQuery))))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := utils.ErrorsResponse{
			ErrorType: "not_found_error",
			Details:   "path not found",
		}
		utils.Encode(w, r, http.StatusNotFound, response)
	})
}
