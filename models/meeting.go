package models

import (
	"time"
)

type Meeting struct {
	ID               uint         `gorm:"primarykey" json:"id"`
	CreatedAt        time.Time    `json:"createdAt"`
	UpdatedAt        time.Time    `json:"updatedAt"`
	Title            string       `json:"title"`
	FirstName        string       `json:"firstName"`
	LastName         string       `json:"lastName"`
	Email            string       `json:"email"`
	CountryCode      string       `json:"countryCode"`
	Phone            string       `json:"phone"`
	MeetingCode      string       `gorm:"unique" json:"meetingCode"`
	Remark           string       `json:"remark"`
	ScheduleDateTime time.Time    `json:"scheduleDateTime"`
	Password         string       `json:"password"`
	AgentID          uint         `json:"agentId"`
	Agent            *User        `json:"agent,omitempty"`
	Chat             *[]Chat      `gorm:"foreignKey:MeetingID" json:"chat,omitempty"`
	IsEnded          bool         `gorm:"default:false" json:"isEnded"`
	IsStarted        bool         `gorm:"default:false" json:"isStarted"`
	Recordings       *[]Recording `gorm:"foreignKey:MeetingID" json:"recording,omitempty"`
	KycStatus        *string      `json:"kycStatus"`
	Location         *string      `json:"location"`
	Documents        *[]Document  `gorm:"foreignKey:MeetingID" json:"documents,omitempty"`
	IsDeleted        bool         `gorm:"default:false" json:"isDeleted"`
	Answers          *[]Answer    `gorm:"foreignKey:MeetingID" json:"answers,omitempty"`
	Remarks          *string      `json:"remarks"`
}
