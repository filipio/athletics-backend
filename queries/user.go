package queries

import (
	"net/http"

	"gorm.io/gorm"
)

func GetUsersQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return db.Preload("Roles")
}

func GetUserQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return GetByIdQuery(db.Preload("Roles"), r)
}
