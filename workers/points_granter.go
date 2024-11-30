package workers

import (
	"context"
	"reflect"
	"time"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
	args "github.com/filipio/athletics-backend/workerargs"
	"github.com/riverqueue/river"
	"gorm.io/gorm"
)

type PointsGranterWorker struct {
	river.WorkerDefaults[args.PointsGranterArgs]
}

func (w *PointsGranterWorker) Work(ctx context.Context, job *river.Job[args.PointsGranterArgs]) error {
	db := ctx.Value(utils.DbContextKey).(*gorm.DB)
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

	// // iterate over answers
	for _, answer := range answers {
		if reflect.DeepEqual(answer.Content.JSON, correctAnswer.JSON) {
			updateMap["points"] = question.Points
		} else {
			updateMap["points"] = 0
		}
		db.Model(&answer).Updates(updateMap)
	}

	return nil
}
