package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"vkyc-backend/configs"
	"vkyc-backend/models"
	"vkyc-backend/utils"

	"github.com/gin-gonic/gin"
)

func Auth(roles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("Authorization")

		tokenString := strings.ReplaceAll(authorizationHeader, "Bearer ", "")

		claims, err := utils.ParseToken(tokenString)

		if utils.CheckError(err, http.StatusUnauthorized, c, "Unauthorized") {
			fmt.Println("Unable to parse token")
			return
		}

		var user models.User

		result := configs.DB.First(&user, models.User{
			ID:     claims.UserID,
			Status: utils.ActiveStatus,
		})

		if utils.CheckError(result.Error, http.StatusUnauthorized, c, "Unauthorized") {
			fmt.Println("User not found")
			return
		}

		if !slices.Contains(roles, user.Role) {
			utils.CheckError(errors.New("Unauthorized"), http.StatusUnauthorized, c, "Unauthorized")
			fmt.Println("Invalid role")
			return
		}

		c.Set("authUser", user)

		c.Next()
	}
}
