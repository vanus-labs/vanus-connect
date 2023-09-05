package internal

type EventData struct {
	QuestionUser   string `json:"question_user"`
	QuestionAtUser string `json:"question_at_user"`
	Question       string `json:"question"`
	Answer         string `json:"answer"`
	AnswerUser     string `json:"answer_user"`
}
