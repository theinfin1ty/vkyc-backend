package migrations

import (
	"fmt"
	"vkyc-backend/configs"
	"vkyc-backend/models"
)

func RunMigrations() {
	err := configs.DB.AutoMigrate(
		&models.User{},
		&models.Meeting{},
		&models.Chat{},
		&models.Document{},
		&models.OTP{},
		&models.Question{},
		&models.Recording{},
		&models.Configuration{},
		&models.Answer{},
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Database migrations completed")
}
