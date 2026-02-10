package workers

import (
	"context"
	"encoding/json"
	"math"
	"reflect"
	"time"

	"github.com/filipio/athletics-backend/config"
	"github.com/filipio/athletics-backend/models"
	args "github.com/filipio/athletics-backend/workerargs"
	"github.com/riverqueue/river"
)

type PointsGranterWorker struct {
	river.WorkerDefaults[args.PointsGranterArgs]
	deps *config.Dependencies
}

func NewPointsGranterWorker(deps *config.Dependencies) *PointsGranterWorker {
	return &PointsGranterWorker{deps: deps}
}

func (w *PointsGranterWorker) Work(ctx context.Context, job *river.Job[args.PointsGranterArgs]) error {
	db := w.deps.DB
	questionId := job.Args.QuestionID

	var question models.Question
	db.First(&question, questionId)
	if question.ID == 0 {
		return nil
	}

	if question.CorrectAnswer == nil {
		return nil
	}

	var answers []models.Answer
	db.Where("question_id = ?", questionId).Find(&answers)

	correctAnswer := question.CorrectAnswer
	updateMap := map[string]interface{}{
		"points_granted_at": time.Now().UTC(),
	}

	// iterate over answers
	for _, answer := range answers {
		isCorrect := false

		if question.Type == "numeric_value" {
			// Special handling for float comparison with epsilon tolerance
			var answerContent models.NumericValueAnswer
			var correctContent models.NumericValueAnswer

			if err := json.Unmarshal(answer.Content.JSON, &answerContent); err == nil {
				if err := json.Unmarshal(correctAnswer.JSON, &correctContent); err == nil {
					const epsilon = 0.001
					isCorrect = math.Abs(answerContent.Value-correctContent.Value) < epsilon
				}
			}
		} else {
			// Existing exact match comparison for other types
			isCorrect = reflect.DeepEqual(answer.Content.JSON, correctAnswer.JSON)
		}

		if isCorrect {
			updateMap["points"] = question.Points
		} else {
			updateMap["points"] = 0
		}
		db.Model(&answer).Updates(updateMap)
	}

	return nil
}
