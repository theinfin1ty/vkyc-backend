package types

import (
	"sync"
	"vkyc-backend/models"

	"github.com/gorilla/websocket"
)

type MeetingUser struct {
	MeetingCode string          `json:"meetingCode"`
	IsHost      bool            `json:"isHost"`
	Email       string          `json:"email"`
	Conn        *websocket.Conn `json:"-"`
	Video       bool            `json:"video"`
	Audio       bool            `json:"audio"`
	Screen      bool            `json:"screen"`
}

type MeetingRoom struct {
	Meeting          models.Meeting           `json:"meeting"`
	WaitingRoomUsers map[string]*MeetingUser  `json:"waitingRoomUsers"`
	MeetingRoomUsers map[string]*MeetingUser  `json:"meetingRoomUsers"`
	Chat             []map[string]interface{} `json:"chat"`

	MeetingMutex          sync.RWMutex
	MeetingRoomUsersMutex sync.RWMutex
	WaitingRoomUsersMutex sync.RWMutex
	ChatMutex             sync.RWMutex
}

type Message struct {
	Type    *string `json:"type"`
	Status  *string `json:"status"`
	Message *string `json:"message"`
}

type BroadcastMessage struct {
	Type        string       `json:"type"`
	MeetingRoom *MeetingRoom `json:"meetingRoom"`
}

type BroadcastParameters struct {
	Msg           *Message
	MeetingRoom   *MeetingRoom
	Conn          *websocket.Conn
	MessageType   string
	ToMeetingRoom bool
	ToWaitingRoom bool
	ToSelf        bool
}
