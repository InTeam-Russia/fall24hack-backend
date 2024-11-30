package ml

import "go.uber.org/zap"

type MockService struct {
	logger *zap.Logger
}

func NewMockService(logger *zap.Logger) Service {
	return &MockService{logger}
}

func (s *MockService) OnAnswer(text string) error {
	s.logger.Info("ml.MockService.OnAnswer called")
	return nil
}

func (s *MockService) OnQuestion(text string) error {
	s.logger.Info("ml.MockService.OnQuestion called")
	return nil
}
