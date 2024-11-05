package models

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	AppModel
	Email    string   `json:"email" validate:"required,email" gorm:"not null;unique"`
	Password string   `json:"password" validate:"required,min=6" gorm:"not null"`
	Roles    []Role   `json:"roles" gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
	Answers  []Answer `json:"answers,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}

	u.Password = string(hashedPasswordBytes)
	return nil
}

func GetUsersQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return db.Preload("Roles")
}

func GetUserQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return GetByIdQuery(db.Preload("Roles"), r)
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func BuildUserResponse(model User) UserResponse {
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
