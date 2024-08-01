package models

import (
	"time"
)

type Chat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Message   string    `json:"message"`
	Sender    string    `json:"sender"`
	MeetingID uint      `json:"meetingID"`
	Meeting   *Meeting  `json:"meeting,omitempty"`
	IsAgent   bool      `gorm:"default:false" json:"isAgent"`
}
