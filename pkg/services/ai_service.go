package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"zhku-oj-server/pkg/ai/client"
)

// AIService AI业务服务
type AIService struct {
	clientManager *client.Manager
}

// NewAIService 创建AI业务服务
func NewAIService(clientManager *client.Manager) *AIService {
	return &AIService{
		clientManager: clientManager,
	}
}

// TestCaseGenerationRequest 测试用例生成请求
type TestCaseGenerationRequest struct {
	ProblemID          string `json:"problem_id"`
	ProblemTitle       string `json:"problem_title"`
	ProblemDescription string `json:"problem_description"`
	InputFormat        string `json:"input_format"`
	OutputFormat       string `json:"output_format"`
	SampleInput        string `json:"sample_input"`
	SampleOutput       string `json:"sample_output"`
	Constraints        string `json:"constraints"`
	TestCaseCount      int    `json:"test_case_count"` // 需要生成的测试用例数量
}

// TestCase 测试用例
type TestCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	Description    string `json:"description"`
	Type           string `json:"type"` // "basic", "boundary", "extreme"
}

// TestCaseGenerationResponse 测试用例生成响应
type TestCaseGenerationResponse struct {
	ProblemID   string     `json:"problem_id"`
	TestCases   []TestCase `json:"test_cases"`
	GeneratedBy string     `json:"generated_by"`
	GeneratedAt time.Time  `json:"generated_at"`
	TokenUsage  int        `json:"token_usage"`
}

// GenerateTestCases 生成测试用例
func (s *AIService) GenerateTestCases(ctx context.Context, req *TestCaseGenerationRequest) (*TestCaseGenerationResponse, error) {
	if s.clientManager == nil {
		return nil, fmt.Errorf("AI服务未初始化")
	}

	// 获取测试用例生成场景的客户端
	aiClient, err := s.clientManager.GetClientByScenario("testcase_generation")
	if err != nil {
		return nil, fmt.Errorf("获取AI客户端失败: %w", err)
	}

	// 构建提示词
	prompt := s.buildTestCasePrompt(req)

	// 调用AI生成
	response, err := aiClient.Generate(ctx, &client.GenerateRequest{
		Prompt:      prompt,
		Temperature: 0.7,
		MaxTokens:   4000,
		Format:      "json",
	})

	if err != nil {
		return nil, fmt.Errorf("AI生成失败: %w", err)
	}

	// 解析响应
	testCases, err := s.parseTestCaseResponse(response.Content)
	if err != nil {
		return nil, fmt.Errorf("解析AI响应失败: %w", err)
	}

	return &TestCaseGenerationResponse{
		ProblemID:   req.ProblemID,
		TestCases:   testCases,
		GeneratedBy: "AI",
		GeneratedAt: time.Now(),
		TokenUsage:  response.Usage.TotalTokens,
	}, nil
}

// buildTestCasePrompt 构建测试用例生成提示词
func (s *AIService) buildTestCasePrompt(req *TestCaseGenerationRequest) string {
	testCaseCount := req.TestCaseCount
	if testCaseCount <= 0 {
		testCaseCount = 5 // 默认生成5个测试用例
	}

	prompt := fmt.Sprintf(`
请为以下编程题目生成 %d 个高质量的测试用例：

题目信息：
- 题目ID: %s
- 题目标题: %s
- 题目描述: %s
- 输入格式: %s
- 输出格式: %s
- 样例输入: %s
- 样例输出: %s
- 约束条件: %s

要求：
1. 生成 %d 个测试用例，包含以下类型：
   - 基础测试用例（basic）：验证基本功能
   - 边界测试用例（boundary）：测试边界条件
   - 极端测试用例（extreme）：测试极限情况
2. 每个测试用例必须包含输入、期望输出和描述
3. 确保测试用例的正确性和多样性
4. 输入输出格式必须严格遵循题目要求

返回JSON格式：
{
  "test_cases": [
    {
      "input": "具体的输入数据",
      "expected_output": "对应的期望输出",
      "description": "测试用例描述",
      "type": "basic|boundary|extreme"
    }
  ]
}

请确保返回的是有效的JSON格式。
`, testCaseCount, req.ProblemID, req.ProblemTitle, req.ProblemDescription,
		req.InputFormat, req.OutputFormat, req.SampleInput, req.SampleOutput,
		req.Constraints, testCaseCount)

	return prompt
}

// parseTestCaseResponse 解析AI响应
func (s *AIService) parseTestCaseResponse(content string) ([]TestCase, error) {
	var response struct {
		TestCases []TestCase `json:"test_cases"`
	}

	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	if len(response.TestCases) == 0 {
		return nil, fmt.Errorf("AI未生成任何测试用例")
	}

	// 验证测试用例
	for i, testCase := range response.TestCases {
		if testCase.Input == "" {
			return nil, fmt.Errorf("测试用例 %d 缺少输入", i+1)
		}
		if testCase.ExpectedOutput == "" {
			return nil, fmt.Errorf("测试用例 %d 缺少期望输出", i+1)
		}
		if testCase.Type == "" {
			response.TestCases[i].Type = "basic" // 默认类型
		}
		if testCase.Description == "" {
			response.TestCases[i].Description = fmt.Sprintf("测试用例 %d", i+1)
		}
	}

	return response.TestCases, nil
}

