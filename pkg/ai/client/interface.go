package client

import (
	"context"
	"time"
)

// AIClient 通用AI客户端接口
type AIClient interface {
	// Chat 通用聊天接口
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)

	// Generate 通用生成接口
	Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)

	// GetModels 获取可用模型列表
	GetModels(ctx context.Context) ([]Model, error)

	// Close 关闭客户端
	Close() error
}

// ChatRequest 聊天请求
type ChatRequest struct {
	Messages       []Message        `json:"messages"`                // 消息列表
	Model          string           `json:"model,omitempty"`         // 模型名称
	Temperature    float64          `json:"temperature,omitempty"`   // 温度参数 0-1
	MaxTokens      int              `json:"max_tokens,omitempty"`    // 最大token数
	TopP           float64          `json:"top_p,omitempty"`         // Top-p采样
	Stream         bool             `json:"stream,omitempty"`        // 是否流式返回
	Stop           []string         `json:"stop,omitempty"`          // 停止词
	ResponseFormat *ResponseFormat  `json:"response_format,omitempty"` // 响应格式
	Timeout        time.Duration    `json:"-"`                       // 请求超时时间
}

// ChatResponse 聊天响应
type ChatResponse struct {
	ID      string    `json:"id"`      // 响应ID
	Content string    `json:"content"` // 响应内容
	Model   string    `json:"model"`   // 使用模型
	Usage   Usage     `json:"usage"`   // 使用情况
	Created time.Time `json:"created"` // 创建时间
}

// GenerateRequest 生成请求
type GenerateRequest struct {
	Prompt      string        `json:"prompt"`                // 提示词
	Model       string        `json:"model,omitempty"`       // 模型名称
	Temperature float64       `json:"temperature,omitempty"` // 温度参数
	MaxTokens   int           `json:"max_tokens,omitempty"`  // 最大token数
	TopP        float64       `json:"top_p,omitempty"`       // Top-p采样
	Format      string        `json:"format,omitempty"`      // 输出格式 (json, text, etc.)
	Schema      interface{}   `json:"schema,omitempty"`      // JSON Schema (当format为json时)
	Timeout     time.Duration `json:"-"`                     // 请求超时时间
}

// GenerateResponse 生成响应
type GenerateResponse struct {
	ID       string                 `json:"id"`                 // 响应ID
	Content  string                 `json:"content"`            // 生成内容
	Data     interface{}            `json:"data"`               // 结构化数据(当format为json时)
	Model    string                 `json:"model"`              // 使用的模型
	Usage    Usage                  `json:"usage"`              // 使用情况
	Created  time.Time              `json:"created"`            // 创建时间
	Metadata map[string]interface{} `json:"metadata,omitempty"` // 元数据
}

// Message 消息
type Message struct {
	Role    string `json:"role"`    // 角色: system, user, assistant
	Content string `json:"content"` // 内容
}

// Model 模型信息
type Model struct {
	ID           string   `json:"id"`           // 模型ID
	Name         string   `json:"name"`         // 模型名称
	Description  string   `json:"description"`  // 模型描述
	MaxTokens    int      `json:"max_tokens"`   // 最大token数
	Capabilities []string `json:"capabilities"` // 能力列表
}

// Usage 使用情况
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`     // 提示词token数
	CompletionTokens int `json:"completion_tokens"` // 完成token数
	TotalTokens      int `json:"total_tokens"`      // 总token数
}

// ResponseFormat 响应格式
type ResponseFormat struct {
	Type string `json:"type"` // "text" 或 "json_object"
}

// ClientConfig 客户端配置
type ClientConfig struct {
	Provider   string            `json:"provider"`    // 提供商名称
	APIKey     string            `json:"api_key"`     // API密钥
	BaseURL    string            `json:"base_url"`    // 基础URL
	Model      string            `json:"model"`       // 默认模型
	Timeout    time.Duration     `json:"timeout"`     // 默认超时时间
	MaxRetries int               `json:"max_retries"` // 最大重试次数
	Headers    map[string]string `json:"headers"`     // 自定义请求头
	Proxy      string            `json:"proxy"`       // 代理地址
}

// ClientFactory 客户端工厂接口  负责处理ai客户端创建
type ClientFactory interface {

	// CreateClient 创建AI客户端
	CreateClient(config *ClientConfig) (AIClient, error)

	// GetSupportedProviders 获取支持的提供商列表
	GetSupportedProviders() []string
}

// Error 错误类型
type Error struct {
	Code    string `json:"code"`    // 错误代码
	Message string `json:"message"` // 错误消息
	Type    string `json:"type"`    // 错误类型
}

func (e *Error) Error() string {
	return e.Message
}

// 预定义错误类型
const (
	ErrorTypeAuth       = "auth_error"       // 认证错误
	ErrorTypeQuota      = "quota_error"      // 配额错误
	ErrorTypeRate       = "rate_error"       // 频率限制错误
	ErrorTypeModel      = "model_error"      // 模型错误
	ErrorTypeValidation = "validation_error" // 验证错误
	ErrorTypeNetwork    = "network_error"    // 网络错误
	ErrorTypeTimeout    = "timeout_error"    // 超时错误
	ErrorTypeUnknown    = "unknown_error"    // 未知错误
)

// 预定义错误代码
const (
	ErrorCodeInvalidAPIKey     = "invalid_api_key"
	ErrorCodeQuotaExceeded     = "quota_exceeded"
	ErrorCodeRateLimitExceeded = "rate_limit_exceeded"
	ErrorCodeModelNotFound     = "model_not_found"
	ErrorCodeInvalidRequest    = "invalid_request"
	ErrorCodeTimeout           = "timeout"
	ErrorCodeNetworkError      = "network_error"
)
