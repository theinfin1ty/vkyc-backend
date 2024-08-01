package routes

import (
	"vkyc-backend/controllers"
	"vkyc-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func QuestionRoutes(app *gin.Engine, apiGroup *gin.RouterGroup) {
	apiGroup.POST("/question", middlewares.Auth([]string{"admin"}), controllers.CreateQuestion)

	apiGroup.GET("/question", middlewares.Auth([]string{"admin"}), controllers.ListQuestions)
	apiGroup.GET("/question/:questionId", middlewares.Auth([]string{"admin"}), controllers.GetQuestion)

	apiGroup.PATCH("/question/:questionId", middlewares.Auth([]string{"admin"}), controllers.UpdateQuestion)

	apiGroup.DELETE("/question/:questionId", middlewares.Auth([]string{"admin"}), controllers.DeleteQuestion)
}
