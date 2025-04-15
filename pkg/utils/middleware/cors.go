/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/24 17:31
@Name: cors.go
@Description: CORS middleware for handling cross-origin requests
*/

package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CorsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 设置允许的源，可以根据需要修改为具体的域名
		// 生产环境建议不要使用"*"，而是配置具体的允许域名
		allowOrigin := "*"
		if origin != "" {
			allowOrigin = origin
		}

		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers",
			"Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Expose-Headers",
			"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 24小时

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
