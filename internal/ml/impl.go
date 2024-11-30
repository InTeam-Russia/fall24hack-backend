package ml

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"

	"go.uber.org/zap"
)

type APIService struct {
	logger  *zap.Logger
	baseURL *url.URL
}

func NewAPIService(logger *zap.Logger, baseURL string) Service {
	url, err := url.Parse(baseURL)
	if err != nil {
		panic(err)
	}

	return &APIService{logger: logger, baseURL: url}
}

func (s *APIService) OnAnswer(text string) error {
	s.logger.Debug("ml.APIService.OnAnswer called")

	relativeURL, err := url.Parse("/on_new_answer")
	if err != nil {
		return err
	}

	jsonData := fmt.Sprintf(`{"text": "%s"}`, text)
	url := s.baseURL.ResolveReference(relativeURL)

	_, err = http.Post(url.String(), "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return err
	}

	s.logger.Debug("ml.APIService.OnAnswer http")

	return nil
}

func (s *APIService) OnQuestion(text string) error {
	s.logger.Debug("ml.APIService.OnQuestion called")

	relativeURL, err := url.Parse("/on_new_question")
	if err != nil {
		return err
	}

	jsonData := fmt.Sprintf(`{"text": "%s"}`, text)
	url := s.baseURL.ResolveReference(relativeURL)

	_, err = http.Post(url.String(), "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return err
	}

	return nil
}
