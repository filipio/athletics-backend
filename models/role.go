package models

type Role struct {
	AppModel
	Name  string `json:"name" validate:"required,oneof=Admin User" gorm:"not null;unique"`
	Users []User `json:"users" gorm:"many2many:user_roles;"`
}
