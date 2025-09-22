package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"zhku-oj-server/pkg/ai/client"
	"zhku-oj-server/pkg/ai/config"
)

// 示例：如何在主函数中集成AI配置和客户端管理
func main() {
	// 0. 检查环境变量
	fmt.Println("=== 检查环境变量 ===")
	if err := config.CheckAllRequired(); err != nil {
		log.Fatal(err)
	}

	// 1. 创建配置管理器
	fmt.Println("\n=== 初始化AI配置管理器 ===")
	configManager, err := config.NewConfigManager("../../../conf/ai.yaml")
	if err != nil {
		log.Fatalf("创建配置管理器失败: %v", err)
	}

	// 验证配置加载
	aiConfig := configManager.GetConfig()
	fmt.Printf("默认提供商: %s\n", aiConfig.DefaultProvider)
	fmt.Printf("启用的提供商: %v\n", configManager.GetEnabledProviders())

	// 2. 创建客户端工厂
	fmt.Println("\n=== 创建客户端工厂 ===")
	factory := client.NewFactory(configManager)
	fmt.Printf("支持的提供商: %v\n", factory.GetSupportedProviders())

	// 3. 创建客户端管理器
	fmt.Println("\n=== 创建客户端管理器 ===")
	clientManager := client.NewManager(factory, configManager)

	// 4. 测试连接
	fmt.Println("\n=== 测试提供商连接 ===")
	testResults := clientManager.TestAllConnections()
	for provider, err := range testResults {
		if err != nil {
			fmt.Printf("❌ %s: %v\n", provider, err)
		} else {
			fmt.Printf("✅ %s: 连接正常\n", provider)
		}
	}

	// 5. 获取默认客户端并测试
	fmt.Println("\n=== 测试默认客户端 ===")
	defaultClient, err := clientManager.GetDefaultClient()
	if err != nil {
		log.Printf("获取默认客户端失败: %v", err)
		return
	}

	// 测试获取模型列表
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	models, err := defaultClient.GetModels(ctx)
	if err != nil {
		log.Printf("获取模型列表失败: %v", err)
	} else {
		fmt.Printf("可用模型:\n")
		for _, model := range models {
			fmt.Printf("  - %s: %s (最大tokens: %d)\n", model.ID, model.Name, model.MaxTokens)
		}
	}

	// 6. 测试生成功能
	fmt.Println("\n=== 测试生成功能 ===")
	testGenerate(defaultClient)

	// 7. 测试场景化客户端
	fmt.Println("\n=== 测试场景化客户端 ===")
	testScenarioClient(clientManager)

	// 8. 测试故障转移
	fmt.Println("\n=== 测试故障转移 ===")
	testFallback(clientManager)

	// 9. 清理资源
	fmt.Println("\n=== 清理资源 ===")
	if err := clientManager.Close(); err != nil {
		log.Printf("关闭客户端管理器失败: %v", err)
	}

	fmt.Println("示例程序执行完成")
}

// testGenerate 测试生成功能
func testGenerate(aiClient client.AIClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &client.GenerateRequest{
		Prompt:      "请生成一个简单的Hello World程序的测试用例",
		Model:       "qwen-turbo",
		Temperature: 0.7,
		MaxTokens:   1000,
		Format:      "json",
	}

	response, err := aiClient.Generate(ctx, req)
	if err != nil {
		log.Printf("生成请求失败: %v", err)
		return
	}

	fmt.Printf("生成结果:\n")
	fmt.Printf("  ID: %s\n", response.ID)
	fmt.Printf("  模型: %s\n", response.Model)
	fmt.Printf("  内容: %s\n", response.Content)
	fmt.Printf("  Token使用: %d/%d/%d (提示/完成/总计)\n", 
		response.Usage.PromptTokens, 
		response.Usage.CompletionTokens, 
		response.Usage.TotalTokens)
}

// testScenarioClient 测试场景化客户端
func testScenarioClient(manager *client.Manager) {
	// 测试测试用例生成场景
	scenarioClient, err := manager.GetClientByScenario("testcase_generation")
	if err != nil {
		log.Printf("获取场景客户端失败: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req := &client.GenerateRequest{
		Prompt:      "为排序算法生成测试用例",
		Temperature: 0.7,
		MaxTokens:   2000,
		Format:      "json",
	}

	response, err := scenarioClient.Generate(ctx, req)
	if err != nil {
		log.Printf("场景生成请求失败: %v", err)
		return
	}

	fmt.Printf("场景生成结果:\n")
	fmt.Printf("  场景: testcase_generation\n")
	fmt.Printf("  内容长度: %d 字符\n", len(response.Content))
	fmt.Printf("  Token使用: %d\n", response.Usage.TotalTokens)
}

// testFallback 测试故障转移
func testFallback(manager *client.Manager) {
	err := manager.ExecuteWithFallback("testcase_generation", func(aiClient client.AIClient) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req := &client.GenerateRequest{
			Prompt:      "测试故障转移功能",
			Temperature: 0.5,
			MaxTokens:   500,
		}

		response, err := aiClient.Generate(ctx, req)
		if err != nil {
			return err
		}

		fmt.Printf("故障转移测试成功:\n")
		fmt.Printf("  内容: %s\n", response.Content[:min(100, len(response.Content))])
		return nil
	})

	if err != nil {
		log.Printf("故障转移测试失败: %v", err)
	}
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 业务服务示例
type TestCaseService struct {
	clientManager *client.Manager
}

// NewTestCaseService 创建测试用例服务
func NewTestCaseService(clientManager *client.Manager) *TestCaseService {
	return &TestCaseService{
		clientManager: clientManager,
	}
}

// GenerateTestCases 生成测试用例
func (s *TestCaseService) GenerateTestCases(problemDescription string) (*TestCaseSet, error) {
	// 使用场景化客户端
	aiClient, err := s.clientManager.GetClientByScenario("testcase_generation")
	if err != nil {
		return nil, fmt.Errorf("获取AI客户端失败: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	prompt := fmt.Sprintf(`
请为以下编程题目生成测试用例：

题目描述：
%s

要求：
1. 生成至少5个测试用例
2. 包含边界情况和极端情况
3. 返回JSON格式
4. 每个测试用例包含输入和期望输出

返回格式：
{
  "test_cases": [
    {
      "input": "输入数据",
      "expected_output": "期望输出",
      "description": "用例描述"
    }
  ]
}
`, problemDescription)

	req := &client.GenerateRequest{
		Prompt:      prompt,
		Temperature: 0.7,
		MaxTokens:   4000,
		Format:      "json",
	}

	response, err := aiClient.Generate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("AI生成失败: %w", err)
	}

	// 这里应该解析JSON响应为TestCaseSet结构
	// 为了示例简化，直接返回基本信息
	return &TestCaseSet{
		ProblemID:   "example",
		TestCases:   []TestCase{}, // 实际应该解析response.Content
		GeneratedBy: "AI",
		CreatedAt:   time.Now(),
	}, nil
}

// TestCaseSet 测试用例集合
type TestCaseSet struct {
	ProblemID   string     `json:"problem_id"`
	TestCases   []TestCase `json:"test_cases"`
	GeneratedBy string     `json:"generated_by"`
	CreatedAt   time.Time  `json:"created_at"`
}

// TestCase 测试用例
type TestCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	Description    string `json:"description"`
}
