package types

type AnswerInput struct {
	MeetingCode string `json:"meetingCode"`
	QuestionID  uint   `json:"questionID"`
	Answer      string `json:"answer"`
}
