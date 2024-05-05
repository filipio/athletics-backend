package models

import (
	"time"
)

type WithID interface {
	GetID() uint
}

type AppModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m AppModel) GetID() uint {
	return m.ID
}
