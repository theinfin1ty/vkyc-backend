package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
	"vkyc-backend/configs"
	"vkyc-backend/models"
)

var fromEmail = configs.GetEnvVariable("FROM_EMAIL")
var smtpPort = configs.GetEnvVariable("SMTP_PORT")
var smtpHost = configs.GetEnvVariable("SMTP_HOST")
var smtpAuth = smtp.PlainAuth("", fromEmail, configs.GetEnvVariable("FROM_EMAIL_PASSWORD"), smtpHost)

func SendMeetingEmail(meeting *models.Meeting, authUser *models.User) error {
	var companyNameConfig, companyLogoConfig models.Configuration
	companyName := ""
	companyLogo := ""

	result := configs.DB.First(&companyNameConfig, models.Configuration{
		Key:    CompanyName,
		Status: ActiveStatus,
	})

	if result.Error == nil {
		companyName = companyNameConfig.Value
	}

	result = configs.DB.First(&companyLogoConfig, models.Configuration{
		Key:    CompanyLogo,
		Status: ActiveStatus,
	})

	if result.Error == nil {
		companyLogo = companyLogoConfig.Value
	}

	to := []string{meeting.Email}

	t, err := template.ParseFiles("templates/meeting.html")

	if err != nil {
		return err
	}

	var userBody, agentBody bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	userBody.Write([]byte(fmt.Sprintf("Subject: Confirmation of Your Scheduled vKYC Appointment with %s\n%s\n\n", companyName, mimeHeaders)))

	err = t.Execute(&userBody, map[string]string{
		"ScheduleDateTime": meeting.ScheduleDateTime.Format(DateFormat),
		"FirstName":        meeting.FirstName,
		"Link":             fmt.Sprintf("%s/public/call/%s?pswd=%s", configs.GetEnvVariable("WEB_URL"), meeting.MeetingCode, meeting.Password),
		"MeetingCode":      meeting.MeetingCode,
		"Password":         meeting.Password,
		"AgentName":        authUser.FirstName,
		"Company":          companyName,
		"CompanyLogo":      companyLogo,
	})

	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, smtpAuth, fromEmail, to, userBody.Bytes())

	if err != nil {
		return err
	}

	to = []string{authUser.Email}

	agentBody.Write([]byte(fmt.Sprintf("Subject: Confirmation of Your Scheduled vKYC Appointment with %s\n%s\n\n", companyName, mimeHeaders)))

	err = t.Execute(&agentBody, map[string]string{
		"ScheduleDateTime": meeting.ScheduleDateTime.Format(DateFormat),
		"FirstName":        authUser.FirstName,
		"Link":             fmt.Sprintf("%s/agent/call/%s?pswd=%s", configs.GetEnvVariable("WEB_URL"), meeting.MeetingCode, meeting.Password),
		"MeetingCode":      meeting.MeetingCode,
		"Password":         meeting.Password,
		"AgentName":        authUser.FirstName,
		"Company":          companyName,
		"CompanyLogo":      companyLogo,
	})

	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, smtpAuth, fromEmail, to, agentBody.Bytes())

	if err != nil {
		return err
	}

	return nil
}

func SendOTPEmail(otp string, user *models.User) error {
	var companyNameConfig, companyLogoConfig models.Configuration
	companyName := ""
	companyLogo := ""

	result := configs.DB.First(&companyNameConfig, models.Configuration{
		Key:    CompanyName,
		Status: ActiveStatus,
	})

	if result.Error == nil {
		companyName = companyNameConfig.Value
	}

	result = configs.DB.First(&companyLogoConfig, models.Configuration{
		Key:    CompanyLogo,
		Status: ActiveStatus,
	})

	if result.Error == nil {
		companyLogo = companyLogoConfig.Value
	}

	to := []string{user.Email}

	t, err := template.ParseFiles("templates/forgotPassword.html")

	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Reset Your Password - Secure OTP Inside\n%s\n\n", mimeHeaders)))

	err = t.Execute(&body, map[string]string{
		"OTP":         otp,
		"FirstName":   user.FirstName,
		"Company":     companyName,
		"CompanyLogo": companyLogo,
	})

	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, smtpAuth, fromEmail, to, body.Bytes())

	if err != nil {
		return err
	}

	return nil
}

func SendInviteEmail(user *models.User, password string) error {
	var companyNameConfig, companyLogoConfig models.Configuration
	companyName := ""
	companyLogo := ""

	result := configs.DB.First(&companyNameConfig, models.Configuration{
		Key:    CompanyName,
		Status: ActiveStatus,
	})

	if result.Error == nil {
		companyName = companyNameConfig.Value
	}

	result = configs.DB.First(&companyLogoConfig, models.Configuration{
		Key:    CompanyLogo,
		Status: ActiveStatus,
	})

	if result.Error == nil {
		companyLogo = companyLogoConfig.Value
	}

	to := []string{user.Email}

	t, err := template.ParseFiles("templates/invitation.html")

	if err != nil {
		return err
	}

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s - Your vKYC Access Credentials\n%s\n\n", companyName, mimeHeaders)))

	err = t.Execute(&body, map[string]string{
		"FirstName":   user.FirstName,
		"Email":       user.Email,
		"Password":    password,
		"Link":        configs.GetEnvVariable("WEB_URL"),
		"Company":     companyName,
		"CompanyLogo": companyLogo,
	})

	if err != nil {
		return err
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, smtpAuth, fromEmail, to, body.Bytes())

	if err != nil {
		return err
	}

	return nil
}
