package dto

import (
	"user-system/models"
	"user-system/utils"
)

// UserRegisterRequest 注册请求结构体
type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Name     string `json:"name" binding:"required"`
}

// UserLoginRequest 登录请求结构体
type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// UserResponse 用户响应结构体（隐藏敏感信息）
type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"` // string 类型
	UpdatedAt string `json:"updated_at"`
}

func BuildUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: utils.FormatTime(user.CreatedAt),
		UpdatedAt: utils.FormatTime(user.UpdatedAt),
	}
}
