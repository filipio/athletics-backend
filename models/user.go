package models

type User struct {
	AppModel
	Email    string `json:"email" validate:"required,email" gorm:"not null;unique"`
	Password string `json:"password" validate:"required,min=6" gorm:"not null"`
	Roles    []Role `json:"roles" gorm:"many2many:user_roles;"`
}
