package controllers

import (
	"net/http"
	"strconv"

	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
)

func PublishEvent(deps *config.Dependencies) utils.HandlerWithError {
	return utils.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		eventID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return utils.AppValidationError{
				FieldPath: "id",
				AppError:  utils.AppError{Message: "invalid event ID"},
			}
		}

		db := deps.DB

		var event models.Event
		if err := db.First(&event, eventID).Error; err != nil {
			return utils.RecordNotFoundError{AppError: utils.AppError{Message: "event not found"}}
		}

		if event.Status == "published" {
			utils.Encode(w, r, http.StatusOK, event.BuildResponse())
			return nil
		}

		event.Status = "published"
		if err := db.Save(&event).Error; err != nil {
			return err
		}

		utils.Encode(w, r, http.StatusOK, event.BuildResponse())
		return nil
	})
}

func UnpublishEvent(deps *config.Dependencies) utils.HandlerWithError {
	return utils.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		eventID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return utils.AppValidationError{
				FieldPath: "id",
				AppError:  utils.AppError{Message: "invalid event ID"},
			}
		}

		db := deps.DB

		var event models.Event
		if err := db.First(&event, eventID).Error; err != nil {
			return utils.RecordNotFoundError{AppError: utils.AppError{Message: "event not found"}}
		}

		if event.Status == "draft" {
			utils.Encode(w, r, http.StatusOK, event.BuildResponse())
			return nil
		}

		event.Status = "draft"
		if err := db.Save(&event).Error; err != nil {
			return err
		}

		utils.Encode(w, r, http.StatusOK, event.BuildResponse())
		return nil
	})
}
