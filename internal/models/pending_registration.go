package models

import (
	"time"

	"github.com/google/uuid"
)

const PendingRegistrationExpirationHours = 24

type PendingRegistration struct {
	AppModel
	Email             string    `json:"email" validate:"required,email" gorm:"not null;index"`
	Username          string    `json:"username" validate:"required" gorm:"not null"`
	PasswordHash      string    `json:"-" gorm:"not null"`
	VerificationToken string    `json:"-" gorm:"not null;index"`
	ExpiresAt         time.Time `json:"expires_at" gorm:"not null;index"`
	Verified          bool      `json:"verified" gorm:"not null;default:false"`
}

func (pr *PendingRegistration) GenerateVerificationToken() error {
	pr.VerificationToken = uuid.New().String()
	pr.ExpiresAt = time.Now().Add(time.Hour * PendingRegistrationExpirationHours)
	return nil
}

func (pr PendingRegistration) IsExpired() bool {
	return time.Now().After(pr.ExpiresAt)
}
