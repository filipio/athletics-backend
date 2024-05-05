package controllers

import (
	"log"
	"net/http"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

func GetAll[T any](db *gorm.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var records []T

			result := db.Find(&records)
			if result.Error != nil {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}

			if err := utils.Encode(w, r, http.StatusOK, records); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}

func Get[T any](db *gorm.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var record T
			id := utils.IntPathValue(r, "id")

			dbResult := db.First(&record, id)

			if dbResult.Error != nil && dbResult.Error.Error() != utils.NotFoundError {
				http.Error(w, dbResult.Error.Error(), http.StatusInternalServerError)
				return
			}

			if dbResult.RowsAffected == 0 {
				if err := utils.Encode(w, r, http.StatusBadRequest, utils.RecordNotFoundResponse); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				if err := utils.Encode(w, r, http.StatusOK, record); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

		})
}

func Create[T models.WithID](db *gorm.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			record, validationError, err := utils.Decode[T](r)

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

			dbResult := db.Create(&record)
			log.Print(record.GetID())

			if dbResult.Error != nil {
				http.Error(w, dbResult.Error.Error(), http.StatusInternalServerError)
				return
			}

			db.First(&record, record.GetID()) // this is done to make dates format properly (gorm issue)

			if err := utils.Encode(w, r, http.StatusOK, record); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}

func Update[T any](db *gorm.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			record, validationError, err := utils.Decode[T](r)

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

			id := utils.IntPathValue(r, "id")
			dbResult := db.Model(&record).Where("id = ?", id).Select("*").Omit("id", "created_at").Updates(&record)

			if dbResult.Error != nil {
				log.Print("error occurred", dbResult.Error.Error())
				http.Error(w, dbResult.Error.Error(), http.StatusInternalServerError)
				return
			}

			if dbResult.RowsAffected == 0 {
				if err := utils.Encode(w, r, http.StatusBadRequest, utils.RecordNotFoundResponse); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			} else {
				db.First(&record, id)
				if err := utils.Encode(w, r, http.StatusOK, record); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		})
}

func Delete[T any](db *gorm.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var record T
			id := utils.IntPathValue(r, "id")

			dbResult := db.Delete(&record, id)

			if dbResult.Error != nil {
				http.Error(w, dbResult.Error.Error(), http.StatusInternalServerError)
				return
			}

			if dbResult.RowsAffected == 0 {
				if err := utils.Encode(w, r, http.StatusBadRequest, utils.RecordNotFoundResponse); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		})
}
