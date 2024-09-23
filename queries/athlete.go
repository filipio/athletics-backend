package queries

import (
	"net/http"

	"gorm.io/gorm"
)

func GetAthletesQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return db.Preload("Disciplines")
}
