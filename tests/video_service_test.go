package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"user-system/dto"
	"user-system/models"
	"user-system/service"
)

func TestCreateVideo(t *testing.T) {
	db := SetupTestDB()
	videoService := service.NewVideoService(db)

	req := dto.VideoRequest{
		Title: "Test Video",
		URL:   "http://example.com/test.mp4",
	}

	video, err := videoService.Create(req)
	assert.Nil(t, err)
	assert.Equal(t, "Test Video", video.Title)
	assert.Equal(t, "http://example.com/test.mp4", video.URL)
}

func TestGetVideos(t *testing.T) {
	db := SetupTestDB()
	videoService := service.NewVideoService(db)

	// 插入多条测试数据
	for i := 1; i <= 15; i++ {
		db.Create(&models.Video{
			Title: "Video " + string(rune(i)),
			URL:   "http://example.com/video.mp4",
		})
	}

	total, videos, err := videoService.GetVideos(1, 10)
	assert.Nil(t, err)
	assert.Equal(t, int64(15), total)
	assert.Equal(t, 10, len(videos))
}

func TestGetVideoByID(t *testing.T) {
	db := SetupTestDB()
	videoService := service.NewVideoService(db)

	video := models.Video{Title: "Find Me", URL: "http://example.com/found.mp4"}
	db.Create(&video)

	result, err := videoService.GetByID(video.ID)
	assert.Nil(t, err)
	assert.Equal(t, video.Title, result.Title)
}

func TestUpdateVideo(t *testing.T) {
	db := SetupTestDB()
	videoService := service.NewVideoService(db)

	video := models.Video{Title: "Old Title", URL: "http://old.mp4"}
	db.Create(&video)

	req := dto.VideoRequest{
		Title: "New Title",
		URL:   "http://new.mp4",
	}
	updated, err := videoService.Update(video.ID, req)
	assert.Nil(t, err)
	assert.Equal(t, "New Title", updated.Title)
	assert.Equal(t, "http://new.mp4", updated.URL)
}

func TestDeleteVideo(t *testing.T) {
	db := SetupTestDB()
	videoService := service.NewVideoService(db)

	video := models.Video{Title: "To Delete", URL: "http://delete.mp4"}
	db.Create(&video)

	err := videoService.Delete(video.ID)
	assert.Nil(t, err)

	var count int64
	db.Model(&models.Video{}).Where("id = ?", video.ID).Count(&count)
	assert.Equal(t, int64(0), count)
}
