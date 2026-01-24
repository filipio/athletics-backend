package controllers

import (
	"net/http"
	"time"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
)

func Logout() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := models.Db(r)

			sessionID, ok := r.Context().Value(utils.SessionIDContextKey).(string)
			if !ok || sessionID == "" {
				return utils.SessionExpiredError{}
			}

			now := time.Now()
			if err := db.Model(&models.RefreshToken{}).
				Where("session_id = ?", sessionID).
				Update("revoked_at", now).Error; err != nil {
				return err
			}

			if err := utils.Encode(w, r, http.StatusOK, utils.AnyMap{
				"message": "logged out successfully",
			}); err != nil {
				return err
			}

			return nil
		})
}
