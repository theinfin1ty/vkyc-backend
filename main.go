package main

import (
	"fmt"
	"vkyc-backend/configs"
	"vkyc-backend/middlewares"
	"vkyc-backend/migrations"
	"vkyc-backend/routes"
	"vkyc-backend/seeds"
	"vkyc-backend/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	if configs.GetEnvVariable("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}

	app.Use(cors.New(corsConfig))
	app.Use(middlewares.RateLimiter())

	migrations.RunMigrations()

	seeds.CreateDefaultAdmin()

	utils.StartScheduler()

	routes.Routes(app)

	fmt.Printf("Server is running on port %s \n", configs.GetEnvVariable("PORT"))

	err := app.Run("0.0.0.0:" + configs.GetEnvVariable("PORT"))

	if err != nil {
		fmt.Println(err)
	}
}
