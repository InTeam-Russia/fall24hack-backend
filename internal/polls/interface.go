package polls

type Repo interface {
	GetUncompletedPolls(pageIndex int, pageSize int, userId int64) ([]Model, error)
	AddAnswer(userId int64, pollId int64, text string) error
}
