package ml

import "go.uber.org/zap"

type MockService struct {
	logger *zap.Logger
}

func NewMockService(logger *zap.Logger) Service {
	return &MockService{logger}
}

func (s *MockService) OnAnswer(userId int64, text string) error {
	s.logger.Info("ml.MockService.OnAnswer called")
	return nil
}

func (s *MockService) OnQuestion(text string) error {
	s.logger.Info("ml.MockService.OnQuestion called")
	return nil
}

func (s *MockService) OnCreateUser(userId int64) error {
	s.logger.Info("ml.MockService.OnCreateUser called")
	return nil
}
