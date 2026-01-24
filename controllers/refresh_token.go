package controllers

import (
	"net/http"
	"time"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RefreshTokenPayload struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

func (payload RefreshTokenPayload) Validate(r *http.Request) error {
	return nil
}

func RefreshToken() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			payload, err := utils.DecodeAndValidate[RefreshTokenPayload](r)
			if err != nil {
				return err
			}

			db := models.Db(r)

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
				return utils.InvalidRefreshTokenError{}
			}

			var user models.User
			db.First(&user, foundToken.UserID)
			if user.GetID() == 0 {
				return utils.LoginError{}
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

			if err := utils.Encode(w, r, http.StatusOK, tokenPair); err != nil {
				return err
			}

			return nil
		})
}
