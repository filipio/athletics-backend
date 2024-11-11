package models

type Discipline struct {
	AppModel
	Name     string    `json:"name" gorm:"not null"`
	Type     string    `json:"type" gorm:"not null"`
	Athletes []Athlete `json:"athletes" gorm:"many2many:athletes_disciplines;constraint:OnDelete:CASCADE"`
}

type DisciplineResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (m Discipline) BuildResponse() any {
	return DisciplineResponse{
		ID:   m.ID,
		Name: m.Name,
		Type: m.Type,
	}
}
