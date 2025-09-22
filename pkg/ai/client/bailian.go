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

// DashScopeClient 阿里云百炼DashScope客户端实现
type DashScopeClient struct {
	config     *ClientConfig
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

// NewDashScopeClient 创建百炼DashScope客户端
func NewDashScopeClient(config *ClientConfig) (*DashScopeClient, error) {
	if config.APIKey == "" {
		return nil, &Error{
			Code:    "missing_api_key",
			Message: "API密钥不能为空",
			Type:    ErrorTypeAuth,
		}
	}

	client := &DashScopeClient{
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
func (c *DashScopeClient) Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error) {
	// 设置默认值
	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.7 // 默认温度
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 1000 // 默认最大token数
	}

	topP := req.TopP
	if topP == 0 {
		topP = 0.9 // 默认top_p值
	}

	// 构建DashScope API请求格式
	dashscopeReq := map[string]interface{}{
		"model": req.Model,
		"input": map[string]interface{}{
			"messages": req.Messages,
		},
		"parameters": map[string]interface{}{
			"temperature": temperature,
			"max_tokens":  maxTokens,
			"top_p":       topP,
			"stream":      req.Stream,
		},
	}

	// 发送请求
	respData, err := c.sendRequest(ctx, "/api/v1/services/aigc/text-generation/generation", dashscopeReq)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var dashscopeResp struct {
		RequestID string `json:"request_id"`
		Output    struct {
			Text         string `json:"text"`
			FinishReason string `json:"finish_reason"`
		} `json:"output"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
			TotalTokens  int `json:"total_tokens"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(respData, &dashscopeResp); err != nil {
		return nil, &Error{
			Code:    "parse_error",
			Message: fmt.Sprintf("解析响应失败: %v", err),
			Type:    ErrorTypeUnknown,
		}
	}

	if dashscopeResp.Output.Text == "" {
		return nil, &Error{
			Code:    "empty_response",
			Message: "AI返回空响应",
			Type:    ErrorTypeUnknown,
		}
	}

	return &ChatResponse{
		ID:      dashscopeResp.RequestID,
		Content: dashscopeResp.Output.Text,
		Model:   req.Model,
		Usage: Usage{
			PromptTokens:     dashscopeResp.Usage.InputTokens,
			CompletionTokens: dashscopeResp.Usage.OutputTokens,
			TotalTokens:      dashscopeResp.Usage.TotalTokens,
		},
		Created: time.Now(),
	}, nil
}

// Generate 单次对话接口
func (c *DashScopeClient) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {
	// 设置默认值（与Chat方法保持一致）
	temperature := req.Temperature
	if temperature == 0 {
		temperature = 0.7 // 默认温度
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 1000 // 默认最大token数
	}

	topP := req.TopP
	if topP == 0 {
		topP = 0.9 // 默认top_p值
	}

	// 转换为Chat格式
	chatReq := &ChatRequest{
		Messages: []Message{
			{Role: "user", Content: req.Prompt},
		},
		Model:       req.Model,
		Temperature: temperature,
		MaxTokens:   maxTokens,
		TopP:        topP,
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
func (c *DashScopeClient) GetModels(ctx context.Context) ([]Model, error) {
	// DashScope支持的模型列表
	models := []Model{
		{
			ID:           "qwen-turbo",
			Name:         "通义千问Turbo",
			Description:  "阿里云通义千问Turbo版本",
			MaxTokens:    6000,
			Capabilities: []string{"chat", "generate", "json"},
		},
		{
			ID:           "qwen-plus",
			Name:         "通义千问Plus",
			Description:  "阿里云通义千问Plus版本",
			MaxTokens:    30000,
			Capabilities: []string{"chat", "generate", "json"},
		},
		{
			ID:           "qwen-max",
			Name:         "通义千问Max",
			Description:  "阿里云通义千问Max版本",
			MaxTokens:    6000,
			Capabilities: []string{"chat", "generate", "json"},
		},
	}

	return models, nil
}

// Close 关闭客户端
func (c *DashScopeClient) Close() error {
	// DashScope客户端无需特殊关闭操作
	return nil
}

// sendRequest 发送HTTP请求的通用方法
func (c *DashScopeClient) sendRequest(ctx context.Context, endpoint string, data interface{}) ([]byte, error) {
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
	req.Header.Set("X-DashScope-SSE", "disable") // 禁用SSE流式输出

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
