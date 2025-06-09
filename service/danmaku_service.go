package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"time"
	"user-system/config"
)

type Client struct {
	Conn   *websocket.Conn
	LiveID uint
	UserID uint
	Name   string
}

type BroadcastMessage struct {
	UserID  uint   `json:"user_id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	LiveID  uint   `json:"live_id"`
}

type DanmakuHub struct {
	LiveID     uint
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan BroadcastMessage
	Mutex      sync.Mutex
}

var hubRegistry = make(map[uint]*DanmakuHub)
var hubMutex sync.Mutex

func GetDanmakuHub(liveID uint) *DanmakuHub {
	hubMutex.Lock()
	defer hubMutex.Unlock()
	if hub, exists := hubRegistry[liveID]; exists {
		return hub
	}
	hub := &DanmakuHub{
		LiveID:     liveID,
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan BroadcastMessage),
	}
	go hub.run()
	hubRegistry[liveID] = hub
	return hub
}

func (h *DanmakuHub) run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client] = true
			h.Mutex.Unlock()
			go h.handleClient(client)
		case client := <-h.Unregister:
			h.Mutex.Lock()
			delete(h.Clients, client)
			h.Mutex.Unlock()
			client.Conn.Close()
		case msg := <-h.Broadcast:
			data, _ := json.Marshal(msg)
			h.Mutex.Lock()
			for c := range h.Clients {
				_ = c.Conn.WriteMessage(websocket.TextMessage, data)
			}
			h.Mutex.Unlock()
			SaveDanmakuToRedis(msg.UserID, msg.LiveID, msg.Content)
		}
	}
}

func (h *DanmakuHub) handleClient(c *Client) {
	defer func() {
		h.Unregister <- c
	}()
	for {
		var msg struct {
			Content string `json:"content"`
		}
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			break
		}
		h.Broadcast <- BroadcastMessage{
			UserID:  c.UserID,
			Name:    c.Name,
			Content: msg.Content,
			LiveID:  c.LiveID,
		}
	}
}

func SaveDanmakuToRedis(userID uint, liveID uint, content string) {
	key := fmt.Sprintf("danmaku:%d", liveID)
	msg := fmt.Sprintf("[%d] %s", userID, content)
	_ = config.RedisClient.RPush(context.Background(), key, msg).Err()
	config.RedisClient.Expire(context.Background(), key, time.Hour*6)
}
