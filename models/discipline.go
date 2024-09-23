package models

type Discipline struct {
	AppModel
	Name     string    `json:"name" gorm:"not null"`
	Athletes []Athlete `json:"athletes" gorm:"many2many:athletes_disciplines;constraint:OnDelete:CASCADE"`
}
