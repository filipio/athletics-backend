package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/filipio/athletics-backend/email"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RequestVerificationPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
}

func (payload RequestVerificationPayload) Validate(r *http.Request) error {
	return nil
}

type VerifyEmailPayload struct {
	Token string `json:"token" validate:"required"`
}

func (payload VerifyEmailPayload) Validate(r *http.Request) error {
	return nil
}

func RequestVerification() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			payload, err := utils.DecodeAndValidate[RequestVerificationPayload](r)
			if err != nil {
				return err
			}

			db := models.Db(r)

			var existingUser models.User
			db.Where("email = ?", payload.Email).First(&existingUser)
			if existingUser.ID != 0 {
				return utils.EmailAlreadyExistsError{}
			}

			var rateLimit models.EmailVerificationRateLimit
			db.FirstOrCreate(&rateLimit, models.EmailVerificationRateLimit{Email: payload.Email})

			if !rateLimit.CanRequestVerification() {
				blockedUntilStr := ""
				if rateLimit.BlockedUntil != nil {
					blockedUntilStr = rateLimit.BlockedUntil.Format(time.RFC3339)
				}
				return utils.EmailVerificationRateLimitError{
					BlockedUntil: &blockedUntilStr,
				}
			}

			db.Where("email = ? AND verified = ?", payload.Email, false).Delete(&models.PendingRegistration{})

			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), 10)
			if err != nil {
				return err
			}

			pendingReg := models.PendingRegistration{
				Email:        payload.Email,
				Username:     payload.Username,
				PasswordHash: string(hashedPassword),
			}

			if err := pendingReg.GenerateVerificationToken(); err != nil {
				return err
			}

			err = db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(&pendingReg).Error; err != nil {
					return err
				}

				rateLimit.IncrementRequestCount()
				if err := tx.Save(&rateLimit).Error; err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				return err
			}

			ctx := context.Background()
			emailErr := email.SendVerificationEmail(ctx, email.VerificationEmailParams{
				To:                payload.Email,
				VerificationToken: pendingReg.VerificationToken,
			})

			if emailErr != nil {
				return utils.EmailSendError{
					OriginalError: emailErr,
				}
			}

			utils.Encode(w, r, http.StatusOK, utils.AnyMap{
				"message": "Verification email sent. Please check your inbox.",
			})

			return nil
		})
}

func VerifyEmail() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			payload, err := utils.DecodeAndValidate[VerifyEmailPayload](r)
			if err != nil {
				return err
			}

			db := models.Db(r)

			var pendingReg models.PendingRegistration
			db.Where("verification_token = ? AND verified = ?", payload.Token, false).First(&pendingReg)

			if pendingReg.ID == 0 {
				return utils.InvalidVerificationTokenError{}
			}

			if pendingReg.IsExpired() {
				return utils.InvalidVerificationTokenError{}
			}

			user := models.User{
				Email:               pendingReg.Email,
				Username:            pendingReg.Username,
				Password:            pendingReg.PasswordHash,
				SkipPasswordHashing: true,
			}

			var role models.Role

			err = db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(&user).Error; err != nil {
					return err
				}

				if err := tx.Where("name = ?", utils.UserRole).First(&role).Error; err != nil {
					return err
				}

				if err := tx.Model(&user).Association("Roles").Append(&role); err != nil {
					return err
				}

				if err := tx.Delete(&pendingReg).Error; err != nil {
					return err
				}

				if err := tx.Where("email = ?", pendingReg.Email).Delete(&models.EmailVerificationRateLimit{}).Error; err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				return err
			}

			tokenPair, err := generateTokenPair(user, db, nil)
			if err != nil {
				return err
			}

			utils.Encode(w, r, http.StatusOK, tokenPair)

			return nil
		})
}
