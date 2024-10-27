package responses

import (
	"github.com/filipio/athletics-backend/models"
)

type DisciplineResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func BuildDisciplineResponse(model models.Discipline) DisciplineResponse {
	return DisciplineResponse{
		ID:   model.ID,
		Name: model.Name,
		Type: model.Type,
	}
}
