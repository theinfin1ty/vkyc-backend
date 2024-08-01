package routes

import (
	"vkyc-backend/controllers"
	"vkyc-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func KYCRoutes(app *gin.Engine, apiGroup *gin.RouterGroup) {
	apiGroup.POST("/kyc", middlewares.Auth([]string{"agent"}), controllers.DocumentVerification)
	apiGroup.POST("/kyc/answer", middlewares.Auth([]string{"agent"}), controllers.SubmitAnswer)
	apiGroup.GET("/kyc/question", middlewares.Auth([]string{"agent"}), controllers.GetKYCQuestions)
	apiGroup.PATCH("/kyc/:meetingCode", middlewares.Auth([]string{"agent"}), controllers.CompleteKycProcess)
}
