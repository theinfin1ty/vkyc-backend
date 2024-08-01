package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"vkyc-backend/configs"
	"vkyc-backend/models"
)

func SendMeetingSMS(meeting *models.Meeting, authUser *models.User) error {
	var apiConfig, apiKeyConfig, senderConfig, templateConfig models.Configuration

	result := configs.DB.First(&apiConfig, models.Configuration{
		Key:    "SMS_API",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.First(&apiKeyConfig, models.Configuration{
		Key:    "SMS_API_KEY",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.First(&senderConfig, models.Configuration{
		Key:    "SMS_SENDER",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.First(&templateConfig, models.Configuration{
		Key:    "SMS_TEMPLATE",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	if apiKeyConfig.Value == "" || apiConfig.Value == "" || senderConfig.Value == "" || templateConfig.Value == "" {
		return fmt.Errorf("SMS Configuration not found in configuration")
	}

	userPhone := fmt.Sprintf("%s%s", meeting.CountryCode, meeting.Phone)

	recipients := []map[string]string{
		{
			"mobiles": userPhone,
			"var":     meeting.Password,
		},
	}

	if authUser.Phone != "" && authUser.CountryCode != "" {
		authUserPhone := fmt.Sprintf("%s%s", authUser.CountryCode, authUser.Phone)

		recipients = append(recipients, map[string]string{
			"mobiles": authUserPhone,
			"var":     meeting.Password,
		})
	}

	data, err := json.Marshal(map[string]interface{}{
		"template_id": templateConfig.Value,
		"short_url":   "0",
		"recipients":  recipients,
	})

	if err != nil {
		return err
	}

	dataBuffer := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", apiConfig.Value, dataBuffer)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authKey", apiKeyConfig.Value)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to send sms")
	}

	return nil
}

func SendOTPSMS(otp string, authUser *models.User) error {
	var apiConfig, apiKeyConfig, senderConfig, templateConfig models.Configuration
	var recipients []map[string]string

	result := configs.DB.First(&apiConfig, models.Configuration{
		Key:    "SMS_API",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.First(&apiKeyConfig, models.Configuration{
		Key:    "SMS_API_KEY",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.First(&senderConfig, models.Configuration{
		Key:    "SMS_SENDER",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.First(&templateConfig, models.Configuration{
		Key:    "SMS_TEMPLATE",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	if apiKeyConfig.Value == "" || apiConfig.Value == "" || senderConfig.Value == "" || templateConfig.Value == "" {
		return fmt.Errorf("SMS Configuration not found in configuration")
	}

	if authUser.Phone != "" && authUser.CountryCode != "" {
		authUserPhone := fmt.Sprintf("%s%s", authUser.CountryCode, authUser.Phone)

		recipients = append(recipients, map[string]string{
			"mobiles": authUserPhone,
			"var":     otp,
		})
	}

	data, err := json.Marshal(map[string]interface{}{
		"template_id": templateConfig.Value,
		"short_url":   "0",
		"recipients":  recipients,
	})

	if err != nil {
		return err
	}

	dataBuffer := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", apiConfig.Value, dataBuffer)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authKey", apiKeyConfig.Value)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to send sms")
	}

	return nil
}

func SendInviteSMS(user *models.User, password string) error {
	var apiConfig, apiKeyConfig, senderConfig, templateConfig models.Configuration
	var recipients []map[string]string

	result := configs.DB.First(&apiConfig, models.Configuration{
		Key:    "SMS_API",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.First(&apiKeyConfig, models.Configuration{
		Key:    "SMS_API_KEY",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.First(&senderConfig, models.Configuration{
		Key:    "SMS_SENDER",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	result = configs.DB.First(&templateConfig, models.Configuration{
		Key:    "SMS_TEMPLATE",
		Status: ActiveStatus,
	})

	if result.Error != nil {
		return result.Error
	}

	if apiKeyConfig.Value == "" || apiConfig.Value == "" || senderConfig.Value == "" || templateConfig.Value == "" {
		return fmt.Errorf("SMS Configuration not found in configuration")
	}

	if user.Phone != "" && user.CountryCode != "" {
		userPhone := fmt.Sprintf("%s%s", user.CountryCode, user.Phone)

		recipients = append(recipients, map[string]string{
			"mobiles": userPhone,
			"var":     password,
		})
	}

	data, err := json.Marshal(map[string]interface{}{
		"template_id": templateConfig.Value,
		"short_url":   "0",
		"recipients":  recipients,
	})

	if err != nil {
		return err
	}

	dataBuffer := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", apiConfig.Value, dataBuffer)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authKey", apiKeyConfig.Value)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to send sms")
	}

	return nil
}
