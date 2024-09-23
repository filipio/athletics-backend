package models

import "gorm.io/datatypes"

type Athlete struct {
	AppModel
	FirstName   *string         `json:"first_name"`
	LastName    *string         `json:"last_name"`
	Birthday    *datatypes.Date `json:"birthday"`
	Country     string          `json:"country"`
	Gender      string          `json:"gender" gorm:"not null"`
	Disciplines []Discipline    `json:"disciplines" gorm:"many2many:athletes_disciplines;constraint:OnDelete:CASCADE"`
}
