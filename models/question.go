package models

import (
	"fmt"
	"net/http"

	"gorm.io/datatypes"
)

type Question struct {
	AppModel
	EventID       uint            `json:"event_id" gorm:"not null" validate:"required,id_of=event"`
	Content       string          `json:"content" gorm:"not null" validate:"required"`
	CorrectAnswer *datatypes.JSON `json:"correct_answer"`
	Type          string          `json:"type" gorm:"not null" validate:"oneof=athlete country"`
	Answers       []Answer        `json:"answers,omitempty" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
}

func (m Question) Validate(r *http.Request) error {
	fmt.Println("questions validation is executed")
	return nil
}
