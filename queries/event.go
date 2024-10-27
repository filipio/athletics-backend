package queries

import (
	"net/http"

	"gorm.io/gorm"
)

func GetEventsQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	queryFunctions := []func(db *gorm.DB) *gorm.DB{}
	queryParams := r.URL.Query()

	if queryParams.Get("active") == "true" {
		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Where("NOW() < deadline")
		})
	}

	return db.Scopes(queryFunctions...)
}
