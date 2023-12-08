package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func RateLimit() func(c *gin.Context) {
	bucket := ratelimit.NewBucket(2, 100)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			c.JSON(http.StatusOK, gin.H{
				"msg": "rate limit...",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
