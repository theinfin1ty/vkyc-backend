package models

import (
	"time"
)

type User struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	FirstName   string     `json:"firstName"`
	LastName    string     `json:"lastName"`
	Email       string     `gorm:"unique" json:"email"`
	CountryCode string     `json:"countryCode"`
	Phone       string     `json:"phone"`
	Password    string     `json:"-"`
	Role        string     `gorm:"default:agent" json:"role"`
	Meetings    *[]Meeting `gorm:"foreignKey:AgentID" json:"meetings,omitempty"`
	Status      string     `gorm:"default:active" json:"status"`
	DOB         *time.Time `json:"dob"`
	Department  *string    `json:"department"`
	OTPs        *[]OTP     `gorm:"foreignKey:UserID" json:"otps,omitempty"`
	IsDeleted   bool       `gorm:"default:false" json:"isDeleted"`
}
