package middleware

import (
	"net/http"
	"user-system/response"

	"github.com/gin-gonic/gin"
	"log"
)

// Recovery 捕获 panic，返回 500 错误，防止服务崩溃
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic recovered: %v", err)
				response.Fail(c, 500, "服务器内部错误")
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
