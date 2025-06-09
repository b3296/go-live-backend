package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Claims 自定义 token 的内容
// 包含用户 ID、邮箱及标准字段

type Claims struct {
	UserID uint
	Email  string
	Name   string
	jwt.RegisteredClaims
}

var secret = []byte("secret-key") // 密钥

// GenerateToken 创建 token
func GenerateToken(userID uint, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(720 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken 验证并解析 token
func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
