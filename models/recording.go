package models

import "time"

type Recording struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	MeetingID uint      `json:"meetingID"`
	Meeting   *Meeting  `json:"meeting,omitempty"`
	Recording string    `json:"recording"`
}
