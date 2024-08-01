package models

import "time"

type Answer struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
	MeetingID  uint      `json:"meetingID"`
	Meeting    *Meeting  `json:"meeting,omitempty"`
	QuestionID uint      `json:"questionID"`
	Question   *Question `json:"question,omitempty"`
	Answer     string    `json:"answer"`
}
