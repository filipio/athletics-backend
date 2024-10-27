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

	return db.Scopes(queryFunctions...)
}
