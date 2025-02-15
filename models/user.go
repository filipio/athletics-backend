package models

import (
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	AppModel
	Username string   `json:"username" validate:"required" gorm:"not null;default:'no_name'"`
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

func (m User) GetAllQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return db.Preload("Roles")
}

func (m User) GetQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	return GetByIdQuery(db.Preload("Roles"), r)
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Email     string    `json:"email"`
	Roles     []string  `json:"roles"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m User) BuildResponse() any {
	roles := make([]string, len(m.Roles))
	for i, role := range m.Roles {
		roles[i] = role.Name
	}

	return UserResponse{
		ID:        m.ID,
		Email:     m.Email,
		Roles:     roles,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
