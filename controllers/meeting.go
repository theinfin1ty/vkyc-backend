package controllers

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"vkyc-backend/configs"
	"vkyc-backend/models"
	"vkyc-backend/types"
	"vkyc-backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateMeeting(c *gin.Context) {
	var body types.MeetingInput
	var agent models.User
	var existingMeeting models.Meeting
	var meetingTimeGapConfig models.Configuration

	authUser, _ := c.Get("authUser")

	user := authUser.(models.User)

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	if user.Role == utils.AdminRole {
		result := configs.DB.First(&agent, models.User{
			ID:        uint(*body.AgentID),
			Status:    utils.ActiveStatus,
			IsDeleted: false,
		})

		if utils.CheckError(result.Error, http.StatusNotFound, c, "Agent Not Found") {
			return
		}
	} else {
		agent = user
	}

	scheduleDateTime, err := time.Parse(utils.DateLayout, *body.ScheduleDateTime)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	meetingTimeGap := 60

	result := configs.DB.First(&meetingTimeGapConfig, models.Configuration{
		Key:    utils.MeetingTimeGap,
		Status: utils.ActiveStatus,
	})

	if result.Error != nil {
		meetingTimeGap = 60
	}

	if meetingTimeGapConfig.Value != "" {
		meetingTimeGap, err = strconv.Atoi(meetingTimeGapConfig.Value)

		if err != nil || meetingTimeGap <= 0 {
			meetingTimeGap = 60
		}
	}

	fromCheckTime := scheduleDateTime.Add(time.Minute * -time.Duration(meetingTimeGap-1))

	result = configs.DB.
		Model(&models.Meeting{}).
		Where("agent_id = ?", agent.ID).
		Where("schedule_date_time > ?", fromCheckTime).
		Where("schedule_date_time < ?", scheduleDateTime).
		Where("is_deleted = ?", false).
		Find(&existingMeeting)

	if result.RowsAffected > 0 {
		msg := fmt.Sprintf("Meetings can only be scheduled with %d minutes gap", meetingTimeGap)
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse(msg))
		return
	}

	meeting := models.Meeting{
		Title:            body.Title,
		FirstName:        body.FirstName,
		LastName:         body.LastName,
		Email:            body.Email,
		CountryCode:      body.CountryCode,
		Phone:            body.Phone,
		MeetingCode:      uuid.NewString(),
		Remark:           body.Remark,
		ScheduleDateTime: scheduleDateTime,
		Password:         utils.GeneratePassword(4, false),
		IsDeleted:        false,
		AgentID:          agent.ID,
	}

	result = configs.DB.Create(&meeting)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	err = utils.SendMeetingEmail(&meeting, &agent)

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	err = utils.SendMeetingSMS(&meeting, &agent)

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"meeting": meeting,
	}))
}

