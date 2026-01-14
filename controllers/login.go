package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginPayload struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (payload LoginPayload) Validate(r *http.Request) error {
	return nil
}

const jwtTokenExpiration = time.Hour * 24 * 30

func Login() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			loginPayload, decodeErr := utils.DecodeAndValidate[LoginPayload](r)

			if decodeErr != nil {
				return decodeErr
			}

			db := models.Db(r)

			var user models.User
			db.First(&user, "email = ?", loginPayload.Email)

			if user.GetID() == 0 {
				return utils.LoginError{}
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginPayload.Password)); err != nil {
				return utils.LoginError{}
			}

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

			responseBody := utils.AnyMap{
				"token": tokenString,
			}

			if err := utils.Encode(w, r, http.StatusOK, responseBody); err != nil {
				return err
			}

			return nil

		})
}
