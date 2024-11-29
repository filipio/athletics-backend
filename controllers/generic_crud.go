package controllers

import (
	"context"
	"net/http"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

func GetAll[T utils.DbModel]() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := models.Db(r)
			var records []T
			var instance T

			query := instance.GetAllQuery(db, r)
			pageNo, perPage, orderBy, orderDirection := utils.PaginationParams(r)
			queryResult := models.Paginate(query, pageNo, perPage, orderBy, orderDirection).Find(&records)

			if queryResult.Error != nil {
				return queryResult.Error
			}

			var totalCount int64
			totalCountResult := instance.GetAllQuery(db.Model(&instance), r).Count(&totalCount)
			if totalCountResult.Error != nil {
				return totalCountResult.Error
			}

			var responseRecords []any = []any{}
			for _, record := range records {
				responseRecords = append(responseRecords, record.BuildResponse())
			}

			paginatedResponse := utils.BuildPaginatedResponse(responseRecords, totalCount, pageNo, perPage)
			if err := utils.Encode(w, r, http.StatusOK, paginatedResponse); err != nil {
				return err
			}

			return nil
		})
}

func Get[T utils.DbModel]() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := models.Db(r)
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

func Create[T utils.DbModel]() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := models.Db(r)
			record, err := utils.DecodeAndValidate[T](r)

			if err != nil {
				return err
			}

			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := record.BeforeCreateCtx(r.Context(), tx); err != nil {
					return err
				}

				if err := tx.Create(&record).Error; err != nil {
					return err
				}

				if err := record.AfterCreateCtx(r.Context(), tx); err != nil {
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

func Update[T utils.DbModel]() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := models.Db(r)
			record, err := utils.DecodeAndValidate[T](r)

			if err != nil {
				return err
			}

			baseQuery := db.Model(&record)
			query := record.UpdateQuery(baseQuery, r)
			id := utils.IntPathValue(r, "id")

			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := record.BeforeUpdateCtx(context.WithValue(r.Context(), utils.RecordIdContextKey, id), tx); err != nil {
					return err
				}

				queryResult := query.Updates(&record)

				if queryResult.Error != nil {
					return err
				}

				if queryResult.RowsAffected == 0 {
					return utils.RecordNotFoundError{}
				}

				if err := record.AfterUpdateCtx(context.WithValue(r.Context(), utils.RecordIdContextKey, id), tx); err != nil {
					return err
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

func Delete[T utils.DbModel]() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := models.Db(r)
			var record T

			query := record.DeleteQuery(db, r)
			id := utils.IntPathValue(r, "id")

			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := record.BeforeDeleteCtx(context.WithValue(r.Context(), utils.RecordIdContextKey, id), tx); err != nil {
					return err
				}

				queryResult := query.Delete(&record)

				if queryResult.Error != nil {
					return queryResult.Error
				}

				if queryResult.RowsAffected == 0 {
					return utils.RecordNotFoundError{}
				}

				if err := record.AfterDeleteCtx(context.WithValue(r.Context(), utils.RecordIdContextKey, id), tx); err != nil {
					return err
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
