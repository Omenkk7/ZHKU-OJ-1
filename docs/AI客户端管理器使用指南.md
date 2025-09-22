# AI客户端管理器使用指南

## 概述

AI客户端管理器提供了完整的配置驱动的AI客户端管理解决方案，支持多提供商、场景化配置、故障转移等功能。

## 架构图

```
配置文件 (ai.yaml)
    ↓
配置管理器 (config.Manager)
    ↓
客户端工厂 (client.Factory)
    ↓
客户端管理器 (client.Manager)
    ↓
业务服务 (TestCaseService)
```

## 快速开始

### 1. 环境准备

```bash
# 设置环境变量
export DASHSCOPE_API_KEY="your-dashscope-api-key"

# 验证配置文件
cat conf/ai.yaml
```

### 2. 基本使用

```go
package main

import (
    "context"
    "log"
    "zhku-oj-server/pkg/ai/client"
    "zhku-oj-server/pkg/ai/config"
)

func main() {
    // 1. 创建配置管理器
    configManager, err := config.NewConfigManager("conf/ai.yaml")
    if err != nil {
        log.Fatal(err)
    }

    // 2. 创建客户端工厂
    factory := client.NewFactory(configManager)

    // 3. 创建客户端管理器
    clientManager := client.NewManager(factory, configManager)
    defer clientManager.Close()

    // 4. 获取客户端并使用
    aiClient, err := clientManager.GetDefaultClient()
    if err != nil {
        log.Fatal(err)
    }

    // 5. 调用AI服务
    response, err := aiClient.Generate(context.Background(), &client.GenerateRequest{
        Prompt: "生成一个Hello World程序的测试用例",
        Format: "json",
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("AI响应: %s", response.Content)
}
```

## 核心功能

### 1. 配置驱动的客户端创建

```go
// 通过配置文件自动创建客户端
client, err := clientManager.GetClient("dashscope")

// 使用默认提供商
defaultClient, err := clientManager.GetDefaultClient()
```

### 2. 场景化配置

```go
// 根据业务场景获取优化的客户端
testCaseClient, err := clientManager.GetClientByScenario("testcase_generation")
codeReviewClient, err := clientManager.GetClientByScenario("code_review")
```

### 3. 故障转移支持

```go
// 自动故障转移执行
err := clientManager.ExecuteWithFallback("testcase_generation", func(aiClient client.AIClient) error {
    response, err := aiClient.Generate(ctx, request)
    if err != nil {
        return err // 自动尝试下一个提供商
    }
    // 处理响应
    return nil
})
```

### 4. 连接测试

```go
// 测试单个提供商
err := clientManager.TestConnection("dashscope")

// 测试所有启用的提供商
results := clientManager.TestAllConnections()
for provider, err := range results {
    if err != nil {
        log.Printf("❌ %s: %v", provider, err)
    } else {
        log.Printf("✅ %s: 连接正常", provider)
    }
}
```

### 5. 配置热更新

```go
// 配置文件变更时自动重新加载客户端
// 管理器会自动监听配置变更并重新创建客户端

// 手动重新加载
err := clientManager.Reload()
```

## 业务集成示例

### 测试用例生成服务

```go
type TestCaseService struct {
    clientManager *client.Manager
}

func NewTestCaseService(clientManager *client.Manager) *TestCaseService {
    return &TestCaseService{clientManager: clientManager}
}

func (s *TestCaseService) GenerateTestCases(problemDesc string) (*TestCaseSet, error) {
    // 使用场景化客户端
    aiClient, err := s.clientManager.GetClientByScenario("testcase_generation")
    if err != nil {
        return nil, err
    }

    // 构建提示词
    prompt := buildTestCasePrompt(problemDesc)
    
    // 调用AI生成
    response, err := aiClient.Generate(context.Background(), &client.GenerateRequest{
        Prompt:      prompt,
        Temperature: 0.7,
        MaxTokens:   4000,
        Format:      "json",
    })
    
    if err != nil {
        return nil, err
    }

    // 解析响应
    return parseTestCaseResponse(response.Content)
}
```

### 代码审查服务

```go
type CodeReviewService struct {
    clientManager *client.Manager
}

func (s *CodeReviewService) ReviewCode(code string) (*ReviewResult, error) {
    // 使用故障转移确保服务可用性
    var result *ReviewResult
    
    err := s.clientManager.ExecuteWithFallback("code_review", func(aiClient client.AIClient) error {
        response, err := aiClient.Generate(context.Background(), &client.GenerateRequest{
            Prompt:      buildCodeReviewPrompt(code),
            Temperature: 0.3, // 代码审查需要更准确的结果
            MaxTokens:   2000,
            Format:      "json",
        })
        
        if err != nil {
            return err
        }
        
        result, err = parseReviewResponse(response.Content)
        return err
    })
    
    return result, err
}
```

