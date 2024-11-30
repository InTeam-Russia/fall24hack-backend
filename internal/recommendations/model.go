package recommendations

type UserResponse struct {
	Id                    int64  `json:"id"`
	FirstName             string `json:"firstName"`
	LastName              string `json:"lastName"`
	Username              string `json:"username"`
	Email                 string `json:"email"`
	TgLink                string `json:"tgLink"`
	OverlappingPercentage int    `json:"overlappingPercentage"`
}
