package ml

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

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

func (s *APIService) OnAnswer(userId int64, text string) error {
	s.logger.Debug("ml.APIService.OnAnswer called")

	relativeURL, err := url.Parse("/on_new_answer")
	if err != nil {
		return err
	}

	jsonData := fmt.Sprintf(`{"text": "%s", "user_id": %d}`, text, userId)
	url := s.baseURL.ResolveReference(relativeURL)

	s.logger.Debug(url.String())

	_, err = http.Post(url.String(), "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return err
	}

	s.logger.Debug("ml.APIService.OnAnswer http")

	return nil
}

func (s *APIService) OnQuestion(text string) (Cluster, error) {
	s.logger.Debug("ml.APIService.OnQuestion called")

	relativeURL, err := url.Parse("/on_new_question")
	if err != nil {
		return 0, err
	}

	jsonData := fmt.Sprintf(`{"text": "%s"}`, text)
	url := s.baseURL.ResolveReference(relativeURL)

	s.logger.Debug(url.String())

	res, err := http.Post(url.String(), "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	s.logger.Debug("ml.APIService.OnQuestion http")

	var data OnNewQuestionResponse
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Cluster, nil
}

func (s *APIService) OnCreateUser(userId int64) error {
	s.logger.Debug("ml.APIService.OnCreateUser called")

	relativeURL, err := url.Parse("/on_register_user")
	if err != nil {
		return err
	}

	url := s.baseURL.ResolveReference(relativeURL)

	q := url.Query()
	q.Add("user_id", strconv.FormatInt(userId, 10))
	url.RawQuery = q.Encode()

	s.logger.Debug(url.String())

	_, err = http.Post(url.String(), "application/json", bytes.NewBuffer(make([]byte, 0)))
	if err != nil {
		return err
	}

	s.logger.Debug("ml.APIService.OnCreateUser http")

	return nil
}

func (s *APIService) UsersANN(userId int64, neighboursCount int, searchType SearchType) ([]User, error) {
	s.logger.Debug("ml.APIService.UsersANN called")

	relativeURL, err := url.Parse("/users_ann")
	if err != nil {
		return nil, err
	}

	url := s.baseURL.ResolveReference(relativeURL)
	s.logger.Debug(url.String())

	jsonData := fmt.Sprintf(`{
		"user_id": %d,
		"k": %d,
		"search_type": "%s"
	}`, userId, neighboursCount, string(searchType))

	res, err := http.Post(url.String(), "application/json", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		return nil, err
	}

	var data []User
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}
