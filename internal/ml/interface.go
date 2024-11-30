package ml

type Service interface {
	OnAnswer(text string) error
	OnQuestion(text string) error
}
