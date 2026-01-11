package controllers

import (
	"net/http"

	"github.com/filipio/athletics-backend/models"
	"github.com/filipio/athletics-backend/utils"
)

type RankingItem struct {
	Username    string `json:"username"`
	TotalPoints int    `json:"total_points"`
}

type MyRankingResponse struct {
	Place             *int `json:"place"`
	TotalPlaces       int  `json:"total_places"`
	TotalPointsScored int  `json:"total_points_scored"`
	TotalPoints       int  `json:"total_points"`
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

			query = query.Group("users.id, users.username").
				Select("users.username, sum(answers.points) as total_points")

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

func GetMyRanking() utils.HandlerWithError {
	return utils.HandlerWithError(
		func(w http.ResponseWriter, r *http.Request) error {
			currentUser := r.Context().Value(utils.UserContextKey).(models.User)
			queryParams := r.URL.Query()

			var userPointsResult struct {
				TotalPoints int
			}

			userPointsQuery := models.Db(r).
				Model(&models.Answer{}).
				Select("COALESCE(SUM(answers.points), 0) as total_points").
				Where("user_id = ?", currentUser.ID)

			if queryParams.Has("event_id") {
				userPointsQuery = userPointsQuery.
					Joins("JOIN questions ON answers.question_id = questions.id").
					Where("questions.event_id = ?", queryParams.Get("event_id"))
			}

			if err := userPointsQuery.Scan(&userPointsResult).Error; err != nil {
				return err
			}

			var place *int
			if userPointsResult.TotalPoints > 0 {
				var usersWithHigherPoints int64

				countSubquery := models.Db(r).
					Model(&models.Answer{}).
					Select("user_id")

				if queryParams.Has("event_id") {
					countSubquery = countSubquery.
						Joins("JOIN questions ON answers.question_id = questions.id").
						Where("questions.event_id = ?", queryParams.Get("event_id"))
				}

				countSubquery = countSubquery.
					Group("user_id").
					Having("SUM(answers.points) > ?", userPointsResult.TotalPoints)

				if err := countSubquery.Count(&usersWithHigherPoints).Error; err != nil {
					return err
				}

				placeValue := int(usersWithHigherPoints) + 1
				place = &placeValue
			}

			var totalPlaces int64

			totalPlacesQuery := models.Db(r).
				Model(&models.Answer{}).
				Distinct("user_id")

			if queryParams.Has("event_id") {
				totalPlacesQuery = totalPlacesQuery.
					Joins("JOIN questions ON answers.question_id = questions.id").
					Where("questions.event_id = ?", queryParams.Get("event_id"))
			}

			if err := totalPlacesQuery.Count(&totalPlaces).Error; err != nil {
				return err
			}

			var totalAvailablePoints struct {
				TotalPoints int
			}

			totalPointsQuery := models.Db(r).
				Model(&models.Question{}).
				Select("COALESCE(SUM(questions.points), 0) as total_points")

			if queryParams.Has("event_id") {
				totalPointsQuery = totalPointsQuery.
					Where("event_id = ?", queryParams.Get("event_id"))
			}

			if err := totalPointsQuery.Scan(&totalAvailablePoints).Error; err != nil {
				return err
			}

			response := MyRankingResponse{
				Place:             place,
				TotalPlaces:       int(totalPlaces),
				TotalPointsScored: userPointsResult.TotalPoints,
				TotalPoints:       totalAvailablePoints.TotalPoints,
			}

			if err := utils.Encode(w, r, http.StatusOK, response); err != nil {
				return err
			}

			return nil
		})
}
