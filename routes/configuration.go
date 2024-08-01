package routes

import (
	"vkyc-backend/controllers"
	"vkyc-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func ConfigurationRoutes(app *gin.Engine, apiGroup *gin.RouterGroup) {
	apiGroup.GET("/config", middlewares.Auth([]string{"admin"}), controllers.GetConfiguration)
	apiGroup.GET("/config/public", controllers.GetPublicConfiguration)

	apiGroup.POST("/config", middlewares.Auth([]string{"admin"}), controllers.UpdateConfiguration)
}
