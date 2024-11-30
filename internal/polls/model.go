package polls

type PollType string

const (
	FREE  PollType = "FREE"
	RADIO PollType = "RADIO"
)

type Model struct {
	Id       int64    `json:"id"`
	Text     string   `json:"text"`
	Type     PollType `json:"type"`
	AuthorID int64    `json:"authorId"`
	Cluster  int      `json:"cluster"`
	Answers  []string `json:"answers"`
}

type CreateModel struct {
	Text     string   `json:"text" binding:"required"`
	Type     PollType `json:"type" binding:"required"`
	AuthorID int64    `json:"authorId" binding:"required"`
	Cluster  int      `json:"cluster" binding:"required"`
}

type OutModel struct {
	Id       int64    `json:"id"`
	Text     string   `json:"text"`
	Type     PollType `json:"type"`
	AuthorID int64    `json:"authorId"`
	Cluster  int      `json:"cluster"`
	Answers  []string `json:"answers"`
}
