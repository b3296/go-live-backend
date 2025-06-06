package main

import (
	"log"
	"user-system/config"
	"user-system/routes"
)

func main() {
	// 初始化配置、数据库、Redis
	config.Init()

	config.SetupLogger()
	// 设置 Gin 路由
	r := routes.SetupRouter()

	routes.SetupWebSocketRoutes(r)

	log.Println("Server started at http://localhost:8888")
	log.Println("WebSocket endpoint at ws://localhost:8888/ws/:room_id")

	// 启动服务监听端口 8888
	r.Run(":8888")
}
