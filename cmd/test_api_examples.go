package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"zhku-oj-server/pkg/ai/client"
	"zhku-oj-server/pkg/utils"
)

func main() {
	// 初始化AI服务
	configDir := "conf"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = "../conf"
	}

	log.Println("正在初始化AI服务...")
	_, err := utils.InitAI(configDir)
	if err != nil {
		log.Fatalf("❌ AI服务初始化失败: %v", err)
	}
	log.Println("🎉 AI服务初始化完成")

	// 获取默认客户端
	aiManager := utils.GetAIClientManager()
	defaultClient, err := aiManager.GetDefaultClient()
	if err != nil {
		log.Fatalf("❌ 获取AI客户端失败: %v", err)
	}

	fmt.Println("🧪 测试API文档中的示例代码")
	fmt.Println(strings.Repeat("=", 60))

	// 测试基础Chat示例
	fmt.Println("\n📋 1. 基础Chat示例")
	basicChatExample(defaultClient)

	// 测试多轮对话
	fmt.Println("\n📋 2. 多轮对话示例")
	multiTurnChatExample(defaultClient)

	// 测试JSON响应
	fmt.Println("\n📋 3. JSON响应示例")
	jsonResponseExample(defaultClient)

	// 测试基础Generate
	fmt.Println("\n📋 4. 基础Generate示例")
	basicGenerateExample(defaultClient)

	// 测试JSON生成
	fmt.Println("\n📋 5. JSON生成示例")
	jsonGenerateExample(defaultClient)

	// 测试错误处理
	fmt.Println("\n📋 6. 错误处理示例")
	errorHandlingExample(defaultClient)

	// 测试参数验证
	fmt.Println("\n📋 7. 参数验证示例")
	parameterValidationExample()

	// 测试超时控制
	fmt.Println("\n📋 8. 超时控制示例")
	timeoutExample(defaultClient)

	// 测试代码生成场景
	fmt.Println("\n📋 9. 代码生成场景")
	codeGenerationExample(defaultClient)

	fmt.Println("\n🎉 所有示例测试完成！")
}

// 基础对话示例
func basicChatExample(aiClient client.AIClient) {
	chatReq := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "你好，请用一句话介绍Go语言"},
		},
		Model: "qwen-plus",
	}

	response, err := aiClient.Chat(context.Background(), chatReq)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		return
	}

	fmt.Printf("✅ 成功!\n")
	fmt.Printf("   回复: %s\n", response.Content)
	fmt.Printf("   Token使用: %d\n", response.Usage.TotalTokens)
}

// 多轮对话示例
func multiTurnChatExample(aiClient client.AIClient) {
	chatReq := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "system", Content: "你是一个编程助手"},
			{Role: "user", Content: "什么是递归？"},
			{Role: "assistant", Content: "递归是函数调用自身的编程技术，用于解决可以分解为相似子问题的问题。"},
			{Role: "user", Content: "能给个简单例子吗？"},
		},
		Model:       "qwen-plus",
		Temperature: 0.3,
		MaxTokens:   300,
	}

	response, err := aiClient.Chat(context.Background(), chatReq)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		return
	}

	fmt.Printf("✅ 成功!\n")
	fmt.Printf("   回复: %s\n", response.Content)
}

// JSON响应示例
func jsonResponseExample(aiClient client.AIClient) {
	chatReq := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "生成一个用户信息的JSON示例，包含姓名、年龄、邮箱"},
		},
		Model: "qwen-plus",
		ResponseFormat: &client.ResponseFormat{
			Type: "json_object",
		},
	}

	response, err := aiClient.Chat(context.Background(), chatReq)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		return
	}

	fmt.Printf("✅ 成功!\n")
	fmt.Printf("   JSON响应: %s\n", response.Content)
}

