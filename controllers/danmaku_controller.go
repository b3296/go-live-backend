package controllers

import (
	"fmt"
	"net/http"
	"user-system/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type DanmakuMessage struct {
	VideoID uint   `json:"video_id"`
	Content string `json:"content"`
}

// 简化实现，全局广播池（后续可拆为 map[video_id][]conn）
var danmakuRooms = make(map[uint][]*websocket.Conn)

func HandleDanmakuWebSocket(c *gin.Context) {
	videoID := c.Query("video_id")
	userID := c.GetUint("user_id")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// 简单解析 video_id
	var vid uint
	fmt.Sscanf(videoID, "%d", &vid)
	danmakuRooms[vid] = append(danmakuRooms[vid], conn)

	for {
		var msg DanmakuMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("read error:", err)
			break
		}

		// 发送给该视频的所有连接
		for _, client := range danmakuRooms[msg.VideoID] {
			_ = client.WriteJSON(gin.H{
				"user_id":  userID,
				"content":  msg.Content,
				"video_id": msg.VideoID,
			})
		}

		// 可选：存入 Redis
		go service.SaveDanmakuToRedis(userID, msg.VideoID, msg.Content)
	}
}
