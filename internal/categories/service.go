package categories

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/kylerjohnsondev/quiz-app-api/internal/utils"
)

const cacheTTL = 24 * time.Hour

type Service interface {
	FetchCategories(ctx context.Context) ([]Category, error)
}

type svc struct {
	mu       sync.RWMutex
	cached   []Category
	cacheExp time.Time
}

func NewService() Service {
	return &svc{}
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (s *svc) FetchCategories(ctx context.Context) ([]Category, error) {
	// Return cached result if still valid.
	s.mu.RLock()
	if time.Now().Before(s.cacheExp) && len(s.cached) > 0 {
		result := make([]Category, len(s.cached))
		copy(result, s.cached)
		s.mu.RUnlock()
		return result, nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()
	// Double-check after acquiring write lock (another goroutine may have refreshed).
	if time.Now().Before(s.cacheExp) && len(s.cached) > 0 {
		result := make([]Category, len(s.cached))
		copy(result, s.cached)
		return result, nil
	}

	quizApiConfig := utils.GetQuizApiConfig()
	categoriesUrl, pathJoinError := url.JoinPath(quizApiConfig.BaseUrl, "categories")
	if pathJoinError != nil {
		log.Print(pathJoinError)
		return nil, pathJoinError
	}

	parsedUrl, err := url.Parse(categoriesUrl)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	req, createRequestError := http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
	if createRequestError != nil {
		return nil, createRequestError
	}
	req.Header.Set("X-Api-Key", quizApiConfig.ApiKey)

	resp, requestError := http.DefaultClient.Do(req)
	if requestError != nil {
		return nil, requestError
	}
	defer resp.Body.Close()

	body, readAllError := io.ReadAll(resp.Body)
	if readAllError != nil {
		log.Print(readAllError)
		return nil, readAllError
	}

	categories := []Category{}
	unmarshallingError := json.Unmarshal(body, &categories)
	if unmarshallingError != nil {
		return nil, unmarshallingError
	}

	s.cached = categories
	s.cacheExp = time.Now().Add(cacheTTL)
	return categories, nil
}
