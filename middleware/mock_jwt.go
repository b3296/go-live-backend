// middleware/mock_jwt.go
package middleware

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func MockJWTWebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.Query("user_id")
		name := c.Query("name")

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			userID = 0
		}

		// 注入模拟身份信息
		c.Set("user_id", uint(userID))
		c.Set("name", name)
		c.Next()
	}
}
