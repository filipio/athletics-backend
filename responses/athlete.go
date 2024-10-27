package responses

import (
	"github.com/filipio/athletics-backend/models"
)

type AthleteResponse struct {
	ID          uint     `json:"id"`
	FirstName   *string  `json:"first_name"`
	LastName    *string  `json:"last_name"`
	Gender      string   `json:"gender"`
	Country     *string  `json:"country"`
	Disciplines []string `json:"disciplines"`
}

func BuildAthleteResponse(model models.Athlete) AthleteResponse {
	disciplines := make([]string, len(model.Disciplines))
	for i, discipline := range model.Disciplines {
		disciplines[i] = discipline.Name
	}
	athleteResponse := AthleteResponse{
		ID:          model.ID,
		FirstName:   model.FirstName,
		LastName:    model.LastName,
		Gender:      model.Gender,
		Country:     model.Country,
		Disciplines: disciplines,
	}

	return athleteResponse
}
