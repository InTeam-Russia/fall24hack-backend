package ml

type User struct {
	Id                    int64 `json:"id"`
	OverlappingPercentage int   `json:"overlapping_percentage"`
}

type OnNewQuestionResponse struct {
	Cluster int `json:"cluster"`
}
