package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"user-system/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth 中间件：验证 token 合法性
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// 拆分 "Bearer xxx" 格式
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization format must be Bearer {token}"})
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := utils.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}
		fmt.Println(claims)
		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("name", claims.Name)

		c.Next()
	}
}

// JWTAuthWebSocket 处理 WebSocket 的 JWT 校验
func JWTAuthWebSocket() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing token"})
			return
		}
		claims, err := utils.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("name", claims.Name)
		c.Next()
	}
}
