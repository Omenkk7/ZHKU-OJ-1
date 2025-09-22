package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"zhku-oj-server/pkg/ai/client"
	"zhku-oj-server/pkg/utils"
)

// TestCase 测试用例结构
type TestCase struct {
	Input    interface{} `json:"input"`
	Expected interface{} `json:"expected"`
	Name     string      `json:"name"`
	Desc     string      `json:"description"`
}

// TestSuite 测试套件
type TestSuite struct {
	FunctionName string     `json:"function_name"`
	TestCases    []TestCase `json:"test_cases"`
	Language     string     `json:"language"`
}

func main() {
	// 初始化AI服务
	configDir := "conf"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = "../conf"
	}

	_, err := utils.InitAI(configDir)
	if err != nil {
		log.Fatalf("AI服务初始化失败: %v", err)
	}

	// 测试函数代码
	functionCode := `
func Add(a, b int) int {
    return a + b
}
`

	// 生成测试用例
	testSuite, err := generateTestCases(functionCode, "Add", "go")
	if err != nil {
		log.Fatalf("生成测试用例失败: %v", err)
	}

	// 打印结果
	fmt.Printf("🎯 为函数 %s 生成的测试用例:\n\n", testSuite.FunctionName)
	for i, testCase := range testSuite.TestCases {
		fmt.Printf("测试用例 %d: %s\n", i+1, testCase.Name)
		fmt.Printf("  描述: %s\n", testCase.Desc)
		fmt.Printf("  输入: %v\n", testCase.Input)
		fmt.Printf("  期望: %v\n\n", testCase.Expected)
	}

	// 测试简单文本生成
	fmt.Println("==================================================")
	fmt.Println("🔧 测试简单文本生成:")
	
	simpleResult, err := generateSimpleText("写一个Go语言的快速排序函数")
	if err != nil {
		log.Printf("简单文本生成失败: %v", err)
	} else {
		fmt.Printf("生成的代码:\n%s\n", simpleResult)
	}
}

// generateTestCases 使用Generate接口生成测试用例
func generateTestCases(functionCode, functionName, language string) (*TestSuite, error) {
	// 获取AI客户端
	aiManager := utils.GetAIClientManager()
	if aiManager == nil {
		return nil, fmt.Errorf("AI服务不可用")
	}

	defaultClient, err := aiManager.GetDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("获取AI客户端失败: %v", err)
	}

	// 构建提示词
	prompt := fmt.Sprintf(`请为以下%s函数生成完整的测试用例，要求：
1. 包含正常情况、边界情况、异常情况
2. 每个测试用例包含输入、期望输出、名称和描述
3. 返回JSON格式，严格按照以下结构：

{
  "function_name": "%s",
  "language": "%s", 
  "test_cases": [
    {
      "name": "测试用例名称",
      "description": "测试用例描述",
      "input": {"a": 1, "b": 2},
      "expected": 3
    }
  ]
}

函数代码：
%s

请直接返回JSON，不要包含其他文字说明。`, language, functionName, language, functionCode)

	// 创建Generate请求
	genReq := &client.GenerateRequest{
		Prompt:      prompt,
		Model:       "qwen-plus",
		Temperature: 0.3, // 较低温度确保输出稳定
		MaxTokens:   2000,
		TopP:        0.9,
		Format:      "json", // 指定JSON格式
	}

	// 调用Generate接口
	ctx := context.Background()
	response, err := defaultClient.Generate(ctx, genReq)
	if err != nil {
		return nil, fmt.Errorf("AI生成失败: %v", err)
	}

	fmt.Printf("🤖 AI原始响应:\n%s\n\n", response.Content)

	// 解析JSON响应
	var testSuite TestSuite
	if err := json.Unmarshal([]byte(response.Content), &testSuite); err != nil {
		return nil, fmt.Errorf("解析JSON响应失败: %v", err)
	}

	return &testSuite, nil
}

// generateSimpleText 使用Generate接口进行简单文本生成
func generateSimpleText(prompt string) (string, error) {
	// 获取AI客户端
	aiManager := utils.GetAIClientManager()
	if aiManager == nil {
		return "", fmt.Errorf("AI服务不可用")
	}

	defaultClient, err := aiManager.GetDefaultClient()
	if err != nil {
		return "", fmt.Errorf("获取AI客户端失败: %v", err)
	}

	// 创建Generate请求
	genReq := &client.GenerateRequest{
		Prompt:      prompt,
		Model:       "qwen-plus",
		Temperature: 0.7,
		MaxTokens:   1000,
		TopP:        0.9,
		Format:      "text", // 文本格式
	}

	// 调用Generate接口
	ctx := context.Background()
	response, err := defaultClient.Generate(ctx, genReq)
	if err != nil {
		return "", fmt.Errorf("AI生成失败: %v", err)
	}

	return response.Content, nil
}
