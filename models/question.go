package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/utils"
	"github.com/filipio/athletics-backend/workerargs"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Question struct {
	AppModel
	EventID       uint            `json:"event_id" gorm:"not null" validate:"required,id_of=event"`
	Content       string          `json:"content" gorm:"not null" validate:"required"`
	CorrectAnswer *datatypes.JSON `json:"correct_answer"`
	Type          string          `json:"type" gorm:"not null" validate:"oneof=athlete country"`
	Answers       []Answer        `json:"answers,omitempty" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
}

func (m Question) Validate(r *http.Request) error {
	fmt.Println("questions validation is executed")

	db := Db(r)
	var event Event
	db.First(&event, m.EventID)

	if event.Deadline.Before(time.Now().UTC()) {
		return utils.AppValidationError{
			FieldPath: "event_id",
			AppError:  utils.AppError{Message: "event is already closed"},
		}
	}
	return nil
}

func GetQuestionsQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	queryFunctions := []func(db *gorm.DB) *gorm.DB{getByIds(r)}
	queryParams := r.URL.Query()

	if queryParams.Has("event_id") {
		eventId := queryParams.Get("event_id")
		queryFunctions = append(queryFunctions, func(db *gorm.DB) *gorm.DB {
			return db.Where("event_id = ?", eventId)
		})
	}

	return db.Scopes(queryFunctions...)
}

func (m Question) BeforeUpdateCtx(ctx context.Context, tx *gorm.DB) error {
	workersClient := ctx.Value(utils.WorkersContextKey).(*config.InsertWorkerClient)
	if _, err := workersClient.InsertTx(tx, workerargs.SortArgs{Strings: []string{"h", "l", "o", "a"}}); err != nil {
		return err
	}
	return nil
}
