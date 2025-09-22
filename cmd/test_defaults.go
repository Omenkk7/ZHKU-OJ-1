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
	if aiManager == nil {
		log.Fatal("AI服务不可用")
	}

	defaultClient, err := aiManager.GetDefaultClient()
	if err != nil {
		log.Fatalf("获取默认客户端失败: %v", err)
	}

	fmt.Println("🧪 测试Chat方法的默认值设置...")

	// 测试1: 不传递任何参数（使用默认值）
	fmt.Println("\n📋 测试1: 使用默认值")
	chatReq1 := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "简单回答：1+1等于几？"},
		},
		Model: "qwen-plus",
		// 不设置 Temperature, MaxTokens, TopP - 应该使用默认值
	}

	response1, err := defaultClient.Chat(context.Background(), chatReq1)
	if err != nil {
		log.Printf("❌ 测试1失败: %v", err)
	} else {
		fmt.Printf("✅ 测试1成功！\n")
		fmt.Printf("   回复: %s\n", response1.Content)
		fmt.Printf("   Token使用: %d\n", response1.Usage.TotalTokens)
	}

	// 测试2: 显式设置参数
	fmt.Println("\n📋 测试2: 显式设置参数")
	chatReq2 := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "简单回答：2+2等于几？"},
		},
		Model:       "qwen-plus",
		Temperature: 0.1, // 低温度，更确定性的回答
		MaxTokens:   50,  // 限制token数
		TopP:        0.5, // 较低的top_p
	}

	response2, err := defaultClient.Chat(context.Background(), chatReq2)
	if err != nil {
		log.Printf("❌ 测试2失败: %v", err)
	} else {
		fmt.Printf("✅ 测试2成功！\n")
		fmt.Printf("   回复: %s\n", response2.Content)
		fmt.Printf("   Token使用: %d\n", response2.Usage.TotalTokens)
	}

	// 测试3: 部分使用默认值
	fmt.Println("\n📋 测试3: 部分使用默认值")
	chatReq3 := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "简单回答：3+3等于几？"},
		},
		Model:     "qwen-plus",
		MaxTokens: 100, // 只设置MaxTokens，其他使用默认值
		// Temperature 和 TopP 使用默认值
	}

	response3, err := defaultClient.Chat(context.Background(), chatReq3)
	if err != nil {
		log.Printf("❌ 测试3失败: %v", err)
	} else {
		fmt.Printf("✅ 测试3成功！\n")
		fmt.Printf("   回复: %s\n", response3.Content)
		fmt.Printf("   Token使用: %d\n", response3.Usage.TotalTokens)
	}

	// 测试4: Generate接口的默认值
	fmt.Println("\n📋 测试4: Generate接口默认值")
	genReq := &client.GenerateRequest{
		Prompt: "写一个简单的Hello World程序",
		Model:  "qwen-plus",
		// 不设置其他参数，使用默认值
	}

	genResponse, err := defaultClient.Generate(context.Background(), genReq)
	if err != nil {
		log.Printf("❌ 测试4失败: %v", err)
	} else {
		fmt.Printf("✅ 测试4成功！\n")
		fmt.Printf("   生成内容: %s\n", genResponse.Content[:100]+"...") // 只显示前100字符
		fmt.Printf("   Token使用: %d\n", genResponse.Usage.TotalTokens)
	}

	fmt.Println("\n🎉 默认值测试完成！")
	fmt.Println("\n📊 默认值设置:")
	fmt.Println("   - Temperature: 0.7 (如果未设置)")
	fmt.Println("   - MaxTokens: 1000 (如果未设置)")
	fmt.Println("   - TopP: 0.9 (如果未设置)")
	fmt.Println("   - Stream: false (默认)")
}
