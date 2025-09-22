# AI配置管理文档

## 概述

AI配置管理系统提供了灵活的多提供商配置、场景化配置和动态配置管理能力。

## 配置文件结构

### 基本结构
```yaml
default_provider: "qianfan"    # 默认AI提供商
providers:                     # AI提供商配置
  qianfan: { ... }
  openai: { ... }
scenarios:                     # 场景化配置
  testcase_generation: { ... }
  code_review: { ... }
global:                        # 全局配置
  enable_cache: true
  cache_ttl: 10m
```

## 提供商配置

### 支持的提供商

| 提供商 | 标识 | 说明 | 状态 |
|--------|------|------|------|
| 阿里云百炼 | dashscope | 阿里云百炼大模型平台 | ✅ 推荐 |
| OpenAI | openai | GPT系列模型 | 💤 暂时禁用 |
| 智谱AI | glm | GLM系列模型 | 💤 暂时禁用 |
| 讯飞星火 | spark | 星火认知大模型 | 💤 暂时禁用 |

### 提供商配置参数

```yaml
providers:
  dashscope:
    provider: "dashscope"            # 提供商标识
    api_key: "${DASHSCOPE_API_KEY}"  # API密钥（支持环境变量）
    base_url: "https://dashscope.aliyuncs.com"  # API基础URL
    model: "qwen-turbo"              # 默认模型
    timeout: 30s                     # 请求超时时间
    max_retries: 3                   # 最大重试次数
    enabled: true                    # 是否启用
    headers:                         # 自定义请求头
      User-Agent: "ZHKU-OJ-Server/1.0"
    proxy: "http://proxy:8080"       # 代理地址（可选）
```

### 环境变量支持

配置文件中可以使用环境变量：
- `${DASHSCOPE_API_KEY}` - 阿里云百炼API密钥
- `${OPENAI_API_KEY}` - OpenAI API密钥（暂时禁用）
- `${GLM_API_KEY}` - 智谱AI API密钥（暂时禁用）
- `${SPARK_API_KEY}` - 讯飞星火API密钥（暂时禁用）

如果环境变量未设置，对应的提供商将自动禁用。

## 场景化配置

### 预定义场景

| 场景 | 标识 | 说明 | 推荐提供商 |
|------|------|------|------------|
| 测试用例生成 | testcase_generation | 生成编程题测试用例 | dashscope |
| 代码审查 | code_review | 代码质量分析 | dashscope |
| 题目生成 | problem_generation | 生成编程题目 | dashscope |
| 聊天助手 | chat_assistant | AI对话助手 | dashscope |

### 场景配置参数

```yaml
scenarios:
  testcase_generation:
    primary_provider: "dashscope"         # 主要提供商
    fallback_providers: []                # 暂时不使用备用提供商
    config:
      temperature: 0.7    # 创造性参数 (0-1)
      max_tokens: 4000    # 最大token数
      top_p: 0.9         # 核采样参数
      format: "json"      # 输出格式: json/text
```

### 参数说明

- **temperature**: 控制输出的随机性和创造性
  - 0.0: 最确定性的输出
  - 1.0: 最随机的输出
  - 推荐值: 0.3-0.8

- **max_tokens**: 限制生成的最大token数量
  - 根据具体需求设置
  - 测试用例生成: 3000-5000
  - 代码审查: 1500-2500

- **top_p**: 核采样参数，控制候选词的范围
  - 0.1: 只考虑概率最高的10%的词
  - 1.0: 考虑所有词
  - 推荐值: 0.8-0.95

- **format**: 输出格式
  - "text": 纯文本格式
  - "json": JSON格式（需要在提示词中说明）

## 全局配置

```yaml
global:
  enable_cache: true          # 启用缓存
  cache_ttl: 10m             # 缓存过期时间
  enable_metrics: true        # 启用指标收集
  enable_retry: true          # 启用重试机制
  default_timeout: 30s        # 默认超时时间
  max_concurrency: 10         # 最大并发数
  rate_limit_rps: 100         # 每秒请求限制
```

## 配置管理器使用

### 基本使用

```go
import "zhku-oj-server/pkg/ai/config"

// 创建配置管理器
manager, err := config.NewManager("conf/ai.yaml")
if err != nil {
    log.Fatal(err)
}

// 获取配置
aiConfig := manager.GetConfig()

// 获取特定提供商配置
dashscopeConfig, err := manager.GetProviderConfig("dashscope")
if err != nil {
    log.Fatal(err)
}

// 获取场景配置
scenarioConfig, err := manager.GetScenarioConfig("testcase_generation")
if err != nil {
    log.Fatal(err)
}
```

### 配置监听

```go
// 添加配置变更监听器
manager.AddWatcher(func(newConfig *config.AIConfig) {
    log.Println("配置已更新")
    // 重新初始化相关组件
})

// 重新加载配置
if err := manager.Reload(); err != nil {
    log.Printf("重新加载配置失败: %v", err)
}
```

### 获取启用的提供商

```go
// 获取所有启用的提供商
enabledProviders := manager.GetEnabledProviders()
fmt.Printf("启用的提供商: %v\n", enabledProviders)
```

## 配置文件位置

### 默认位置
- `conf/ai.yaml` - 主配置文件
- `conf/ai.local.yaml` - 本地开发配置（可选）
- `conf/ai.prod.yaml` - 生产环境配置（可选）

### 环境特定配置
可以根据环境使用不同的配置文件：
```bash
# 开发环境
export AI_CONFIG_PATH="conf/ai.local.yaml"

# 生产环境  
export AI_CONFIG_PATH="conf/ai.prod.yaml"
```

## 最佳实践

### 1. 环境变量管理
```bash
# 设置环境变量
export DASHSCOPE_API_KEY="your-dashscope-api-key"
# export OPENAI_API_KEY="your-openai-api-key"  # 暂时不需要

# 或使用 .env 文件
echo "DASHSCOPE_API_KEY=your-key" >> .env
```

### 2. 提供商选择策略
- **主要场景**: 使用阿里云百炼，成本较低、响应较快
- **备用方案**: 目前专注于百炼，后续可扩展其他提供商
- **特殊需求**: 根据具体需求调整模型参数

### 3. 参数调优
- **测试用例生成**: temperature=0.7, 平衡创造性和准确性
- **代码审查**: temperature=0.3, 注重准确性和一致性
- **创意生成**: temperature=0.8, 提高创造性

### 4. 监控和日志
- 启用指标收集监控API调用情况
- 记录配置变更日志
- 监控各提供商的成功率和响应时间

### 5. 安全考虑
- API密钥使用环境变量，不要硬编码
- 定期轮换API密钥
- 限制API调用频率，避免超出配额
- 在生产环境中禁用不必要的提供商

## 故障排除

### 常见问题

1. **提供商不可用**
   - 检查API密钥是否正确
   - 检查网络连接
   - 检查提供商服务状态

2. **配置加载失败**
   - 检查配置文件语法
   - 检查文件路径是否正确
   - 检查文件权限

3. **环境变量未设置**
   - 检查环境变量是否正确设置
   - 检查变量名是否匹配

### 调试模式
```go
// 启用详细日志
manager.SetLogLevel("debug")

// 验证配置
if err := manager.ValidateConfig(); err != nil {
    log.Printf("配置验证失败: %v", err)
}
```
