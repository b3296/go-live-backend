package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Success 成功响应（带可选 data）
func Success(c *gin.Context, msg string, data ...interface{}) {
	resp := gin.H{
		"code": 200,
		"msg":  msg,
	}
	if len(data) > 0 {
		resp["data"] = data[0]
	} else {
		resp["data"] = nil
	}
	c.JSON(http.StatusOK, resp)
}

// Fail 失败响应（带可选 data）
func Fail(c *gin.Context, code int, msg string, data ...interface{}) {
	resp := gin.H{
		"code": code,
		"msg":  msg,
	}
	if len(data) > 0 {
		resp["data"] = data[0]
	} else {
		resp["data"] = nil
	}
	c.JSON(http.StatusOK, resp)
}

// SrsHookSuccess 返回符合 SRS HTTP Hook 要求的成功响应
func SrsHookSuccess(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
	})
}
