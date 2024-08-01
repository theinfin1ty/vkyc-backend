package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

var clientLimiters = make(map[string]*rate.Limiter)

// var Limiter = rate.NewLimiter(rate.Limit(10), 10)

func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientID := strings.Join([]string{c.ClientIP(), c.Request.UserAgent(), c.GetHeader("Accept-Language")}, ":")

		limiter, ok := clientLimiters[clientID]

		if !ok {
			limiter = rate.NewLimiter(rate.Limit(100), 50)
			clientLimiters[clientID] = limiter
		}

		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Rate limit exceeded",
			})
			return
		}
		c.Next()
	}
}
