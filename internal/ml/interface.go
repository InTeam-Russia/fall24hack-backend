package ml

type SearchType string

const (
	CODIRECTIONAL SearchType = "codirectional"
	OPPOSITE      SearchType = "opposite"
)

type Cluster = int

type Service interface {
	OnAnswer(userId int64, text string) error
	OnQuestion(text string) (Cluster, error)
	OnCreateUser(userId int64) error
	UsersANN(userId int64, neighboursCount int, searchType SearchType) ([]User, error)
}
