# AI测试用例生成功能设计文档

## 1. 功能概述

### 1.1 功能描述
开发一个AI驱动的测试用例生成功能，通过调用大模型API（主要使用阿里云千炼），根据题目描述自动生成高质量的测试用例。

### 1.2 核心目标
- **自动化生成**：根据题目描述自动生成多样化的测试用例
- **质量保证**：确保生成的测试用例语法正确、逻辑合理
- **覆盖度完整**：包含基础、边界、极端等多种类型的测试用例
- **高区分度**：生成的测试用例能够有效区分不同水平的解答
- **可扩展性**：支持多种大模型提供商，便于后续扩展

### 1.3 业务价值
- 减少出题老师的工作量
- 提高测试用例的质量和覆盖度
- 标准化测试用例生成流程
- 提升OJ平台的整体质量

## 2. 需求分析

### 2.1 功能需求

#### 2.1.1 核心功能
1. **测试用例生成**
   - 根据题目描述生成测试用例
   - 支持指定生成数量和类型
   - 支持批量生成

2. **质量验证**
   - 语法正确性验证
   - 逻辑合理性验证
   - 答案正确性验证

3. **质量评估**
   - 测试用例覆盖度评估
   - 难度分布评估
   - 区分度评估

4. **结果管理**
   - 生成结果存储
   - 生成历史查询
   - 结果导出

#### 2.1.2 扩展功能
1. **多模型支持**
   - 阿里云千炼
   - OpenAI GPT系列
   - 其他大模型

2. **提示词优化**
   - 提示词模板管理
   - 动态提示词构建
   - 提示词效果评估

3. **智能优化**
   - 根据反馈优化生成策略
   - 学习历史生成模式
   - 自适应参数调整

### 2.2 非功能需求

#### 2.2.1 性能需求
- 单次生成响应时间 < 30秒
- 支持并发生成请求
- 系统可用性 > 99%

#### 2.2.2 安全需求
- API密钥安全存储
- 请求数据加密传输
- 访问权限控制

#### 2.2.3 可维护性需求
- 模块化设计，低耦合
- 完善的日志记录
- 错误处理和恢复机制

## 3. 系统架构设计

### 3.1 整体架构

```
┌─────────────────────────────────────────────────────────────┐
│                    前端界面层                                │
├─────────────────────────────────────────────────────────────┤
│                    API接口层                                │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  测试用例API    │  │   其他API       │                  │
│  └─────────────────┘  └─────────────────┘                  │
├─────────────────────────────────────────────────────────────┤
│                    业务服务层                               │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │ 测试用例服务    │  │   其他服务      │                  │
│  └─────────────────┘  └─────────────────┘                  │
├─────────────────────────────────────────────────────────────┤
│                    AI工具层                                 │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────┐ │
│  │   AI客户端      │  │   提示词工具    │  │  响应处理   │ │
│  └─────────────────┘  └─────────────────┘  └─────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    数据访问层                               │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │   MongoDB       │  │     缓存        │                  │
│  └─────────────────┘  └─────────────────┘                  │
├─────────────────────────────────────────────────────────────┤
│                    外部服务层                               │
│  ┌─────────────────┐  ┌─────────────────┐                  │
│  │  阿里云千炼     │  │   其他LLM       │                  │
│  └─────────────────┘  └─────────────────┘                  │
└─────────────────────────────────────────────────────────────┘
```

### 3.2 包结构设计

#### 3.2.1 AI工具包（通用层）
```
pkg/ai/
├── client/                    # AI客户端抽象层
│   ├── interface.go          # 通用AI客户端接口
│   ├── qianfan.go           # 阿里云千炼实现
│   ├── openai.go            # OpenAI实现
│   ├── factory.go           # 客户端工厂
│   └── types.go             # 通用类型定义
├── prompt/                   # 提示词工具
│   ├── template.go          # 模板引擎
│   ├── builder.go           # 构建器
│   └── types.go             # 提示词类型
└── processor/                # 响应处理器
    ├── parser.go            # 解析器
    ├── validator.go         # 验证器
    └── types.go             # 处理器类型
```

