# AI测试用例生成技术实现方案

## 1. 技术选型

### 1.1 大模型选择
- **主要选择**：阿里云千炼（百度文心一言）
  - 中文支持好
  - API稳定性高
  - 成本相对较低
  - 有专门的代码生成能力

- **备选方案**：
  - OpenAI GPT-3.5/4
  - 智谱AI GLM
  - 讯飞星火

### 1.2 技术栈
- **后端**：Go 1.23+
- **数据库**：MongoDB（已有）
- **缓存**：Redis（可选）
- **HTTP客户端**：标准库 net/http
- **JSON处理**：标准库 encoding/json
- **配置管理**：viper（已有）
- **日志**：logrus（已有）

## 2. 核心技术实现

### 2.1 AI客户端设计模式

#### 2.1.1 策略模式 + 工厂模式
```go
// 策略接口
type AIClient interface {
    Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error)
}

// 工厂接口
type ClientFactory interface {
    CreateClient(provider string, config *Config) (AIClient, error)
}

// 具体实现
type QianfanClient struct {
    apiKey  string
    baseURL string
    client  *http.Client
}

type OpenAIClient struct {
    apiKey  string
    baseURL string
    client  *http.Client
}
```

#### 2.1.2 适配器模式
```go
// 统一请求格式
type UnifiedRequest struct {
    Prompt      string
    Model       string
    Temperature float64
    MaxTokens   int
}

// 各厂商适配器
type QianfanAdapter struct{}
func (a *QianfanAdapter) AdaptRequest(req *UnifiedRequest) *QianfanRequest
func (a *QianfanAdapter) AdaptResponse(resp *QianfanResponse) *UnifiedResponse
```

### 2.2 提示词工程

#### 2.2.1 模板系统
```go
type PromptTemplate struct {
    Name        string
    Template    string
    Variables   []string
    Examples    []PromptExample
    Metadata    map[string]interface{}
}

type PromptBuilder struct {
    templates map[string]*PromptTemplate
    engine    *template.Template
}

func (pb *PromptBuilder) Build(templateName string, data interface{}) (string, error)
```

#### 2.2.2 提示词模板示例
```
你是一个专业的算法题测试用例生成专家。请根据以下题目信息生成高质量的测试用例。

## 题目信息
- 标题：{{.Title}}
- 描述：{{.Description}}
- 输入格式：{{.InputFormat}}
- 输出格式：{{.OutputFormat}}
- 约束条件：{{.Constraints}}
- 难度：{{.Difficulty}}

## 示例用例
{{range .Examples}}
输入：{{.Input}}
输出：{{.Output}}
{{end}}

## 生成要求
1. 生成{{.Count}}个测试用例
2. 包含以下类型：{{range .RequiredTypes}}{{.}} {{end}}
3. 确保测试用例的正确性和多样性
4. 覆盖边界条件和极端情况

## 输出格式
请严格按照以下JSON格式输出：
```json
{
  "test_cases": [
    {
      "input": "具体输入数据",
      "expected_output": "期望输出结果", 
      "type": "basic|boundary|extreme|corner",
      "description": "测试用例说明",
      "difficulty": 1-10
    }
  ]
}
```

请开始生成测试用例：
```

### 2.3 响应解析策略

#### 2.3.1 多层解析
```go
type ResponseParser struct {
    jsonParser   *JSONParser
    textParser   *TextParser
    regexParser  *RegexParser
}

func (rp *ResponseParser) Parse(content string) (*ParsedResult, error) {
    // 1. 尝试JSON解析
    if result, err := rp.jsonParser.Parse(content); err == nil {
        return result, nil
    }
    
    // 2. 尝试文本解析
    if result, err := rp.textParser.Parse(content); err == nil {
        return result, nil
    }
    
    // 3. 正则表达式兜底
    return rp.regexParser.Parse(content)
}
```

