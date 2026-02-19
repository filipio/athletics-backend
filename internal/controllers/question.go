package controllers

import (
	"net/http"

	"github.com/filipio/athletics-backend/pkg/config"
	"github.com/filipio/athletics-backend/internal/models"
	"github.com/filipio/athletics-backend/pkg/httpio"
	"github.com/filipio/athletics-backend/internal/workers/args"
	"gorm.io/gorm"
)

func CreateQuestion(deps *config.Dependencies) httpio.HandlerWithError {
	return httpio.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		db := deps.DB
		question, err := httpio.DecodeAndValidate[models.Question](r, db)
		if err != nil {
			return err
		}

		// POST-specific validation: CorrectAnswer must be nil on creation
		if question.CorrectAnswer != nil {
			return httpio.AppValidationError{
				FieldPath: "correct_answer",
				AppError:  httpio.AppError{Message: "must be null for POST requests"},
			}
		}

		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&question).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}

		db.First(&question, question.GetID())
		response := question.BuildResponse()

		if err := httpio.Encode(w, r, http.StatusCreated, response); err != nil {
			return err
		}

		return nil
	})
}

func UpdateQuestion(deps *config.Dependencies) httpio.HandlerWithError {
	return httpio.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		db := deps.DB
		question, err := httpio.DecodeAndValidate[models.Question](r, db)
		if err != nil {
			return err
		}

		baseQuery := db.Model(&question)
		query := question.UpdateQuery(baseQuery, r)
		id := httpio.IntPathValue(r, "id")

		if err := db.Transaction(func(tx *gorm.DB) error {
			queryResult := query.Updates(&question)
			if queryResult.Error != nil {
				return queryResult.Error
			}
			if queryResult.RowsAffected == 0 {
				return httpio.RecordNotFoundError{}
			}

			// Business logic: Queue PointsGranter worker (moved from BeforeUpdateCtx hook)
			if _, err := deps.Workers.InsertTx(tx, args.PointsGranterArgs{
				QuestionID: uint(id),
			}); err != nil {
				return err
			}

			return nil
		}); err != nil {
			return err
		}

		db.First(&question, id)
		response := question.BuildResponse()

		if err := httpio.Encode(w, r, http.StatusOK, response); err != nil {
			return err
		}

		return nil
	})
}
