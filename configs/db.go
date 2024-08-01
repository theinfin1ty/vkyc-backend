package configs

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() *gorm.DB {
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Etc/UTC", GetEnvVariable("DB_HOST"), GetEnvVariable("DB_USERNAME"), GetEnvVariable("DB_PASSWORD"), GetEnvVariable("DB_NAME"), GetEnvVariable("DB_PORT"))
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Silent),
	// })
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		GetEnvVariable("DB_USERNAME"),
		GetEnvVariable("DB_PASSWORD"),
		GetEnvVariable("DB_HOST"),
		GetEnvVariable("DB_PORT"),
		GetEnvVariable("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()

	if err != nil {
		panic(err)
	}

	err = sqlDB.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database")

	return db
}

var DB = ConnectDB()