#### 2.3.2 容错机制
```go
type FaultTolerantParser struct {
    parsers []Parser
    cleaner *ContentCleaner
}

func (ftp *FaultTolerantParser) Parse(content string) (*ParsedResult, error) {
    // 内容清理
    cleaned := ftp.cleaner.Clean(content)
    
    // 多解析器尝试
    for _, parser := range ftp.parsers {
        if result, err := parser.Parse(cleaned); err == nil {
            return result, nil
        }
    }
    
    return nil, errors.New("all parsers failed")
}
```

### 2.4 质量验证机制

#### 2.4.1 多维度验证
```go
type ValidationPipeline struct {
    validators []Validator
}

type Validator interface {
    Validate(ctx context.Context, testCase *TestCase, problem *Problem) *ValidationResult
}

// 语法验证器
type SyntaxValidator struct{}
func (sv *SyntaxValidator) Validate(ctx context.Context, testCase *TestCase, problem *Problem) *ValidationResult

// 逻辑验证器  
type LogicValidator struct{}
func (lv *LogicValidator) Validate(ctx context.Context, testCase *TestCase, problem *Problem) *ValidationResult

// 答案验证器
type AnswerValidator struct {
    judgeService JudgeService
}
func (av *AnswerValidator) Validate(ctx context.Context, testCase *TestCase, problem *Problem) *ValidationResult
```

#### 2.4.2 验证规则
```go
type ValidationRule struct {
    Name        string
    Description string
    Validator   func(*TestCase, *Problem) error
    Weight      float64
}

var DefaultValidationRules = []ValidationRule{
    {
        Name: "input_format_check",
        Validator: func(tc *TestCase, p *Problem) error {
            return validateInputFormat(tc.Input, p.InputFormat)
        },
        Weight: 0.3,
    },
    {
        Name: "constraint_check", 
        Validator: func(tc *TestCase, p *Problem) error {
            return validateConstraints(tc.Input, p.Constraints)
        },
        Weight: 0.4,
    },
    {
        Name: "output_correctness",
        Validator: func(tc *TestCase, p *Problem) error {
            return validateOutputCorrectness(tc, p)
        },
        Weight: 0.3,
    },
}
```

### 2.5 缓存策略

#### 2.5.1 多级缓存
```go
type CacheManager struct {
    memoryCache *MemoryCache
    redisCache  *RedisCache
    dbCache     *DatabaseCache
}

func (cm *CacheManager) Get(key string) (*CachedResult, error) {
    // L1: 内存缓存
    if result, found := cm.memoryCache.Get(key); found {
        return result, nil
    }
    
    // L2: Redis缓存
    if result, err := cm.redisCache.Get(key); err == nil {
        cm.memoryCache.Set(key, result)
        return result, nil
    }
    
    // L3: 数据库缓存
    if result, err := cm.dbCache.Get(key); err == nil {
        cm.redisCache.Set(key, result)
        cm.memoryCache.Set(key, result)
        return result, nil
    }
    
    return nil, ErrCacheNotFound
}
```

#### 2.5.2 缓存键设计
```go
func GenerateCacheKey(req *GenerationRequest) string {
    h := sha256.New()
    h.Write([]byte(req.ProblemID))
    h.Write([]byte(req.Description))
    h.Write([]byte(fmt.Sprintf("%d", req.Count)))
    h.Write([]byte(strings.Join(req.RequiredTypes, ",")))
    return fmt.Sprintf("testcase:%x", h.Sum(nil))
}
```

### 2.6 错误处理和重试机制

#### 2.6.1 指数退避重试
```go
type RetryConfig struct {
    MaxRetries      int
    InitialDelay    time.Duration
    MaxDelay        time.Duration
    BackoffFactor   float64
    RetryableErrors []string
}

func (rc *RetryConfig) ShouldRetry(err error, attempt int) bool {
    if attempt >= rc.MaxRetries {
        return false
    }
    
    for _, retryableErr := range rc.RetryableErrors {
        if strings.Contains(err.Error(), retryableErr) {
            return true
        }
    }
    
    return false
}

func (rc *RetryConfig) GetDelay(attempt int) time.Duration {
    delay := time.Duration(float64(rc.InitialDelay) * math.Pow(rc.BackoffFactor, float64(attempt)))
    if delay > rc.MaxDelay {
        delay = rc.MaxDelay
    }
    return delay
}
```