func ListMeetings(c *gin.Context) {
	var meetings []models.Meeting
	var result *gorm.DB
	var agentId int
	// var count int64

	// if c.Query("page") == "" || c.Query("size") == "" {
	// 	c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Page and Size are required parameters"))
	// 	return
	// }

	meetingType := c.Query("type")

	agentIdParam := c.Query("agentId")

	if agentIdParam != "" {
		agentId, _ = strconv.Atoi(agentIdParam)
	}

	// size, err := strconv.Atoi(c.Query("size"))

	// if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
	// 	return
	// }

	// page, err := strconv.Atoi(c.Query("page"))

	// if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
	// 	return
	// }

	// skip := page * (size - 1)

	authUser, _ := c.Get("authUser")

	user := authUser.(models.User)

	switch meetingType {
	case "pending":
		if user.Role == utils.AdminRole {
			if agentId > 0 {
				result = configs.DB.
					Preload("Agent").
					Where("agent_id = ?", uint(agentId)).
					Where("is_ended = ?", false).
					Where("is_deleted = ?", false).
					Order("schedule_date_time desc").
					Find(&meetings)
			} else {
				result = configs.DB.
					Preload("Agent").
					Where("is_ended = ?", false).
					Where("is_deleted = ?", false).
					Order("schedule_date_time desc").
					Find(&meetings)
			}
		} else {
			result = configs.DB.
				Preload("Agent").
				Where("agent_id = ?", user.ID).
				Where("is_ended = ?", false).
				Where("is_deleted = ?", false).
				Order("schedule_date_time desc").
				Find(&meetings)
		}
		// .Limit(size).Offset(skip)

		if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		// result = configs.DB.
		// 	Model(&models.Meeting{}).
		// 	Where("agent_id = ? and is_ended = ? and is_deleted = ?", user.ID, false, false).
		// 	Count(&count)

		// if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		// 	return
		// }
	case "finished":
		if user.Role == utils.AdminRole {
			if agentId > 0 {
				result = configs.DB.
					Preload("Agent").
					Where("agent_id = ?", uint(agentId)).
					Where("is_ended = ?", true).
					Where("is_deleted = ?", false).
					Order("schedule_date_time desc").
					Find(&meetings)
			} else {
				result = configs.DB.
					Preload("Agent").
					Where("is_ended = ?", true).
					Where("is_deleted = ?", false).
					Order("schedule_date_time desc").
					Find(&meetings)
			}
		} else {
			result = configs.DB.
				Preload("Agent").
				Where("agent_id = ?", user.ID).
				Where("is_ended = ?", true).
				Where("is_deleted = ?", false).
				Order("schedule_date_time desc").
				Find(&meetings)
		}
		// .Limit(size).Offset(skip)

		if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		// result = configs.DB.
		// 	Model(&models.Meeting{}).
		// 	Where("agent_id = ? and is_ended = ? and is_deleted = ?", user.ID, true, false).
		// 	Count(&count)

		// if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		// 	return
		// }
	default:
		if user.Role == utils.AdminRole {
			if agentId > 0 {
				result = configs.DB.
					Preload("Agent").
					Where("agent_id = ?", uint(agentId)).
					Where("is_deleted = ?", false).
					Order("schedule_date_time desc").
					Find(&meetings)
			} else {
				result = configs.DB.
					Preload("Agent").
					Where("is_deleted = ?", false).
					Order("schedule_date_time desc").
					Find(&meetings)
			}
		} else {
			result = configs.DB.
				Preload("Agent").
				Where("agent_id = ?", user.ID).
				Where("is_deleted = ?", false).
				Order("schedule_date_time desc").
				Find(&meetings)
		}
		// .Limit(size).Offset(skip)

		if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		// result = configs.DB.
		// 	Model(&models.Meeting{}).
		// 	Where("agent_id = ?", user.ID).
		// 	Count(&count)

		// if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		// 	return
		// }
	}

	// pagination := types.Pagination{
	// 	PageSize:    size,
	// 	CurrentPage: page,
	// 	TotalPages:  int(count) / size,
	// }

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"meetings": meetings,
		// "pagination": pagination,
	}))
}

func GetMeeting(c *gin.Context) {
	var meeting models.Meeting
	var meetingCode string
	var meetingQuery models.Meeting

	meetingId, err := strconv.Atoi(c.Param("meetingId"))

	if err != nil {
		meetingCode = c.Param("meetingId")
		_, err := uuid.Parse(meetingCode)
		if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
			return
		}
	}

	if meetingId != 0 {
		meetingQuery.ID = uint(meetingId)
		meetingQuery.IsDeleted = false
	}

	if meetingCode != "" {
		meetingQuery.MeetingCode = meetingCode
		meetingQuery.IsDeleted = false
	}

	result := configs.DB.
		Preload("Agent").
		Preload("Chat").
		Preload("Recordings").
		Preload("Documents").
		Preload("Answers").
		Preload("Answers.Question").
		First(&meeting, meetingQuery)

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting Not Found") {
		return
	}

	if meeting.Documents != nil {
		for i := range *meeting.Documents {
			document := &(*meeting.Documents)[i]
			if document.Type == "ocr" && document.APIResponse != nil {
				decodedResponsedData, err := base64.StdEncoding.DecodeString(*document.APIResponse)

				if utils.CheckError(err, http.StatusNotFound, c, "Meeting Not Found") {
					return
				}

				decryptedResponseData, err := utils.DecryptData(decodedResponsedData)

				if utils.CheckError(err, http.StatusNotFound, c, "Meeting Not Found") {
					return
				}

				decryptedResponse := string(decryptedResponseData)
				document.APIResponse = &decryptedResponse
			}
		}
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"meeting": meeting,
	}))
}

