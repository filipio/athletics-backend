package models

import "time"

type Event struct {
	AppModel
	Name        string     `json:"name" gorm:"not null" validate:"required"`
	Description *string    `json:"description"`
	Deadline    time.Time  `json:"deadline" gorm:"not null" validate:"required"`
	Questions   []Question `json:"questions,omitempty" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}
