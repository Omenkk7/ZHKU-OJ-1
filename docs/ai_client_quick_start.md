# AI客户端快速入门指南

## 🚀 快速开始

### 1. 初始化AI服务

```go
package main

import (
    "context"
    "fmt"
    "log"
    "zhku-oj-server/pkg/ai/client"
    "zhku-oj-server/pkg/utils"
)

func main() {
    // 初始化AI服务
    _, err := utils.InitAI("conf")
    if err != nil {
        log.Fatalf("AI服务初始化失败: %v", err)
    }

    // 获取默认客户端
    aiManager := utils.GetAIClientManager()
    aiClient, err := aiManager.GetDefaultClient()
    if err != nil {
        log.Fatalf("获取AI客户端失败: %v", err)
    }

    // 现在可以使用aiClient了
}
```

### 2. 基础Chat调用

```go
// 最简单的调用 - 使用默认参数
func simpleChat(aiClient client.AIClient) {
    chatReq := &client.ChatRequest{
        Messages: []client.Message{
            {Role: "user", Content: "你好，请介绍一下Go语言"},
        },
        Model: "qwen-plus",
        // 自动使用默认值: Temperature=0.7, MaxTokens=1000, TopP=0.9
    }

    response, err := aiClient.Chat(context.Background(), chatReq)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("回复: %s\n", response.Content)
    fmt.Printf("Token使用: %d\n", response.Usage.TotalTokens)
}
```

### 3. 基础Generate调用

```go
// 简单的文本生成
func simpleGenerate(aiClient client.AIClient) {
    genReq := &client.GenerateRequest{
        Prompt: "写一个计算斐波那契数列的函数",
        Model:  "qwen-plus",
        // 自动使用默认值: Temperature=0.7, MaxTokens=1000, TopP=0.9
    }

    response, err := aiClient.Generate(context.Background(), genReq)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("生成内容: %s\n", response.Content)
}
```

## 📋 常用场景

### 1. 代码生成

```go
func generateCode(aiClient client.AIClient, language, description string) (string, error) {
    prompt := fmt.Sprintf("请用%s语言实现：%s\n只返回代码，不要解释。", language, description)

    genReq := &client.GenerateRequest{
        Prompt:      prompt,
        Model:       "qwen-plus",
        Temperature: 0.2, // 低温度确保代码准确性
        MaxTokens:   1000,
    }

    response, err := aiClient.Generate(context.Background(), genReq)
    if err != nil {
        return "", err
    }

    return response.Content, nil
}

// 使用示例
code, err := generateCode(aiClient, "Go", "实现一个栈数据结构")
```

### 2. JSON格式输出

```go
func generateJSON(aiClient client.AIClient) {
    genReq := &client.GenerateRequest{
        Prompt: "生成一个包含用户信息的JSON示例，包含姓名、年龄、邮箱",
        Model:  "qwen-plus",
        Format: "json", // 指定JSON格式
    }

    response, err := aiClient.Generate(context.Background(), genReq)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("JSON: %s\n", response.Content)
    
    // 可以进一步解析JSON
    // var data map[string]interface{}
    // json.Unmarshal([]byte(response.Content), &data)
}
```

### 3. 多轮对话

```go
func multiTurnChat(aiClient client.AIClient) {
    chatReq := &client.ChatRequest{
        Messages: []client.Message{
            {Role: "system", Content: "你是一个编程助手"},
            {Role: "user", Content: "什么是递归？"},
            {Role: "assistant", Content: "递归是函数调用自身的编程技术..."},
            {Role: "user", Content: "能给个例子吗？"},
        },
        Model:       "qwen-plus",
        Temperature: 0.3, // 较低温度，更准确的回答
        MaxTokens:   500,
    }

    response, err := aiClient.Chat(context.Background(), chatReq)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("回复: %s\n", response.Content)
}
```

### 4. 测试用例生成

```go
func generateTestCases(aiClient client.AIClient, functionCode string) {
    prompt := fmt.Sprintf(`
为以下函数生成测试用例，返回JSON格式：
%s

返回格式：
{
  "test_cases": [
    {"input": "输入", "expected": "期望输出", "description": "描述"}
  ]
}`, functionCode)

    genReq := &client.GenerateRequest{
        Prompt:      prompt,
        Model:       "qwen-plus",
        Format:      "json",
        Temperature: 0.3,
    }

    response, err := aiClient.Generate(context.Background(), genReq)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }

    fmt.Printf("测试用例: %s\n", response.Content)
}
```

