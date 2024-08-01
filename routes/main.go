package routes

import "github.com/gin-gonic/gin"

func Routes(app *gin.Engine) {
	apiGroup := app.Group("/api")
	AuthRoutes(app, apiGroup)
	MeetingRoutes(app, apiGroup)
	UserRoutes(app, apiGroup)
	KYCRoutes(app, apiGroup)
	QuestionRoutes(app, apiGroup)
	ConfigurationRoutes(app, apiGroup)
}
