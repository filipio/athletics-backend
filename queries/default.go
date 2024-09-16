package queries

import (
	"net/http"

	"gorm.io/gorm"
)

type BuildQueryFunc func(db *gorm.DB, r *http.Request) *gorm.DB

func DefaultQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return db
}
