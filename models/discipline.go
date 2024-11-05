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

func BuildDisciplineResponse(model Discipline) DisciplineResponse {
	return DisciplineResponse{
		ID:   model.ID,
		Name: model.Name,
		Type: model.Type,
	}
}
