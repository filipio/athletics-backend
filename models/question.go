package models

import (
	"net/http"

	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

type Question struct {
	AppModel
	EventID       uint              `json:"event_id" gorm:"not null" validate:"required,id_of=event"`
	Content       string            `json:"content" gorm:"not null" validate:"required"`
	CorrectAnswer *AnswerOfQuestion `json:"correct_answer"`
	Type          string            `json:"type" gorm:"not null" validate:"oneof=athlete athletes_three country countries_three numeric_value"`
	Answers       []Answer          `json:"answers,omitempty" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
	Points        uint              `json:"points" gorm:"not null;default:1" validate:"required,gte=1"`
}

func (m Question) Validate(db *gorm.DB) error {
	if m.CorrectAnswer != nil {
		if err := m.CorrectAnswer.Validate(m.Type); err != nil {
			return utils.AppValidationError{
				FieldPath: "correct_answer",
				AppError:  utils.AppError{Message: err.Error()},
			}
		}
	}

	var event Event
	db.First(&event, m.EventID)

	return event.IsPresent()
}

func (m Question) GetAllQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = getByIds(db, r)
	queryParams := r.URL.Query()

	if queryParams.Has("event_id") {
		db = db.Where("event_id = ?", queryParams.Get("event_id"))
	}

	return db
}

func (m Question) BuildResponse() any {
	return m
}
