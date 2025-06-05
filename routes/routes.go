package routes

import (
	"github.com/gin-contrib/cors"
	"time"
	"user-system/config"
	"user-system/controllers"
	"user-system/middleware"
	"user-system/service"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置所有路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// ✅ 开启 CORS 支持
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 前端地址
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.POST("/api/register", controllers.Register)
	r.POST("/api/login", controllers.Login)

	// 需要鉴权的接口组
	auth := r.Group("/api")
	auth.Use(middleware.JWTAuth())
	auth.GET("/profile", controllers.Profile)

	video := r.Group("/api/video")
	video.Use(middleware.JWTAuth()) // 需要登录权限
	{
		video.POST("/", controllers.CreateVideo)
		video.GET("/", controllers.GetVideos)
		video.GET("/:id", controllers.GetVideo)
		video.PUT("/:id", controllers.UpdateVideo)
		video.DELETE("/:id", controllers.DeleteVideo)
	}

	// 注入 CommentController
	commentService := service.NewCommentService(config.DB, config.RedisClient)
	commentController := &controllers.CommentController{
		CommentService: commentService,
	}

	comment := r.Group("/api/comment")
	comment.Use(middleware.JWTAuth())
	{
		comment.POST("/", commentController.Create)
		comment.GET("/", commentController.List)
	}

	live := r.Group("/api/live")
	liveController := controllers.NewLiveController(service.NewLiveService(config.DB))
	live.POST("/create", middleware.JWTAuth(), liveController.CreateLive)
	live.GET("/list", middleware.JWTAuth(), liveController.List)
	live.POST("/start", liveController.StartLive)
	live.POST("/stop", liveController.StopLive) // 可用于 SRS 回调或后台操作

	return r
}
