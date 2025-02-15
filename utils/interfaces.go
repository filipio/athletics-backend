package utils

import (
	"context"
	"net/http"

	"gorm.io/gorm"
)

type Validatable interface {
	Validate(*http.Request) error
}

type DbModel interface {
	Validatable
	GetID() uint
	BeforeCreateCtx(context.Context, *gorm.DB) error
	AfterCreateCtx(context.Context, *gorm.DB) error
	BeforeUpdateCtx(context.Context, *gorm.DB) error
	AfterUpdateCtx(context.Context, *gorm.DB) error
	BeforeDeleteCtx(context.Context, *gorm.DB) error
	AfterDeleteCtx(context.Context, *gorm.DB) error
	GetAllQuery(*gorm.DB, *http.Request) *gorm.DB
	GetQuery(*gorm.DB, *http.Request) *gorm.DB
	UpdateQuery(*gorm.DB, *http.Request) *gorm.DB
	DeleteQuery(*gorm.DB, *http.Request) *gorm.DB
	BuildResponse() any
}
