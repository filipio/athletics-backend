package controllers

import (
	"net/http"

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

			result := db.First(&record, id)

			if result.Error != nil && result.Error.Error() != "record not found" {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}

			if result.RowsAffected == 0 {
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

func CreateOrUpdate[T any](db *gorm.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			decodedStruct, validationError, err := utils.Decode[T](r)

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

			result := db.Save(&decodedStruct)
			if result.Error != nil {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}

			if err := utils.Encode(w, r, http.StatusOK, decodedStruct); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		})
}

func Delete[T any](db *gorm.DB) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var typeInstance T
			id := utils.IntPathValue(r, "id")

			result := db.Delete(&typeInstance, id)

			if result.Error != nil {
				http.Error(w, result.Error.Error(), http.StatusInternalServerError)
				return
			}

			if result.RowsAffected == 0 {
				if err := utils.Encode(w, r, http.StatusBadRequest, utils.RecordNotFoundResponse); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		})
}
