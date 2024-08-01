package utils

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"vkyc-backend/configs"
	"vkyc-backend/models"
	"vkyc-backend/types"

	"github.com/gorilla/websocket"
)

var MeetingRooms = make(map[string]*types.MeetingRoom)
var MeetingRoomsMutex sync.RWMutex

var TempMeetingUsers = make(map[string]*websocket.Conn)

var Upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func SelfBroadcast(conn *websocket.Conn, msg *types.Message) {
	err := conn.WriteJSON(msg)
	if err != nil {
		fmt.Println(err)
	}
}

func Broadcast(params types.BroadcastParameters) {
	broadcastMessage := types.BroadcastMessage{
		Type:        params.MessageType,
		MeetingRoom: params.MeetingRoom,
	}

	if params.ToMeetingRoom {
		params.MeetingRoom.MeetingRoomUsersMutex.RLock()
		for _, connectedUser := range params.MeetingRoom.MeetingRoomUsers {
			if !params.ToSelf && (params.Conn != nil && connectedUser.Conn == params.Conn) {
				continue
			}
			if params.Msg != nil && *params.Msg.Type == "rtc" {
				err := connectedUser.Conn.WriteJSON(params.Msg)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				err := connectedUser.Conn.WriteJSON(broadcastMessage)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		params.MeetingRoom.MeetingRoomUsersMutex.RUnlock()
	}

	if params.ToWaitingRoom {
		params.MeetingRoom.WaitingRoomUsersMutex.RLock()
		for _, connectedUser := range params.MeetingRoom.WaitingRoomUsers {
			if !params.ToSelf && (params.Conn != nil && connectedUser.Conn == params.Conn) {
				continue
			}
			if params.Msg != nil && *params.Msg.Type == "rtc" {
				err := connectedUser.Conn.WriteJSON(params.Msg)
				if err != nil {
					fmt.Println(err)
				}
			} else {
				err := connectedUser.Conn.WriteJSON(broadcastMessage)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		params.MeetingRoom.WaitingRoomUsersMutex.RUnlock()
	}
}

func CheckActiveConnection(meetingRoomUser *types.MeetingUser, conn *websocket.Conn, token string) bool {
	messageType := "ping"
	SelfBroadcast(meetingRoomUser.Conn, &types.Message{
		Type: &messageType,
	})

	TempMeetingUsers[token] = conn

	pongReceived := false
	checks := 5

	for i := 0; i < checks; i++ {
		time.Sleep(1 * time.Second)

		_, exists := TempMeetingUsers[token]

		if !exists {
			pongReceived = true
		}
	}

	delete(TempMeetingUsers, token)
	return pongReceived
}

func HandlePong(token string) {
	delete(TempMeetingUsers, token)
}

func HandleChat(msg types.Message, meetingRoom *types.MeetingRoom, token string) {
	meetingRoom.MeetingRoomUsersMutex.RLock()
	user, userExists := meetingRoom.MeetingRoomUsers[token]
	meetingRoom.MeetingRoomUsersMutex.RUnlock()

	if !userExists {
		return
	}

	meetingRoom.ChatMutex.Lock()
	meetingRoom.Chat = append(meetingRoom.Chat, map[string]interface{}{
		"sender":      user.Email,
		"senderToken": token,
		"message":     *msg.Message,
		"isAgent":     user.IsHost,
		"createdAt":   time.Now().UTC().Format(DateLayout),
		"meetingID":   meetingRoom.Meeting.ID,
	})
	meetingRoom.ChatMutex.Unlock()

	Broadcast(types.BroadcastParameters{
		Msg:           &msg,
		MeetingRoom:   meetingRoom,
		Conn:          user.Conn,
		MessageType:   "chat",
		ToSelf:        true,
		ToMeetingRoom: true,
		ToWaitingRoom: true,
	})

	result := configs.DB.Create(&models.Chat{
		Message:   *msg.Message,
		Sender:    user.Email,
		MeetingID: meetingRoom.Meeting.ID,
		IsAgent:   user.IsHost,
	})

	if result.Error != nil {
		fmt.Println(result.Error)
	}
}

func HandleControls(msg types.Message, meetingRoom *types.MeetingRoom, token string) {
	var user *types.MeetingUser

	meetingRoom.WaitingRoomUsersMutex.RLock()
	waitingRoomUser, existsInWaitingRoom := meetingRoom.WaitingRoomUsers[token]
	meetingRoom.WaitingRoomUsersMutex.RUnlock()

	meetingRoom.MeetingRoomUsersMutex.RLock()
	meetingRoomUser, existsInMeetingRoom := meetingRoom.MeetingRoomUsers[token]
	meetingRoom.MeetingRoomUsersMutex.RUnlock()

	if !existsInMeetingRoom && !existsInWaitingRoom {
		return
	}

	if existsInWaitingRoom {
		user = waitingRoomUser
	}

	if existsInMeetingRoom {
		user = meetingRoomUser
	}

	switch *msg.Type {
	case "video":
		user.Video = *msg.Status == "on"
	case "audio":
		user.Audio = *msg.Status == "on"
	case "screen":
		user.Screen = *msg.Status == "on"
	default:
		return
	}

	Broadcast(types.BroadcastParameters{
		Msg:           &msg,
		MeetingRoom:   meetingRoom,
		Conn:          user.Conn,
		MessageType:   "controls",
		ToSelf:        true,
		ToMeetingRoom: true,
		ToWaitingRoom: true,
	})
}

func HandleMeetingRoom(msg types.Message, meetingRoom *types.MeetingRoom, token string) {
	meetingRoom.MeetingRoomUsersMutex.RLock()
	user, userExists := meetingRoom.MeetingRoomUsers[token]
	meetingRoom.MeetingRoomUsersMutex.RUnlock()

	if !userExists {
		return
	}

	switch *msg.Status {
	case "end":
		if !user.IsHost {
			return
		}

		err := updateMeetingKYCStatus(&meetingRoom.Meeting)

		if err != nil {
			fmt.Println(err)
		}

		meetingRoom.MeetingMutex.Lock()
		meetingRoom.Meeting.IsEnded = true
		meetingRoom.MeetingMutex.Unlock()

		// result := configs.DB.Save(&meetingRoom.Meeting)

		// if result.Error != nil {
		// 	fmt.Println(result.Error)
		// }

		Broadcast(types.BroadcastParameters{
			Msg:           &msg,
			MeetingRoom:   meetingRoom,
			Conn:          user.Conn,
			MessageType:   *msg.Status,
			ToSelf:        true,
			ToMeetingRoom: true,
			ToWaitingRoom: true,
		})

		meetingRoom.MeetingRoomUsersMutex.RLock()
		for _, connectedUser := range meetingRoom.MeetingRoomUsers {
			connectedUser.Conn.Close()
		}
		meetingRoom.MeetingRoomUsersMutex.RUnlock()

		meetingRoom.WaitingRoomUsersMutex.RLock()
		for _, connectedUser := range meetingRoom.WaitingRoomUsers {
			connectedUser.Conn.Close()
		}
		meetingRoom.WaitingRoomUsersMutex.RUnlock()

		MeetingRoomsMutex.Lock()
		delete(MeetingRooms, meetingRoom.Meeting.MeetingCode)
		MeetingRoomsMutex.Unlock()
	case "leave":
		user.Conn.Close()

		meetingRoom.MeetingRoomUsersMutex.Lock()
		delete(meetingRoom.MeetingRoomUsers, token)
		meetingRoom.MeetingRoomUsersMutex.Unlock()

		Broadcast(types.BroadcastParameters{
			Msg:           &msg,
			MeetingRoom:   meetingRoom,
			Conn:          nil,
			MessageType:   *msg.Status,
			ToSelf:        false,
			ToMeetingRoom: true,
			ToWaitingRoom: true,
		})
	default:
		return
	}
}

func HandleWaitingRoom(msg types.Message, meetingRoom *types.MeetingRoom, token string) {
	meetingRoom.MeetingRoomUsersMutex.RLock()
	user, userExists := meetingRoom.MeetingRoomUsers[token]
	meetingRoom.MeetingRoomUsersMutex.RUnlock()

	meetingRoom.WaitingRoomUsersMutex.RLock()
	waitingUser, waitingUserExists := meetingRoom.WaitingRoomUsers[token]
	meetingRoom.WaitingRoomUsersMutex.RUnlock()

	switch *msg.Status {
	case "admit":
		if !(userExists && user.IsHost) {
			return
		}

		meetingRoom.WaitingRoomUsersMutex.RLock()
		waitingRoomUser, waitingRoomUserExists := meetingRoom.WaitingRoomUsers[*msg.Message]
		meetingRoom.WaitingRoomUsersMutex.RUnlock()

		if waitingRoomUserExists {
			meetingRoom.MeetingRoomUsersMutex.Lock()
			meetingRoom.MeetingRoomUsers[*msg.Message] = waitingRoomUser
			meetingRoom.MeetingRoomUsersMutex.Unlock()

			meetingRoom.WaitingRoomUsersMutex.Lock()
			delete(meetingRoom.WaitingRoomUsers, *msg.Message)
			meetingRoom.WaitingRoomUsersMutex.Unlock()
		}

		meetingRoom.MeetingMutex.Lock()
		meetingRoom.Meeting.IsStarted = true
		meetingRoom.MeetingMutex.Unlock()

		result := configs.DB.Save(&meetingRoom.Meeting)

		if result.Error != nil {
			fmt.Println(result.Error)
		}

		Broadcast(types.BroadcastParameters{
			Msg:           &msg,
			MeetingRoom:   meetingRoom,
			Conn:          user.Conn,
			MessageType:   *msg.Status,
			ToSelf:        true,
			ToMeetingRoom: true,
			ToWaitingRoom: true,
		})
	case "leave":
		if !waitingUserExists {
			return
		}

		waitingUser.Conn.Close()
		meetingRoom.WaitingRoomUsersMutex.Lock()
		delete(meetingRoom.WaitingRoomUsers, token)
		meetingRoom.WaitingRoomUsersMutex.Unlock()

		Broadcast(types.BroadcastParameters{
			Msg:           &msg,
			MeetingRoom:   meetingRoom,
			Conn:          nil,
			MessageType:   *msg.Status,
			ToSelf:        false,
			ToMeetingRoom: true,
			ToWaitingRoom: true,
		})
	default:
		return
	}
}

func HandleLocation(msg types.Message, meetingRoom *types.MeetingRoom, token string) {
	meetingRoom.MeetingRoomUsersMutex.RLock()
	user, userExists := meetingRoom.MeetingRoomUsers[token]
	meetingRoom.MeetingRoomUsersMutex.RUnlock()

	if !userExists {
		return
	}

	if *msg.Status == "request" {
		Broadcast(types.BroadcastParameters{
			Msg:           &msg,
			MeetingRoom:   meetingRoom,
			Conn:          user.Conn,
			MessageType:   "location",
			ToSelf:        false,
			ToMeetingRoom: true,
			ToWaitingRoom: false,
		})
		return
	}

	meetingRoom.MeetingMutex.Lock()
	meetingRoom.Meeting.Location = msg.Message
	meetingRoom.MeetingMutex.Unlock()

	result := configs.DB.Save(&meetingRoom.Meeting)

	if result.Error != nil {
		fmt.Println(result.Error)
	}

	Broadcast(types.BroadcastParameters{
		Msg:           &msg,
		MeetingRoom:   meetingRoom,
		Conn:          user.Conn,
		MessageType:   "location",
		ToSelf:        false,
		ToMeetingRoom: true,
		ToWaitingRoom: false,
	})
}

func updateMeetingKYCStatus(meeting *models.Meeting) error {
	var passFacematchCount,
		passLivenessCount,
		passOCRCount,
		incorrectAnswersCount,
		signatureCount,
		correctAnswersCount int64
	var remark string

	kycStatus := RejectedStatus

	result := configs.DB.
		Model(&models.Document{}).
		Where("meeting_id = ?", meeting.ID).
		Where("type = ?", "facematch").
		Where("api_status = ?", PassStatus).
		Count(&passFacematchCount)

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.
		Model(&models.Document{}).
		Where("meeting_id = ?", meeting.ID).
		Where("type = ?", "liveness").
		Where("api_status = ?", PassStatus).
		Count(&passLivenessCount)

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.
		Model(&models.Document{}).
		Where("meeting_id = ?", meeting.ID).
		Where("type = ?", "ocr").
		Where("api_status = ?", "Success").
		Count(&passOCRCount)

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.
		Model(&models.Document{}).
		Where("meeting_id = ?", meeting.ID).
		Where("type = ?", "signature").
		Count(&signatureCount)

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.
		Model(&models.Answer{}).
		Where("meeting_id = ?", meeting.ID).
		Where("answer = 'incorrect'").
		Count(&incorrectAnswersCount)

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.
		Model(&models.Answer{}).
		Where("meeting_id = ?", meeting.ID).
		Where("answer = 'correct'").
		Count(&correctAnswersCount)

	if result.Error != nil {
		return result.Error
	}

	if correctAnswersCount == 0 &&
		incorrectAnswersCount == 0 &&
		signatureCount == 0 &&
		passOCRCount == 0 &&
		passLivenessCount == 0 &&
		passFacematchCount == 0 {
		remark = FailedAfterConnectingRemark
	}

	if remark == "" {
		if incorrectAnswersCount > 0 {
			remark = QAFailedRemark
		} else if signatureCount == 0 {
			remark = SignFailedRemark
		} else if passOCRCount == 0 {
			remark = OCRFailedRemark
		} else if passLivenessCount == 0 {
			remark = LivenessFailedRemark
		} else if passFacematchCount == 0 {
			remark = FacematchFailedRemark
		} else {
			kycStatus = ApprovedStatus
		}
	}

	result = configs.DB.Model(&models.Meeting{}).Where("id = ?", meeting.ID).Updates(models.Meeting{
		KycStatus: &kycStatus,
		Remarks:   &remark,
		IsEnded:   true,
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}
