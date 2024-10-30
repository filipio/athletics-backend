package queries

import (
	"net/http"

	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

type BuildQueryFunc func(db *gorm.DB, r *http.Request) *gorm.DB

func DefaultQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return db
}

func GetByIdQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	id := utils.IntPathValue(r, "id")
	return db.Where("id = ?", id)
}

func DefaultUpdateQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return baseUpdateQuery(GetByIdQuery(db, r))
}

func baseUpdateQuery(db *gorm.DB) *gorm.DB {
	return db.Select("*").Omit("id", "created_at")
}

func Paginate(db *gorm.DB, pageNo int, perPage int, orderBy string, orderDirection string) *gorm.DB {
	return db.Offset((pageNo - 1) * perPage).Limit(perPage).Order(orderBy + " " + orderDirection)
}
