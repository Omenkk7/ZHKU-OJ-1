package config

import (
	"fmt"
	"log"
	"os"
)

// EnvChecker 环境变量检查器
type EnvChecker struct {
	requiredVars map[string]string // 环境变量名 -> 描述
}

// NewEnvChecker 创建环境变量检查器
func NewEnvChecker() *EnvChecker {
	return &EnvChecker{
		requiredVars: make(map[string]string),
	}
}

// AddRequired 添加必需的环境变量
func (e *EnvChecker) AddRequired(varName, description string) {
	e.requiredVars[varName] = description
}

// CheckAll 检查所有必需的环境变量
func (e *EnvChecker) CheckAll() error {
	var missingVars []string
	
	for varName, description := range e.requiredVars {
		value := os.Getenv(varName)
		if value == "" {
			missingVars = append(missingVars, fmt.Sprintf("%s (%s)", varName, description))
		} else {
			log.Printf("✅ %s 已设置 (长度: %d)", varName, len(value))
		}
	}
	
	if len(missingVars) > 0 {
		return fmt.Errorf("缺少必需的环境变量:\n%v\n\n请设置这些环境变量后重试", missingVars)
	}
	
	return nil
}

// CheckProvider 检查特定提供商的环境变量
func CheckProvider(provider string) error {
	var envKey string
	var description string
	
	switch provider {
	case "dashscope":
		envKey = "DASHSCOPE_API_KEY"
		description = "阿里云百炼API密钥"
	case "openai":
		envKey = "OPENAI_API_KEY"
		description = "OpenAI API密钥"
	case "glm":
		envKey = "GLM_API_KEY"
		description = "智谱AI API密钥"
	case "spark":
		envKey = "SPARK_API_KEY"
		description = "讯飞星火API密钥"
	default:
		return fmt.Errorf("未知的提供商: %s", provider)
	}
	
	apiKey := os.Getenv(envKey)
	if apiKey == "" {
		return fmt.Errorf("请设置 %s 环境变量 (%s)", envKey, description)
	}
	
	log.Printf("✅ %s 已设置 (长度: %d)", envKey, len(apiKey))
	return nil
}

// CheckDashScope 检查DashScope环境变量
func CheckDashScope() error {
	return CheckProvider("dashscope")
}

// CheckOpenAI 检查OpenAI环境变量
func CheckOpenAI() error {
	return CheckProvider("openai")
}

// SetupExample 设置示例环境变量（仅用于测试）
func SetupExample() {
	// 仅用于开发测试，生产环境不应使用
	if os.Getenv("DASHSCOPE_API_KEY") == "" {
		log.Println("⚠️  警告: 未检测到 DASHSCOPE_API_KEY，设置示例值（仅用于测试）")
		os.Setenv("DASHSCOPE_API_KEY", "example-api-key-for-testing")
	}
}

// PrintEnvStatus 打印环境变量状态
func PrintEnvStatus() {
	fmt.Println("=== 环境变量状态 ===")
	
	providers := map[string]string{
		"DASHSCOPE_API_KEY": "阿里云百炼",
		"OPENAI_API_KEY":    "OpenAI",
		"GLM_API_KEY":       "智谱AI",
		"SPARK_API_KEY":     "讯飞星火",
	}
	
	for envVar, providerName := range providers {
		value := os.Getenv(envVar)
		if value != "" {
			fmt.Printf("✅ %s (%s): 已设置 (长度: %d)\n", envVar, providerName, len(value))
		} else {
			fmt.Printf("❌ %s (%s): 未设置\n", envVar, providerName)
		}
	}
}

// ValidateAPIKey 验证API密钥格式（基本检查）
func ValidateAPIKey(provider, apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("API密钥不能为空")
	}
	
	// 基本长度检查
	minLength := 10
	switch provider {
	case "dashscope":
		minLength = 20 // DashScope API密钥通常较长
	case "openai":
		minLength = 40 // OpenAI API密钥格式: sk-...
		if len(apiKey) > 10 && !startsWith(apiKey, "sk-") {
			return fmt.Errorf("OpenAI API密钥应以 'sk-' 开头")
		}
	}
	
	if len(apiKey) < minLength {
		return fmt.Errorf("API密钥长度过短，期望至少 %d 个字符，实际 %d 个字符", minLength, len(apiKey))
	}
	
	return nil
}

// startsWith 检查字符串是否以指定前缀开头
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

// GetRequiredEnvVars 获取所有必需的环境变量列表
func GetRequiredEnvVars() map[string]string {
	return map[string]string{
		"DASHSCOPE_API_KEY": "阿里云百炼API密钥",
		// 其他提供商暂时注释
		// "OPENAI_API_KEY":    "OpenAI API密钥",
		// "GLM_API_KEY":       "智谱AI API密钥",
		// "SPARK_API_KEY":     "讯飞星火API密钥",
	}
}

// CheckAllRequired 检查所有必需的环境变量
func CheckAllRequired() error {
	checker := NewEnvChecker()
	
	for envVar, description := range GetRequiredEnvVars() {
		checker.AddRequired(envVar, description)
	}
	
	return checker.CheckAll()
}
