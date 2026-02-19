package controllers

import (
	"net/http"
	"time"

	"github.com/filipio/athletics-backend/pkg/config"
	"github.com/filipio/athletics-backend/internal/models"
	"github.com/filipio/athletics-backend/pkg/httpio"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RefreshTokenPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (payload RefreshTokenPayload) Validate(db *gorm.DB) error {
	return nil
}

func RefreshToken(deps *config.Dependencies) httpio.HandlerWithError {
	return httpio.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := deps.DB
			payload, err := httpio.DecodeAndValidate[RefreshTokenPayload](r, db)
			if err != nil {
				return err
			}

			var refreshTokens []models.RefreshToken
			db.Where("expires_at > ? AND revoked_at IS NULL", time.Now()).Find(&refreshTokens)

			var foundToken *models.RefreshToken
			for i := range refreshTokens {
				if err := bcrypt.CompareHashAndPassword([]byte(refreshTokens[i].TokenHash), []byte(payload.RefreshToken)); err == nil {
					foundToken = &refreshTokens[i]
					break
				}
			}

			if foundToken == nil {
				return httpio.InvalidRefreshTokenError{}
			}

			var user models.User
			db.First(&user, foundToken.UserID)
			if user.GetID() == 0 {
				return httpio.LoginError{}
			}

			sessionID := foundToken.SessionID

			var tokenPair TokenPair
			err = db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Delete(foundToken).Error; err != nil {
					return err
				}

				var err error
				tokenPair, err = generateTokenPair(user, tx, &sessionID)
				return err
			})
			if err != nil {
				return err
			}

			if err := httpio.Encode(w, r, http.StatusOK, tokenPair); err != nil {
				return err
			}

			return nil
		})
}