func UpdateMeeting(c *gin.Context) {
	var body types.MeetingInput
	var meeting models.Meeting
	var agent models.User
	meetingId, err := strconv.Atoi(c.Param("meetingId"))

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	err = c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	authUser, _ := c.Get("authUser")

	user := authUser.(models.User)

	result := configs.DB.First(&meeting, models.Meeting{
		ID:        uint(meetingId),
		IsDeleted: false,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting Not Found") {
		return
	}

	if user.Role == utils.AdminRole && *body.AgentID != 0 && uint(*body.AgentID) != meeting.AgentID {
		result := configs.DB.First(&agent, models.User{
			ID:        uint(*body.AgentID),
			Status:    utils.ActiveStatus,
			IsDeleted: false,
		})

		if utils.CheckError(result.Error, http.StatusNotFound, c, "Agent Not Found") {
			return
		}

		meeting.AgentID = agent.ID
	}

	meeting.Title = body.Title
	meeting.FirstName = body.FirstName
	meeting.LastName = body.LastName
	meeting.Email = body.Email
	meeting.CountryCode = body.CountryCode
	meeting.Phone = body.Phone
	meeting.Remark = body.Remark

	result = configs.DB.Save(&meeting)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "Meeting Updated",
	}))
}

func RescheduleMeeting(c *gin.Context) {
	var body map[string]string
	var meeting models.Meeting
	var existingMeeting models.Meeting
	var meetingTimeGapConfig models.Configuration

	meetingId, err := strconv.Atoi(c.Param("meetingId"))

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	err = c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	result := configs.DB.Preload("Agent").First(&meeting, models.Meeting{
		ID:        uint(meetingId),
		IsDeleted: false,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting Not Found") {
		return
	}

	scheduleDateTime, err := time.Parse(utils.DateLayout, body["scheduleDateTime"])

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	meetingTimeGap := 60

	result = configs.DB.First(&meetingTimeGapConfig, models.Configuration{
		Key:    utils.MeetingTimeGap,
		Status: utils.ActiveStatus,
	})

	if result.Error != nil {
		meetingTimeGap = 60
	}

	if meetingTimeGapConfig.Value != "" {
		meetingTimeGap, err = strconv.Atoi(meetingTimeGapConfig.Value)

		if err != nil || meetingTimeGap <= 0 {
			meetingTimeGap = 60
		}
	}

	fromCheckTime := scheduleDateTime.Add(time.Minute * -time.Duration(meetingTimeGap-1))

	result = configs.DB.
		Model(&models.Meeting{}).
		Where("agent_id = ?", meeting.AgentID).
		Where("schedule_date_time >= ?", fromCheckTime).
		Where("schedule_date_time <= ?", scheduleDateTime).
		Where("is_deleted = ?", false).
		Where("id != ?", meeting.ID).
		Find(&existingMeeting)

	if result.RowsAffected > 0 {
		msg := fmt.Sprintf("Meetings can only be scheduled with %d minutes gap", meetingTimeGap)
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse(msg))
		return
	}

	meeting.ScheduleDateTime = scheduleDateTime

	result = configs.DB.Save(&meeting)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	err = utils.SendMeetingEmail(&meeting, meeting.Agent)

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	err = utils.SendMeetingSMS(&meeting, meeting.Agent)

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "Meeting Rescheduled",
	}))
}

func ResendMeetingInvite(c *gin.Context) {
	var meeting models.Meeting

	meetingId, err := strconv.Atoi(c.Param("meetingId"))

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	result := configs.DB.Preload("Agent").First(&meeting, models.Meeting{
		ID:        uint(meetingId),
		IsDeleted: false,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting Not Found") {
		return
	}

	err = utils.SendMeetingEmail(&meeting, meeting.Agent)

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	err = utils.SendMeetingSMS(&meeting, meeting.Agent)

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"meeting": meeting,
	}))
}

func DeleteMeeting(c *gin.Context) {
	var meeting models.Meeting

	meetingId, err := strconv.Atoi(c.Param("meetingId"))

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	result := configs.DB.First(&meeting, models.Meeting{
		ID: uint(meetingId),
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting Not Found") {
		return
	}

	meeting.IsDeleted = true

	result = configs.DB.Save(&meeting)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "Meeting Deleted",
	}))
}

func CheckMeetingCode(c *gin.Context) {
	var meeting models.Meeting
	meetingCode := c.Param("meetingId")

	result := configs.DB.Preload("Agent").Where(&models.Meeting{
		MeetingCode: meetingCode,
		IsDeleted:   false,
		// IsEnded:     false,
	}).First(&meeting)

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Invalid Meeting Code") {
		return
	}

	if meeting.IsEnded {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Meeting Ended"))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "Valid Meeting Code",
	}))
}

