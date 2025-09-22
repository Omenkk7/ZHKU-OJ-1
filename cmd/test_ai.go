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
	// 获取配置目录路径
	configDir := "conf"
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		configDir = "../conf"
	}

	// 初始化AI服务
	_, err := utils.InitAI(configDir)
	if err != nil {
		log.Fatalf("AI服务初始化失败: %v", err)
	}

	// 获取AI客户端管理器
	aiManager := utils.GetAIClientManager()
	if aiManager == nil {
		log.Fatal("AI服务不可用")
	}

	// 获取默认客户端（qwen-plus）
	defaultClient, err := aiManager.GetDefaultClient()
	if err != nil {
		log.Fatalf("获取默认客户端失败: %v", err)
	}

	fmt.Println("🤖 测试qwen-plus模型单次对话...")

	// 创建聊天请求
	chatReq := &client.ChatRequest{
		Messages: []client.Message{
			{Role: "user", Content: "你好，请简单介绍一下你自己，不超过50字"},
		},
		Model:       "qwen-plus",
		Temperature: 0.7,
		MaxTokens:   100,
		TopP:        0.9, // 设置有效的top_p值
	}

	// 发送聊天请求
	ctx := context.Background()
	response, err := defaultClient.Chat(ctx, chatReq)
	if err != nil {
		log.Fatalf("聊天请求失败: %v", err)
	}

	// 打印结果
	fmt.Printf("✅ 对话成功！\n")
	fmt.Printf("📝 AI回复: %s\n", response.Content)
	fmt.Printf("🔧 使用模型: %s\n", response.Model)
	fmt.Printf("📊 Token使用: 输入=%d, 输出=%d, 总计=%d\n", 
		response.Usage.PromptTokens, 
		response.Usage.CompletionTokens, 
		response.Usage.TotalTokens)
	fmt.Printf("⏰ 响应时间: %v\n", response.Created)

	fmt.Println("\n🎉 qwen-plus模型和单次对话功能测试成功！")
}
