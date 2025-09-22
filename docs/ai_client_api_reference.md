# AI客户端接口文档

## 📋 概述

本文档详细说明了AI客户端的Chat和Generate方法的请求体和返回体结构，以及使用示例。

## 🔧 接口定义

### 1. Chat方法 - 多轮对话接口

#### 1.1 方法签名
```go
Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
```

#### 1.2 请求体结构 (ChatRequest)

```go
type ChatRequest struct {
    Messages       []Message        `json:"messages"`                // 消息列表 (必填)
    Model          string           `json:"model,omitempty"`         // 模型名称 (必填)
    Temperature    float64          `json:"temperature,omitempty"`   // 温度参数 0-1 (可选，默认0.7)
    MaxTokens      int              `json:"max_tokens,omitempty"`    // 最大token数 (可选，默认1000)
    TopP           float64          `json:"top_p,omitempty"`         // Top-p采样 (可选，默认0.9)
    Stream         bool             `json:"stream,omitempty"`        // 是否流式返回 (可选，默认false)
    Stop           []string         `json:"stop,omitempty"`          // 停止词 (可选)
    ResponseFormat *ResponseFormat  `json:"response_format,omitempty"` // 响应格式 (可选)
    Timeout        time.Duration    `json:"-"`                       // 请求超时时间 (可选)
}

// 消息结构
type Message struct {
    Role    string `json:"role"`    // 角色: "system", "user", "assistant"
    Content string `json:"content"` // 消息内容
}

// 响应格式
type ResponseFormat struct {
    Type string `json:"type"` // "text" 或 "json_object"
}
```

#### 1.3 返回体结构 (ChatResponse)

```go
type ChatResponse struct {
    ID      string    `json:"id"`      // 响应ID
    Content string    `json:"content"` // 响应内容
    Model   string    `json:"model"`   // 使用的模型
    Usage   Usage     `json:"usage"`   // Token使用情况
    Created time.Time `json:"created"` // 创建时间
}

// Token使用情况
type Usage struct {
    PromptTokens     int `json:"prompt_tokens"`     // 提示词token数
    CompletionTokens int `json:"completion_tokens"` // 完成token数
    TotalTokens      int `json:"total_tokens"`      // 总token数
}
```

### 2. Generate方法 - 单次生成接口

#### 2.1 方法签名
```go
Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
```

#### 2.2 请求体结构 (GenerateRequest)

```go
type GenerateRequest struct {
    Prompt      string        `json:"prompt"`                // 提示词 (必填)
    Model       string        `json:"model,omitempty"`       // 模型名称 (必填)
    Temperature float64       `json:"temperature,omitempty"` // 温度参数 (可选，默认0.7)
    MaxTokens   int           `json:"max_tokens,omitempty"`  // 最大token数 (可选，默认1000)
    TopP        float64       `json:"top_p,omitempty"`       // Top-p采样 (可选，默认0.9)
    Format      string        `json:"format,omitempty"`      // 输出格式 (可选，默认"text")
    Schema      interface{}   `json:"schema,omitempty"`      // JSON Schema (当format为json时)
    Timeout     time.Duration `json:"-"`                     // 请求超时时间 (可选)
}
```

#### 2.3 返回体结构 (GenerateResponse)

```go
type GenerateResponse struct {
    ID       string                 `json:"id"`                 // 响应ID
    Content  string                 `json:"content"`            // 生成内容
    Data     interface{}            `json:"data"`               // 结构化数据(当format为json时)
    Model    string                 `json:"model"`              // 使用的模型
    Usage    Usage                  `json:"usage"`              // Token使用情况
    Created  time.Time              `json:"created"`            // 创建时间
    Metadata map[string]interface{} `json:"metadata,omitempty"` // 元数据
}
```

## 🎯 参数详解

