package controllers

import (
	"net/http"

	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"github.com/filipio/athletics-backend/workerargs"
	"gorm.io/gorm"
)

func CreateQuestion(deps *config.Dependencies) utils.HandlerWithError {
	return utils.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		db := deps.DB
		question, err := utils.DecodeAndValidate[models.Question](r, db)
		if err != nil {
			return err
		}

		// POST-specific validation: CorrectAnswer must be nil on creation
		if question.CorrectAnswer != nil {
			return utils.AppValidationError{
				FieldPath: "correct_answer",
				AppError:  utils.AppError{Message: "must be null for POST requests"},
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

		if err := utils.Encode(w, r, http.StatusCreated, response); err != nil {
			return err
		}

		return nil
	})
}

func UpdateQuestion(deps *config.Dependencies) utils.HandlerWithError {
	return utils.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		db := deps.DB
		question, err := utils.DecodeAndValidate[models.Question](r, db)
		if err != nil {
			return err
		}

		baseQuery := db.Model(&question)
		query := question.UpdateQuery(baseQuery, r)
		id := utils.IntPathValue(r, "id")

		if err := db.Transaction(func(tx *gorm.DB) error {
			queryResult := query.Updates(&question)
			if queryResult.Error != nil {
				return queryResult.Error
			}
			if queryResult.RowsAffected == 0 {
				return utils.RecordNotFoundError{}
			}

			// Business logic: Queue PointsGranter worker (moved from BeforeUpdateCtx hook)
			if _, err := deps.Workers.InsertTx(tx, workerargs.PointsGranterArgs{
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

		if err := utils.Encode(w, r, http.StatusOK, response); err != nil {
			return err
		}

		return nil
	})
}
