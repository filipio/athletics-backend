package queries

import (
	"net/http"

	"gorm.io/gorm"
)

func DefaultQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return db
}
