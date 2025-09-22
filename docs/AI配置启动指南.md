# AI配置启动指南

## 概述

本文档详细说明如何在主函数中启动和配置AI系统，包括配置加载、客户端初始化和服务集成。

## 1. 环境准备

### 1.1 设置环境变量

在启动应用之前，需要设置相应的API密钥：

```bash
# Linux/macOS
export DASHSCOPE_API_KEY="your-dashscope-api-key"
# export OPENAI_API_KEY="your-openai-api-key"  # 暂时不需要

# Windows
set DASHSCOPE_API_KEY=your-dashscope-api-key
# set OPENAI_API_KEY=your-openai-api-key  # 暂时不需要

# 或者创建 .env 文件
echo "DASHSCOPE_API_KEY=your-dashscope-api-key" > .env
# echo "OPENAI_API_KEY=your-openai-api-key" >> .env  # 暂时不需要
```

### 环境变量验证

系统提供了环境变量检查工具：

```go
import "zhku-oj-server/pkg/ai/config"

// 检查DashScope环境变量
if err := config.CheckDashScope(); err != nil {
    log.Fatal(err)
}

// 检查所有必需的环境变量
if err := config.CheckAllRequired(); err != nil {
    log.Fatal(err)
}

// 打印环境变量状态
config.PrintEnvStatus()
```

### 1.2 配置文件准备

确保配置文件存在：
```bash
# 检查配置文件
ls -la conf/ai.yaml

# 如果不存在，复制示例配置
cp conf/ai.yaml.example conf/ai.yaml
```

## 2. 主函数集成

### 2.1 基本集成方式

在 `cmd/main.go` 中集成AI配置：

```go
package main

import (
    "log"
    "zhku-oj-server/pkg/ai/config"
    "zhku-oj-server/pkg/ai/client"
    // ... 其他导入
)

func main() {
    // 1. 加载现有配置
    appConfig, err := utils.LoadConfig("conf/config.yaml")
    if err != nil {
        log.Fatalf("加载应用配置失败: %v", err)
    }

    // 2. 初始化AI配置管理器
    aiConfigManager, err := initAIConfig()
    if err != nil {
        log.Fatalf("初始化AI配置失败: %v", err)
    }

    // 3. 初始化AI客户端管理器
    aiClientManager, err := initAIClientManager(aiConfigManager)
    if err != nil {
        log.Fatalf("初始化AI客户端失败: %v", err)
    }

    // 4. 初始化业务服务（传入AI客户端管理器）
    services := initServices(appConfig, aiClientManager)

    // 5. 启动服务器
    startServer(appConfig, services)
}

// 初始化AI配置
func initAIConfig() (*config.ConfigManager, error) {
    log.Println("正在初始化AI配置...")

    // 检查必需的环境变量
    if err := config.CheckAllRequired(); err != nil {
        return nil, fmt.Errorf("环境变量检查失败: %w", err)
    }

    // 创建AI配置管理器
    aiConfigPath := getAIConfigPath()
    manager, err := config.NewConfigManager(aiConfigPath)
    if err != nil {
        return nil, err
    }

    // 验证配置
    aiConfig := manager.GetConfig()
    log.Printf("AI配置加载成功，默认提供商: %s", aiConfig.DefaultProvider)
    log.Printf("启用的提供商: %v", manager.GetEnabledProviders())

    return manager, nil
}

// 初始化AI客户端管理器
func initAIClientManager(configManager *config.ConfigManager) (*client.Manager, error) {
    log.Println("正在初始化AI客户端管理器...")
    
    // 创建客户端工厂
    factory := client.NewFactory(configManager)
    
    // 创建客户端管理器
    clientManager := client.NewManager(factory, configManager)
    
    // 测试默认客户端连接
    if err := testAIConnection(clientManager); err != nil {
        log.Printf("AI连接测试失败: %v", err)
        // 注意：这里可以选择是否要因为AI连接失败而终止启动
        // return nil, err
    } else {
        log.Println("AI连接测试成功")
    }

    return clientManager, nil
}

// 获取AI配置文件路径
func getAIConfigPath() string {
    // 优先使用环境变量指定的路径
    if path := os.Getenv("AI_CONFIG_PATH"); path != "" {
        return path
    }
    
    // 根据环境选择配置文件
    env := os.Getenv("GO_ENV")
    switch env {
    case "production":
        return "conf/ai.prod.yaml"
    case "development":
        return "conf/ai.local.yaml"
    default:
        return "conf/ai.yaml"
    }
}

// 测试AI连接
func testAIConnection(clientManager *client.Manager) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // 获取默认客户端
    client, err := clientManager.GetDefaultClient()
    if err != nil {
        return fmt.Errorf("获取默认AI客户端失败: %w", err)
    }

    // 简单的连接测试
    _, err = client.GetModels(ctx)
    if err != nil {
        return fmt.Errorf("AI连接测试失败: %w", err)
    }

    return nil
}

// 初始化业务服务
func initServices(appConfig *Config, aiClientManager *client.Manager) *Services {
    return &Services{
        // 传入AI客户端管理器到需要AI功能的服务
        TestCaseService: NewTestCaseService(aiClientManager),
        // ... 其他服务
    }
}
```

