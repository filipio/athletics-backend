package controllers

import (
	"net/http"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
)

type RankingItem struct {
	ID          uint   `json:"user_id"`
	Email       string `json:"email"`
	TotalPoints int    `json:"total_points"`
}

func GetRanking() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			query := models.Db(r).
				Model(&models.Answer{}).
				Joins("JOIN users ON answers.user_id = users.id")

			queryParams := r.URL.Query()
			if queryParams.Has("event_id") {
				query = query.Joins("JOIN questions ON answers.question_id = questions.id").
					Where("questions.event_id = ?", queryParams.Get("event_id"))
			}

			query = query.Group("users.id, users.email").
				Select("users.id, users.email, sum(answers.points) as total_points")

			orderBy := "total_points"
			orderDirection := "desc"
			paginatedResponse, err := models.Paginate[RankingItem](query, r, &utils.PaginationParams{
				OrderBy:        orderBy,
				OrderDirection: orderDirection,
			}, nil, &models.Answer{})

			if err != nil {
				return err
			}

			if err := utils.Encode(w, r, http.StatusOK, paginatedResponse); err != nil {
				return err
			}

			return nil
		})
}
