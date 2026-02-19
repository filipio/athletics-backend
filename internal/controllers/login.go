package controllers

import (
	"net/http"

	"github.com/filipio/athletics-backend/pkg/config"
	"github.com/filipio/athletics-backend/internal/models"
	"github.com/filipio/athletics-backend/pkg/httpio"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (payload LoginPayload) Validate(db *gorm.DB) error {
	return nil
}

func Login(deps *config.Dependencies) httpio.HandlerWithError {
	return httpio.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := deps.DB
			loginPayload, decodeErr := httpio.DecodeAndValidate[LoginPayload](r, db)

			if decodeErr != nil {
				return decodeErr
			}

			var user models.User
			db.First(&user, "email = ?", loginPayload.Email)

			if user.GetID() == 0 {
				return httpio.LoginError{}
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginPayload.Password)); err != nil {
				return httpio.LoginError{}
			}

			tokenPair, err := generateTokenPair(user, db, nil)
			if err != nil {
				return err
			}

			if err := httpio.Encode(w, r, http.StatusOK, tokenPair); err != nil {
				return err
			}

			return nil

		})
}
