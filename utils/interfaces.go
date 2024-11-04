package utils

import (
	"net/http"
)

type DbModel interface {
	GetID() uint
	Validate(*http.Request) error
}