### 2.2 服务结构定义

```go
// Services 服务集合
type Services struct {
    TestCaseService *TestCaseService
    // ... 其他服务
}

// TestCaseService 测试用例服务
type TestCaseService struct {
    aiClientManager *client.Manager
    // ... 其他依赖
}

func NewTestCaseService(aiClientManager *client.Manager) *TestCaseService {
    return &TestCaseService{
        aiClientManager: aiClientManager,
    }
}
```

## 3. 高级配置

### 3.1 配置热更新

```go
func setupConfigWatcher(aiConfigManager *config.Manager, aiClientManager *client.Manager) {
    // 添加配置变更监听器
    aiConfigManager.AddWatcher(func(newConfig *config.AIConfig) {
        log.Println("检测到AI配置变更，正在重新初始化客户端...")
        
        // 重新初始化客户端管理器
        if err := aiClientManager.Reload(); err != nil {
            log.Printf("重新加载AI客户端失败: %v", err)
        } else {
            log.Println("AI客户端重新加载成功")
        }
    })

    // 启动配置文件监控（可选）
    go func() {
        ticker := time.NewTicker(30 * time.Second)
        defer ticker.Stop()
        
        for range ticker.C {
            if err := aiConfigManager.Reload(); err != nil {
                log.Printf("重新加载AI配置失败: %v", err)
            }
        }
    }()
}
```

### 3.2 优雅关闭

```go
func setupGracefulShutdown(aiClientManager *client.Manager) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-c
        log.Println("正在优雅关闭...")
        
        // 关闭AI客户端
        if err := aiClientManager.Close(); err != nil {
            log.Printf("关闭AI客户端失败: %v", err)
        }
        
        os.Exit(0)
    }()
}
```

## 4. 错误处理策略

### 4.1 启动时错误处理

```go
func initAIWithFallback() (*client.Manager, error) {
    // 尝试加载AI配置
    aiConfigManager, err := config.NewManager("conf/ai.yaml")
    if err != nil {
        log.Printf("加载AI配置失败: %v，使用默认配置", err)
        
        // 使用默认配置
        defaultConfig := config.DefaultAIConfig()
        aiConfigManager = config.NewManagerWithConfig(defaultConfig)
    }

    // 初始化客户端管理器
    factory := client.NewFactory(aiConfigManager)
    clientManager := client.NewManager(factory, aiConfigManager)

    // 测试连接，如果失败则禁用AI功能
    if err := testAIConnection(clientManager); err != nil {
        log.Printf("AI服务不可用: %v，AI功能将被禁用", err)
        return client.NewDisabledManager(), nil // 返回禁用的管理器
    }

    return clientManager, nil
}
```

### 4.2 运行时错误处理

```go
func (s *TestCaseService) GenerateWithFallback(req *GenerationRequest) (*TestCaseSet, error) {
    // 尝试使用AI生成
    result, err := s.generateWithAI(req)
    if err != nil {
        log.Printf("AI生成失败: %v，使用备用方案", err)
        
        // 备用方案：使用规则生成或返回错误
        return s.generateWithRules(req)
    }
    
    return result, nil
}
```

## 5. 监控和日志

