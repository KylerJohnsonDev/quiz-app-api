package questions

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
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
	FetchQuestions(ctx context.Context) ([]Question, error)
}

// Because *svc has a FetchQuestions method with the same signature as Service.FetchQuestions,
// *svc implements Service. No function fields or assignments are needed.
type svc struct {
}

func NewService() Service {
	return &svc{}
}

const QUIZ_API_URL string = "https://quizapi.io/api/v1/questions"

func (s *svc) FetchQuestions(ctx context.Context) ([]Question, error) {
	api_key := os.Getenv("QUIZ_APP_API_KEY")

	params := url.Values{}
	params.Set("apiKey", api_key)
	params.Set("limit", "10")

	parsedUrl, err0 := url.Parse(QUIZ_API_URL)
	if err0 != nil {
		log.Fatal(err0.Error())
		return nil, err0
	}

	parsedUrl.RawQuery = params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		log.Fatal(err2.Error())
		return nil, err2
	}

	questions := []Question{}
	err3 := json.Unmarshal(body, &questions)
	if err3 != nil {
		log.Fatal(err3.Error())
		return nil, err3
	}
	return questions, nil
}
