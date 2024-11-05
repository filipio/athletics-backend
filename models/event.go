package models

import (
	"net/http"
	"time"

	"gorm.io/gorm"
)

type Event struct {
	AppModel
	Name        string     `json:"name" gorm:"not null" validate:"required"`
	Description *string    `json:"description"`
	Deadline    time.Time  `json:"deadline" gorm:"not null" validate:"required"`
	Questions   []Question `json:"questions,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}

func GetEventsQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	queryFunctions := []func(db *gorm.DB) *gorm.DB{getByIds(r)}
	queryParams := r.URL.Query()

	if queryParams.Get("active") == "true" {
		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Where("NOW() < deadline")
		})
	}

	return db.Scopes(queryFunctions...)
}
