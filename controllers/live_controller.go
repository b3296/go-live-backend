package controllers

import (
	"net/http"
	"strconv"
	"user-system/dto"
	"user-system/response"
	"user-system/service"

	"github.com/gin-gonic/gin"
)

type LiveController struct {
	LiveService *service.LiveService
}

func NewLiveController(liveService *service.LiveService) *LiveController {
	return &LiveController{LiveService: liveService}
}

func (ctl *LiveController) List(c *gin.Context) {
	var req dto.LiveListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	lives, err := ctl.LiveService.GetList(req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "获直播列表论失败", err.Error())
		return
	}

	// 构建返回结构
	resp := dto.BuildLiveResponseList(lives)
	response.Success(c, "获取成功", resp)
}

func (c *LiveController) CreateLive(ctx *gin.Context) {
	var req dto.CreateLiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	uid := ctx.GetUint("userID")
	live, err := c.LiveService.CreateLive(uid, req.Title, req.CoverURL)
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "创建失败", err)
		return
	}

	response.Success(ctx, "success", gin.H{
		"live_id":    live.ID,
		"stream_key": live.StreamKey,
	})
}

// SRS 推流成功后回调：通知后端直播开始
func (c *LiveController) StartLive(ctx *gin.Context) {
	var req dto.LiveStartRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}

	err := c.LiveService.StartLive(req.Stream)
	if err != nil {
		response.Fail(ctx, http.StatusInternalServerError, "启动直播失败", err.Error())
		return
	}

	response.Success(ctx, "直播开始成功")
}

func (c *LiveController) StopLive(ctx *gin.Context) {
	var req dto.StopLiveRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Fail(ctx, http.StatusBadRequest, "参数错误")
		return
	}
	streamKey := req.StreamKey
	forced, _ := strconv.ParseBool(ctx.Query("forced"))
	crashed, _ := strconv.ParseBool(ctx.Query("crashed"))

	if err := c.LiveService.StopLive(streamKey, forced, crashed); err != nil {
		response.Fail(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response.Success(ctx, "success")
}
