package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"user-system/utils"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil {
			bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
			requestBody = bodyBytes
			// 重新赋值 request body 否则后续读取会报错
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		// 响应捕获器
		respBody := &bytes.Buffer{}
		writer := &bodyWriter{body: respBody, ResponseWriter: c.Writer}
		c.Writer = writer

		// 执行请求
		c.Next()

		// 记录响应数据和耗时
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		utils.Log("route").Infof(
			`
==================== 接口访问日志 ====================
| Method  : %s
| Path    : %s
| Status  : %d
| Latency : %s
| Request : %s
| Response: %s
======================================================
`, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), latency, string(requestBody), respBody.String(),
		)
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *bodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b) // 写入缓冲
	return w.ResponseWriter.Write(b)
}
