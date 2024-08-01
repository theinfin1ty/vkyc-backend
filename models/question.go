package models

import "time"

type Question struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Question  string    `json:"question"`
	IsDeleted bool      `gorm:"default:false" json:"isDeleted"`
	Status    string    `gorm:"default:active" json:"status"`
}
