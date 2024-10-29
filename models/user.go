package models

import (
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
