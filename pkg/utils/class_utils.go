package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// GenerateClassCode 生成班级代码
// 格式: 院系缩写+年份+随机4位数字，如CS20240001
func GenerateClassCode(department string) string {
	// 获取院系缩写（取每个单词的首字母）
	words := strings.Fields(department)
	prefix := ""
	for _, word := range words {
		if len(word) > 0 {
			prefix += strings.ToUpper(string(word[0]))
		}
	}

	// 如果没有提取到缩写，使用默认值
	if prefix == "" {
		prefix = "CL" // Class的缩写
	}

	// 如果缩写太长，截取前两个字符
	if len(prefix) > 2 {
		prefix = prefix[:2]
	}

	// 获取当前年份
	year := time.Now().Year()

	// 生成4位随机数
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(9000) + 1000 // 1000-9999之间的随机数

	// 组合成班级代码
	return fmt.Sprintf("%s%d%04d", prefix, year, randomNum)
}
