package controllers

import (
	"net/http"
	"strconv"
	"user-system/service"
	"user-system/utils"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type DanmakuMessage struct {
	LiveID  uint   `json:"live_id"`
	Content string `json:"content"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleDanmakuWebSocket(c *gin.Context) {
	liveIDStr := c.Query("live_id")
	userID := c.GetUint("user_id")
	name := c.GetString("Name") // 从 JWT middleware 获取昵称

	liveID, err := strconv.ParseUint(liveIDStr, 10, 64)
	if err != nil || liveID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid live_id"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		utils.Log("danmaku").Errorf("WebSocket upgrade error:s%", err)
		return
	}

	hub := service.GetDanmakuHub(uint(liveID))
	hub.Register <- &service.Client{
		Conn:   conn,
		LiveID: uint(liveID),
		UserID: userID,
		Name:   name,
	}
}
