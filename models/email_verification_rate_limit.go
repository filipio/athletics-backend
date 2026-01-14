package models

import (
	"time"
)

const (
	MaxVerificationRequests   = 3
	VerificationWindowMinutes = 15
)

type EmailVerificationRateLimit struct {
	AppModel
	Email         string     `json:"email" gorm:"not null;unique;index"`
	RequestCount  int        `json:"request_count" gorm:"not null;default:0"`
	LastRequestAt *time.Time `json:"last_request_at"`
	BlockedUntil  *time.Time `json:"blocked_until" gorm:"index"`
}

func (evrl *EmailVerificationRateLimit) IsBlocked() bool {
	if evrl.BlockedUntil == nil {
		return false
	}
	return time.Now().Before(*evrl.BlockedUntil)
}

func (evrl *EmailVerificationRateLimit) CanRequestVerification() bool {
	if evrl.IsBlocked() {
		return false
	}

	now := time.Now()

	// Reset counter if window has passed
	if evrl.LastRequestAt != nil && now.Sub(*evrl.LastRequestAt) > time.Duration(VerificationWindowMinutes)*time.Minute {
		evrl.RequestCount = 0
	}

	return evrl.RequestCount < MaxVerificationRequests
}

func (evrl *EmailVerificationRateLimit) IncrementRequestCount() {
	now := time.Now()
	evrl.RequestCount++
	evrl.LastRequestAt = &now

	if evrl.RequestCount >= MaxVerificationRequests {
		blockUntil := now.Add(time.Duration(VerificationWindowMinutes) * time.Minute)
		evrl.BlockedUntil = &blockUntil
	}
}
