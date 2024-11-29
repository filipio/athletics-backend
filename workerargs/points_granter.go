package workerargs

type PointsGranterArgs struct {
	QuestionID uint `json:"id"`
}

func (PointsGranterArgs) Kind() string { return "points_granter" }
