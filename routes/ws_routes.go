package routes

import (
	"user-system/controllers"
	"user-system/middleware"

	"github.com/gin-gonic/gin"
)

func SetupWebSocketRoutes(r *gin.Engine) {
	r.GET("/ws/danmaku", middleware.JWTAuthWebSocket(), controllers.HandleDanmakuWebSocket)
}