## ⚙️ 参数配置

### 默认值说明

| 参数 | 默认值 | 说明 |
|------|--------|------|
| **Temperature** | `0.7` | 控制输出随机性，0-1之间 |
| **MaxTokens** | `1000` | 最大生成token数 |
| **TopP** | `0.9` | 核采样参数，0-1之间 |
| **Stream** | `false` | 是否启用流式响应 |

### 参数调优建议

```go
// 快速响应场景
fastConfig := &client.ChatRequest{
    Model:       "qwen-turbo",  // 更快的模型
    Temperature: 0.3,           // 较低温度
    MaxTokens:   200,           // 限制长度
}

// 高质量输出场景
qualityConfig := &client.ChatRequest{
    Model:       "qwen-plus",   // 更强的模型
    Temperature: 0.7,           // 平衡创造性
    MaxTokens:   1000,          // 允许更长输出
}

// 确定性输出场景
deterministicConfig := &client.ChatRequest{
    Model:       "qwen-plus",
    Temperature: 0.1,           // 很低的温度
    MaxTokens:   500,
}
```

## 🔍 错误处理

### 基础错误处理

```go
response, err := aiClient.Chat(context.Background(), chatReq)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "timeout"):
        fmt.Println("请求超时，请稍后重试")
    case strings.Contains(err.Error(), "model"):
        fmt.Println("模型名称无效")
    case strings.Contains(err.Error(), "quota"):
        fmt.Println("配额已用完")
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
    return
}
```

### 参数验证

```go
func validateRequest(req *client.ChatRequest) error {
    if req.Model == "" {
        return fmt.Errorf("模型名称不能为空")
    }
    if len(req.Messages) == 0 {
        return fmt.Errorf("消息列表不能为空")
    }
    if req.Temperature < 0 || req.Temperature > 1 {
        return fmt.Errorf("温度参数必须在0-1之间")
    }
    return nil
}
```

### 重试机制

```go
func chatWithRetry(aiClient client.AIClient, req *client.ChatRequest, maxRetries int) (*client.ChatResponse, error) {
    var response *client.ChatResponse
    var err error

    for i := 0; i < maxRetries; i++ {
        response, err = aiClient.Chat(context.Background(), req)
        if err == nil {
            return response, nil
        }

        fmt.Printf("第%d次尝试失败: %v\n", i+1, err)
        if i < maxRetries-1 {
            time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
        }
    }

    return nil, fmt.Errorf("所有重试都失败了: %v", err)
}
```

## 🎯 最佳实践

### 1. 提示词设计
- **明确具体**: 提供清晰、具体的指令
- **示例引导**: 提供期望输出的示例
- **格式约束**: 明确指定输出格式要求

### 2. 参数选择
- **创意任务**: Temperature 0.7-0.9
- **准确任务**: Temperature 0.1-0.3
- **代码生成**: Temperature 0.2, 较低的随机性

### 3. 性能优化
- **复用客户端**: 避免频繁创建新的客户端实例
- **合理设置MaxTokens**: 根据实际需要设置，避免浪费
- **并发控制**: 合理控制并发请求数量

### 4. 错误处理
- **总是检查错误**: 每次调用都要处理可能的错误
- **实现重试**: 对于重要请求实现重试机制
- **记录日志**: 记录请求和响应用于调试

## 📚 完整示例

```go
package main

import (
    "context"
    "fmt"
    "log"
    "zhku-oj-server/pkg/ai/client"
    "zhku-oj-server/pkg/utils"
)

func main() {
    // 初始化
    _, err := utils.InitAI("conf")
    if err != nil {
        log.Fatalf("初始化失败: %v", err)
    }

    aiManager := utils.GetAIClientManager()
    aiClient, err := aiManager.GetDefaultClient()
    if err != nil {
        log.Fatalf("获取客户端失败: %v", err)
    }

    // 简单对话
    chatReq := &client.ChatRequest{
        Messages: []client.Message{
            {Role: "user", Content: "写一个Hello World程序"},
        },
        Model: "qwen-plus",
    }

    response, err := aiClient.Chat(context.Background(), chatReq)
    if err != nil {
        log.Fatalf("请求失败: %v", err)
    }

    fmt.Printf("AI回复: %s\n", response.Content)
    fmt.Printf("Token使用: %d\n", response.Usage.TotalTokens)
}
```

---

*更多详细信息请参考 [完整API文档](./ai_client_api_reference.md)*
