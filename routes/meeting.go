package routes

import (
	"vkyc-backend/controllers"
	"vkyc-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func MeetingRoutes(app *gin.Engine, apiGroup *gin.RouterGroup) {
	apiGroup.POST("/meeting", middlewares.Auth([]string{"agent", "admin"}), controllers.CreateMeeting)
	apiGroup.POST("/meeting/join", controllers.JoinMeeting)
	apiGroup.POST("/meeting/agent/join", middlewares.Auth([]string{"agent"}), controllers.AgentJoinMeeting)
	apiGroup.POST("/meeting/recording/:meetingId", middlewares.Auth([]string{"agent", "admin"}), controllers.SaveMeetingRecording)

	apiGroup.GET("/meeting", middlewares.Auth([]string{"agent", "admin"}), controllers.ListMeetings)
	apiGroup.GET("/meeting/ws", controllers.MeetingEvents)
	apiGroup.GET("/meeting/:meetingId", middlewares.Auth([]string{"agent", "admin"}), controllers.GetMeeting)
	apiGroup.GET("/meeting/:meetingId/check", controllers.CheckMeetingCode)
	apiGroup.GET("/meeting/:meetingId/chat", controllers.GetMeetingChat)
	apiGroup.GET("/meeting/:meetingId/resend", middlewares.Auth([]string{"agent", "admin"}), controllers.ResendMeetingInvite)
	apiGroup.GET("/meeting/media/:mediaId", middlewares.Auth([]string{"agent", "admin"}), controllers.GetMeetingMedia)

	apiGroup.PATCH("/meeting/:meetingId", middlewares.Auth([]string{"agent", "admin"}), controllers.UpdateMeeting)
	apiGroup.PATCH("/meeting/:meetingId/reschedule", middlewares.Auth([]string{"agent", "admin"}), controllers.RescheduleMeeting)

	apiGroup.DELETE("/meeting/:meetingId", middlewares.Auth([]string{"agent", "admin"}), controllers.DeleteMeeting)
	apiGroup.DELETE("/meeting/media/:mediaId", middlewares.Auth([]string{"agent", "admin"}), controllers.DeleteMeetingMedia)
}
