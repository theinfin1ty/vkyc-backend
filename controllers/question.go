package controllers

import (
	"net/http"
	"slices"
	"strconv"
	"vkyc-backend/configs"
	"vkyc-backend/models"
	"vkyc-backend/types"
	"vkyc-backend/utils"

	"github.com/gin-gonic/gin"
)

func CreateQuestion(c *gin.Context) {
	var body types.QuestionInput
	var existingQuestion models.Question

	err := c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	if !slices.Contains([]string{utils.ActiveStatus, utils.InactiveStatus}, body.Status) {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Invalid Status"))
		return
	}

	result := configs.DB.First(&existingQuestion, models.Question{
		Question:  body.Question,
		IsDeleted: false,
	})

	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Question already exists"))
		return
	}

	question := models.Question{
		Question:  body.Question,
		Status:    body.Status,
		IsDeleted: false,
	}

	result = configs.DB.Create(&question)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"question": question,
	}))
}

func ListQuestions(c *gin.Context) {
	var questions []models.Question

	result := configs.DB.
		Model(&models.Question{}).
		Where("is_deleted = ?", false).
		Find(&questions)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"questions": questions,
	}))
}

func GetQuestion(c *gin.Context) {
	var question models.Question

	questionId, err := strconv.Atoi(c.Param("questionId"))

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	result := configs.DB.First(&question, models.Question{
		ID:        uint(questionId),
		IsDeleted: false,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Question Not Found") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"question": question,
	}))
}

func UpdateQuestion(c *gin.Context) {
	var body types.QuestionInput
	var question, existingQuestion models.Question

	questionId, err := strconv.Atoi(c.Param("questionId"))

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	err = c.BindJSON(&body)

	if utils.CheckError(err, http.StatusBadRequest, c, "Validation Failed") {
		return
	}

	result := configs.DB.First(&question, models.Question{
		ID:        uint(questionId),
		IsDeleted: false,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Question Not Found") {
		return
	}

	result = configs.DB.
		Model(&models.Question{}).
		Where("is_deleted = ?", false).
		Where("id != ?", questionId).
		Where("question = ?", body.Question).
		Find(&existingQuestion)

	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, utils.BadRequestResponse("Question already exists"))
		return
	}

	question.Question = body.Question
	question.Status = body.Status

	result = configs.DB.Save(&question)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"question": question,
	}))
}

func DeleteQuestion(c *gin.Context) {
	var question models.Question

	questionId, err := strconv.Atoi(c.Param("questionId"))

	if utils.CheckError(err, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	result := configs.DB.First(&question, models.Question{
		ID:        uint(questionId),
		IsDeleted: false,
	})

	if utils.CheckError(result.Error, http.StatusNotFound, c, "Question Not Found") {
		return
	}

	question.IsDeleted = true

	result = configs.DB.Save(&question)

	if utils.CheckError(result.Error, http.StatusInternalServerError, c, "Internal Server Error") {
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse(&gin.H{
		"message": "Question Deleted",
	}))
}
