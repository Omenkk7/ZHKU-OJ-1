/*
@Author: urmsone urmsone@163.com
@Date: 2025-09-22
@Name: ai.go
@Description: AI服务工具类，提供AI客户端管理器的初始化和全局访问
*/

package utils

import (
	"fmt"
	"log"
	"sync"
	"zhku-oj-server/pkg/ai/client"
	"zhku-oj-server/pkg/ai/config"
)

var (
	aiClientManager *client.Manager
	aiOnce          sync.Once
	aiInitError     error
)



// InitAI 初始化AI服务
// 采用单例模式，确保只初始化一次
func InitAI(configDir string) (*client.Manager, error) {
	aiOnce.Do(func() {
		log.Println("正在初始化AI服务...")

		// 1. 检查环境变量
		if err := config.CheckAllRequired(); err != nil {
			aiInitError = fmt.Errorf("环境变量检查失败: %w", err)
			log.Printf("❌ %v", aiInitError)
			return
		}
		log.Println("✅ 环境变量检查通过")

		// 2. 创建AI配置管理器
		aiConfigPath := configDir + "/ai.yaml"
		configManager, err := config.NewConfigManager(aiConfigPath)
		if err != nil {
			aiInitError = fmt.Errorf("创建AI配置管理器失败: %w", err)
			log.Printf("❌ %v", aiInitError)
			return
		}
		log.Println("✅ AI配置管理器创建成功")

		// 3. 创建客户端工厂
		factory := client.NewFactory(configManager)
		log.Println("✅ AI客户端工厂创建成功")

		// 4. 创建客户端管理器
		clientManager := client.NewManager(factory, configManager)
		log.Println("✅ AI客户端管理器创建成功")

		// 5. 测试连接
		log.Println("正在测试AI服务连接...")
		testResults := clientManager.TestAllConnections()
		hasWorkingProvider := false

		for provider, err := range testResults {
			if err != nil {
				log.Printf("❌ AI提供商 %s 连接失败: %v", provider, err)
			} else {
				log.Printf("✅ AI提供商 %s 连接成功", provider)
				hasWorkingProvider = true
			}
		}

		if !hasWorkingProvider {
			aiInitError = fmt.Errorf("没有可用的AI提供商")
			log.Printf("❌ %v", aiInitError)
			return
		}

		// 6. 验证配置
		aiConfig := configManager.GetConfig()
		log.Printf("✅ AI配置加载成功，默认提供商: %s", aiConfig.DefaultProvider)
		log.Printf("✅ 启用的提供商: %v", configManager.GetEnabledProviders())

		// 7. 保存全局实例
		aiClientManager = clientManager
		log.Println("🎉 AI服务初始化完成")
	})

	if aiInitError != nil {
		return nil, aiInitError
	}

	return aiClientManager, nil
}

// GetAIClientManager 获取AI客户端管理器实例
// 如果未初始化，返回nil
func GetAIClientManager() *client.Manager {
	return aiClientManager
}

// IsAIAvailable 检查AI服务是否可用
func IsAIAvailable() bool {
	if aiClientManager == nil {
		return false
	}

	// 测试默认客户端连接
	err := aiClientManager.TestConnection("")
	return err == nil
}

// GetAIStatus 获取AI服务状态
func GetAIStatus() map[string]interface{} {
	if aiClientManager == nil {
		return map[string]interface{}{
			"available":    false,
			"initialized": false,
			"error":       "AI服务未初始化",
		}
	}

	// 测试所有提供商连接
	testResults := aiClientManager.TestAllConnections()
	enabledProviders := aiClientManager.GetEnabledProviders()

	status := map[string]interface{}{
		"available":         len(enabledProviders) > 0,
		"initialized":       true,
		"enabled_providers": enabledProviders,
		"provider_status":   make(map[string]string),
	}

	for provider, err := range testResults {
		if err != nil {
			status["provider_status"].(map[string]string)[provider] = "error: " + err.Error()
		} else {
			status["provider_status"].(map[string]string)[provider] = "ok"
		}
	}

	return status
}

// MustGetAIClientManager 获取AI客户端管理器实例
// 如果未初始化，会panic
func MustGetAIClientManager() *client.Manager {
	if aiClientManager == nil {
		panic("AI服务未初始化，请先调用 utils.InitAI()")
	}
	return aiClientManager
}

// TryInitAI 尝试初始化AI服务
// 如果失败，不会panic，而是返回错误
func TryInitAI(configDir string) error {
	_, err := InitAI(configDir)
	return err
}

// ResetAI 重置AI服务（主要用于测试）
func ResetAI() {
	aiOnce = sync.Once{}
	aiClientManager = nil
	aiInitError = nil
}

// GetAIInitError 获取AI初始化错误
func GetAIInitError() error {
	return aiInitError
}

// AIServiceInfo AI服务信息
type AIServiceInfo struct {
	Available        bool              `json:"available"`
	Initialized      bool              `json:"initialized"`
	EnabledProviders []string          `json:"enabled_providers"`
	ProviderStatus   map[string]string `json:"provider_status"`
	DefaultProvider  string            `json:"default_provider"`
	InitError        string            `json:"init_error,omitempty"`
}

// GetAIServiceInfo 获取AI服务详细信息
func GetAIServiceInfo() *AIServiceInfo {
	info := &AIServiceInfo{
		Available:      false,
		Initialized:    aiClientManager != nil,
		ProviderStatus: make(map[string]string),
	}

	if aiInitError != nil {
		info.InitError = aiInitError.Error()
		return info
	}

	if aiClientManager == nil {
		info.InitError = "AI服务未初始化"
		return info
	}

	// 获取详细状态
	testResults := aiClientManager.TestAllConnections()
	enabledProviders := aiClientManager.GetEnabledProviders()

	info.Available = len(enabledProviders) > 0
	info.EnabledProviders = enabledProviders

	for provider, err := range testResults {
		if err != nil {
			info.ProviderStatus[provider] = "error: " + err.Error()
		} else {
			info.ProviderStatus[provider] = "ok"
		}
	}

	// 获取默认提供商
	if configManager := aiClientManager.GetConfigManager(); configManager != nil {
		if aiConfig := configManager.GetConfig(); aiConfig != nil {
			info.DefaultProvider = aiConfig.DefaultProvider
		}
	}

	return info
}

// LogAIStatus 打印AI服务状态日志
func LogAIStatus() {
	info := GetAIServiceInfo()
	
	if !info.Initialized {
		log.Printf("🔴 AI服务状态: 未初始化 - %s", info.InitError)
		return
	}

	if !info.Available {
		log.Printf("🟡 AI服务状态: 已初始化但不可用")
		for provider, status := range info.ProviderStatus {
			log.Printf("   - %s: %s", provider, status)
		}
		return
	}

	log.Printf("🟢 AI服务状态: 正常运行")
	log.Printf("   - 默认提供商: %s", info.DefaultProvider)
	log.Printf("   - 启用提供商: %v", info.EnabledProviders)
	for provider, status := range info.ProviderStatus {
		if status == "ok" {
			log.Printf("   - %s: ✅ 连接正常", provider)
		} else {
			log.Printf("   - %s: ❌ %s", provider, status)
		}
	}
}
