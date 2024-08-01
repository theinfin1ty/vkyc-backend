package routes

import (
	"vkyc-backend/controllers"
	"vkyc-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoutes(app *gin.Engine, apiGroup *gin.RouterGroup) {
	apiGroup.POST("/user", middlewares.Auth([]string{"admin"}), controllers.AddUser)

	apiGroup.GET("/user", middlewares.Auth([]string{"admin"}), controllers.ListUsers)
	apiGroup.GET("/user/dashboard", middlewares.Auth([]string{"admin", "agent"}), controllers.Dashboard)
	apiGroup.GET("/user/:id", middlewares.Auth([]string{"admin"}), controllers.GetUser)

	apiGroup.PATCH("/user", middlewares.Auth([]string{"admin", "agent"}), controllers.UpdateProfile)
	apiGroup.PATCH("/user/:id", middlewares.Auth([]string{"admin"}), controllers.UpdateUser)
	apiGroup.PATCH("/user/:id/status", middlewares.Auth([]string{"admin"}), controllers.UpdateUserStatus)

	apiGroup.DELETE("/user/:id", middlewares.Auth([]string{"admin"}), controllers.DeleteUser)
}
