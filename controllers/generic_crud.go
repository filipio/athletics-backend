package controllers

import (
	"net/http"

	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

func GetAll[T utils.DbModel](deps *config.Dependencies) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := deps.DB
			var instance T

			query := instance.GetAllQuery(db, r)

			var totalCount int64

			totalCountResult := query.Model(&instance).Count(&totalCount)
			if totalCountResult.Error != nil {
				return totalCountResult.Error
			}

			paginationParams := utils.BuildPaginationParams(r)

			var records []T
			queryResult := models.PaginateQuery(query, paginationParams).Find(&records)
			if queryResult.Error != nil {
				return queryResult.Error
			}

			var responseRecords []any = make([]any, len(records))
			for i, record := range records {
				responseRecords[i] = record.BuildResponse()
			}

			paginationResponse := utils.BuildPaginatedResponse(responseRecords, totalCount, paginationParams)
			if err := utils.Encode(w, r, http.StatusOK, paginationResponse); err != nil {
				return err
			}

			return nil
		})
}

func Get[T utils.DbModel](deps *config.Dependencies) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := deps.DB
			var record T

			query := record.GetQuery(db, r)
			queryResult := query.First(&record)

			if queryResult.Error != nil {
				if queryResult.Error.Error() != "record not found" {
					return queryResult.Error
				} else {
					return utils.RecordNotFoundError{}
				}
			}

			response := record.BuildResponse()

			if err := utils.Encode(w, r, http.StatusOK, response); err != nil {
				return err
			}

			return nil

		})
}

func Create[T utils.DbModel](deps *config.Dependencies) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := deps.DB
			record, err := utils.DecodeAndValidate[T](r, db)

			if err != nil {
				return err
			}

			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := tx.Create(&record).Error; err != nil {
					return err
				}

				return nil
			}); err != nil {
				return err
			}

			db.First(&record, record.GetID())
			response := record.BuildResponse()

			if err := utils.Encode(w, r, http.StatusOK, response); err != nil {
				return err
			}

			return nil
		})
}

func Update[T utils.DbModel](deps *config.Dependencies) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := deps.DB
			record, err := utils.DecodeAndValidate[T](r, db)

			if err != nil {
				return err
			}

			baseQuery := db.Model(&record)
			query := record.UpdateQuery(baseQuery, r)
			id := utils.IntPathValue(r, "id")

			if err := db.Transaction(func(tx *gorm.DB) error {
				queryResult := query.Updates(&record)

				if queryResult.Error != nil {
					return queryResult.Error
				}

				if queryResult.RowsAffected == 0 {
					return utils.RecordNotFoundError{}
				}

				return nil
			}); err != nil {
				return err
			}

			db.First(&record, id)
			response := record.BuildResponse()

			if err := utils.Encode(w, r, http.StatusOK, response); err != nil {
				return err
			}

			return nil
		})
}

func Delete[T utils.DbModel](deps *config.Dependencies) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := deps.DB
			var record T

			query := record.DeleteQuery(db, r)

			if err := db.Transaction(func(tx *gorm.DB) error {
				queryResult := query.Delete(&record)

				if queryResult.Error != nil {
					return queryResult.Error
				}

				if queryResult.RowsAffected == 0 {
					return utils.RecordNotFoundError{}
				}

				return nil
			}); err != nil {
				return err
			}

			if err := utils.Encode(w, r, http.StatusOK, utils.AnyMap{}); err != nil {
				return err
			}

			return nil
		})
}
