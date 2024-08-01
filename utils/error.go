package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckError(err error, status int, c *gin.Context, message string) bool {
	if err != nil {
		fmt.Println(message, err)
		switch status {
		case http.StatusBadRequest:
			c.AbortWithStatusJSON(http.StatusBadRequest, BadRequestResponse(message))
			return true
		case http.StatusNotFound:
			c.AbortWithStatusJSON(http.StatusNotFound, NotFoundResponse(message))
			return true
		case http.StatusInternalServerError:
			c.AbortWithStatusJSON(http.StatusInternalServerError, InternalServerErrorResponse(err))
			return true
		case http.StatusUnauthorized:
			c.AbortWithStatusJSON(http.StatusUnauthorized, UnauthorizedResponse(err))
			return true
		default:
			return false
		}
	}
	return false
}
