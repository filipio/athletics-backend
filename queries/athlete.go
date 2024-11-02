package queries

import (
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func GetAthletesQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = db.Preload("Disciplines")

	queryFunctions := []func(db *gorm.DB) *gorm.DB{getByIds(r)}
	queryParams := r.URL.Query()

	if queryParams.Has("search") {
		searchTerm := "%" + strings.ToLower(queryParams.Get("search")) + "%"
		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Where("first_name || ' ' || last_name LIKE ? OR last_name || ' ' || first_name LIKE ?", searchTerm, searchTerm)
		})
	}

	if queryParams.Has("discipline_ids") {
		disciplineIds := strings.Split(queryParams.Get("discipline_ids"), ",")

		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN athletes_disciplines ON athletes_disciplines.athlete_id = athletes.id").
				Where("athletes_disciplines.discipline_id IN (?)", disciplineIds)
		})
	}

	if queryParams.Has("country") {
		country := strings.ToUpper(queryParams.Get("country"))
		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Where("country = ?", country)
		})
	}

	if queryParams.Has("gender") {
		gender := strings.ToLower(queryParams.Get("gender"))
		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Where("gender = ?", gender)
		})
	}

	// below is the way to fetch the user from the context
	// user := r.Context().Value(utils.UserContextKey).(models.User)

	return db.Scopes(queryFunctions...)
}
