package controllers

import (
	"encoding/base64"
	"net/http"
	"strings"
	"time"
	"vkyc-backend/configs"
	"vkyc-backend/models"
	"vkyc-backend/types"
	"vkyc-backend/utils"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var body types.LoginInput
	var user models.User

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	body.Email = strings.ToLower(body.Email)

	result := configs.DB.First(&user, models.User{
		Email: body.Email,
	})

	if utils.CheckError(result.Error, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	password, err := base64.StdEncoding.DecodeString(body.Password)

	if utils.CheckError(err, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	err = utils.CheckPasswordHash(string(password), user.Password)

	if utils.CheckError(err, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	tokens, err := utils.GenerateAuthTokens(&user)

	if utils.CheckError(err, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"user":   user,
		"tokens": tokens,
	}))
}

func ForgotPassword(c *gin.Context) {
	var body types.ForgotPasswordInput
	var user models.User

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	body.Email = strings.ToLower(body.Email)

	result := configs.DB.First(&user, models.User{
		Email: body.Email,
	})

	if utils.CheckError(result.Error, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	otp := utils.GeneratePassword(6, false)

	otpRecord := models.OTP{
		OTP:       otp,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * time.Minute),
	}

	result = configs.DB.Save(&otpRecord)

	if utils.CheckError(result.Error, http.StatusBadRequest, c, "Unable to save OTP") {
		return
	}

	err = utils.SendOTPEmail(otp, &user)

	if utils.CheckError(err, http.StatusBadRequest, c, "Unable to send OTP") {
		return
	}

	err = utils.SendOTPSMS(otp, &user)

	if utils.CheckError(err, http.StatusBadRequest, c, "Unable to send OTP") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "OTP sent successfully",
	}))
}

func ResetPassword(c *gin.Context) {
	var body types.ResetPasswordInput
	var otpRecord models.OTP
	var user models.User

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	body.Email = strings.ToLower(body.Email)

	result := configs.DB.First(&user, models.User{
		Email: body.Email,
	})

	if utils.CheckError(result.Error, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	result = configs.DB.First(&otpRecord, "otp = ? AND user_id = ? AND expires_at >= ?", body.OTP, user.ID, time.Now())

	if utils.CheckError(result.Error, http.StatusBadRequest, c, "Invalid OTP") {
		return
	}

	decodedPassword, err := base64.StdEncoding.DecodeString(body.Password)

	if utils.CheckError(err, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	body.Password, err = utils.GeneratePasswordHash(string(decodedPassword))

	if utils.CheckError(err, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	user.Password = body.Password

	result = configs.DB.Save(&user)

	if utils.CheckError(result.Error, http.StatusBadRequest, c, "Unable to save user") {
		return
	}

	result = configs.DB.Delete(&models.OTP{}, "user_id = ?", user.ID)

	if utils.CheckError(result.Error, http.StatusBadRequest, c, "Unable to delete OTP") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "Password reset successfully",
	}))
}

func ChangePassword(c *gin.Context) {
	var body types.ChangePasswordInput

	authUser, _ := c.Get("authUser")

	user := authUser.(models.User)

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	currentPassword, err := base64.StdEncoding.DecodeString(body.CurrentPassword)

	if utils.CheckError(err, http.StatusBadRequest, c, "Invalid Password") {
		return
	}

	err = utils.CheckPasswordHash(string(currentPassword), user.Password)

	if utils.CheckError(err, http.StatusBadRequest, c, "Invalid Password") {
		return
	}

	newPassword, err := base64.StdEncoding.DecodeString(body.NewPassword)

	if utils.CheckError(err, http.StatusBadRequest, c, "Invalid Password") {
		return
	}

	body.NewPassword, err = utils.GeneratePasswordHash(string(newPassword))

	if utils.CheckError(err, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	user.Password = body.NewPassword

	result := configs.DB.Save(&user)

	if utils.CheckError(result.Error, http.StatusBadRequest, c, "Unable to save user") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "Password changed successfully",
	}))
}
