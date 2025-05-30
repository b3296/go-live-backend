package tests

import (
	"testing"
	"user-system/dto"
	"user-system/models"
	"user-system/service"

	"github.com/stretchr/testify/assert"
)

func TestCreateComment(t *testing.T) {
	db := SetupTestDB()

	// 初始化 service
	commentService := service.NewCommentService(db)

	// 插入模拟用户和视频
	user := models.User{Email: "testEmail", Password: "hashed"}
	video := models.Video{Title: "test video", URL: "http://example.com/video.mp4"}

	db.Create(&user)
	db.Create(&video)

	// 构造请求
	req := dto.CommentRequest{
		Content: "hello test comment",
		VideoID: video.ID,
	}

	// 调用 service 创建评论
	comment, err := commentService.Create(req, user.ID)
	assert.Nil(t, err)
	assert.Equal(t, req.Content, comment.Content)
	assert.Equal(t, user.ID, comment.UserID)
	assert.Equal(t, video.ID, comment.VideoID)
}

func TestCreateComment_Anonymous(t *testing.T) {
	db := SetupTestDB()
	commentService := service.NewCommentService(db)

	// 构造未登录用户请求
	req := dto.CommentRequest{
		Content: "anonymous should fail",
		VideoID: 1,
	}

	// 用户ID为 0，表示未登录
	_, err := commentService.Create(req, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "未登录用户不允许评论", err.Error())
}
