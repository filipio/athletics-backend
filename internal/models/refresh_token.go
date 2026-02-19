package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RefreshToken struct {
	AppModel
	UserID    uint       `json:"user_id" validate:"required" gorm:"not null;index"`
	TokenHash string     `json:"-" gorm:"not null;index"`
	SessionID string     `json:"session_id" gorm:"not null;unique;type:uuid;index"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null;index"`
	RevokedAt *time.Time `json:"revoked_at" gorm:"index"`
}

func (rt *RefreshToken) GenerateToken() string {
	return uuid.New().String()
}

func (rt *RefreshToken) GenerateSessionID() string {
	return uuid.New().String()
}

// HashToken hashes a token string using bcrypt
func (rt *RefreshToken) HashToken(token string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(token), 10)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func (rt RefreshToken) IsValid() bool {
	if rt.RevokedAt != nil {
		return false
	}
	return time.Now().Before(rt.ExpiresAt)
}

func (rt *RefreshToken) Revoke(db *gorm.DB) error {
	now := time.Now()
	rt.RevokedAt = &now
	return db.Model(rt).Update("revoked_at", now).Error
}
