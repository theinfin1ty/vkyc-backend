package utils

import (
	"vkyc-backend/types"

	"github.com/gin-gonic/gin"
)

func SuccessResponse(data *gin.H) types.Response {
	var successResponse types.Response
	successResponse.Status = 200
	successResponse.Message = "Success"
	successResponse.Data = data
	return successResponse
}

func BadRequestResponse(err string) types.Response {
	var errorResponse types.Response
	errorResponse.Status = 400
	errorResponse.Message = err
	return errorResponse
}

func NotFoundResponse(err string) types.Response {
	var errorResponse types.Response
	errorResponse.Status = 404
	errorResponse.Message = err
	return errorResponse
}

func InternalServerErrorResponse(err error) types.Response {
	var errorResponse types.Response
	errorResponse.Status = 500
	errorResponse.Message = err.Error()
	return errorResponse
}

func UnauthorizedResponse(err error) types.Response {
	var errorResponse types.Response
	errorResponse.Status = 401
	errorResponse.Message = "Unauthorized"
	return errorResponse
}
