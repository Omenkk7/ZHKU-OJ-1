package config

import "time"

// AIConfig AI总配置
type AIConfig struct {
	DefaultProvider string                    `yaml:"default_provider" json:"default_provider"`
	Providers       map[string]*ClientConfig  `yaml:"providers" json:"providers"`
	Scenarios       map[string]*ScenarioConfig `yaml:"scenarios,omitempty" json:"scenarios,omitempty"`
	Global          *GlobalConfig             `yaml:"global,omitempty" json:"global,omitempty"`
}

// ClientConfig 客户端配置
type ClientConfig struct {
	Provider   string            `yaml:"provider" json:"provider"`
	APIKey     string            `yaml:"api_key" json:"api_key"`
	BaseURL    string            `yaml:"base_url" json:"base_url"`
	Model      string            `yaml:"model" json:"model"`
	Timeout    time.Duration     `yaml:"timeout" json:"timeout"`
	MaxRetries int               `yaml:"max_retries" json:"max_retries"`
	Headers    map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Proxy      string            `yaml:"proxy,omitempty" json:"proxy,omitempty"`
	Enabled    bool              `yaml:"enabled" json:"enabled"`
}

// ScenarioConfig 场景配置
type ScenarioConfig struct {
	PrimaryProvider   string   `yaml:"primary_provider" json:"primary_provider"`
	FallbackProviders []string `yaml:"fallback_providers,omitempty" json:"fallback_providers,omitempty"`
	Config            struct {
		Temperature float64 `yaml:"temperature,omitempty" json:"temperature,omitempty"`
		MaxTokens   int     `yaml:"max_tokens,omitempty" json:"max_tokens,omitempty"`
		TopP        float64 `yaml:"top_p,omitempty" json:"top_p,omitempty"`
		Format      string  `yaml:"format,omitempty" json:"format,omitempty"` // text, json
	} `yaml:"config,omitempty" json:"config,omitempty"`
}

// GlobalConfig 全局配置
type GlobalConfig struct {
	EnableCache     bool          `yaml:"enable_cache" json:"enable_cache"`
	CacheTTL        time.Duration `yaml:"cache_ttl" json:"cache_ttl"`
	EnableMetrics   bool          `yaml:"enable_metrics" json:"enable_metrics"`
	EnableRetry     bool          `yaml:"enable_retry" json:"enable_retry"`
	DefaultTimeout  time.Duration `yaml:"default_timeout" json:"default_timeout"`
	MaxConcurrency  int           `yaml:"max_concurrency" json:"max_concurrency"`
	RateLimitRPS    int           `yaml:"rate_limit_rps" json:"rate_limit_rps"`
}

// DefaultAIConfig 返回默认配置
func DefaultAIConfig() *AIConfig {
	return &AIConfig{
		DefaultProvider: "dashscope",
		Providers: map[string]*ClientConfig{
			"dashscope": {
				Provider:   "dashscope",
				// APIKey 将从环境变量 DASHSCOPE_API_KEY 中读取
				BaseURL:    "https://dashscope.aliyuncs.com",
				Model:      "qwen-plus",
				Timeout:    30 * time.Second,
				MaxRetries: 3,
				Enabled:    true,
			},
			// 其他提供商暂时注释
			// "openai": {
			// 	Provider:   "openai",
			// 	// APIKey 将从环境变量 OPENAI_API_KEY 中读取
			// 	BaseURL:    "https://api.openai.com",
			// 	Model:      "gpt-3.5-turbo",
			// 	Timeout:    30 * time.Second,
			// 	MaxRetries: 3,
			// 	Enabled:    false,
			// },
		},
		Scenarios: map[string]*ScenarioConfig{
			"testcase_generation": {
				PrimaryProvider:   "dashscope",
				FallbackProviders: []string{}, // 暂时不使用备用提供商
				Config: struct {
					Temperature float64 `yaml:"temperature,omitempty" json:"temperature,omitempty"`
					MaxTokens   int     `yaml:"max_tokens,omitempty" json:"max_tokens,omitempty"`
					TopP        float64 `yaml:"top_p,omitempty" json:"top_p,omitempty"`
					Format      string  `yaml:"format,omitempty" json:"format,omitempty"`
				}{
					Temperature: 0.7,
					MaxTokens:   4000,
					TopP:        0.9,
					Format:      "json",
				},
			},
		},
		Global: &GlobalConfig{
			EnableCache:     true,
			CacheTTL:        10 * time.Minute,
			EnableMetrics:   true,
			EnableRetry:     true,
			DefaultTimeout:  30 * time.Second,
			MaxConcurrency:  10,
			RateLimitRPS:    100,
		},
	}
}
