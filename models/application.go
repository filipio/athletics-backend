package models

import (
	"fmt"
	"net/http"
	"time"
)

type AppModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m AppModel) GetID() uint {
	return m.ID
}

func (m AppModel) Validate(r *http.Request) error {
	fmt.Println("validating itself!")
	return nil
}