#### 3.2.2 测试用例业务包
```
pkg/testcase/
├── generator/                # 生成器
│   ├── interface.go         # 生成器接口
│   ├── ai_generator.go      # AI生成器
│   └── rule_generator.go    # 规则生成器
├── validator/                # 验证器
│   ├── interface.go         # 验证器接口
│   ├── syntax_validator.go  # 语法验证
│   └── logic_validator.go   # 逻辑验证
├── evaluator/                # 评估器
│   ├── interface.go         # 评估器接口
│   ├── quality_evaluator.go # 质量评估
│   └── coverage_evaluator.go# 覆盖度评估
└── types.go                  # 业务类型定义
```

#### 3.2.3 API和服务层
```
pkg/app/api-server/
├── dto/
│   └── dto_testcase.go      # 测试用例DTO
├── service/
│   └── service_testcase.go  # 测试用例服务
└── server/
    └── server_testcase.go   # 测试用例API
```

#### 3.2.4 数据层
```
pkg/
├── models/
│   └── mongo_testcase.go    # 测试用例数据模型
└── dao/
    └── dao_testcase.go      # 测试用例数据访问
```

### 3.3 核心流程设计

#### 3.3.1 测试用例生成流程
```
用户请求 → API接口 → 服务层 → AI生成器 → 提示词构建 → LLM调用 → 响应解析 → 质量验证 → 结果评估 → 数据存储 → 返回结果
```

#### 3.3.2 质量保证流程
```
生成测试用例 → 语法验证 → 逻辑验证 → 答案验证 → 覆盖度检查 → 质量评分 → 不合格重新生成 → 合格返回结果
```

## 4. 接口设计

### 4.1 REST API设计

#### 4.1.1 生成测试用例
```
POST /api/v1/testcase/generate
Content-Type: application/json

Request:
{
  "problem_id": "string",
  "problem_title": "string", 
  "description": "string",
  "input_format": "string",
  "output_format": "string",
  "constraints": "string",
  "examples": [
    {
      "input": "string",
      "output": "string"
    }
  ],
  "difficulty": "string",
  "tags": ["string"],
  "count": 10,
  "types": ["basic", "boundary", "extreme"],
  "llm_provider": "qianfan",
  "options": {
    "validate": true,
    "max_retries": 3,
    "timeout": 30
  }
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": "string",
    "problem_id": "string",
    "test_cases": [
      {
        "input": "string",
        "expected_output": "string",
        "type": "basic",
        "description": "string",
        "difficulty": 5
      }
    ],
    "quality": "high",
    "coverage": {
      "has_basic_cases": true,
      "has_boundary_cases": true,
      "has_extreme_cases": true,
      "coverage_score": 85
    },
    "statistics": {
      "total_count": 10,
      "type_distribution": {},
      "difficulty_distribution": {},
      "avg_difficulty": 5.2
    },
    "generated_at": "2025-09-21T22:00:00Z",
    "generated_by": "ai-generator",
    "llm_provider": "qianfan"
  }
}
```

#### 4.1.2 查询生成历史
```
GET /api/v1/testcase/history?problem_id={id}&page={page}&size={size}

Response:
{
  "code": 200,
  "message": "success", 
  "data": {
    "total": 100,
    "page": 1,
    "size": 10,
    "items": [
      {
        "id": "string",
        "problem_id": "string",
        "quality": "high",
        "test_case_count": 10,
        "generated_at": "2025-09-21T22:00:00Z",
        "generated_by": "user123"
      }
    ]
  }
}
```

#### 4.1.3 验证测试用例
```
POST /api/v1/testcase/validate

Request:
{
  "problem_id": "string",
  "test_cases": [
    {
      "input": "string",
      "expected_output": "string"
    }
  ]
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "results": [
      {
        "is_valid": true,
        "syntax_errors": [],
        "logic_errors": [],
        "warnings": [],
        "score": 95
      }
    ],
    "overall_valid": true,
    "overall_score": 92
  }
}
```

### 4.2 内部接口设计

#### 4.2.1 AI客户端接口
```go
type AIClient interface {
    Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
    Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
    GetModels(ctx context.Context) ([]Model, error)
}
```

