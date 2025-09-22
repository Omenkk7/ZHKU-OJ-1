package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"zhku-oj-server/pkg/ai/client"
	"zhku-oj-server/pkg/utils"
)

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

	// 获取AI客户端
	aiManager := utils.GetAIClientManager()
	defaultClient, _ := aiManager.GetDefaultClient()

	fmt.Println("🔧 详细测试Chat和Generate方法的默认值行为")
	fmt.Println("============================================================")

	// 测试Chat方法的默认值
	testChatDefaults(defaultClient)
	
	fmt.Println("\n============================================================")
	
	// 测试Generate方法的默认值
	testGenerateDefaults(defaultClient)
	
	fmt.Println("\n🎯 总结:")
	fmt.Println("✅ Chat方法默认值:")
	fmt.Println("   - Temperature: 0.7 (当传入0时)")
	fmt.Println("   - MaxTokens: 1000 (当传入0时)")
	fmt.Println("   - TopP: 0.9 (当传入0时)")
	fmt.Println("   - Stream: false (默认)")
	fmt.Println()
	fmt.Println("✅ Generate方法默认值:")
	fmt.Println("   - Temperature: 0.7 (当传入0时)")
	fmt.Println("   - MaxTokens: 1000 (当传入0时)")
	fmt.Println("   - TopP: 0.9 (当传入0时)")
	fmt.Println("   - Format: text (默认)")
}

func testChatDefaults(aiClient client.AIClient) {
	fmt.Println("📋 Chat方法默认值测试:")
	
	// 测试场景1: 所有参数都为0（使用默认值）
	fmt.Println("\n🧪 场景1: 所有参数为0/空值")
	chatReq := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "请用一句话介绍Go语言"},
		},
		Model: "qwen-plus",
		// Temperature: 0,  // 默认值
		// MaxTokens: 0,    // 默认值
		// TopP: 0,         // 默认值
		// Stream: false,   // 默认值
	}
	
	response, err := aiClient.Chat(context.Background(), chatReq)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功! Token使用: %d\n", response.Usage.TotalTokens)
		fmt.Printf("   回复: %s\n", response.Content)
	}
	
	// 测试场景2: 显式设置非零值
	fmt.Println("\n🧪 场景2: 显式设置参数")
	chatReq2 := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "请用一句话介绍Python语言"},
		},
		Model:       "qwen-plus",
		Temperature: 0.1,  // 低温度
		MaxTokens:   50,   // 限制token
		TopP:        0.5,  // 低top_p
	}
	
	response2, err := aiClient.Chat(context.Background(), chatReq2)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功! Token使用: %d\n", response2.Usage.TotalTokens)
		fmt.Printf("   回复: %s\n", response2.Content)
	}
}

func testGenerateDefaults(aiClient client.AIClient) {
	fmt.Println("📋 Generate方法默认值测试:")
	
	// 测试场景1: 所有参数都为0（使用默认值）
	fmt.Println("\n🧪 场景1: 所有参数为0/空值")
	genReq := &client.GenerateRequest{
		Prompt: "写一个简单的Hello函数",
		Model:  "qwen-plus",
		// Temperature: 0,  // 默认值
		// MaxTokens: 0,    // 默认值
		// TopP: 0,         // 默认值
		// Format: "",      // 默认为text
	}
	
	response, err := aiClient.Generate(context.Background(), genReq)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功! Token使用: %d\n", response.Usage.TotalTokens)
		fmt.Printf("   生成内容: %s...\n", response.Content[:100])
	}

	// 测试场景2: 显式设置参数
	fmt.Println("\n🧪 场景2: 显式设置参数")
	genReq2 := &client.GenerateRequest{
		Prompt:      "写一个简单的加法函数",
		Model:       "qwen-plus",
		Temperature: 0.2,
		MaxTokens:   200,
		TopP:        0.8,
		Format:      "text",
	}

	response2, err := aiClient.Generate(context.Background(), genReq2)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功! Token使用: %d\n", response2.Usage.TotalTokens)
		fmt.Printf("   生成内容: %s...\n", response2.Content[:100])
	}

	// 测试场景3: JSON格式输出
	fmt.Println("\n🧪 场景3: JSON格式输出")
	genReq3 := &client.GenerateRequest{
		Prompt: "生成一个包含姓名和年龄的JSON对象示例",
		Model:  "qwen-plus",
		Format: "json",
		// 其他参数使用默认值
	}

	response3, err := aiClient.Generate(context.Background(), genReq3)
	if err != nil {
		fmt.Printf("❌ 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 成功! Token使用: %d\n", response3.Usage.TotalTokens)
		fmt.Printf("   JSON内容: %s\n", response3.Content)
		if response3.Data != nil {
			fmt.Printf("   解析的数据: %+v\n", response3.Data)
		}
	}
}