// 基础文本生成示例
func basicGenerateExample(aiClient client.AIClient) {
	genReq := &client.GenerateRequest{
		Prompt: "写一个计算两个数最大公约数的函数（用Go语言）",
		Model:  "qwen-plus",
	}

	response, err := aiClient.Generate(context.Background(), genReq)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		return
	}

	fmt.Printf("✅ 成功!\n")
	fmt.Printf("   生成内容: %s...\n", response.Content[:200])
	fmt.Printf("   Token使用: %d\n", response.Usage.TotalTokens)
}

// JSON格式生成示例
func jsonGenerateExample(aiClient client.AIClient) {
	genReq := &client.GenerateRequest{
		Prompt: "生成一个包含3个测试用例的JSON对象，每个测试用例包含输入和期望输出",
		Model:  "qwen-plus",
		Format: "json",
	}

	response, err := aiClient.Generate(context.Background(), genReq)
	if err != nil {
		fmt.Printf("❌ 错误: %v\n", err)
		return
	}

	fmt.Printf("✅ 成功!\n")
	fmt.Printf("   JSON内容: %s\n", response.Content)
}

// 错误处理示例
func errorHandlingExample(aiClient client.AIClient) {
	// 故意使用无效的模型名称来触发错误
	chatReq := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "测试"},
		},
		Model: "invalid-model-name",
	}

	response, err := aiClient.Chat(context.Background(), chatReq)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "timeout"):
			fmt.Println("⏰ 请求超时，请稍后重试")
		case strings.Contains(err.Error(), "invalid model") || strings.Contains(err.Error(), "model"):
			fmt.Println("🚫 模型名称无效")
		case strings.Contains(err.Error(), "quota exceeded"):
			fmt.Println("📊 配额已用完")
		case strings.Contains(err.Error(), "invalid parameter"):
			fmt.Println("⚠️ 参数错误")
		default:
			fmt.Printf("❓ 未知错误: %v\n", err)
		}
		return
	}

	fmt.Printf("✅ 意外成功: %s\n", response.Content)
}

// 参数验证示例
func parameterValidationExample() {
	fmt.Println("🔍 测试参数验证...")

	// 测试空模型名称
	req1 := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "测试"},
		},
		Model: "", // 空模型名称
	}
	if err := validateChatRequest(req1); err != nil {
		fmt.Printf("✅ 验证成功: %v\n", err)
	}

	// 测试无效温度
	req2 := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "测试"},
		},
		Model:       "qwen-plus",
		Temperature: 1.5, // 无效温度
	}
	if err := validateChatRequest(req2); err != nil {
		fmt.Printf("✅ 验证成功: %v\n", err)
	}
}

// 参数验证函数
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

// 超时控制示例
func timeoutExample(aiClient client.AIClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	chatReq := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "写一个详细的算法分析"},
		},
		Model:   "qwen-plus",
		Timeout: 8 * time.Second,
	}

	response, err := aiClient.Chat(ctx, chatReq)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			fmt.Println("⏰ 请求超时")
		} else {
			fmt.Printf("❌ 其他错误: %v\n", err)
		}
		return
	}

	fmt.Printf("✅ 成功: %s...\n", response.Content[:100])
}

// 代码生成场景示例
func codeGenerationExample(aiClient client.AIClient) {
	code, err := generateCode(aiClient, "Go", "实现一个简单的栈数据结构")
	if err != nil {
		fmt.Printf("❌ 代码生成失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 代码生成成功!\n")
	fmt.Printf("   生成的代码: %s...\n", code[:200])
}

// 代码生成函数
func generateCode(aiClient client.AIClient, language, description string) (string, error) {
	prompt := fmt.Sprintf("请用%s语言实现：%s\n只返回代码，不要解释。", language, description)

	genReq := &client.GenerateRequest{
		Prompt:      prompt,
		Model:       "qwen-plus",
		Temperature: 0.2,
		MaxTokens:   800,
	}

	response, err := aiClient.Generate(context.Background(), genReq)
	if err != nil {
		return "", err
	}

	return response.Content, nil
}