func JoinMeeting(c *gin.Context) {
	var body types.JoinMeetingInput
	var meeting models.Meeting

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	result := configs.DB.Preload("Agent").First(&meeting, models.Meeting{
		MeetingCode: body.MeetingCode,
		IsDeleted:   false,
		// Password:    body.Password,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting Not Found") {
		return
	}

	if meeting.Password != body.Password {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid Password"))
		return
	}

	if meeting.IsEnded {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Meeting Ended"))
		return
	}

	tokenString := fmt.Sprintf("%s#%s#%s", meeting.MeetingCode, "false", meeting.Email)

	tokenBytes, err := utils.EncryptData([]byte(tokenString))

	token := base64.StdEncoding.EncodeToString(tokenBytes)

	if utils.CheckError(err, http.StatusNotFound, c, "Meeting Not Found or Invalid Password") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"token": token,
	}))
}

func AgentJoinMeeting(c *gin.Context) {
	var body types.JoinMeetingInput
	var meeting models.Meeting

	authUser, _ := c.Get("authUser")

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	result := configs.DB.Preload("Agent").First(&meeting, models.Meeting{
		MeetingCode: body.MeetingCode,
		Password:    body.Password,
		AgentID:     authUser.(models.User).ID,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting not found") {
		return
	}

	tokenString := fmt.Sprintf("%s#%s#%s", meeting.MeetingCode, "true", meeting.Agent.Email)

	tokenBytes, err := utils.EncryptData([]byte(tokenString))

	token := base64.StdEncoding.EncodeToString(tokenBytes)

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"token": token,
	}))
}

