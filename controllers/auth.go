package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const jwtTokenExpiration = time.Hour * 24 * 30

func Register(db *gorm.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			user, validationError, err := utils.Decode[models.User](r)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if validationError != nil {
				if err := utils.Encode(w, r, http.StatusBadRequest, validationError); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			// Hash the password
			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

			if err != nil {
				if err := utils.Encode(w, r, http.StatusBadRequest, utils.HashError); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			user.Password = string(hash)
			var role models.Role

			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(&user).Error; err != nil {
					return err
				}

				if err := tx.Where("name = ?", utils.UserRole).First(&role).Error; err != nil {
					return err
				}

				if err := tx.Model(&user).Association("Roles").Append(&role); err != nil {
					return err
				}

				if err := tx.First(&user, user.GetID()).Error; err != nil {
					return err
				}

				return nil
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := utils.Encode(w, r, http.StatusOK, utils.AnyMap{}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		})
}

func Login(db *gorm.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			bodyUser, validationError, decodeErr := utils.Decode[models.User](r)

			if decodeErr != nil {
				http.Error(w, decodeErr.Error(), http.StatusBadRequest)
				return
			}

			if validationError != nil {
				if err := utils.Encode(w, r, http.StatusBadRequest, validationError); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			var user models.User

			db.First(&user, "email = ?", bodyUser.Email)

			if user.GetID() == 0 {
				if err := utils.Encode(w, r, http.StatusUnauthorized, utils.AuthError); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(bodyUser.Password)); err != nil {
				if err := utils.Encode(w, r, http.StatusUnauthorized, utils.AuthError); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			var roles []models.Role

			if err := db.Model(&user).Association("Roles").Find(&roles); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var roleNames []string = make([]string, len(roles))
			for i, role := range roles {
				roleNames[i] = role.Name
			}

			// Generate a JWT token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"sub":   user.ID,
				"exp":   time.Now().Add(jwtTokenExpiration).Unix(),
				"roles": roleNames,
			})

			tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_SECRET")))

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			responseBody := utils.AnyMap{
				"token": tokenString,
			}

			if err := utils.Encode(w, r, http.StatusOK, responseBody); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		})
}
