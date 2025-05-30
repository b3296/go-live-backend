package controllers

import (
	"net/http"
	"strconv"
	"user-system/config"
	"user-system/dto"
	"user-system/response"
	"user-system/service"

	"github.com/gin-gonic/gin"
)

func getVideoService() *service.VideoService {
	return service.NewVideoService(config.DB)
}

func CreateVideo(c *gin.Context) {
	var req dto.VideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	videoService := getVideoService()
	video, err := videoService.Create(req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "创建失败", err.Error())
		return
	}

	response.Success(c, "创建成功", dto.BuildVideoResponse(video))
}

func GetVideos(c *gin.Context) {
	// 读取分页参数，默认1页，每页10条
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	videoService := getVideoService()
	total, videos, err := videoService.GetVideos(page, pageSize)
	if err != nil {
		response.Fail(c, 500, "查询失败", err.Error())
		return
	}

	// 构造返回结构体
	respData := struct {
		Total int64               `json:"total"`
		List  []dto.VideoResponse `json:"list"`
	}{
		Total: total,
		List:  dto.BuildVideoResponseList(videos),
	}

	response.Success(c, "查询成功", respData)
}

func GetVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	videoService := getVideoService()
	video, err := videoService.GetByID(uint(id))
	if err != nil {
		response.Fail(c, http.StatusNotFound, "视频不存在", err.Error())
		return
	}

	response.Success(c, "获取成功", dto.BuildVideoResponse(video))
}

func UpdateVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	var req dto.VideoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, "参数错误", err.Error())
		return
	}

	videoService := getVideoService()
	video, err := videoService.Update(uint(id), req)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, "更新失败", err.Error())
		return
	}

	response.Success(c, "更新成功", dto.BuildVideoResponse(video))
}

func DeleteVideo(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	videoService := getVideoService()
	if err := videoService.Delete(uint(id)); err != nil {
		response.Fail(c, http.StatusInternalServerError, "删除失败", err.Error())
		return
	}

	response.Success(c, "删除成功")
}
