package models

import (
	"gorm.io/datatypes"
)

type Question struct {
	AppModel
	EventID       uint            `json:"event_id" gorm:"not null" validate:"required,id_of=event"`
	Content       string          `json:"content" gorm:"not null" validate:"required"`
	CorrectAnswer *datatypes.JSON `json:"correct_answer"`
	Type          string          `json:"type" gorm:"not null" validate:"oneof=athlete country"`
}
