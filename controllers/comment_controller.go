package controllers

import (
	"net/http"
	"user-system/dto"
	"user-system/models"
	"user-system/response"
	"user-system/service"

	"github.com/gin-gonic/gin"
)

type CommentController struct {
	CommentService *service.CommentService
}

// 创建评论
func (cc *CommentController) Create(c *gin.Context) {
	var req dto.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	var count int64
	cc.CommentService.DB.Model(&models.Video{}).Where("id = ?", req.VideoID).Count(&count)
	if count == 0 {
		response.Fail(c, http.StatusBadRequest, "评论失败：视频不存在")
		return
	}

	// 当前登录用户ID
	userID := c.GetUint("userID")

	comment, err := cc.CommentService.Create(req, userID)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "评论失败", err.Error())
		return
	}

	response.Success(c, "评论成功", dto.BuildCommentResponse(comment))
}

// 获取评论分页
func (ctl *CommentController) List(c *gin.Context) {
	var req dto.CommentListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	comments, err := ctl.CommentService.GetComments(req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获取评论失败", err.Error())
		return
	}

	// 构建返回结构
	resp := dto.BuildCommentResponseList(comments)
	response.Success(c, "获取成功", resp)
}
