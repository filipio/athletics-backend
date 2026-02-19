package controllers

import (
	"net/http"
	"strconv"

	"github.com/filipio/athletics-backend/pkg/config"
	"github.com/filipio/athletics-backend/internal/models"
	"github.com/filipio/athletics-backend/pkg/httpio"
)

func PublishEvent(deps *config.Dependencies) httpio.HandlerWithError {
	return httpio.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		eventID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return httpio.AppValidationError{
				FieldPath: "id",
				AppError:  httpio.AppError{Message: "invalid event ID"},
			}
		}

		db := deps.DB

		var event models.Event
		if err := db.First(&event, eventID).Error; err != nil {
			return httpio.RecordNotFoundError{AppError: httpio.AppError{Message: "event not found"}}
		}

		if event.Status == "published" {
			httpio.Encode(w, r, http.StatusOK, event.BuildResponse())
			return nil
		}

		event.Status = "published"
		if err := db.Save(&event).Error; err != nil {
			return err
		}

		httpio.Encode(w, r, http.StatusOK, event.BuildResponse())
		return nil
	})
}

func UnpublishEvent(deps *config.Dependencies) httpio.HandlerWithError {
	return httpio.HandlerWithError(func(w http.ResponseWriter, r *http.Request) error {
		eventID, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			return httpio.AppValidationError{
				FieldPath: "id",
				AppError:  httpio.AppError{Message: "invalid event ID"},
			}
		}

		db := deps.DB

		var event models.Event
		if err := db.First(&event, eventID).Error; err != nil {
			return httpio.RecordNotFoundError{AppError: httpio.AppError{Message: "event not found"}}
		}

		if event.Status == "draft" {
			httpio.Encode(w, r, http.StatusOK, event.BuildResponse())
			return nil
		}

		event.Status = "draft"
		if err := db.Save(&event).Error; err != nil {
			return err
		}

		httpio.Encode(w, r, http.StatusOK, event.BuildResponse())
		return nil
	})
}