### 5.1 添加监控指标

```go
import "github.com/prometheus/client_golang/prometheus"

var (
    aiRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "ai_requests_total",
            Help: "Total number of AI requests",
        },
        []string{"provider", "scenario", "status"},
    )
    
    aiRequestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "ai_request_duration_seconds",
            Help: "AI request duration in seconds",
        },
        []string{"provider", "scenario"},
    )
)

func init() {
    prometheus.MustRegister(aiRequestsTotal)
    prometheus.MustRegister(aiRequestDuration)
}
```

### 5.2 结构化日志

```go
import "github.com/sirupsen/logrus"

func logAIOperation(provider, scenario string, duration time.Duration, err error) {
    fields := logrus.Fields{
        "component": "ai",
        "provider":  provider,
        "scenario":  scenario,
        "duration":  duration.Seconds(),
    }
    
    if err != nil {
        fields["error"] = err.Error()
        logrus.WithFields(fields).Error("AI操作失败")
    } else {
        logrus.WithFields(fields).Info("AI操作成功")
    }
}
```

## 6. 部署注意事项

### 6.1 Docker部署

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# 复制配置文件
COPY --from=builder /app/conf ./conf
COPY --from=builder /app/main .

# 设置环境变量
ENV GO_ENV=production
ENV AI_CONFIG_PATH=conf/ai.prod.yaml

CMD ["./main"]
```

### 6.2 Kubernetes部署

```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: zhku-oj-server
spec:
  template:
    spec:
      containers:
      - name: app
        image: zhku-oj-server:latest
        env:
        - name: DASHSCOPE_API_KEY
          valueFrom:
            secretKeyRef:
              name: ai-secrets
              key: dashscope-api-key
        # - name: OPENAI_API_KEY  # 暂时不需要
        #   valueFrom:
        #     secretKeyRef:
        #       name: ai-secrets
        #       key: openai-api-key
        - name: GO_ENV
          value: "production"
        volumeMounts:
        - name: config
          mountPath: /root/conf
      volumes:
      - name: config
        configMap:
          name: ai-config
```

## 7. 完整示例

### 7.1 完整的main.go示例

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "time"
    
    "zhku-oj-server/pkg/ai/config"
    "zhku-oj-server/pkg/ai/client"
    "zhku-oj-server/pkg/utils"
)

func main() {
    log.Println("正在启动ZHKU OJ服务器...")

    // 1. 加载应用配置
    appConfig, err := utils.LoadConfig("conf/config.yaml")
    if err != nil {
        log.Fatalf("加载应用配置失败: %v", err)
    }

    // 2. 初始化AI系统
    aiClientManager, err := initAISystem()
    if err != nil {
        log.Fatalf("初始化AI系统失败: %v", err)
    }
    defer aiClientManager.Close()

    // 3. 初始化业务服务
    services := initServices(appConfig, aiClientManager)

    // 4. 设置优雅关闭
    setupGracefulShutdown(aiClientManager)

    // 5. 启动服务器
    log.Println("ZHKU OJ服务器启动成功")
    startServer(appConfig, services)
}

func initAISystem() (*client.Manager, error) {
    log.Println("正在初始化AI系统...")

    // 加载AI配置
    aiConfigManager, err := config.NewManager("conf/ai.yaml")
    if err != nil {
        return nil, err
    }

    // 创建客户端管理器
    factory := client.NewFactory(aiConfigManager)
    clientManager := client.NewManager(factory, aiConfigManager)

    // 测试连接
    if err := testAIConnection(clientManager); err != nil {
        log.Printf("AI连接测试失败: %v", err)
    } else {
        log.Println("AI系统初始化成功")
    }

    return clientManager, nil
}

func testAIConnection(clientManager *client.Manager) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    client, err := clientManager.GetDefaultClient()
    if err != nil {
        return err
    }

    _, err = client.GetModels(ctx)
    return err
}

func setupGracefulShutdown(aiClientManager *client.Manager) {
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    go func() {
        <-c
        log.Println("正在优雅关闭...")
        aiClientManager.Close()
        os.Exit(0)
    }()
}
```

这个启动指南提供了完整的AI配置集成方案，包括错误处理、监控、部署等各个方面的考虑。
