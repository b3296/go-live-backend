package service

import (
	"errors"
	"fmt"
	"user-system/dto"
	"user-system/models"
	"user-system/utils"

	"gorm.io/gorm"
)

type VideoService struct {
	DB *gorm.DB
}

func NewVideoService(db *gorm.DB) *VideoService {
	return &VideoService{DB: db}
}

func (s *VideoService) Create(req dto.VideoRequest) (*models.Video, error) {
	video := &models.Video{}

	// 传指针给 CopyStructExt
	if err := utils.CopyStructExt(&req, video, nil); err != nil {
		fmt.Println("CopyStructExt Create error:", err)
		return nil, err
	}

	if err := s.DB.Create(video).Error; err != nil {
		return nil, err
	}

	return video, nil
}

func (s *VideoService) GetVideos(page, pageSize int) (int64, []models.Video, error) {
	var videos []models.Video
	var total int64

	db := s.DB.Model(&models.Video{})
	if err := db.Count(&total).Error; err != nil {
		return 0, nil, err
	}

	offset := (page - 1) * pageSize
	if err := db.Order("id desc").Offset(offset).Limit(pageSize).Find(&videos).Error; err != nil {
		return 0, nil, err
	}

	return total, videos, nil
}

func (s *VideoService) GetByID(id uint) (*models.Video, error) {
	var video models.Video
	if err := s.DB.First(&video, id).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

func (s *VideoService) Update(id uint, req dto.VideoRequest) (*models.Video, error) {
	video, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if err := utils.CopyStructExt(&req, video, nil); err != nil {
		fmt.Println("CopyStructExt Update error:", err)
		return nil, err
	}

	if err := s.DB.Save(video).Error; err != nil {
		return nil, err
	}

	return video, nil
}

func (s *VideoService) Delete(id uint) error {
	result := s.DB.Delete(&models.Video{}, id)
	if result.RowsAffected == 0 {
		return errors.New("记录不存在")
	}
	return result.Error
}