func MeetingEvents(c *gin.Context) {
	conn, err := utils.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	go func() {
		var meeting models.Meeting

		meetingCode := c.Query("meetingCode")
		token := c.Query("token")
		userState := c.Query("userState")

		if meetingCode == "" || token == "" {
			return
		}

		decodedToken, err := base64.StdEncoding.DecodeString(token)

		if err != nil {
			fmt.Println("Invalid token!")
			conn.Close()
			return
		}

		decryptedTokenBytes, err := utils.DecryptData([]byte(decodedToken))

		decryptedToken := string(decryptedTokenBytes)

		if err != nil || decryptedToken == "" {
			fmt.Println("Invalid token!")
			conn.Close()
			return
		}

		decryptedTokenDetails := strings.Split(decryptedToken, "#")

		if decryptedTokenDetails[0] != meetingCode {
			fmt.Println("Invalid token")
			conn.Close()
			return
		}

		result := configs.DB.Preload("Agent").First(&meeting, models.Meeting{
			MeetingCode: meetingCode,
			IsEnded:     false,
			IsDeleted:   false,
		})

		if result.Error != nil {
			fmt.Println("Meeting not found")
			conn.Close()
			return
		}

		utils.MeetingRoomsMutex.RLock()
		meetingRoom, meetingRoomExists := utils.MeetingRooms[meetingCode]
		utils.MeetingRoomsMutex.RUnlock()

		if !meetingRoomExists {
			utils.MeetingRoomsMutex.Lock()
			var waitingRoomUsers = make(map[string]*types.MeetingUser)
			var meetingRoomUsers = make(map[string]*types.MeetingUser)
			meetingRoom = &types.MeetingRoom{
				Meeting:          meeting,
				WaitingRoomUsers: waitingRoomUsers,
				MeetingRoomUsers: meetingRoomUsers,
			}
			utils.MeetingRooms[meetingCode] = meetingRoom
			utils.MeetingRoomsMutex.Unlock()
		}

		meetingRoom.WaitingRoomUsersMutex.RLock()
		_, waitingRoomUserExists := meetingRoom.WaitingRoomUsers[token]
		meetingRoom.WaitingRoomUsersMutex.RUnlock()

		meetingRoom.MeetingRoomUsersMutex.RLock()
		meetingRoomUser, meetingRoomUserExists := meetingRoom.MeetingRoomUsers[token]
		meetingRoom.MeetingRoomUsersMutex.RUnlock()

		if waitingRoomUserExists {
			meetingRoom.WaitingRoomUsersMutex.Lock()
			meetingRoom.WaitingRoomUsers[token].Conn.Close()
			meetingRoom.WaitingRoomUsers[token].Conn = conn
			meetingRoom.WaitingRoomUsersMutex.Unlock()

			utils.Broadcast(types.BroadcastParameters{
				Msg:           nil,
				MeetingRoom:   meetingRoom,
				Conn:          conn,
				MessageType:   "waiting",
				ToSelf:        true,
				ToWaitingRoom: true,
				ToMeetingRoom: true,
			})
		} else if meetingRoomUserExists {
			if userState == "admit" {
				meetingRoom.MeetingRoomUsersMutex.Lock()
				meetingRoom.MeetingRoomUsers[token].Conn.Close()
				meetingRoom.MeetingRoomUsers[token].Conn = conn
				meetingRoom.MeetingRoomUsersMutex.Unlock()

				utils.Broadcast(types.BroadcastParameters{
					Msg:           nil,
					MeetingRoom:   meetingRoom,
					Conn:          conn,
					MessageType:   "rejoin",
					ToSelf:        false,
					ToWaitingRoom: false,
					ToMeetingRoom: true,
				})
			} else {
				activeConnection := utils.CheckActiveConnection(meetingRoomUser, conn, token)

				if activeConnection {
					messageType := "existingUser"
					utils.SelfBroadcast(conn, &types.Message{
						Type: &messageType,
					})
					conn.Close()
					return
				} else {
					meetingRoomUser.Conn.Close()
					meetingRoomUser.Conn = conn
				}
			}
		} else {
			user := &types.MeetingUser{
				MeetingCode: meetingCode,
				IsHost:      decryptedTokenDetails[1] == "true",
				Email:       decryptedTokenDetails[2],
				Conn:        conn,
			}

			if decryptedTokenDetails[1] == "true" {
				meetingRoom.MeetingRoomUsersMutex.Lock()
				meetingRoom.MeetingRoomUsers[token] = user
				meetingRoom.MeetingRoomUsersMutex.Unlock()

				utils.Broadcast(types.BroadcastParameters{
					Msg:           nil,
					MeetingRoom:   meetingRoom,
					Conn:          conn,
					MessageType:   "waiting",
					ToSelf:        true,
					ToWaitingRoom: false,
					ToMeetingRoom: true,
				})
			} else {
				meetingRoom.WaitingRoomUsersMutex.Lock()
				meetingRoom.WaitingRoomUsers[token] = user
				meetingRoom.WaitingRoomUsersMutex.Unlock()

				utils.Broadcast(types.BroadcastParameters{
					Msg:           nil,
					MeetingRoom:   meetingRoom,
					Conn:          conn,
					MessageType:   "waiting",
					ToSelf:        true,
					ToWaitingRoom: true,
					ToMeetingRoom: true,
				})
			}
		}

		conn.SetCloseHandler(func(code int, text string) error {
			meetingRoom.WaitingRoomUsersMutex.Lock()
			delete(meetingRoom.WaitingRoomUsers, token)
			meetingRoom.WaitingRoomUsersMutex.Unlock()

			meetingRoom.MeetingRoomUsersMutex.Lock()
			delete(meetingRoom.MeetingRoomUsers, token)
			meetingRoom.MeetingRoomUsersMutex.Unlock()

			// utils.Broadcast(utils.BroadcastParameters{
			// 	Msg:           nil,
			// 	MeetingRoom:   meetingRoom,
			// 	Conn:          nil,
			// 	MessageType:   "disconnect",
			// 	ToMeetingRoom: true,
			// 	ToWaitingRoom: false,
			// 	ToSelf:        false,
			// })

			conn.Close()

			return nil
		})

		for {
			var msg types.Message
			err := conn.ReadJSON(&msg)

			if err != nil {
				fmt.Println(err.Error())
				break
			}

			if msg.Type == nil {
				continue
			}

			switch *msg.Type {
			case "rtc":
				utils.Broadcast(types.BroadcastParameters{
					Msg:           &msg,
					MeetingRoom:   meetingRoom,
					Conn:          conn,
					MessageType:   "rtc",
					ToSelf:        false,
					ToMeetingRoom: true,
					ToWaitingRoom: true,
				})
			case "chat":
				utils.HandleChat(msg, meetingRoom, token)
			case "video":
				utils.HandleControls(msg, meetingRoom, token)
			case "audio":
				utils.HandleControls(msg, meetingRoom, token)
			case "screen":
				utils.HandleControls(msg, meetingRoom, token)
			case "meetingRoom":
				utils.HandleMeetingRoom(msg, meetingRoom, token)
			case "waitingRoom":
				utils.HandleWaitingRoom(msg, meetingRoom, token)
			case "location":
				utils.HandleLocation(msg, meetingRoom, token)
			case "pong":
				utils.HandlePong(token)
				fmt.Println("Received pong message")
			default:
				fmt.Println("No case matched")
			}
		}
	}()
}

