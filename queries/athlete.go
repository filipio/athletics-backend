package queries

import (
	"fmt"
	"net/http"
	"strings"

	"gorm.io/gorm"
)

func GetAthletesQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = db.Preload("Disciplines")

	queryFunctions := []func(db *gorm.DB) *gorm.DB{}
	queryParams := r.URL.Query()

	if queryParams.Has("search") {
		searchTerm := "%" + strings.ToLower(queryParams.Get("search")) + "%"
		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Where("first_name || ' ' || last_name LIKE ? OR last_name || ' ' || first_name LIKE ?", searchTerm, searchTerm)
		})
	}

	// filter by discipline_ids passed in query params, using disciplines relationship
	if queryParams.Has("discipline_ids") {
		fmt.Println(queryParams.Get("discipline_ids"))
		disciplineIds := strings.Split(queryParams.Get("discipline_ids"), ",")

		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Joins("JOIN athletes_disciplines ON athletes_disciplines.athlete_id = athletes.id").
				Where("athletes_disciplines.discipline_id IN (?)", disciplineIds)
		})
	}

	// below is the way to fetch the user from the context
	// user := r.Context().Value(utils.UserContextKey).(models.User)

	return db.Scopes(queryFunctions...)
}
