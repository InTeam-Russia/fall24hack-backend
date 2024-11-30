package ml

type Service interface {
	OnAnswer(userId int64, text string) error
	OnQuestion(text string) error
	OnCreateUser(userId int64) error
}
