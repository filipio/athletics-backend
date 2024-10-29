package models

import (
	"net/http"

	"github.com/filipio/athletics-backend/utils"
	"gorm.io/datatypes"
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
