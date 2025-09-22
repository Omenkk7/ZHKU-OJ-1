package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
) 

// ConfigManager AI配置管理器
type ConfigManager struct {
	config     *AIConfig
	configPath string
	mutex      sync.RWMutex
	watchers   []func(*AIConfig)
}

// NewConfigManager 创建配置管理器
func NewConfigManager(configPath string) (*ConfigManager, error) {
	manager := &ConfigManager{
		configPath: configPath,
		watchers:   make([]func(*AIConfig), 0),
	}

	if err := manager.Load(); err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	return manager, nil
}

// NewManagerWithConfig 使用现有配置创建管理器
func NewManagerWithConfig(config *AIConfig) *ConfigManager {
	return &ConfigManager{
		config:   config,
		watchers: make([]func(*AIConfig), 0),
	}
}

// Load 加载配置文件
func (m *ConfigManager) Load() error {
	if m.configPath == "" {
		m.config = DefaultAIConfig()
		return nil
	}

	data, err := ioutil.ReadFile(m.configPath)
	if err != nil {
		// 如果文件不存在，使用默认配置
		if os.IsNotExist(err) {
			m.config = DefaultAIConfig()
			return m.Save() // 保存默认配置到文件
		}
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config AIConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 先设置配置，然后处理环境变量
	m.config = &config

	// 处理环境变量
	if err := m.resolveEnvVars(m.config); err != nil {
		return fmt.Errorf("处理环境变量失败: %w", err)
	}

	// 从环境变量加载API密钥
	if err := m.loadAPIKeysFromEnv(); err != nil {
		return fmt.Errorf("加载API密钥失败: %w", err)
	}

	// 验证配置
	if err := m.validateConfig(m.config); err != nil {
		return fmt.Errorf("配置验证失败: %w", err)
	}

	// 通知观察者
	m.notifyWatchers(m.config)

	return nil
}

// Save 保存配置到文件
func (m *ConfigManager) Save() error {
	if m.configPath == "" {
		return fmt.Errorf("未指定配置文件路径")
	}

	m.mutex.RLock()
	config := m.config
	m.mutex.RUnlock()

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	if err := ioutil.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// GetConfig 获取配置
func (m *ConfigManager) GetConfig() *AIConfig {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	
	// 返回配置的副本，避免外部修改
	return m.copyConfig(m.config)
}

// GetProviderConfig 获取指定提供商配置
func (m *ConfigManager) GetProviderConfig(provider string) (*ClientConfig, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if provider == "" {
		provider = m.config.DefaultProvider
	}

	config, exists := m.config.Providers[provider]
	if !exists {
		return nil, fmt.Errorf("提供商配置不存在: %s", provider)
	}

	if !config.Enabled {
		return nil, fmt.Errorf("提供商已禁用: %s", provider)
	}

	return config, nil
}

// GetScenarioConfig 获取场景配置
func (m *ConfigManager) GetScenarioConfig(scenario string) (*ScenarioConfig, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	config, exists := m.config.Scenarios[scenario]
	if !exists {
		return nil, fmt.Errorf("场景配置不存在: %s", scenario)
	}

	return config, nil
}

// GetEnabledProviders 获取启用的提供商列表
func (m *ConfigManager) GetEnabledProviders() []string {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var providers []string
	for name, config := range m.config.Providers {
		if config.Enabled {
			providers = append(providers, name)
		}
	}
	return providers
}

// AddWatcher 添加配置变更观察者
func (m *ConfigManager) AddWatcher(watcher func(*AIConfig)) {
	m.mutex.Lock()
	m.watchers = append(m.watchers, watcher)
	m.mutex.Unlock()
}

// Reload 重新加载配置
func (m *ConfigManager) Reload() error {
	return m.Load()
}

// resolveEnvVars 处理环境变量
func (m *ConfigManager) resolveEnvVars(config *AIConfig) error {
	for name, providerConfig := range config.Providers {
		// 处理API密钥
		if strings.HasPrefix(providerConfig.APIKey, "${") && strings.HasSuffix(providerConfig.APIKey, "}") {
			envVar := strings.Trim(providerConfig.APIKey, "${}")
			envValue := os.Getenv(envVar)
			if envValue == "" {
				// 如果环境变量未设置，禁用该提供商
				providerConfig.Enabled = false
				continue
			}
			providerConfig.APIKey = envValue
		}

		// 处理其他可能包含环境变量的字段
		if strings.Contains(providerConfig.BaseURL, "${") {
			providerConfig.BaseURL = os.ExpandEnv(providerConfig.BaseURL)
		}
		if strings.Contains(providerConfig.Proxy, "${") {
			providerConfig.Proxy = os.ExpandEnv(providerConfig.Proxy)
		}

		// 设置提供商名称
		providerConfig.Provider = name
	}
	return nil
}

// validateConfig 验证配置
func (m *ConfigManager) validateConfig(config *AIConfig) error {
	if config.DefaultProvider == "" {
		return fmt.Errorf("必须指定默认提供商")
	}

	defaultConfig, exists := config.Providers[config.DefaultProvider]
	if !exists {
		return fmt.Errorf("默认提供商配置不存在: %s", config.DefaultProvider)
	}

	if !defaultConfig.Enabled {
		return fmt.Errorf("默认提供商已禁用: %s", config.DefaultProvider)
	}

	// 验证每个提供商配置
	for name, providerConfig := range config.Providers {
		if providerConfig.Enabled {
			if providerConfig.APIKey == "" {
				return fmt.Errorf("提供商 %s 缺少API密钥", name)
			}
			if providerConfig.BaseURL == "" {
				return fmt.Errorf("提供商 %s 缺少基础URL", name)
			}
			if providerConfig.Timeout <= 0 {
				providerConfig.Timeout = 30 * time.Second
			}
			if providerConfig.MaxRetries < 0 {
				providerConfig.MaxRetries = 3
			}
		}
	}
	

	// 验证场景配置
	for scenarioName, scenarioConfig := range config.Scenarios {
		if _, exists := config.Providers[scenarioConfig.PrimaryProvider]; !exists {
			return fmt.Errorf("场景 %s 的主提供商不存在: %s", scenarioName, scenarioConfig.PrimaryProvider)
		}

		for _, fallbackProvider := range scenarioConfig.FallbackProviders {
			if _, exists := config.Providers[fallbackProvider]; !exists {
				return fmt.Errorf("场景 %s 的备用提供商不存在: %s", scenarioName, fallbackProvider)
			}
		}
	}

	return nil
}

// copyConfig 复制配置
func (m *ConfigManager) copyConfig(config *AIConfig) *AIConfig {
	// 简单的深拷贝实现
	data, _ := yaml.Marshal(config)
	var copy AIConfig
	yaml.Unmarshal(data, &copy)
	return &copy
}

// notifyWatchers 通知观察者
func (m *ConfigManager) notifyWatchers(config *AIConfig) {
	for _, watcher := range m.watchers {
		go watcher(config)
	}
}

// loadAPIKeysFromEnv 从环境变量加载API密钥
func (m *ConfigManager) loadAPIKeysFromEnv() error {
	for providerName, providerConfig := range m.config.Providers {
		var envKey string

		// 根据提供商名称确定环境变量名
		switch providerName {
		case "dashscope":
			envKey = "DASHSCOPE_API_KEY"
		case "openai":
			envKey = "OPENAI_API_KEY"
		case "glm":
			envKey = "GLM_API_KEY"
		case "spark":
			envKey = "SPARK_API_KEY"
		default:
			// 对于未知提供商，使用大写的提供商名称 + _API_KEY
			envKey = strings.ToUpper(providerName) + "_API_KEY"
		}

		// 从环境变量读取API密钥
		apiKey := os.Getenv(envKey)
		if apiKey == "" {
			if providerConfig.Enabled {
				return fmt.Errorf("提供商 %s 已启用但未设置环境变量 %s", providerName, envKey)
			}
			// 如果提供商未启用，跳过API密钥检查
			log.Printf("提供商 %s 未启用，跳过API密钥检查", providerName)
			continue
		}

		// 设置API密钥
		providerConfig.APIKey = apiKey

		log.Printf("已从环境变量 %s 加载提供商 %s 的API密钥", envKey, providerName)
	}

	return nil
}
