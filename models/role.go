package models

type Role struct {
	AppModel
	Name  string `json:"name" validate:"required,one_of=admin user organizer" gorm:"not null;unique"`
	Users []User `json:"users" gorm:"many2many:user_roles;constraint:OnDelete:CASCADE"`
}
