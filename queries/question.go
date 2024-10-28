package queries

import (
	"net/http"

	"gorm.io/gorm"
)

func GetQuestionsQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	queryFunctions := []func(db *gorm.DB) *gorm.DB{}
	queryParams := r.URL.Query()

	if queryParams.Has("event_id") {
		eventId := queryParams.Get("event_id")
		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Where("event_id = ?", eventId)
		})
	}

	return db.Scopes(queryFunctions...)
}
