package models

import "time"

type Configuration struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Name        *string   `json:"name"`
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description *string   `json:"description"`
	Status      string    `gorm:"default:active" json:"status"`
}
