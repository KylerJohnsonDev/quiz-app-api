package categories

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/kylerjohnsondev/quiz-app-api/internal/utils"
)

type Service interface {
	FetchCategories(ctx context.Context) ([]Category, error)
}

type svc struct {
}

func NewService() Service {
	return &svc{}
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (s *svc) FetchCategories(ctx context.Context) ([]Category, error) {
	quizApiConfig := utils.GetQuizApiConfig()
	categoriesUrl, pathJoinError := url.JoinPath(quizApiConfig.BaseUrl, "categories")
	if pathJoinError != nil {
		log.Fatal(pathJoinError.Error())
	}

	params := url.Values{}
	params.Set("apiKey", quizApiConfig.ApiKey)

	parsedUrl, err := url.Parse(categoriesUrl)
	if err != nil {
		log.Fatal(err.Error())
	}

	parsedUrl.RawQuery = params.Encode()

	req, createRequestError := http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
	if createRequestError != nil {
		return nil, createRequestError
	}

	resp, requestError := http.DefaultClient.Do(req)
	if requestError != nil {
		return nil, requestError
	}
	defer resp.Body.Close()

	body, readAllError := io.ReadAll(resp.Body)
	if readAllError != nil {
		log.Fatal(readAllError.Error())
		return nil, readAllError
	}

	categories := []Category{}
	unmarshallingError := json.Unmarshal(body, &categories)
	if unmarshallingError != nil {
		return nil, unmarshallingError
	}

	return categories, nil
}