#### 2.6.2 熔断器模式
```go
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    failures    int
    lastFailure time.Time
    state       CircuitState
    mutex       sync.RWMutex
}

type CircuitState int

const (
    StateClosed CircuitState = iota
    StateOpen
    StateHalfOpen
)

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mutex.RLock()
    state := cb.state
    cb.mutex.RUnlock()
    
    switch state {
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.setState(StateHalfOpen)
            return cb.callAndUpdate(fn)
        }
        return ErrCircuitBreakerOpen
    case StateHalfOpen:
        return cb.callAndUpdate(fn)
    default: // StateClosed
        return cb.callAndUpdate(fn)
    }
}
```

## 3. 性能优化策略

### 3.1 并发处理
```go
type ConcurrentGenerator struct {
    workerPool   *WorkerPool
    rateLimiter  *RateLimiter
    semaphore    chan struct{}
}

func (cg *ConcurrentGenerator) GenerateBatch(requests []*GenerationRequest) ([]*TestCaseSet, error) {
    results := make([]*TestCaseSet, len(requests))
    errors := make([]error, len(requests))
    
    var wg sync.WaitGroup
    for i, req := range requests {
        wg.Add(1)
        go func(index int, request *GenerationRequest) {
            defer wg.Done()
            
            // 获取信号量
            cg.semaphore <- struct{}{}
            defer func() { <-cg.semaphore }()
            
            // 限流
            cg.rateLimiter.Wait()
            
            // 生成
            result, err := cg.generate(request)
            results[index] = result
            errors[index] = err
        }(i, req)
    }
    
    wg.Wait()
    return results, combineErrors(errors)
}
```

### 3.2 流式处理
```go
type StreamProcessor struct {
    inputChan  chan *GenerationRequest
    outputChan chan *GenerationResult
    workers    int
}

func (sp *StreamProcessor) Start() {
    for i := 0; i < sp.workers; i++ {
        go sp.worker()
    }
}

func (sp *StreamProcessor) worker() {
    for req := range sp.inputChan {
        result := sp.processRequest(req)
        sp.outputChan <- result
    }
}
```

### 3.3 资源池管理
```go
type ResourcePool struct {
    clients chan AIClient
    factory ClientFactory
    config  *PoolConfig
}

func (rp *ResourcePool) Get() (AIClient, error) {
    select {
    case client := <-rp.clients:
        return client, nil
    default:
        return rp.factory.CreateClient(rp.config)
    }
}

func (rp *ResourcePool) Put(client AIClient) {
    select {
    case rp.clients <- client:
    default:
        client.Close()
    }
}
```

## 4. 监控和日志

### 4.1 指标收集
```go
type Metrics struct {
    GenerationCount    prometheus.Counter
    GenerationDuration prometheus.Histogram
    ErrorCount         prometheus.Counter
    QualityScore       prometheus.Gauge
}

func (m *Metrics) RecordGeneration(duration time.Duration, quality float64, err error) {
    m.GenerationCount.Inc()
    m.GenerationDuration.Observe(duration.Seconds())
    m.QualityScore.Set(quality)
    
    if err != nil {
        m.ErrorCount.Inc()
    }
}
```

