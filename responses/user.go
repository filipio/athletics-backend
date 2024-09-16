package responses

import (
	"time"

	"github.com/filipio/athletics-backend/models"
)

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func BuildUserResponse(model models.User) UserResponse {
	roles := make([]string, len(model.Roles))
	for i, role := range model.Roles {
		roles[i] = role.Name
	}

	return UserResponse{
		ID:        model.ID,
		Email:     model.Email,
		Roles:     roles,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
