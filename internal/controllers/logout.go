package controllers

import (
	"net/http"
	"time"

	"github.com/filipio/athletics-backend/pkg/config"
	"github.com/filipio/athletics-backend/internal/models"
	"github.com/filipio/athletics-backend/pkg/httpio"
)

func Logout(deps *config.Dependencies) httpio.HandlerWithError {
	return httpio.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := deps.DB

			sessionID, ok := r.Context().Value(httpio.SessionIDContextKey).(string)
			if !ok || sessionID == "" {
				return httpio.SessionExpiredError{}
			}

			now := time.Now()
			if err := db.Model(&models.RefreshToken{}).
				Where("session_id = ?", sessionID).
				Update("revoked_at", now).Error; err != nil {
				return err
			}

			if err := httpio.Encode(w, r, http.StatusOK, httpio.AnyMap{
				"message": "logged out successfully",
			}); err != nil {
				return err
			}

			return nil
		})
}