## 配置管理

### 提供商配置

```yaml
providers:
  dashscope:
    provider: "dashscope"
    api_key: "${DASHSCOPE_API_KEY}"
    base_url: "https://dashscope.aliyuncs.com"
    model: "qwen-turbo"
    timeout: 30s
    max_retries: 3
    enabled: true
```

### 场景配置

```yaml
scenarios:
  testcase_generation:
    primary_provider: "dashscope"
    fallback_providers: []
    config:
      temperature: 0.7
      max_tokens: 4000
      format: "json"
```

## 错误处理

### 错误类型

```go
// 检查错误类型
if err != nil {
    if aiErr, ok := err.(*client.Error); ok {
        switch aiErr.Type {
        case client.ErrorTypeAuth:
            // 处理认证错误
        case client.ErrorTypeQuota:
            // 处理配额错误
        case client.ErrorTypeNetwork:
            // 处理网络错误
        }
    }
}
```

### 重试机制

```go
// 客户端会根据配置自动重试
// max_retries: 3 表示最多重试3次

// 自定义重试逻辑
func retryableGenerate(aiClient client.AIClient, req *client.GenerateRequest) (*client.GenerateResponse, error) {
    maxRetries := 3
    for i := 0; i < maxRetries; i++ {
        response, err := aiClient.Generate(context.Background(), req)
        if err == nil {
            return response, nil
        }
        
        // 检查是否为可重试错误
        if aiErr, ok := err.(*client.Error); ok {
            if aiErr.Type == client.ErrorTypeNetwork || aiErr.Type == client.ErrorTypeTimeout {
                time.Sleep(time.Duration(i+1) * time.Second)
                continue
            }
        }
        
        return nil, err
    }
    return nil, fmt.Errorf("重试次数已用完")
}
```

## 监控和日志

### 状态监控

```go
// 获取提供商状态
status, err := clientManager.GetProviderStatus("dashscope")
if err != nil {
    log.Printf("获取状态失败: %v", err)
} else {
    log.Printf("提供商状态: %+v", status)
}

// 获取启用的提供商
enabledProviders := clientManager.GetEnabledProviders()
log.Printf("启用的提供商: %v", enabledProviders)
```

### 使用统计

```go
// 在实际使用中记录统计信息
type AIUsageStats struct {
    Provider     string
    Scenario     string
    TokensUsed   int
    RequestTime  time.Duration
    Success      bool
}

func trackUsage(provider, scenario string, usage client.Usage, duration time.Duration, success bool) {
    stats := AIUsageStats{
        Provider:    provider,
        Scenario:    scenario,
        TokensUsed:  usage.TotalTokens,
        RequestTime: duration,
        Success:     success,
    }
    // 发送到监控系统
    sendToMonitoring(stats)
}
```

## 最佳实践

### 1. 资源管理

```go
// 在应用启动时创建管理器
func initAI() *client.Manager {
    configManager, _ := config.NewConfigManager("conf/ai.yaml")
    factory := client.NewFactory(configManager)
    return client.NewManager(factory, configManager)
}

// 在应用关闭时清理资源
func cleanup(manager *client.Manager) {
    if err := manager.Close(); err != nil {
        log.Printf("关闭AI客户端管理器失败: %v", err)
    }
}
```

### 2. 超时控制

```go
// 为每个请求设置合适的超时时间
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

response, err := aiClient.Generate(ctx, request)
```

### 3. 错误处理

```go
// 优雅处理AI服务不可用的情况
func generateWithFallback(manager *client.Manager, prompt string) (string, error) {
    // 尝试AI生成
    err := manager.ExecuteWithFallback("testcase_generation", func(aiClient client.AIClient) error {
        response, err := aiClient.Generate(context.Background(), &client.GenerateRequest{
            Prompt: prompt,
        })
        if err == nil {
            // 处理成功响应
            return nil
        }
        return err
    })
    
    if err != nil {
        // AI服务不可用时的降级处理
        return generateFallbackContent(prompt), nil
    }
    
    return "", nil
}
```

## 运行示例

```bash
# 运行示例程序
cd pkg/ai/example
go run main.go

# 预期输出
=== 初始化AI配置管理器 ===
默认提供商: dashscope
启用的提供商: [dashscope]

=== 创建客户端工厂 ===
支持的提供商: [dashscope openai glm spark]

=== 测试提供商连接 ===
✅ dashscope: 连接正常

=== 测试生成功能 ===
生成结果:
  ID: req_xxx
  模型: qwen-turbo
  内容: {...}
  Token使用: 50/200/250 (提示/完成/总计)
```

这个管理器系统提供了完整的配置驱动的AI客户端管理能力，支持多种使用场景和高可用性需求。
