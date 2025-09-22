# AI工具包 - 通用AI能力抽象层

## 包结构

```
pkg/ai/
├── client/                   # AI客户端抽象层
│   ├── interface.go          # AI客户端通用接口
│   ├── qianfan.go           # 阿里云千炼实现
│   ├── openai.go            # OpenAI实现
│   └── types.go             # 通用AI类型定义
├── prompt/                   # 提示词工具
│   ├── template.go          # 提示词模板引擎
│   ├── builder.go           # 提示词构建器
│   └── types.go             # 提示词相关类型
└── processor/                # AI响应处理器
    ├── parser.go            # 通用解析器
    ├── validator.go         # 通用验证器
    └── types.go             # 处理器类型定义
```

## 核心流程

### 1. 完整的AI调用流程

```
业务请求 → 提示词构建 → AI客户端调用 → 响应解析 → 结果验证 → 返回业务数据
```

### 2. 详细流程说明

#### 步骤1: 业务请求
```go
// 业务层发起请求
businessData := &TestCaseRequest{
    ProblemTitle: "两数之和",
    Description: "给定一个整数数组...",
    Count: 10,
}
```

#### 步骤2: 提示词构建
```go
// prompt.Builder 接口
type Builder interface {
    Build(templateName string, data interface{}) (string, error)
}

// 使用流程
promptBuilder := prompt.NewBuilder()
promptText, err := promptBuilder.Build("testcase_generation", businessData)
```

#### 步骤3: AI客户端调用
```go
// client.AIClient 接口
type AIClient interface {
    Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
}

// 使用流程
aiClient := client.NewQianfanClient(config)
response, err := aiClient.Generate(ctx, &GenerateRequest{
    Prompt: promptText,
    Format: "json",
})
```

#### 步骤4: 响应解析
```go
// processor.Parser 接口
type Parser interface {
    Parse(content string, target interface{}) error
}

// 使用流程
parser := processor.NewJSONParser()
var testCases []TestCase
err := parser.Parse(response.Content, &testCases)
```

#### 步骤5: 结果验证
```go
// processor.Validator 接口
type Validator interface {
    Validate(data interface{}) (*ValidationResult, error)
}

// 使用流程
validator := processor.NewTestCaseValidator()
result, err := validator.Validate(testCases)
```

## 接口抽象设计

### 1. AI客户端接口 (client/interface.go)
```go
type AIClient interface {
    // 通用生成接口
    Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
    
    // 通用聊天接口
    Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    
    // 获取模型列表
    GetModels(ctx context.Context) ([]Model, error)
}
```

### 2. 提示词构建接口 (prompt/types.go)
```go
type Builder interface {
    // 构建提示词
    Build(templateName string, data interface{}) (string, error)
    
    // 加载模板
    LoadTemplate(name string, template string) error
}

type TemplateEngine interface {
    // 渲染模板
    Render(template string, data interface{}) (string, error)
}
```

### 3. 响应处理接口 (processor/types.go)
```go
type Parser interface {
    // 解析响应内容
    Parse(content string, target interface{}) error
}

type Validator interface {
    // 验证数据有效性
    Validate(data interface{}) (*ValidationResult, error)
}
```

## 多层交互的抽象接口

### 业务层 ↔ AI工具层
```go
// 业务层定义自己的服务接口
type TestCaseGeneratorService interface {
    Generate(req *TestCaseRequest) (*TestCaseResponse, error)
}

// 业务层实现依赖AI工具层的接口
type testCaseService struct {
    aiClient      client.AIClient      // 依赖抽象接口
    promptBuilder prompt.Builder       // 依赖抽象接口  
    parser        processor.Parser     // 依赖抽象接口
    validator     processor.Validator  // 依赖抽象接口
}
```

### AI工具层内部交互
```go
// 各组件通过接口交互，而不是具体实现
type AIService struct {
    client    AIClient    // 接口依赖
    builder   Builder     // 接口依赖
    parser    Parser      // 接口依赖
    validator Validator   // 接口依赖
}
```

## 设计优势

### 1. 解耦性
- 各层通过接口交互，具体实现可以随时替换
- 业务层不依赖具体的AI提供商
- 各组件职责单一，易于维护

### 2. 可测试性
- 所有依赖都是接口，可以轻松Mock
- 单元测试可以独立进行
- 集成测试可以使用真实或模拟的实现

### 3. 可扩展性
- 新增AI提供商只需实现AIClient接口
- 新增解析器只需实现Parser接口
- 新增验证器只需实现Validator接口

### 4. 可复用性
- AI工具包完全通用，可用于任何业务场景
- 提示词工具可用于各种文本生成场景
- 响应处理器可处理各种AI返回格式

## 使用示例

```go
// 依赖注入方式组装
func NewTestCaseService(
    aiClient client.AIClient,
    promptBuilder prompt.Builder,
    parser processor.Parser,
    validator processor.Validator,
) TestCaseGeneratorService {
    return &testCaseService{
        aiClient:      aiClient,
        promptBuilder: promptBuilder,
        parser:        parser,
        validator:     validator,
    }
}

// 业务逻辑实现
func (s *testCaseService) Generate(req *TestCaseRequest) (*TestCaseResponse, error) {
    // 1. 构建提示词
    prompt, err := s.promptBuilder.Build("testcase", req)
    if err != nil {
        return nil, err
    }
    
    // 2. 调用AI
    response, err := s.aiClient.Generate(ctx, &GenerateRequest{
        Prompt: prompt,
        Format: "json",
    })
    if err != nil {
        return nil, err
    }
    
    // 3. 解析响应
    var testCases []TestCase
    err = s.parser.Parse(response.Content, &testCases)
    if err != nil {
        return nil, err
    }
    
    // 4. 验证结果
    validationResult, err := s.validator.Validate(testCases)
    if err != nil {
        return nil, err
    }
    
    return &TestCaseResponse{
        TestCases: testCases,
        Quality:   validationResult.Quality,
    }, nil
}
```

这种设计确保了各层之间通过抽象接口交互，实现了高内聚、低耦合的架构。
