package controllers

import (
	"net/http"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"github.com/filipio/athletics-backend/utils/app_errors"
	"gorm.io/gorm"
)

func GetAll[T any](db *gorm.DB) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			var records []T

			result := db.Find(&records)
			if result.Error != nil {
				return result.Error
			}

			if err := utils.Encode(w, r, http.StatusOK, records); err != nil {
				return err
			}

			return nil
		})
}

func Get[T any](db *gorm.DB) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			var record T
			id := utils.IntPathValue(r, "id")

			dbResult := db.First(&record, id)

			if dbResult.Error != nil && dbResult.Error.Error() != "record not found" {
				return dbResult.Error
			}

			if dbResult.RowsAffected == 0 {
				return app_errors.RecordNotFoundError{}
			} else {
				if err := utils.Encode(w, r, http.StatusOK, record); err != nil {
					return err
				}
			}

			return nil

		})
}

func Create[T models.WithID](db *gorm.DB) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			// user, _ := r.Context().Value(utils.UserContextKey).(models.User)
			// fmt.Println(user.Email)

			record, err := utils.Decode[T](r)

			if err != nil {
				return err
			}

			dbResult := db.Create(&record)

			if dbResult.Error != nil {
				return dbResult.Error
			}

			db.First(&record, record.GetID()) // this is done to make dates format properly (gorm issue)

			if err := utils.Encode(w, r, http.StatusOK, record); err != nil {
				return err
			}

			return nil
		})
}

func Update[T any](db *gorm.DB) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			record, err := utils.Decode[T](r)

			if err != nil {
				return err
			}

			id := utils.IntPathValue(r, "id")
			dbResult := db.Model(&record).Where("id = ?", id).Select("*").Omit("id", "created_at").Updates(&record)

			if dbResult.Error != nil {
				return dbResult.Error
			}

			if dbResult.RowsAffected == 0 {
				return app_errors.RecordNotFoundError{}
			} else {
				db.First(&record, id)
				if err := utils.Encode(w, r, http.StatusOK, record); err != nil {
					return err
				}
			}

			return nil
		})
}

func Delete[T any](db *gorm.DB) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			var record T
			id := utils.IntPathValue(r, "id")

			dbResult := db.Delete(&record, id)

			if dbResult.Error != nil {
				return dbResult.Error
			}

			if dbResult.RowsAffected == 0 {
				return app_errors.RecordNotFoundError{}
			}

			return nil
		})
}
