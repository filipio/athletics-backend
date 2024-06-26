package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"github.com/filipio/athletics-backend/utils/app_errors"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const jwtTokenExpiration = time.Hour * 24 * 30

func Register(db *gorm.DB) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			user, err := utils.Decode[models.User](r)

			if err != nil {
				return err
			}

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
				return err
			}

			utils.Encode(w, r, http.StatusOK, utils.AnyMap{})

			return nil
		})
}

func Login(db *gorm.DB) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			bodyUser, decodeErr := utils.Decode[models.User](r)

			if decodeErr != nil {
				return decodeErr
			}

			var user models.User
			db.First(&user, "email = ?", bodyUser.Email)

			if user.GetID() == 0 {
				return app_errors.LoginError{}
			}

			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(bodyUser.Password)); err != nil {
				return app_errors.LoginError{}
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
				"sub":   user.ID,
				"exp":   time.Now().Add(jwtTokenExpiration).Unix(),
				"roles": roleNames,
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