### 通用参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| **Model** | string | ✅ | - | 模型名称，如"qwen-plus", "qwen-turbo" |
| **Temperature** | float64 | ❌ | 0.7 | 控制输出随机性，0-1之间，越高越随机 |
| **MaxTokens** | int | ❌ | 1000 | 最大生成token数，控制输出长度 |
| **TopP** | float64 | ❌ | 0.9 | 核采样参数，0-1之间，控制词汇多样性 |
| **Timeout** | time.Duration | ❌ | - | 请求超时时间 |

### Chat专用参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| **Messages** | []Message | ✅ | - | 对话消息列表 |
| **Stream** | bool | ❌ | false | 是否启用流式响应 |
| **Stop** | []string | ❌ | - | 停止词列表 |
| **ResponseFormat** | *ResponseFormat | ❌ | - | 响应格式控制 |

### Generate专用参数

| 参数名 | 类型 | 必填 | 默认值 | 说明 |
|--------|------|------|--------|------|
| **Prompt** | string | ✅ | - | 生成提示词 |
| **Format** | string | ❌ | "text" | 输出格式："text", "json" |
| **Schema** | interface{} | ❌ | - | JSON Schema定义 |

## 🚀 使用示例

### Chat方法示例

#### 基础对话
```go
package main

import (
    "context"
    "fmt"
    "zhku-oj-server/pkg/ai/client"
)

func basicChatExample(aiClient client.AIClient) {
    // 最简单的调用 - 使用默认值
    chatReq := &client.ChatRequest{
        Messages: []client.Message{
            {Role: "user", Content: "你好，请介绍一下Go语言"},
        },
        Model: "qwen-plus",
        // Temperature: 0.7 (默认)
        // MaxTokens: 1000 (默认)
        // TopP: 0.9 (默认)
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

#### 多轮对话
```go
func multiTurnChatExample(aiClient client.AIClient) {
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

#### JSON格式响应
```go
func jsonResponseExample(aiClient client.AIClient) {
    chatReq := &client.ChatRequest{
        Messages: []client.Message{
            {Role: "user", Content: "生成一个用户信息的JSON示例"},
        },
        Model: "qwen-plus",
        ResponseFormat: &client.ResponseFormat{
            Type: "json_object",
        },
    }
    
    response, err := aiClient.Chat(context.Background(), chatReq)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    
    fmt.Printf("JSON响应: %s\n", response.Content)
}
```

### Generate方法示例

#### 基础文本生成
```go
func basicGenerateExample(aiClient client.AIClient) {
    // 使用默认值的简单调用
    genReq := &client.GenerateRequest{
        Prompt: "写一个计算斐波那契数列的函数",
        Model:  "qwen-plus",
        // Temperature: 0.7 (默认)
        // MaxTokens: 1000 (默认)
        // Format: "text" (默认)
    }
    
    response, err := aiClient.Generate(context.Background(), genReq)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    
    fmt.Printf("生成内容: %s\n", response.Content)
    fmt.Printf("Token使用: %d\n", response.Usage.TotalTokens)
}
```

#### 自定义参数生成
```go
func customGenerateExample(aiClient client.AIClient) {
    genReq := &client.GenerateRequest{
        Prompt:      "写一个简单的排序算法",
        Model:       "qwen-plus",
        Temperature: 0.2,  // 低温度，更确定性
        MaxTokens:   800,  // 限制长度
        TopP:        0.8,  // 适中的多样性
        Format:      "text",
    }
    
    response, err := aiClient.Generate(context.Background(), genReq)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    
    fmt.Printf("生成内容: %s\n", response.Content)
}
```

#### JSON格式生成
```go
func jsonGenerateExample(aiClient client.AIClient) {
    genReq := &client.GenerateRequest{
        Prompt: "生成一个包含测试用例的JSON对象，包含输入和期望输出",
        Model:  "qwen-plus",
        Format: "json", // 指定JSON格式
    }
    
    response, err := aiClient.Generate(context.Background(), genReq)
    if err != nil {
        fmt.Printf("错误: %v\n", err)
        return
    }
    
    fmt.Printf("JSON内容: %s\n", response.Content)
    
    // 可以进一步解析JSON
    // var testCase map[string]interface{}
    // json.Unmarshal([]byte(response.Content), &testCase)
}
```

## 🔍 错误处理

### 常见错误类型

```go
// 错误处理示例
response, err := aiClient.Chat(context.Background(), chatReq)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "timeout"):
        fmt.Println("请求超时，请稍后重试")
    case strings.Contains(err.Error(), "invalid model"):
        fmt.Println("模型名称无效")
    case strings.Contains(err.Error(), "quota exceeded"):
        fmt.Println("配额已用完")
    case strings.Contains(err.Error(), "invalid parameter"):
        fmt.Println("参数错误")
    default:
        fmt.Printf("未知错误: %v\n", err)
    }
    return
}
```

### 参数验证

```go
func validateChatRequest(req *client.ChatRequest) error {
    if req.Model == "" {
        return fmt.Errorf("模型名称不能为空")
    }

    if len(req.Messages) == 0 {
        return fmt.Errorf("消息列表不能为空")
    }

    if req.Temperature < 0 || req.Temperature > 1 {
        return fmt.Errorf("温度参数必须在0-1之间")
    }

    if req.TopP < 0 || req.TopP > 1 {
        return fmt.Errorf("TopP参数必须在0-1之间")
    }

    return nil
}
```

## 🎨 高级用法

### 1. 超时控制

```go
func timeoutExample(aiClient client.AIClient) {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    chatReq := &client.ChatRequest{
        Messages: []client.Message{
            {Role: "user", Content: "写一个复杂的算法"},
        },
        Model:   "qwen-plus",
        Timeout: 25 * time.Second, // 请求级别的超时
    }

    response, err := aiClient.Chat(ctx, chatReq)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("请求超时")
        } else {
            fmt.Printf("其他错误: %v\n", err)
        }
        return
    }

    fmt.Printf("回复: %s\n", response.Content)
}
```

### 2. 批量处理

```go
func batchProcessExample(aiClient client.AIClient) {
    prompts := []string{
        "解释什么是数组",
        "解释什么是链表",
        "解释什么是栈",
        "解释什么是队列",
    }

    var wg sync.WaitGroup
    results := make(chan string, len(prompts))

    for _, prompt := range prompts {
        wg.Add(1)
        go func(p string) {
            defer wg.Done()

            genReq := &client.GenerateRequest{
                Prompt: p,
                Model:  "qwen-plus",
                MaxTokens: 200,
            }

            response, err := aiClient.Generate(context.Background(), genReq)
            if err != nil {
                results <- fmt.Sprintf("错误: %v", err)
                return
            }

            results <- response.Content
        }(prompt)
    }

    wg.Wait()
    close(results)

    for result := range results {
        fmt.Println("结果:", result)
        fmt.Println("---")
    }
}
```

### 3. 重试机制

```go
func retryExample(aiClient client.AIClient) {
    maxRetries := 3
    var response *client.ChatResponse
    var err error

    chatReq := &client.ChatRequest{
        Messages: []client.Message{
            {Role: "user", Content: "你好"},
        },
        Model: "qwen-plus",
    }

    for i := 0; i < maxRetries; i++ {
        response, err = aiClient.Chat(context.Background(), chatReq)
        if err == nil {
            break
        }

        fmt.Printf("第%d次尝试失败: %v\n", i+1, err)
        if i < maxRetries-1 {
            time.Sleep(time.Duration(i+1) * time.Second) // 指数退避
        }
    }

    if err != nil {
        fmt.Printf("所有重试都失败了: %v\n", err)
        return
    }

    fmt.Printf("成功: %s\n", response.Content)
}
```

## 📊 性能优化建议

### 1. 参数调优

```go
// 快速响应场景
fastConfig := &client.ChatRequest{
    Model:       "qwen-turbo",    // 使用更快的模型
    Temperature: 0.3,             // 较低温度，减少随机性
    MaxTokens:   200,             // 限制输出长度
    TopP:        0.8,             // 适中的多样性
}

// 高质量输出场景
qualityConfig := &client.ChatRequest{
    Model:       "qwen-plus",     // 使用更强的模型
    Temperature: 0.7,             // 平衡创造性和准确性
    MaxTokens:   1000,            // 允许更长输出
    TopP:        0.9,             // 更高的多样性
}

// 确定性输出场景
deterministicConfig := &client.ChatRequest{
    Model:       "qwen-plus",
    Temperature: 0.1,             // 很低的温度
    MaxTokens:   500,
    TopP:        0.5,             // 较低的多样性
}
```

### 2. 连接池管理

```go
// 在实际应用中，建议使用连接池
func getAIClient() client.AIClient {
    // 使用单例模式或连接池
    return utils.GetAIClientManager().GetDefaultClient()
}
```

## 🔧 实际应用场景

### 1. 代码生成

```go
func generateCode(language, description string) (string, error) {
    aiClient, err := getAIClient()
    if err != nil {
        return "", err
    }

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
```

### 2. 测试用例生成

```go
func generateTestCases(functionCode string) ([]TestCase, error) {
    aiClient, err := getAIClient()
    if err != nil {
        return nil, err
    }

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
        Prompt: prompt,
        Model:  "qwen-plus",
        Format: "json",
        Temperature: 0.3,
    }

    response, err := aiClient.Generate(context.Background(), genReq)
    if err != nil {
        return nil, err
    }

    // 解析JSON响应
    var result struct {
        TestCases []TestCase `json:"test_cases"`
    }

    if err := json.Unmarshal([]byte(response.Content), &result); err != nil {
        return nil, fmt.Errorf("解析JSON失败: %v", err)
    }

    return result.TestCases, nil
}

type TestCase struct {
    Input       string `json:"input"`
    Expected    string `json:"expected"`
    Description string `json:"description"`
}
```

### 3. 智能问答

```go
func intelligentQA(question string, context []string) (string, error) {
    aiClient, err := getAIClient()
    if err != nil {
        return "", err
    }

    messages := []client.Message{
        {Role: "system", Content: "你是一个智能助手，基于提供的上下文回答问题。"},
    }

    // 添加上下文
    for _, ctx := range context {
        messages = append(messages, client.Message{
            Role:    "user",
            Content: fmt.Sprintf("上下文: %s", ctx),
        })
    }

    // 添加问题
    messages = append(messages, client.Message{
        Role:    "user",
        Content: fmt.Sprintf("问题: %s", question),
    })

    chatReq := &client.ChatRequest{
        Messages:    messages,
        Model:       "qwen-plus",
        Temperature: 0.5,
        MaxTokens:   800,
    }

    response, err := aiClient.Chat(context.Background(), chatReq)
    if err != nil {
        return "", err
    }

    return response.Content, nil
}
```

## 📝 最佳实践

### 1. 提示词设计
- **明确具体**: 提供清晰、具体的指令
- **示例引导**: 提供期望输出的示例
- **格式约束**: 明确指定输出格式要求
- **角色设定**: 为AI设定合适的角色

### 2. 参数选择
- **Temperature**: 创意任务用0.7-0.9，准确任务用0.1-0.3
- **MaxTokens**: 根据实际需要设置，避免浪费
- **TopP**: 通常保持0.9，特殊情况下调整

### 3. 错误处理
- **总是检查错误**: 每次调用都要处理可能的错误
- **实现重试**: 对于重要请求实现重试机制
- **记录日志**: 记录请求和响应用于调试

### 4. 性能考虑
- **复用客户端**: 避免频繁创建新的客户端实例
- **并发控制**: 合理控制并发请求数量
- **缓存结果**: 对于相同请求考虑缓存结果

---

## 📚 相关文档

- [AI服务配置文档](./ai_service_config.md)
- [错误码参考](./error_codes.md)
- [性能调优指南](./performance_tuning.md)

---

*最后更新: 2025-09-22*
```
