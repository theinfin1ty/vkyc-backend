package controllers

import (
	"errors"
	"net/http"
	"slices"
	"strconv"
	"vkyc-backend/configs"
	"vkyc-backend/models"
	"vkyc-backend/types"
	"vkyc-backend/utils"

	"github.com/gin-gonic/gin"
)

func AddUser(c *gin.Context) {
	var body types.AddUserInput
	var user models.User

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	result := configs.DB.Find(&user, models.User{
		Email: body.Email,
	}).Limit(1)

	if utils.CheckError(result.Error, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Email already exists"))
		return
	}

	if utils.CheckError(err, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	password := utils.GeneratePassword(8, true)

	hashedPassword, err := utils.GeneratePasswordHash(password)

	if utils.CheckError(err, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	user = models.User{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Password:  hashedPassword,
		Role:      body.Role,
		// DOB:       time.Now(),
	}

	if body.CountryCode != "" && body.Phone != "" {
		user.CountryCode = body.CountryCode
		user.Phone = body.Phone
	}

	result = configs.DB.Create(&user)

	if utils.CheckError(result.Error, http.StatusBadRequest, c, "Invalid Email or Password") {
		return
	}

	err = utils.SendInviteEmail(&user, password)

	if utils.CheckError(err, http.StatusBadRequest, c, "Error occured while sending invite email") {
		return
	}

	// err = utils.SendInviteSMS(&user, password)

	// if utils.CheckError(err, http.StatusBadRequest, c, "Error occured while sending invite sms") {
	// 	return
	// }

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"user": user,
	}))
}

func ListUsers(c *gin.Context) {
	var users []models.User

	// search := c.Query("search")
	role := c.Query("role")

	// if c.Query("page") == "" || c.Query("size") == "" || c.Query("role") == "" {
	// 	c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Page, Size and Role are required parameters"))
	// 	return
	// }

	// size, err := strconv.Atoi(c.Query("size"))

	// if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
	// 	return
	// }

	// page, err := strconv.Atoi(c.Query("page"))

	// if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
	// 	return
	// }

	// skip := size * (page - 1)

	result := configs.DB.
		Model(&users).
		// Where("first_name ILIKE ? OR last_name ILIKE ?", "%"+search+"%", "%"+search+"%").
		Where("role = ?", role).
		Order("first_name ASC").
		Find(&users)
		// .Limit(size).
		// Offset(skip)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"users": users,
	}))
}

func GetUser(c *gin.Context) {
	var user models.User

	userId, err := strconv.Atoi(c.Param("id"))

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	result := configs.DB.First(&user, models.User{
		ID: uint(userId),
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "User Not Found") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"user": user,
	}))
}

func UpdateUser(c *gin.Context) {
	var body types.UpdateUserInput
	var user models.User

	userId, err := strconv.Atoi(c.Param("id"))

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	err = c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	if body.Email != "" {
		result := configs.DB.
			Model(&models.User{}).
			Where("email = ?", body.Email).
			Where("id != ?", userId).
			Find(&user)

		if result.RowsAffected > 0 {
			c.JSON(http.StatusBadRequest, utils.BadRequestResponse("User with given email already exists"))
			return
		}
	}

	result := configs.DB.First(&user, models.User{
		ID: uint(userId),
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "User Not Found") {
		return
	}

	// dob, err := time.Parse(utils.DateLayout, body.DOB)

	// if utils.CheckError(err, http.StatusBadRequest, c, "Invalid DOB") {
	// 	return
	// }

	user.FirstName = body.FirstName
	user.LastName = body.LastName
	user.Email = body.Email
	user.CountryCode = body.CountryCode
	user.Phone = body.Phone
	user.Department = body.Department
	// user.DOB = dob

	result = configs.DB.Save(&user)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"user": user,
	}))
}

func UpdateProfile(c *gin.Context) {
	var body types.UpdateUserInput
	authUser, _ := c.Get("authUser")

	user := authUser.(models.User)

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	// dob, err := time.Parse(utils.DateLayout, body.DOB)

	// if utils.CheckError(err, http.StatusBadRequest, c, "Invalid DOB") {
	// 	return
	// }

	user.FirstName = body.FirstName
	user.LastName = body.LastName
	user.Email = body.Email
	user.CountryCode = body.CountryCode
	user.Phone = body.Phone
	// user.DOB = dob

	result := configs.DB.Save(&user)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"user": user,
	}))
}

