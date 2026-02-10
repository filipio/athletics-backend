package controllers

import (
	"net/http"

	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	"gorm.io/gorm"
)

func CreateAnswer(deps *config.Dependencies) utils.HandlerWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		db := deps.DB
		answer, err := utils.DecodeAndValidate[models.Answer](r, db)
		if err != nil {
			return err
		}

		// HTTP-level validation: current user must match answer.UserID
		currentUser := r.Context().Value(utils.UserContextKey).(models.User)
		if answer.UserID != currentUser.ID {
			return utils.InvalidUserError{}
		}

		// Check for duplicate answer (user already answered this question)
		var existingAnswer models.Answer
		db.Where("user_id = ? AND question_id = ?", answer.UserID, answer.QuestionID).First(&existingAnswer)
		if existingAnswer.ID != 0 {
			return utils.AppValidationError{
				FieldPath: "question_id",
				AppError:  utils.AppError{Message: "already answered by current user"},
			}
		}

		// Create the answer
		if err := db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&answer).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			return err
		}

		db.First(&answer, answer.GetID())
		response := answer.BuildResponse()

		if err := utils.Encode(w, r, http.StatusCreated, response); err != nil {
			return err
		}

		return nil
	}
}

func UpdateAnswer(deps *config.Dependencies) utils.HandlerWithError {
	return func(w http.ResponseWriter, r *http.Request) error {
		db := deps.DB
		answer, err := utils.DecodeAndValidate[models.Answer](r, db)
		if err != nil {
			return err
		}

		// HTTP-level validation: current user must match answer.UserID
		currentUser := r.Context().Value(utils.UserContextKey).(models.User)
		if answer.UserID != currentUser.ID {
			return utils.InvalidUserError{}
		}

		// Fetch existing answer to check if points already granted
		id := utils.IntPathValue(r, "id")
		var existingAnswer models.Answer
		db.First(&existingAnswer, id)
		if existingAnswer.ID == 0 {
			return utils.RecordNotFoundError{}
		}

		// Check if points already granted
		if existingAnswer.PointsGrantedAt != nil {
			return utils.AppValidationError{
				FieldPath: "question_id",
				AppError:  utils.AppError{Message: "points already granted"},
			}
		}

		// Update the answer
		baseQuery := db.Model(&answer)
		query := answer.UpdateQuery(baseQuery, r)

		if err := db.Transaction(func(tx *gorm.DB) error {
			queryResult := query.Updates(&answer)

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

		db.First(&answer, id)
		response := answer.BuildResponse()

		if err := utils.Encode(w, r, http.StatusOK, response); err != nil {
			return err
		}

		return nil
	}
}