#### 4.2.2 测试用例生成器接口
```go
type Generator interface {
    Generate(ctx context.Context, req *GenerationRequest) (*TestCaseSet, error)
    GenerateWithValidation(ctx context.Context, req *GenerationRequest) (*TestCaseSet, error)
    RegenerateFailedCases(ctx context.Context, req *GenerationRequest, failedCases []TestCase) (*TestCaseSet, error)
}
```

#### 4.2.3 验证器接口
```go
type Validator interface {
    Validate(ctx context.Context, testCase *TestCase, problem *Problem) (*ValidationResult, error)
    ValidateBatch(ctx context.Context, testCases []TestCase, problem *Problem) ([]ValidationResult, error)
}
```

## 5. 数据模型设计

### 5.1 测试用例生成记录
```go
type TestCaseGeneration struct {
    ID          primitive.ObjectID `bson:"_id,omitempty"`
    ProblemID   string            `bson:"problem_id"`
    UserID      string            `bson:"user_id"`
    TestCases   []TestCase        `bson:"test_cases"`
    Quality     string            `bson:"quality"`
    Coverage    Coverage          `bson:"coverage"`
    Statistics  Statistics        `bson:"statistics"`
    LLMProvider string            `bson:"llm_provider"`
    PromptUsed  string            `bson:"prompt_used"`
    GeneratedAt time.Time         `bson:"generated_at"`
    Status      int               `bson:"status"` // 1:成功 0:失败
}
```

### 5.2 测试用例
```go
type TestCase struct {
    Input          string `bson:"input"`
    ExpectedOutput string `bson:"expected_output"`
    Type           string `bson:"type"`
    Description    string `bson:"description"`
    Difficulty     int    `bson:"difficulty"`
    Tags           []string `bson:"tags"`
}
```

### 5.3 质量评估
```go
type QualityAssessment struct {
    OverallQuality string        `bson:"overall_quality"`
    Scores         QualityScores `bson:"scores"`
    Suggestions    []string      `bson:"suggestions"`
    IsAcceptable   bool          `bson:"is_acceptable"`
}
```

## 6. 配置设计

### 6.1 AI配置扩展
```yaml
AI:
  DefaultProvider: "qianfan"
  Providers:
    qianfan:
      APIKey: "your-api-key"
      BaseURL: "https://aip.baidubce.com"
      Model: "ERNIE-Bot-turbo"
      Timeout: 30
      MaxRetries: 3
    openai:
      APIKey: "your-api-key"
      BaseURL: "https://api.openai.com"
      Model: "gpt-3.5-turbo"
      Timeout: 30
      MaxRetries: 3

TestCase:
  DefaultCount: 10
  MaxCount: 50
  DefaultTypes: ["basic", "boundary", "extreme"]
  QualityThreshold: 80
  ValidationEnabled: true
```

## 7. 实施计划

### 7.1 开发阶段
1. **阶段一**：AI工具包开发（1周）
   - AI客户端接口设计
   - 阿里云千炼客户端实现
   - 提示词工具开发

2. **阶段二**：测试用例业务包开发（1周）
   - 生成器实现
   - 验证器实现
   - 评估器实现

3. **阶段三**：API和服务层开发（1周）
   - REST API实现
   - 服务层逻辑
   - 数据访问层

4. **阶段四**：集成测试和优化（1周）
   - 功能测试
   - 性能优化
   - 错误处理完善

### 7.2 测试计划
1. **单元测试**：各模块独立测试
2. **集成测试**：模块间协作测试
3. **性能测试**：并发和响应时间测试
4. **用户测试**：实际使用场景测试

### 7.3 部署计划
1. **开发环境**：本地开发测试
2. **测试环境**：功能验证
3. **生产环境**：正式上线

## 8. 风险评估

### 8.1 技术风险
- LLM API稳定性和可用性
- 生成质量的一致性
- 响应时间的不确定性

### 8.2 业务风险
- 生成的测试用例质量不达标
- 成本控制（API调用费用）
- 用户接受度

### 8.3 风险缓解措施
- 多LLM提供商支持
- 质量验证和重试机制
- 成本监控和限制
- 用户反馈收集和优化
