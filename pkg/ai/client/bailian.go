package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// BailianClient 阿里云千炼客户端实现
type BailianClient struct {
	config     *ClientConfig
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewBailianClient 创建千炼客户端
func NewBailianClient(config *ClientConfig) (*BailianClient, error) {
	if config.APIKey == "" {
		return nil, &Error{
			Code:    "missing_api_key",
			Message: "API密钥不能为空",
			Type:    ErrorTypeAuth,
		}
	}

	client := &BailianClient{
		config:  config,
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}

	return client, nil
}

// Chat 实现聊天接口
func (c *BailianClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 构建千炼API请求格式
	qianfanReq := map[string]interface{}{
		"messages":    req.Messages,
		"temperature": req.Temperature,
		"max_tokens":  req.MaxTokens,
		"stream":      req.Stream,
	}

	if req.Model != "" {
		qianfanReq["model"] = req.Model
	}

	// 发送请求
	respData, err := c.sendRequest(ctx, "/chat/completions", qianfanReq)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var qianfanResp struct {
		ID      string `json:"id"`
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			PromptTokens     int `json:"prompt_tokens"`
			CompletionTokens int `json:"completion_tokens"`
			TotalTokens      int `json:"total_tokens"`
		} `json:"usage"`
		Created int64  `json:"created"`
		Model   string `json:"model"`
	}

	if err := json.Unmarshal(respData, &qianfanResp); err != nil {
		return nil, &Error{
			Code:    "parse_error",
			Message: fmt.Sprintf("解析响应失败: %v", err),
			Type:    ErrorTypeUnknown,
		}
	}

	if len(qianfanResp.Choices) == 0 {
		return nil, &Error{
			Code:    "empty_response",
			Message: "AI返回空响应",
			Type:    ErrorTypeUnknown,
		}
	}

	return &ChatResponse{
		ID:      qianfanResp.ID,
		Content: qianfanResp.Choices[0].Message.Content,
		Model:   qianfanResp.Model,
		Usage: Usage{
			PromptTokens:     qianfanResp.Usage.PromptTokens,
			CompletionTokens: qianfanResp.Usage.CompletionTokens,
			TotalTokens:      qianfanResp.Usage.TotalTokens,
		},
		Created: time.Unix(qianfanResp.Created, 0),
	}, nil
}

// Generate 实现生成接口
func (c *BailianClient) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	// 转换为Chat格式
	chatReq := &ChatRequest{
		Messages: []Message{
			{Role: "user", Content: req.Prompt},
		},
		Model:       req.Model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		Timeout:     req.Timeout,
	}

	chatResp, err := c.Chat(ctx, chatReq)
	if err != nil {
		return nil, err
	}

	genResp := &GenerateResponse{
		ID:      chatResp.ID,
		Content: chatResp.Content,
		Model:   chatResp.Model,
		Usage:   chatResp.Usage,
		Created: chatResp.Created,
	}

	// 如果要求JSON格式，尝试解析
	if req.Format == "json" {
		var data interface{}
		if err := json.Unmarshal([]byte(chatResp.Content), &data); err == nil {
			genResp.Data = data
		}
	}

	return genResp, nil
}

// GetModels 获取可用模型列表
func (c *BailianClient) GetModels(ctx context.Context) ([]Model, error) {
	// 千炼支持的模型列表（示例）
	models := []Model{
		{
			ID:           "ERNIE-Bot-turbo",
			Name:         "文心一言Turbo",
			Description:  "百度文心一言Turbo版本",
			MaxTokens:    4096,
			Capabilities: []string{"chat", "generate", "json"},
		},
		{
			ID:           "ERNIE-Bot",
			Name:         "文心一言",
			Description:  "百度文心一言标准版本",
			MaxTokens:    8192,
			Capabilities: []string{"chat", "generate", "json"},
		},
	}

	return models, nil
}

// Close 关闭客户端
func (c *BailianClient) Close() error {
	// 千炼客户端无需特殊关闭操作
	return nil
}

// sendRequest 发送HTTP请求的通用方法
func (c *BailianClient) sendRequest(ctx context.Context, endpoint string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, &Error{
			Code:    "marshal_error",
			Message: fmt.Sprintf("序列化请求失败: %v", err),
			Type:    ErrorTypeValidation,
		}
	}

	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, &Error{
			Code:    "request_error",
			Message: fmt.Sprintf("创建请求失败: %v", err),
			Type:    ErrorTypeNetwork,
		}
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	// 添加自定义请求头
	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &Error{
			Code:    "network_error",
			Message: fmt.Sprintf("网络请求失败: %v", err),
			Type:    ErrorTypeNetwork,
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &Error{
			Code:    "read_error",
			Message: fmt.Sprintf("读取响应失败: %v", err),
			Type:    ErrorTypeNetwork,
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, &Error{
			Code:    fmt.Sprintf("http_%d", resp.StatusCode),
			Message: fmt.Sprintf("HTTP错误: %d, 响应: %s", resp.StatusCode, string(body)),
			Type:    ErrorTypeNetwork,
		}
	}

	return body, nil
}
