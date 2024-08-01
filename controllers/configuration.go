package controllers

import (
	"net/http"
	"slices"
	"vkyc-backend/configs"
	"vkyc-backend/models"
	"vkyc-backend/types"
	"vkyc-backend/utils"

	"github.com/gin-gonic/gin"
)

func GetConfiguration(c *gin.Context) {
	var configurations []models.Configuration

	result := configs.DB.Find(&configurations)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"configs": configurations,
	}))
}

func GetPublicConfiguration(c *gin.Context) {
	var configurations []models.Configuration

	result := configs.DB.
		Model(&models.Configuration{}).
		Where("status = ?", utils.ActiveStatus).
		Where("`key` IN ?", []string{utils.CompanyName, utils.CompanyLogo, utils.CompanyLogo2}).
		Find(&configurations)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"configs": configurations,
	}))
}

func UpdateConfiguration(c *gin.Context) {
	var body types.ConfigurationInput
	var configuration models.Configuration

	allowedKeys := []string{
		utils.FacematchAPI,
		utils.FacematchAPIKey,
		utils.LivenessAPI,
		utils.LivenessAPIKey,
		utils.OCRAPI,
		utils.OCRAPIKey,
		utils.SMSAPI,
		utils.SMSAPIKey,
		utils.SMSSender,
		utils.SMSTemplate,
		utils.MeetingExpireBeforeStart,
		utils.MeetingExpireAfterStart,
		utils.MeetingTimeGap,
		utils.MinFacematchScore,
		utils.MinLivenessScore,
		utils.CompanyName,
		utils.CompanyLogo,
		utils.CompanyLogo2,
	}

	allowedStatuses := []string{
		utils.ActiveStatus,
		utils.InactiveStatus,
	}

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	if !slices.Contains(allowedKeys, body.Key) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid Key"))
		return
	}

	if !slices.Contains(allowedStatuses, body.Status) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid Status"))
		return
	}

	configs.DB.First(&configuration, models.Configuration{
		Key: body.Key,
	})

	configuration.Name = body.Name
	configuration.Key = body.Key
	configuration.Value = body.Value
	configuration.Description = body.Description
	configuration.Status = body.Status

	result := configs.DB.Save(&configuration)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"config": configuration,
	}))
}
