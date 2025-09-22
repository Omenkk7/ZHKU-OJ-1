package client

import (
	"fmt"
	"zhku-oj-server/pkg/ai/config"
)

// Factory 客户端工厂
type Factory struct {
	configManager *config.ConfigManager
}

// NewFactory 创建客户端工厂
func NewFactory(configManager *config.ConfigManager) *Factory {
	return &Factory{
		configManager: configManager,
	}
}

// CreateClient 根据提供商名称创建客户端
func (f *Factory) CreateClient(provider string) (AIClient, error) {
	// 从配置管理器获取提供商配置
	configClientConfig, err := f.configManager.GetProviderConfig(provider)
	if err != nil {
		return nil, fmt.Errorf("获取提供商配置失败: %w", err)
	}

	// 转换为客户端需要的配置格式
	clientConfig := &ClientConfig{
		Provider:   configClientConfig.Provider,
		APIKey:     configClientConfig.APIKey,
		BaseURL:    configClientConfig.BaseURL,
		Model:      configClientConfig.Model,
		Timeout:    configClientConfig.Timeout,
		MaxRetries: configClientConfig.MaxRetries,
		Headers:    configClientConfig.Headers,
		Proxy:      configClientConfig.Proxy,
	}

	// 根据提供商类型创建对应的客户端
	switch provider {
	case "dashscope":
		return NewDashScopeClient(clientConfig)
	// case "openai":
	//	return NewOpenAIClient(clientConfig)
	// case "glm":
	//	return NewGLMClient(clientConfig)
	// case "spark":
	//	return NewSparkClient(clientConfig)
	default:
		return nil, fmt.Errorf("不支持的AI提供商: %s，目前只支持 dashscope", provider)
	}
}

// CreateClientByScenario 根据场景创建客户端（使用主要提供商）
func (f *Factory) CreateClientByScenario(scenario string) (AIClient, error) {
	// 获取场景配置
	scenarioConfig, err := f.configManager.GetScenarioConfig(scenario)
	if err != nil {
		return nil, fmt.Errorf("获取场景配置失败: %w", err)
	}

	// 使用场景的主要提供商创建客户端
	return f.CreateClient(scenarioConfig.PrimaryProvider)
}

// CreateClientWithFallback 创建客户端，支持故障转移
func (f *Factory) CreateClientWithFallback(scenario string) ([]AIClient, error) {
	// 获取场景配置
	scenarioConfig, err := f.configManager.GetScenarioConfig(scenario)
	if err != nil {
		return nil, fmt.Errorf("获取场景配置失败: %w", err)
	}

	var clients []AIClient

	// 创建主要提供商客户端
	primaryClient, err := f.CreateClient(scenarioConfig.PrimaryProvider)
	if err == nil {
		clients = append(clients, primaryClient)
	}

	// 创建备用提供商客户端
	for _, fallbackProvider := range scenarioConfig.FallbackProviders {
		fallbackClient, err := f.CreateClient(fallbackProvider)
		if err == nil {
			clients = append(clients, fallbackClient)
		}
	}

	if len(clients) == 0 {
		return nil, fmt.Errorf("场景 %s 没有可用的AI提供商", scenario)
	}

	return clients, nil
}

// GetSupportedProviders 获取支持的提供商列表
func (f *Factory) GetSupportedProviders() []string {
	return []string{"dashscope", "openai", "glm", "spark"}
}

// ValidateProvider 验证提供商是否支持
func (f *Factory) ValidateProvider(provider string) error {
	supportedProviders := f.GetSupportedProviders()
	for _, supported := range supportedProviders {
		if provider == supported {
			return nil
		}
	}
	return fmt.Errorf("不支持的提供商: %s，支持的提供商: %v", provider, supportedProviders)
}

// 以下是其他提供商的客户端创建函数（暂时返回错误，待实现）

// NewOpenAIClient 创建OpenAI客户端（暂未实现）
func NewOpenAIClient(config *config.ClientConfig) (AIClient, error) {
	return nil, fmt.Errorf("OpenAI客户端暂未实现")
}

// NewGLMClient 创建GLM客户端（暂未实现）
func NewGLMClient(config *config.ClientConfig) (AIClient, error) {
	return nil, fmt.Errorf("GLM客户端暂未实现")
}

// NewSparkClient 创建Spark客户端（暂未实现）
func NewSparkClient(config *config.ClientConfig) (AIClient, error) {
	return nil, fmt.Errorf("Spark客户端暂未实现")
}
