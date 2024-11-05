package controllers

import (
	"net/http"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
)

func GetAll[T utils.DbModel, V any](buildQuery models.BuildQueryFunc, buildResponse models.BuildResponseFunc[T, V]) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := models.Db(r)
			var records []T

			query := buildQuery(db, r)
			pageNo, perPage, orderBy, orderDirection := utils.PaginationParams(r)
			queryResult := models.Paginate(query, pageNo, perPage, orderBy, orderDirection).Find(&records)

			if queryResult.Error != nil {
				return queryResult.Error
			}

			var instance T
			var totalCount int64
			totalCountResult := buildQuery(db.Model(&instance), r).Count(&totalCount)
			if totalCountResult.Error != nil {
				return totalCountResult.Error
			}

			var responseRecords []V = []V{}
			for _, record := range records {
				responseRecords = append(responseRecords, buildResponse(record))
			}

			paginatedResponse := utils.BuildPaginatedResponse(responseRecords, totalCount, pageNo, perPage)
			if err := utils.Encode(w, r, http.StatusOK, paginatedResponse); err != nil {
				return err
			}

			return nil
		})
}

func Get[T utils.DbModel, V any](buildQuery models.BuildQueryFunc, buildResponse models.BuildResponseFunc[T, V]) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := models.Db(r)
			var record T

			query := buildQuery(db, r)
			queryResult := query.First(&record)

			if queryResult.Error != nil {
				if queryResult.Error.Error() != "record not found" {
					return queryResult.Error
				} else {
					return utils.RecordNotFoundError{}
				}
			}

			response := buildResponse(record)

			if err := utils.Encode(w, r, http.StatusOK, response); err != nil {
				return err
			}

			return nil

		})
}

func Create[T utils.DbModel, V any](buildResponse models.BuildResponseFunc[T, V]) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := models.Db(r)
			record, err := utils.DecodeAndValidate[T](r)

			if err != nil {
				return err
			}

			queryResult := db.Create(&record)

			if queryResult.Error != nil {
				return queryResult.Error
			}

			db.First(&record, record.GetID())
			response := buildResponse(record)

			if err := utils.Encode(w, r, http.StatusOK, response); err != nil {
				return err
			}

			return nil
		})
}

// todo: restrict to only given query
func Update[T utils.DbModel, V any](buildQuery models.BuildQueryFunc, buildResponse models.BuildResponseFunc[T, V]) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := models.Db(r)
			record, err := utils.DecodeAndValidate[T](r)

			if err != nil {
				return err
			}

			baseQuery := db.Model(&record)
			query := buildQuery(baseQuery, r)
			queryResult := query.Updates(&record)

			if queryResult.Error != nil {
				return queryResult.Error
			}

			if queryResult.RowsAffected == 0 {
				return utils.RecordNotFoundError{}
			}

			id := utils.IntPathValue(r, "id")
			db.First(&record, id)
			response := buildResponse(record)

			if err := utils.Encode(w, r, http.StatusOK, response); err != nil {
				return err
			}

			return nil
		})
}

// todo: restrict to only given query
func Delete[T utils.DbModel](buildQuery models.BuildQueryFunc) utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			db := models.Db(r)
			var record T

			query := buildQuery(db, r)
			queryResult := query.Delete(&record)

			if queryResult.Error != nil {
				return queryResult.Error
			}

			if queryResult.RowsAffected == 0 {
				return utils.RecordNotFoundError{}
			}

			if err := utils.Encode(w, r, http.StatusOK, utils.AnyMap{}); err != nil {
				return err
			}

			return nil
		})
}