func GetMeetingChat(c *gin.Context) {
	var meeting models.Meeting
	meetingCode := c.Param("meetingId")

	result := configs.DB.Preload("Chat").First(&meeting, models.Meeting{
		MeetingCode: meetingCode,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting Not Found") {
		return
	}

	utils.MeetingRoomsMutex.RLock()
	meetingRoom, meetingRoomExists := utils.MeetingRooms[meetingCode]
	utils.MeetingRoomsMutex.RUnlock()

	if meetingRoomExists {
		var chats []map[string]interface{}
		for _, chat := range *meeting.Chat {
			chats = append(chats, map[string]interface{}{
				"sender":    chat.Sender,
				"message":   chat.Message,
				"isAgent":   chat.IsAgent,
				"createdAt": chat.CreatedAt,
				"meetingID": chat.MeetingID,
			})
		}
		meetingRoom.ChatMutex.Lock()
		meetingRoom.Chat = chats
		meetingRoom.ChatMutex.Unlock()
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"chat": meeting.Chat,
	}))
}

func SaveMeetingRecording(c *gin.Context) {
	var meeting models.Meeting

	authUser, _ := c.Get("authUser")

	user := authUser.(models.User)

	file, err := c.FormFile("recording")
	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	meetingCode := c.Param("meetingId")

	result := configs.DB.First(&meeting, models.Meeting{
		MeetingCode: meetingCode,
		AgentID:     user.ID,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Meeting not found") {
		return
	}

	fileExt := filepath.Ext(file.Filename)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixMilli(), fileExt)
	savePath := fmt.Sprintf("uploads/%s/recordings/%s", meetingCode, fileName)

	err = c.SaveUploadedFile(file, savePath)
	if utils.CheckError(err, http.StatusInternalServerError, c, "Failed to save recording") {
		return
	}

	recording := models.Recording{
		MeetingID: meeting.ID,
		Recording: savePath,
	}

	result = configs.DB.Create(&recording)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Failed to save recording") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "Recording saved successfully",
	}))
}

func GetMeetingMedia(c *gin.Context) {
	var recording models.Recording
	var document models.Document

	mediaType := c.Query("mediaType")

	mediaId, err := strconv.Atoi(c.Param("mediaId"))

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	if mediaType == "document" {
		result := configs.DB.Preload("Meeting").First(&document, models.Document{
			ID: uint(mediaId),
		})

		if utils.CheckError(result.Error, http.StatusNotFound, c, "Document Not Found") {
			return
		}

		encryptedImageData, err := os.ReadFile(document.Image)

		if utils.CheckError(err, http.StatusInternalServerError, c, "Failed to read document") {
			return
		}

		decryptedImageData, err := utils.DecryptData(encryptedImageData)

		if utils.CheckError(err, http.StatusInternalServerError, c, "Failed to read document") {
			return
		}

		c.Header("Content-Type", "image/jpeg")

		c.Data(http.StatusOK, "image/jpeg", decryptedImageData)
		return
	} else {
		result := configs.DB.Preload("Meeting").First(&recording, models.Recording{
			ID: uint(mediaId),
		})

		if utils.CheckError(result.Error, http.StatusNotFound, c, "Document Not Found") {
			return
		}

		c.File(recording.Recording)
		return
	}
}

func DeleteMeetingMedia(c *gin.Context) {
	var recording models.Recording
	var document models.Document

	mediaType := c.Query("mediaType")
	mediaId, err := strconv.Atoi(c.Param("mediaId"))

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	if mediaType == "document" {
		result := configs.DB.Preload("Meeting").First(&document, models.Document{
			ID: uint(mediaId),
		})

		if utils.CheckError(result.Error, http.StatusNotFound, c, "Document Not Found") {
			return
		}

		err = os.Remove(document.Image)
		if utils.CheckError(err, http.StatusInternalServerError, c, "Failed to delete document") {
			return
		}

		result = configs.DB.Delete(&document)

		if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Failed to delete recording") {
			return
		}

		c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
			"message": "Document deleted successfully",
		}))

		return
	} else {
		result := configs.DB.Preload("Meeting").First(&recording, models.Recording{
			ID: uint(mediaId),
		})

		if utils.CheckError(result.Error, http.StatusNotFound, c, "Recording Not Found") {
			return
		}

		err = os.Remove(recording.Recording)
		if utils.CheckError(err, http.StatusInternalServerError, c, "Failed to delete recording") {
			return
		}

		result = configs.DB.Delete(&recording)

		if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Failed to delete recording") {
			return
		}

		c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
			"message": "Recording deleted successfully",
		}))

		return
	}
}
