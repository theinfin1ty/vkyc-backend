package routes

import (
	"vkyc-backend/controllers"
	"vkyc-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(app *gin.Engine, apiGroup *gin.RouterGroup) {
	apiGroup.POST("/auth/login", controllers.Login)
	apiGroup.POST("/auth/forgot-password", controllers.ForgotPassword)
	apiGroup.POST("/auth/reset-password", controllers.ResetPassword)
	apiGroup.POST("/auth/change-password", middlewares.Auth([]string{"agent", "admin"}), controllers.ChangePassword)
}