### 4.2 结构化日志
```go
type GenerationLogger struct {
    logger *logrus.Logger
}

func (gl *GenerationLogger) LogGeneration(ctx context.Context, req *GenerationRequest, result *TestCaseSet, err error) {
    fields := logrus.Fields{
        "problem_id":   req.ProblemID,
        "count":        req.Count,
        "llm_provider": req.LLMProvider,
        "duration":     time.Since(ctx.Value("start_time").(time.Time)),
    }
    
    if result != nil {
        fields["quality"] = result.Quality
        fields["test_case_count"] = len(result.TestCases)
    }
    
    if err != nil {
        fields["error"] = err.Error()
        gl.logger.WithFields(fields).Error("Generation failed")
    } else {
        gl.logger.WithFields(fields).Info("Generation completed")
    }
}
```

## 5. 安全考虑

### 5.1 API密钥管理
```go
type SecretManager struct {
    vault VaultClient
    cache map[string]*CachedSecret
    mutex sync.RWMutex
}

func (sm *SecretManager) GetAPIKey(provider string) (string, error) {
    sm.mutex.RLock()
    if cached, exists := sm.cache[provider]; exists && !cached.IsExpired() {
        sm.mutex.RUnlock()
        return cached.Value, nil
    }
    sm.mutex.RUnlock()
    
    // 从安全存储获取
    secret, err := sm.vault.GetSecret(fmt.Sprintf("ai/%s/api_key", provider))
    if err != nil {
        return "", err
    }
    
    // 更新缓存
    sm.mutex.Lock()
    sm.cache[provider] = &CachedSecret{
        Value:     secret,
        ExpiresAt: time.Now().Add(time.Hour),
    }
    sm.mutex.Unlock()
    
    return secret, nil
}
```

### 5.2 请求验证
```go
type RequestValidator struct {
    maxPromptLength int
    allowedProviders map[string]bool
    rateLimiter     *RateLimiter
}

func (rv *RequestValidator) Validate(req *GenerationRequest, userID string) error {
    // 检查提示词长度
    if len(req.Description) > rv.maxPromptLength {
        return ErrPromptTooLong
    }
    
    // 检查提供商
    if !rv.allowedProviders[req.LLMProvider] {
        return ErrInvalidProvider
    }
    
    // 检查频率限制
    if !rv.rateLimiter.Allow(userID) {
        return ErrRateLimitExceeded
    }
    
    return nil
}
```

## 6. 测试策略

### 6.1 单元测试
```go
func TestAIGenerator_Generate(t *testing.T) {
    mockClient := &MockAIClient{}
    mockClient.On("Generate", mock.Anything, mock.Anything).Return(&GenerateResponse{
        Content: `{"test_cases": [{"input": "1", "expected_output": "1"}]}`,
    }, nil)
    
    generator := NewAIGenerator(mockClient, nil, nil, nil)
    
    req := &GenerationRequest{
        ProblemID: "test-problem",
        Count:     1,
    }
    
    result, err := generator.Generate(context.Background(), req)
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Len(t, result.TestCases, 1)
}
```

### 6.2 集成测试
```go
func TestTestCaseGenerationFlow(t *testing.T) {
    // 设置测试环境
    testDB := setupTestDB()
    testConfig := loadTestConfig()
    
    // 创建服务
    service := NewTestCaseService(testDB, testConfig)
    
    // 执行测试
    req := &GenerationRequest{
        ProblemID:   "integration-test",
        Description: "Calculate sum of two numbers",
        Count:       5,
    }
    
    result, err := service.Generate(context.Background(), req)
    
    // 验证结果
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Len(t, result.TestCases, 5)
    
    // 验证数据库存储
    stored, err := testDB.GetTestCaseGeneration(result.ID)
    assert.NoError(t, err)
    assert.Equal(t, result.ProblemID, stored.ProblemID)
}
```

### 6.3 性能测试
```go
func BenchmarkTestCaseGeneration(b *testing.B) {
    generator := setupBenchmarkGenerator()
    req := &GenerationRequest{
        ProblemID: "benchmark-test",
        Count:     10,
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := generator.Generate(context.Background(), req)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

这个技术实现方案提供了详细的代码架构和实现策略，确保了系统的可扩展性、可维护性和高性能。