// CodeReviewRequest 代码审查请求
type CodeReviewRequest struct {
	Code     string `json:"code"`
	Language string `json:"language"`
	ProblemID string `json:"problem_id,omitempty"`
}

// CodeReviewResponse 代码审查响应
type CodeReviewResponse struct {
	Issues      []CodeIssue `json:"issues"`
	Suggestions []string    `json:"suggestions"`
	Score       int         `json:"score"` // 代码质量评分 0-100
	GeneratedAt time.Time   `json:"generated_at"`
	TokenUsage  int         `json:"token_usage"`
}

// CodeIssue 代码问题
type CodeIssue struct {
	Type        string `json:"type"`        // "error", "warning", "info"
	Line        int    `json:"line"`        // 行号
	Message     string `json:"message"`     // 问题描述
	Suggestion  string `json:"suggestion"`  // 修改建议
	Severity    string `json:"severity"`    // "high", "medium", "low"
}

// ReviewCode 代码审查
func (s *AIService) ReviewCode(ctx context.Context, req *CodeReviewRequest) (*CodeReviewResponse, error) {
	if s.clientManager == nil {
		return nil, fmt.Errorf("AI服务未初始化")
	}

	// 获取代码审查场景的客户端
	aiClient, err := s.clientManager.GetClientByScenario("code_review")
	if err != nil {
		return nil, fmt.Errorf("获取AI客户端失败: %w", err)
	}

	// 构建提示词
	prompt := s.buildCodeReviewPrompt(req)

	// 调用AI生成
	response, err := aiClient.Generate(ctx, &client.GenerateRequest{
		Prompt:      prompt,
		Temperature: 0.3, // 代码审查需要更准确的结果
		MaxTokens:   2000,
		Format:      "json",
	})

	if err != nil {
		return nil, fmt.Errorf("AI代码审查失败: %w", err)
	}

	// 解析响应
	reviewResult, err := s.parseCodeReviewResponse(response.Content)
	if err != nil {
		return nil, fmt.Errorf("解析代码审查响应失败: %w", err)
	}

	reviewResult.GeneratedAt = time.Now()
	reviewResult.TokenUsage = response.Usage.TotalTokens

	return reviewResult, nil
}

// buildCodeReviewPrompt 构建代码审查提示词
func (s *AIService) buildCodeReviewPrompt(req *CodeReviewRequest) string {
	prompt := fmt.Sprintf(`
请对以下 %s 代码进行详细审查：

代码：
```%s
%s
```

请从以下方面进行审查：
1. 代码正确性和逻辑错误
2. 性能优化建议
3. 代码风格和规范
4. 安全性问题
5. 可维护性和可读性

返回JSON格式：
{
  "issues": [
    {
      "type": "error|warning|info",
      "line": 行号,
      "message": "问题描述",
      "suggestion": "修改建议",
      "severity": "high|medium|low"
    }
  ],
  "suggestions": ["总体建议1", "总体建议2"],
  "score": 85
}

请确保返回有效的JSON格式。
`, req.Language, req.Language, req.Code)

	return prompt
}

// parseCodeReviewResponse 解析代码审查响应
func (s *AIService) parseCodeReviewResponse(content string) (*CodeReviewResponse, error) {
	var response CodeReviewResponse

	if err := json.Unmarshal([]byte(content), &response); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %w", err)
	}

	// 验证评分范围
	if response.Score < 0 {
		response.Score = 0
	} else if response.Score > 100 {
		response.Score = 100
	}

	return &response, nil
}

// IsAvailable 检查AI服务是否可用
func (s *AIService) IsAvailable() bool {
	if s.clientManager == nil {
		return false
	}

	// 测试默认客户端连接
	err := s.clientManager.TestConnection("")
	return err == nil
}

// GetStatus 获取AI服务状态
func (s *AIService) GetStatus() map[string]interface{} {
	if s.clientManager == nil {
		return map[string]interface{}{
			"available": false,
			"error":     "AI服务未初始化",
		}
	}

	// 测试所有提供商连接
	testResults := s.clientManager.TestAllConnections()
	enabledProviders := s.clientManager.GetEnabledProviders()

	status := map[string]interface{}{
		"available":         len(enabledProviders) > 0,
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

// GetClientManager 获取AI客户端管理器（供其他服务使用）
func (s *AIService) GetClientManager() *client.Manager {
	return s.clientManager
}
