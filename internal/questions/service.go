package questions

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/kylerjohnsondev/quiz-app-api/internal/httperror"
	"github.com/kylerjohnsondev/quiz-app-api/internal/utils"
)

type Answers struct {
	AnswerA string `json:"answer_a"`
	AnswerB string `json:"answer_b"`
	AnswerC string `json:"answer_c"`
	AnswerD string `json:"answer_d"`
	AnswerE string `json:"answer_e"`
	AnswerF string `json:"answer_f"`
}

type CorrectAnswers struct {
	AnswerACorrect string `json:"answer_a_correct"`
	AnswerBCorrect string `json:"answer_b_correct"`
	AnswerCCorrect string `json:"answer_c_correct"`
	AnswerDCorrect string `json:"answer_d_correct"`
	AnswerECorrect string `json:"answer_e_correct"`
	AnswerFCorrect string `json:"answer_f_correct"`
}

type Question struct {
	ID                     int            `json:"id"`
	Question               string         `json:"question"`
	Description            string         `json:"description"`
	MultipleCorrectAnswers string         `json:"multiple_correct_answers"`
	Answers                Answers        `json:"answers"`
	CorrectAnswers         CorrectAnswers `json:"correct_answers"`
	Explanation            string         `json:"explanation"`
	Category               string         `json:"category"`
	Difficulty             string         `json:"difficulty"`
}

type Service interface {
	FetchQuestions(ctx context.Context, category string, difficulty string, limit string) ([]Question, error)
}

// Because *svc has a FetchQuestions method with the same signature as Service.FetchQuestions,
// *svc implements Service. No function fields or assignments are needed.
type svc struct {
}

func NewService() Service {
	return &svc{}
}

func (s *svc) FetchQuestions(ctx context.Context, category string, difficulty string, limit string) ([]Question, error) {
	quizApiConfig := utils.GetQuizApiConfig()

	questionsUrl, joinPathError := url.JoinPath(quizApiConfig.BaseUrl, "questions")
	if joinPathError != nil {
		log.Print(joinPathError)
		return nil, joinPathError
	}

	params := url.Values{}
	if category != "" {
		params.Set("category", category)
	}

	if difficulty != "" {
		difficultyLowercase := strings.ToLower(difficulty)
		if difficultyLowercase == "easy" || difficultyLowercase == "medium" || difficultyLowercase == "hard" {
			params.Set("difficulty", difficultyLowercase)
		} else {
			logMessage := fmt.Sprintf("Query parameter not one of 'easy', 'medium', 'hard'. Received %s.", difficulty)
			slog.Warn(logMessage)
		}
	}

	limitInt, convLimitParamErr := strconv.ParseUint(limit, 0, 8) // max 255
	if convLimitParamErr != nil || limitInt == 0 {
		// default to 10 if limit not provided
		params.Set("limit", "10")
	} else {
		params.Set("limit", limit)
	}

	parsedUrl, err0 := url.Parse(questionsUrl)
	if err0 != nil {
		log.Print(err0)
		return nil, err0
	}

	parsedUrl.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", quizApiConfig.ApiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		log.Print(err2)
		return nil, err2
	}

	if resp.StatusCode >= 400 {
		return nil, &httperror.HTTPError{
			StatusCode:  resp.StatusCode,
			Body:        body,
			ContentType: resp.Header.Get("Content-Type"),
		}
	}

	questions := []Question{}
	err3 := json.Unmarshal(body, &questions)
	if err3 != nil {
		log.Print(err3)
		return nil, err3
	}
	return questions, nil
}
