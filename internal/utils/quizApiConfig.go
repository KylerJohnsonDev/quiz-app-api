package utils

import "os"

type QuizApiConfig struct {
	BaseUrl string
	ApiKey  string
}

func GetQuizApiConfig() QuizApiConfig {
	return QuizApiConfig{
		BaseUrl: "https://quizapi.io/api/v1/",
		ApiKey:  os.Getenv("QUIZ_APP_API_KEY"),
	}
}
