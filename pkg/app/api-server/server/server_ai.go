package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
	"zhku-oj-server/pkg/utils"
)

// proxyRequest 通用的代理请求方法
func (s *Server) proxyRequest(c *gin.Context, path string) {
	lg := utils.GetDefaultLogger()
	cfg, _ := utils.LoadConfig("conf/config.yaml")
	difyCfg := cfg.GetDifyConfig()
	apiKey := difyCfg.APIKey     // D从环境变量读取中间层API密钥
	difyAPIURL := difyCfg.APIURL // Dify API地址

	// 1. 构建目标URL
	targetURL := difyAPIURL + path
	if c.Request.URL.RawQuery != "" {
		targetURL += "?" + c.Request.URL.RawQuery
	}
	lg.Info("转发请求至: ", targetURL)

	// 2. 创建转发请求（携带原始请求体）
	proxyReq, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		lg.Error("创建转发请求失败: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// 3. 复制原始请求头（跳过前端传入的Authorization）
	for name, values := range c.Request.Header {
		if name == "Authorization" {
			continue // 跳过前端传入的Authorization
		}
		for _, value := range values {
			proxyReq.Header.Add(name, value)
		}
	}

	// 4. 设置必要的Headers
	proxyReq.Header.Set("Authorization", "Bearer "+apiKey)
	proxyReq.Header.Set("Content-Type", "application/json")

	// 5. 发送请求
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(proxyReq)
	if err != nil {
		lg.Error("转发请求失败: ", err)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Bad Gateway"})
		return
	}
	defer resp.Body.Close()

	// 6. 将响应直接返回给前端
	c.Status(resp.StatusCode)
	for name, values := range resp.Header {
		for _, value := range values {
			c.Writer.Header().Add(name, value)
		}
	}
	_, err = io.Copy(c.Writer, resp.Body)
	if err != nil {
		lg.Error("复制响应体失败: ", err)
	}
}

// Chat 发起对话 /chat-messages
func (s *Server) Chat(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("发起对话......")
	s.proxyRequest(c, "/chat-messages")
}

// GetMessage 查询对话的信息 /message
func (s *Server) GetMessage(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("查询对话的信息......")
	s.proxyRequest(c, "/messages")
}

// GetConversations 查询对话记录 /conversations
func (s *Server) GetConversations(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("查询对话记录......")
	s.proxyRequest(c, "/conversations")
}

// RenameConversation 重命名对话 /conversations/:conversation_id/name
func (s *Server) RenameConversation(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("重命名对话......")

	// 获取路径参数
	conversationID := c.Param("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversation_id is required"})
		return
	}

	s.proxyRequest(c, fmt.Sprintf("/conversations/%s/name", conversationID))
}

// DeleteConversation 删除对话 /conversations/:conversation_id
func (s *Server) DeleteConversation(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("删除对话......")

	// 获取路径参数
	conversationID := c.Param("conversation_id")
	if conversationID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "conversation_id is required"})
		return
	}

	s.proxyRequest(c, fmt.Sprintf("/conversations/%s", conversationID))
}

// 复制请求头
func cloneHeader(src http.Header) http.Header {
	dst := make(http.Header)
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
	return dst
}

// 复制响应头
func copyHeader(dst, src http.Header) {
	for k, vs := range src {
		for _, v := range vs {
			dst.Add(k, v)
		}
	}
}
