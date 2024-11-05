package models

import (
	"net/http"

	"github.com/filipio/athletics-backend/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Answer struct {
	AppModel
	UserID     uint           `json:"user_id" gorm:"not null" validate:"required"`
	QuestionID uint           `json:"question_id" gorm:"not null" validate:"required,id_of=question"`
	Content    datatypes.JSON `json:"content" gorm:"not null" validate:"required"`
}

func (m Answer) Validate(r *http.Request) error {
	currentUser := r.Context().Value(utils.UserContextKey).(User)

	if m.UserID != currentUser.ID {
		return utils.InvalidUserError{}
	}

	return nil
}

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

func UpdateAnswerQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return baseUpdateQuery(GetAnswerQuery(db, r))
}

func onlyCurrentUserRecords(r *http.Request) func(db *gorm.DB) *gorm.DB {
	currentUser := r.Context().Value(utils.UserContextKey).(User)
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", currentUser.ID)
	}
}
