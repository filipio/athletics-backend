package controllers

import (
	"net/http"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/queries"
	"github.com/filipio/athletics-backend/responses"
	"github.com/filipio/athletics-backend/utils"
	"github.com/filipio/athletics-backend/utils/app_errors"
	"gorm.io/gorm"
)

func GetAll[T any, V any](db *gorm.DB, buildQuery queries.BuildQueryFunc, buildResponse responses.BuildResponseFunc[T, V]) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			var records []T

			query := buildQuery(db, r)
			queryResult := query.Find(&records)

			if queryResult.Error != nil {
				return queryResult.Error
			}

			var responses []V
			for _, record := range records {
				responses = append(responses, buildResponse(record))
			}

			if err := utils.Encode(w, r, http.StatusOK, responses); err != nil {
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
			queryResult := db.First(&record, id)

			if queryResult.Error != nil {
				if queryResult.Error.Error() != "record not found" {
					return queryResult.Error
				} else {
					return app_errors.RecordNotFoundError{}
				}
			}

			if err := utils.Encode(w, r, http.StatusOK, record); err != nil {
				return err
			}

			return nil

		})
}

func Create[T models.WithID](db *gorm.DB) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			record, err := utils.Decode[T](r)

			if err != nil {
				return err
			}

			queryResult := db.Create(&record)

			if queryResult.Error != nil {
				return queryResult.Error
			}

			db.First(&record, record.GetID())

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
			queryResult := db.Model(&record).Where("id = ?", id).Select("*").Omit("id", "created_at").Updates(&record)

			if queryResult.Error != nil {
				return queryResult.Error
			}

			if queryResult.RowsAffected == 0 {
				return app_errors.RecordNotFoundError{}
			}

			db.First(&record, id)
			if err := utils.Encode(w, r, http.StatusOK, record); err != nil {
				return err
			}

			return nil
		})
}

func Delete[T any](db *gorm.DB) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			var record T
			id := utils.IntPathValue(r, "id")

			queryResult := db.Delete(&record, id)

			if queryResult.Error != nil {
				return queryResult.Error
			}

			if queryResult.RowsAffected == 0 {
				return app_errors.RecordNotFoundError{}
			}

			if err := utils.Encode(w, r, http.StatusOK, utils.AnyMap{}); err != nil {
				return err
			}

			return nil
		})
}
