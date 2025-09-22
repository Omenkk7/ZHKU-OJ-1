package client

import (
	"context"
	"fmt"
	"sync"
	"time"
	"zhku-oj-server/pkg/ai/config"
)

// Manager 客户端管理器
type Manager struct {
	factory         *Factory
	configManager   *config.ConfigManager
	clients         map[string]AIClient
	scenarioClients map[string][]AIClient // 场景客户端缓存（支持故障转移）
	mutex           sync.RWMutex
	defaultProvider string
}

// NewManager 创建客户端管理器
func NewManager(factory *Factory, configManager *config.ConfigManager) *Manager {
	aiConfig := configManager.GetConfig()
	
	manager := &Manager{
		factory:         factory,
		configManager:   configManager,
		clients:         make(map[string]AIClient),
		scenarioClients: make(map[string][]AIClient),
		defaultProvider: aiConfig.DefaultProvider,
	}

	// 监听配置变更
	configManager.AddWatcher(manager.onConfigChanged)

	return manager
}

// GetClient 获取指定提供商的客户端
func (m *Manager) GetClient(provider string) (AIClient, error) {
	// 如果未指定提供商，使用默认提供商
	if provider == "" {
		provider = m.defaultProvider
	}

	m.mutex.RLock()
	if client, exists := m.clients[provider]; exists {
		m.mutex.RUnlock()
		return client, nil
	}
	m.mutex.RUnlock()

	// 创建新客户端
	client, err := m.factory.CreateClient(provider)
	if err != nil {
		return nil, fmt.Errorf("创建客户端失败: %w", err)
	}

	// 缓存客户端
	m.mutex.Lock()
	m.clients[provider] = client
	m.mutex.Unlock()

	return client, nil
}

// GetDefaultClient 获取默认客户端
func (m *Manager) GetDefaultClient() (AIClient, error) {
	return m.GetClient(m.defaultProvider)
}

// GetClientByScenario 根据场景获取客户端
func (m *Manager) GetClientByScenario(scenario string) (AIClient, error) {
	// 获取场景配置
	scenarioConfig, err := m.configManager.GetScenarioConfig(scenario)
	if err != nil {
		return nil, fmt.Errorf("获取场景配置失败: %w", err)
	}

	// 使用场景的主要提供商
	return m.GetClient(scenarioConfig.PrimaryProvider)
}

// GetClientsWithFallback 获取支持故障转移的客户端列表
func (m *Manager) GetClientsWithFallback(scenario string) ([]AIClient, error) {
	m.mutex.RLock()
	if clients, exists := m.scenarioClients[scenario]; exists {
		m.mutex.RUnlock()
		return clients, nil
	}
	m.mutex.RUnlock()

	// 创建场景客户端列表
	clients, err := m.factory.CreateClientWithFallback(scenario)
	if err != nil {
		return nil, fmt.Errorf("创建场景客户端失败: %w", err)
	}

	// 缓存场景客户端
	m.mutex.Lock()
	m.scenarioClients[scenario] = clients
	m.mutex.Unlock()

	return clients, nil
}

// ExecuteWithFallback 使用故障转移执行AI请求
func (m *Manager) ExecuteWithFallback(scenario string, executor func(AIClient) error) error {
	clients, err := m.GetClientsWithFallback(scenario)
	if err != nil {
		return err
	}

	var lastErr error
	for i, client := range clients {
		err := executor(client)
		if err == nil {
			return nil // 成功执行
		}
		
		lastErr = err
		// 记录失败日志（这里可以集成日志系统）
		fmt.Printf("提供商 %d 执行失败: %v，尝试下一个提供商\n", i, err)
	}

	return fmt.Errorf("所有提供商都执行失败，最后错误: %w", lastErr)
}

// TestConnection 测试客户端连接
func (m *Manager) TestConnection(provider string) error {
	client, err := m.GetClient(provider)
	if err != nil {
		return fmt.Errorf("获取客户端失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 尝试获取模型列表来测试连接
	_, err = client.GetModels(ctx)
	if err != nil {
		return fmt.Errorf("连接测试失败: %w", err)
	}

	return nil
}

// TestAllConnections 测试所有启用的提供商连接
func (m *Manager) TestAllConnections() map[string]error {
	enabledProviders := m.configManager.GetEnabledProviders()
	results := make(map[string]error)

	for _, provider := range enabledProviders {
		results[provider] = m.TestConnection(provider)
	}

	return results
}

// GetEnabledProviders 获取启用的提供商列表
func (m *Manager) GetEnabledProviders() []string {
	return m.configManager.GetEnabledProviders()
}

// GetProviderStatus 获取提供商状态
func (m *Manager) GetProviderStatus(provider string) (*ProviderStatus, error) {
	// 检查配置是否存在
	_, err := m.configManager.GetProviderConfig(provider)
	if err != nil {
		return &ProviderStatus{
			Provider:  provider,
			Available: false,
			Error:     err.Error(),
		}, nil
	}

	// 测试连接
	err = m.TestConnection(provider)
	status := &ProviderStatus{
		Provider:  provider,
		Available: err == nil,
		LastCheck: time.Now(),
	}

	if err != nil {
		status.Error = err.Error()
	}

	return status, nil
}

// Reload 重新加载配置和客户端
func (m *Manager) Reload() error {
	// 重新加载配置
	if err := m.configManager.Reload(); err != nil {
		return fmt.Errorf("重新加载配置失败: %w", err)
	}

	// 清除客户端缓存
	m.mutex.Lock()
	// 关闭现有客户端
	for _, client := range m.clients {
		client.Close()
	}
	m.clients = make(map[string]AIClient)
	m.scenarioClients = make(map[string][]AIClient)
	
	// 更新默认提供商
	aiConfig := m.configManager.GetConfig()
	m.defaultProvider = aiConfig.DefaultProvider
	m.mutex.Unlock()

	return nil
}

// Close 关闭管理器和所有客户端
func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	var lastErr error
	for provider, client := range m.clients {
		if err := client.Close(); err != nil {
			lastErr = fmt.Errorf("关闭客户端 %s 失败: %w", provider, err)
		}
	}

	m.clients = make(map[string]AIClient)
	m.scenarioClients = make(map[string][]AIClient)

	return lastErr
}

// onConfigChanged 配置变更回调
func (m *Manager) onConfigChanged(newConfig *config.AIConfig) {
	fmt.Println("检测到AI配置变更，重新加载客户端...")
	
	if err := m.Reload(); err != nil {
		fmt.Printf("重新加载客户端失败: %v\n", err)
	} else {
		fmt.Println("客户端重新加载成功")
	}
}

// ProviderStatus 提供商状态
type ProviderStatus struct {
	Provider  string    `json:"provider"`
	Available bool      `json:"available"`
	Error     string    `json:"error,omitempty"`
	LastCheck time.Time `json:"last_check"`
}

// GetConfigManager 获取配置管理器
func (m *Manager) GetConfigManager() *config.ConfigManager {
	return m.configManager
}
