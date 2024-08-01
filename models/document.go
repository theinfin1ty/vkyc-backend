package models

import (
	"time"
)

type Document struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	Type           string    `json:"type"`
	Name           *string   `json:"name"`
	APIResponse    *string   `gorm:"type:longtext" json:"apiResponse"`
	APIStatus      *string   `json:"apiStatus"`
	Image          string    `json:"image"`
	Score          *float64  `json:"score"`
	MeetingID      uint      `json:"meetingID"`
	Meeting        *Meeting  `json:"meeting,omitempty"`
	VerificationID *string   `json:"verificationID"`
}
