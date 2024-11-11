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

func (m Answer) GetAllQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = onlyCurrentUserRecords(db, r)
	db = getByIds(db, r)

	queryParams := r.URL.Query()
	if queryParams.Has("question_id") {
		db = db.Where("question_id IN (?)", queryParams.Get("question_id"))
	}

	return db
}

func (m Answer) GetQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = onlyCurrentUserRecords(db, r)
	db = GetByIdQuery(db, r)
	return db
}

func (m Answer) BuildResponse() any {
	return m
}

func onlyCurrentUserRecords(db *gorm.DB, r *http.Request) *gorm.DB {
	currentUser := r.Context().Value(utils.UserContextKey).(User)
	return db.Where("user_id = ?", currentUser.ID)
}
