package queries

import (
	"net/http"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

func GetAnswersQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	queryFunctions := []func(db *gorm.DB) *gorm.DB{onlyCurrentUserRecords(r)}

	queryParams := r.URL.Query()
	if queryParams.Has("question_id") {
		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Where("question_id IN (?)", queryParams.Get("question_id"))
		})
	}

	return db.Scopes(queryFunctions...)
}

func GetAnswerQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = GetByIdQuery(db, r)
	queryFunctions := []func(db *gorm.DB) *gorm.DB{onlyCurrentUserRecords(r)}
	return db.Scopes(queryFunctions...)
}

func onlyCurrentUserRecords(r *http.Request) func(db *gorm.DB) *gorm.DB {
	currentUser := r.Context().Value(utils.UserContextKey).(models.User)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", currentUser.ID)
	}
}
