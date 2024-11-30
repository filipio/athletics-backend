package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/filipio/athletics-backend/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type AnswerOfQuestion struct {
	datatypes.JSON
}

func (*AnswerOfQuestion) typeMapping() map[string]func() any {
	return map[string]func() any{
		"athlete": func() any { return &AthleteAnswer{} },
		"country": func() any { return &AthletesIdsAnswer{} },
	}
}

type AthleteAnswer struct {
	Value string `json:"value" validate:"required"`
}

type AthletesIdsAnswer struct {
	AthleteIdOne   uint `json:"athlete_id_one" validate:"required,id_of=athlete"`
	AthleteIdTwo   uint `json:"athlete_id_two" validate:"required,id_of=athlete"`
	AthleteIdThree uint `json:"athlete_id_three" validate:"required,id_of=athlete"`
}

func (a *AnswerOfQuestion) Validate(questionType string) error {
	constructor, ok := a.typeMapping()[questionType]
	if !ok {
		return errors.New("invalid question type: " + questionType)
	}

	answer := constructor()

	if err := json.Unmarshal(a.JSON, answer); err != nil {
		return errors.New("invalid json")
	}
	if err := utils.Validate(answer); err != nil {
		return errors.New("invalid json")
	}

	return nil
}

type Answer struct {
	AppModel
	UserID          uint             `json:"user_id" gorm:"not null" validate:"required"`
	QuestionID      uint             `json:"question_id" gorm:"not null" validate:"required,id_of=question"`
	Content         AnswerOfQuestion `json:"content" gorm:"not null" validate:"required"`
	Points          uint             `json:"points" gorm:"not null;default:0" validate:"eq=0"` // validation is set so it is not possible to set points manually
	PointsGrantedAt *time.Time       `json:"points_granted_at"`
}

func (m Answer) Validate(r *http.Request) error {
	currentUser := r.Context().Value(utils.UserContextKey).(User)

	if m.UserID != currentUser.ID {
		return utils.InvalidUserError{}
	}

	db := Db(r)

	var question Question
	db.First(&question, m.QuestionID)

	if err := m.Content.Validate(question.Type); err != nil {
		return utils.AppValidationError{
			FieldPath: "content",
			AppError:  utils.AppError{Message: err.Error()},
		}
	}

	var otherAnswer Answer
	db.Where("user_id = ? AND question_id = ?", m.UserID, m.QuestionID).First(&otherAnswer)
	if otherAnswer.ID != 0 {
		if r.Method == http.MethodPost {
			return utils.AppValidationError{
				FieldPath: "question_id",
				AppError:  utils.AppError{Message: "already answered by current user"},
			}
		} else if r.Method == http.MethodPut && otherAnswer.PointsGrantedAt != nil {
			return utils.AppValidationError{
				FieldPath: "question_id",
				AppError:  utils.AppError{Message: "points already granted"},
			}
		}
	}

	var event Event
	db.First(&event, question.EventID)

	return event.IsPresent()
}

func (m Answer) GetAllQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = onlyCurrentUserRecords(db, r)
	db = getByIds(db, r)

	queryParams := r.URL.Query()
	if queryParams.Has("question_id") {
		db = db.Where("question_id = ?", queryParams.Get("question_id"))
	}

	if queryParams.Has("user_id") {
		db = db.Where("user_id = ?", queryParams.Get("user_id"))
	}

	return db
}

func (m Answer) GetQuery(db *gorm.DB, r *http.Request) *gorm.DB {
	db = onlyCurrentUserRecords(db, r)
	db = GetByIdQuery(db, r)
	return db
}

func (m Answer) BuildResponse() any {
	return m
}

func onlyCurrentUserRecords(db *gorm.DB, r *http.Request) *gorm.DB {
	onlyForCurrentUser := r.Context().Value(utils.OnlyCurrentUserContextKey).(bool)
	if onlyForCurrentUser {
		currentUser := r.Context().Value(utils.UserContextKey).(User)
		return db.Where("user_id = ?", currentUser.ID)
	} else {
		return db
	}
}
