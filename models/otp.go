package models

import "time"

type OTP struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	ExpiresAt time.Time `json:"expiresAt"`
	OTP       string    `json:"otp"`
	UserID    uint      `json:"userID"`
}
