package tests

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"
	"time"
	"user-system/config"

	"user-system/controllers"
	"user-system/middleware"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func startTestServer() *httptest.Server {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/ws/danmaku", middleware.MockJWTWebSocket(), controllers.HandleDanmakuWebSocket)
	return httptest.NewServer(r)
}

func createWSConn(t *testing.T, serverURL string, userID int, name string, liveID int) *websocket.Conn {
	u := url.URL{
		Scheme: "ws",
		Host:   serverURL,
		Path:   "/ws/danmaku",
		RawQuery: fmt.Sprintf("user_id=%d&name=%s&live_id=%d",
			userID, name, liveID),
	}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	assert.NoError(t, err)
	return conn
}

func TestDanmaku(t *testing.T) {
	config.InitRedis()
	server := startTestServer()
	defer server.Close()
	serverAddr := server.Listener.Addr().String()

	type userInfo struct {
		Conn   *websocket.Conn
		UserID int
		Name   string
		LiveID int
	}

	users := []userInfo{
		{UserID: 1, Name: "Alice", LiveID: 101},
		{UserID: 2, Name: "Bob", LiveID: 101},
		{UserID: 3, Name: "Carol", LiveID: 102},
		{UserID: 1, Name: "Alice", LiveID: 102},
	}

	for i := range users {
		users[i].Conn = createWSConn(t, serverAddr, users[i].UserID, users[i].Name, users[i].LiveID)
		defer users[i].Conn.Close()
	}

	// 发送弹幕（Alice 在 101）
	sendMsg := map[string]interface{}{
		"live_id": 101,
		"content": "Hello from Alice!",
	}
	assert.NoError(t, users[0].Conn.WriteJSON(sendMsg))

	// 等待广播
	time.Sleep(200 * time.Millisecond)

	// 接收弹幕（101 的用户都应收到）
	type recv struct {
		UserID  int    `json:"user_id"`
		Name    string `json:"name"`
		LiveID  int    `json:"live_id"`
		Content string `json:"content"`
	}

	var wg sync.WaitGroup
	for _, u := range users {
		if u.LiveID == 101 {
			wg.Add(1)
			go func(conn *websocket.Conn, expectedLiveID int) {
				defer wg.Done()
				_, msg, err := conn.ReadMessage()
				assert.NoError(t, err)
				var res recv
				_ = json.Unmarshal(msg, &res)
				assert.Equal(t, expectedLiveID, res.LiveID)
				assert.Equal(t, "Hello from Alice!", res.Content)
			}(u.Conn, u.LiveID)
		}
	}
	wg.Wait()
}