func UpdateUserStatus(c *gin.Context) {
	var user models.User
	var body map[string]string

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	if !slices.Contains([]string{"active", "blocked"}, body["status"]) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid Status"))
		return
	}

	userId, err := strconv.Atoi(c.Param("id"))

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	result := configs.DB.First(&user, models.User{
		ID: uint(userId),
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "User Not Found") {
		return
	}

	user.Status = body["status"]

	result = configs.DB.Save(&user)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "User Status Updated",
	}))
}

func DeleteUser(c *gin.Context) {
	var user models.User

	userId, err := strconv.Atoi(c.Param("id"))

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	result := configs.DB.First(&user, models.User{
		ID: uint(userId),
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "User Not Found") {
		return
	}

	user.IsDeleted = true

	result = configs.DB.Save(&user)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "User Deleted",
	}))
}

func Dashboard(c *gin.Context) {
	var totalMeetings,
		completedMeetings,
		pendingMeetings,
		meetingsEndedAfterConnecting,
		meetingsEndedAfterQA,
		meetingsEndedAfterOCR,
		meetingsEndedAfterSign,
		meetingsEndedAfterFacematch,
		kycApproved,
		kycRejected int64

	authUser, authUserExists := c.Get("authUser")

	if !authUserExists {
		utils.CheckError(errors.New("authUser not found"), http.StatusInternalServerError, c, "User not authenticated")
		return
	}

	user := authUser.(models.User)

	if user.Role == "" {
		utils.CheckError(errors.New("user role not found"), http.StatusInternalServerError, c, "User role not found")
		return
	}

	from := c.Query("from")
	to := c.Query("to")

	// now := time.Now()

	// if from == "" || to == "" {
	// 	from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format(utils.DateLayout)
	// 	to = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999, now.Location()).Format(utils.DateLayout)
	// }

	if user.Role == utils.AdminRole {
		totalMeetingsQuery := configs.DB.Model(&models.Meeting{}).Where("is_deleted = false")

		if from != "" || to != "" {
			totalMeetingsQuery.Where("schedule_date_time between ? and ?", from, to)
		}

		totalMeetingsQuery.Count(&totalMeetings)

		if utils.CheckError(totalMeetingsQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		completedMeetingsQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("is_ended = true").
			Where("is_deleted = false")

		if from != "" || to != "" {
			completedMeetingsQuery.Where("schedule_date_time between ? and ?", from, to)
		}

		completedMeetingsQuery.Count(&completedMeetings)

		if utils.CheckError(completedMeetingsQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		pendingMeetingsQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("is_ended = false").
			Where("is_deleted = false")

		if from != "" || to != "" {
			pendingMeetingsQuery.Where("schedule_date_time between ? and ?", from, to)
		}

		pendingMeetingsQuery.Count(&pendingMeetings)

		if utils.CheckError(pendingMeetingsQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		kycApprovedQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("kyc_status = ?", utils.ApprovedStatus).
			Where("is_deleted = false")
		if from != "" || to != "" {
			kycApprovedQuery.Where("schedule_date_time between ? and ?", from, to)
		}
		kycApprovedQuery.Count(&kycApproved)

		if utils.CheckError(kycApprovedQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		kycRejectedQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("kyc_status = ?", utils.RejectedStatus).
			Where("is_deleted = false")
		if from != "" || to != "" {
			kycRejectedQuery.Where("schedule_date_time between ? and ?", from, to)
		}
		kycRejectedQuery.Count(&kycRejected)

		if utils.CheckError(kycRejectedQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		meetingsEndedAfterQAQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("is_deleted = false").
			Where("remarks = ?", utils.QAFailedRemark)
		if from != "" || to != "" {
			meetingsEndedAfterQAQuery.Where("schedule_date_time between ? and ?", from, to)
		}
		meetingsEndedAfterQAQuery.Count(&meetingsEndedAfterQA)

		if utils.CheckError(meetingsEndedAfterQAQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		meetingsEndedAfterOCRQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("is_deleted = false").
			Where("remarks = ?", utils.OCRFailedRemark)
		if from != "" || to != "" {
			meetingsEndedAfterOCRQuery.Where("schedule_date_time between ? and ?", from, to)
		}
		meetingsEndedAfterOCRQuery.Count(&meetingsEndedAfterOCR)

		if utils.CheckError(meetingsEndedAfterOCRQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		meetingsEndedAfterFacematchQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("is_deleted = false").
			Where("remarks = ?", utils.FacematchFailedRemark)
		if from != "" || to != "" {
			meetingsEndedAfterFacematchQuery.Where("schedule_date_time between ? and ?", from, to)
		}
		meetingsEndedAfterFacematchQuery.Count(&meetingsEndedAfterFacematch)

		if utils.CheckError(meetingsEndedAfterFacematchQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		meetingsEndedAfterSignQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("is_deleted = false").
			Where("remarks = ?", utils.SignFailedRemark)
		if from != "" || to != "" {
			meetingsEndedAfterSignQuery.Where("schedule_date_time between ? and ?", from, to)
		}
		meetingsEndedAfterSignQuery.Count(&meetingsEndedAfterSign)

		if utils.CheckError(meetingsEndedAfterSignQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		meetingsEndedAfterConnectingQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("is_deleted = false").
			Where("remarks = ?", utils.FailedAfterConnectingRemark)
		if from != "" || to != "" {
			meetingsEndedAfterConnectingQuery.Where("schedule_date_time between ? and ?", from, to)
		}
		meetingsEndedAfterConnectingQuery.Count(&meetingsEndedAfterConnecting)

		if utils.CheckError(meetingsEndedAfterConnectingQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
			"totalMeetings":                totalMeetings,
			"completedMeetings":            completedMeetings,
			"pendingMeetings":              pendingMeetings,
			"meetingsEndedAfterConnecting": meetingsEndedAfterConnecting,
			"meetingsEndedAfterQA":         meetingsEndedAfterQA,
			"meetingsEndedAfterOCR":        meetingsEndedAfterOCR,
			"meetingsEndedAfterSign":       meetingsEndedAfterSign,
			"meetingsEndedAfterFacematch":  meetingsEndedAfterFacematch,
			"kycApproved":                  kycApproved,
			"kycRejected":                  kycRejected,
		}))
		return
	}

	if user.Role == utils.AgentRole {
		totalMeetingsQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("agent_id = ?", user.ID).
			Where("is_deleted = false")
		if from != "" || to != "" {
			totalMeetingsQuery.Where("schedule_date_time between ? and ?", from, to)
		}
		totalMeetingsQuery.Count(&totalMeetings)

		if utils.CheckError(totalMeetingsQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		completedMeetingsQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("agent_id = ? ", user.ID).
			Where("is_ended = true").
			Where("is_deleted = false")
		if from != "" || to != "" {
			completedMeetingsQuery.Where("schedule_date_time between ? and ?", from, to)
		}
		completedMeetingsQuery.Count(&completedMeetings)

		if utils.CheckError(completedMeetingsQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		pendingMeetingsQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("agent_id = ?", user.ID).
			Where("is_ended = false").
			Where("is_deleted = false")
		if from != "" || to != "" {
			pendingMeetingsQuery.Where("schedule_date_time between ? and ?", from, to)
		}
		pendingMeetingsQuery.Count(&pendingMeetings)

		if utils.CheckError(pendingMeetingsQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}

		kycRejectedQuery := configs.DB.
			Model(&models.Meeting{}).
			Where("agent_id = ?", user.ID).
			Where("kyc_status = ?", utils.RejectedStatus).
			Where("is_deleted = false")
		if from != "" || to != "" {
			kycRejectedQuery.Where("schedule_date_time between ? and ?", from, to)
		}
		kycRejectedQuery.Count(&kycRejected)

		if utils.CheckError(kycRejectedQuery.Error, http.StatusInternalServerError, c, "Internal Server Error") {
			return
		}
		c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
			"totalMeetings":     totalMeetings,
			"completedMeetings": completedMeetings,
			"pendingMeetings":   pendingMeetings,
			"kycRejected":       kycRejected,
		}))
	}
}
