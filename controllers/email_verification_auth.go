package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/filipio/athletics-backend/email"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"github.com/golang-jwt/jwt/v5"
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

			// Find pending registration by token
			var pendingReg models.PendingRegistration
			db.Where("verification_token = ? AND verified = ?", payload.Token, false).First(&pendingReg)

			if pendingReg.ID == 0 {
				return utils.InvalidVerificationTokenError{}
			}

			if pendingReg.IsExpired() {
				return utils.InvalidVerificationTokenError{}
			}

			// Create user
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

				// Delete pending registration
				if err := tx.Delete(&pendingReg).Error; err != nil {
					return err
				}

				// Delete rate limit
				if err := tx.Where("email = ?", pendingReg.Email).Delete(&models.EmailVerificationRateLimit{}).Error; err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				return err
			}

			// Generate JWT token
			var roles []models.Role
			if err := db.Model(&user).Association("Roles").Find(&roles); err != nil {
				return err
			}

			var roleNames []string = make([]string, len(roles))
			for i, role := range roles {
				roleNames[i] = role.Name
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub":      user.ID,
				"exp":      time.Now().Add(jwtTokenExpiration).Unix(),
				"roles":    roleNames,
				"username": user.Username,
			})

			tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_SECRET")))
			if err != nil {
				return err
			}

			utils.Encode(w, r, http.StatusOK, utils.AnyMap{
				"token": tokenString,
				"user":  user.BuildResponse(),
			})

			return nil
		})
}
