package service

import (
	"errors"
	"user-system/config"
	"user-system/dto"
	"user-system/models"
	"user-system/utils"

	"golang.org/x/crypto/bcrypt"
)

// Register 用户注册逻辑
func Register(req dto.UserRegisterRequest) error {
	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return errors.New("email already exists")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{
		Email:    req.Email,
		Password: string(hashed),
		Name:     req.Name,
	}

	return config.DB.Create(&user).Error
}

// Login 登录并返回 JWT
func Login(req dto.UserLoginRequest) (string, error) {
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	return utils.GenerateToken(user.ID, user.Email)
}

// GetProfile 获取当前用户信息
func GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return dto.BuildUserResponse(user), nil
}

func GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := config.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
