package utils

import (
	"fmt"
	"strconv"
	"time"
	"vkyc-backend/configs"
	"vkyc-backend/models"

	"github.com/go-co-op/gocron"
)

func endExpiredMeetings() {
	fmt.Println("Running cron job...")
	var beforeMeetingStartConfiguration models.Configuration
	var afterMeetingStartConfiguration models.Configuration

	result := configs.DB.Find(&beforeMeetingStartConfiguration, models.Configuration{
		Key:    MeetingExpireBeforeStart,
		Status: ActiveStatus,
	})

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return
	}

	result = configs.DB.Find(&afterMeetingStartConfiguration, models.Configuration{
		Key:    MeetingExpireAfterStart,
		Status: ActiveStatus,
	})

	if result.Error != nil {
		fmt.Println(result.Error.Error())
		return
	}

	if beforeMeetingStartConfiguration.Value != "" {
		var yetToStartMeetings []models.Meeting
		minutes, err := strconv.Atoi(beforeMeetingStartConfiguration.Value)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		currentTime := time.Now().Add(time.Duration(minutes) * time.Minute * -1)

		result := configs.DB.
			Where("schedule_date_time < ? AND is_ended = ? AND is_deleted = ? AND is_started = ?", currentTime, false, false, false).
			Find(&yetToStartMeetings)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
			return
		}

		for _, meeting := range yetToStartMeetings {
			// updateMeetingKYCStatus(&meeting)
			meeting.IsEnded = true

			meetingRoom, meetingRoomExists := MeetingRooms[meeting.MeetingCode]

			if meetingRoomExists {
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
			}

			result = configs.DB.Save(&meeting)
			if result.Error != nil {
				fmt.Println(result.Error.Error())
				// return
			}
		}
	}

	if afterMeetingStartConfiguration.Value != "" {
		var meetings []models.Meeting
		minutes, err := strconv.Atoi(afterMeetingStartConfiguration.Value)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		currentTime := time.Now().Add(time.Duration(minutes) * time.Minute * -1)

		result := configs.DB.
			Where("schedule_date_time < ? AND is_ended = ? AND is_deleted = ? AND is_started = ?", currentTime, false, false, true).
			Find(&meetings)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
			return
		}

		for _, meeting := range meetings {
			updateMeetingKYCStatus(&meeting)

			meetingRoom, meetingRoomExists := MeetingRooms[meeting.MeetingCode]

			if meetingRoomExists {

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
			}
		}

		// result = configs.DB.Save(&meetings)

		if result.Error != nil {
			fmt.Println(result.Error.Error())
			return
		}
	}
}

func StartScheduler() {
	scheduler := gocron.NewScheduler(&time.Location{})

	_, err := scheduler.Every(5).Minute().Do(endExpiredMeetings)

	if err != nil {
		fmt.Println(err.Error())
	}

	scheduler.StartAsync()
}
