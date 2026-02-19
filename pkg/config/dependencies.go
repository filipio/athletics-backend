package config

import (
	"github.com/filipio/athletics-backend/internal/email"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB          *gorm.DB
	Workers     *InsertWorkerClient
	EmailSender email.EmailSender
}

func NewDependencies(db *gorm.DB, workers *InsertWorkerClient, emailSender email.EmailSender) *Dependencies {
	return &Dependencies{
		DB:          db,
		Workers:     workers,
		EmailSender: emailSender,
	}
}
