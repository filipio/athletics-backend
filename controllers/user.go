package controllers

import (
	"net/http"

	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUser(deps *config.Dependencies) utils.HandlerWithError {
	return utils.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		db := deps.DB
		user, err := utils.DecodeAndValidate[models.User](r, db)
		if err != nil {
			return err
		}

		// Business logic: Hash password (moved from BeforeCreate hook)
		if !user.SkipPasswordHashing {
			hashedPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
			if err != nil {
				return err
			}
			user.Password = string(hashedPasswordBytes)
		}

		// Create user
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&user).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}

		db.First(&user, user.GetID())
		response := user.BuildResponse()

		if err := utils.Encode(w, r, http.StatusOK, response); err != nil {
			return err
		}

		return nil
	})
}
