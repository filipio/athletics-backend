package queries

import (
	"net/http"

	"gorm.io/gorm"
)

func GetPokemonsQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	queryFunctions := []func(db *gorm.DB) *gorm.DB{}
	queryParams := r.URL.Query()

	if queryParams.Has("name") {
		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Where("pokemon_name IN (?)", queryParams.Get("name"))
		})
	}

	// TODO: remove below
	// userAny := r.Context().Value(utils.UserContextKey).(models.User)
	// fmt.Println(userAny.ID, userAny.Email)

	return db.Scopes(queryFunctions...)
}
