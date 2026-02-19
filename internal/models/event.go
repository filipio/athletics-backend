package models

import (
	"net/http"
	"time"

	"github.com/filipio/athletics-backend/pkg/httpio"
	"gorm.io/gorm"
)

type Event struct {
	AppModel
	Name        string     `json:"name" gorm:"not null" validate:"required"`
	Description *string    `json:"description"`
	Deadline    time.Time  `json:"deadline" gorm:"not null" validate:"required"`
	Status      string     `json:"status" gorm:"not null;default:draft" validate:"omitempty,oneof=draft published"`
	Questions   []Question `json:"questions,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}

func (m Event) isOrganizerOrAdmin(r *http.Request) bool {
	user, ok := r.Context().Value(httpio.UserContextKey).(User)
	if !ok {
		return false
	}

	for _, role := range user.Roles {
		if role.Name == httpio.OrganizerRole || role.Name == httpio.AdminRole {
			return true
		}
	}

	return false
}

func (m Event) GetAllQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = getByIds(db, r)
	queryParams := r.URL.Query()

	if queryParams.Get("active") == "true" {
		db = db.Where("NOW() < deadline")
	}

	if !m.isOrganizerOrAdmin(r) {
		db = db.Where("status = ?", "published")
	}

	return db
}

func (m Event) GetQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	if !m.isOrganizerOrAdmin(r) {
		db = db.Where("status = ?", "published")
	}

	return db
}

func (m Event) BuildResponse() any {
	return m
}

func (m *Event) IsPresent() error {
	if m.Deadline.Before(time.Now().UTC()) {
		return httpio.AppValidationError{
			FieldPath: "event_id",
			AppError:  httpio.AppError{Message: "event is already closed"},
		}
	}

	return nil
}
