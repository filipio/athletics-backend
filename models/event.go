package models

import (
	"net/http"
	"time"

	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

type Event struct {
	AppModel
	Name        string     `json:"name" gorm:"not null" validate:"required"`
	Description *string    `json:"description"`
	Deadline    time.Time  `json:"deadline" gorm:"not null" validate:"required"`
	Questions   []Question `json:"questions,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}

func (m Event) GetAllQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = getByIds(db, r)
	queryParams := r.URL.Query()

	if queryParams.Get("active") == "true" {
		db = db.Where("NOW() < deadline")
	}

	return db
}

func (m Event) BuildResponse() any {
	return m
}

func (m *Event) IsPresent() error {
	if m.Deadline.Before(time.Now().UTC()) {
		return utils.AppValidationError{
			FieldPath: "event_id",
			AppError:  utils.AppError{Message: "event is already closed"},
		}
	}

	return nil
}
