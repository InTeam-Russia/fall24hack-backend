package polls

type Repo interface {
	GetUncompletedPolls(pageIndex int, pageSize int, userId int64) ([]Model, error)
}
