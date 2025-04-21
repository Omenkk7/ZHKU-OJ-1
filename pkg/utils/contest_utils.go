package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// 竞赛类型常量
const (
	ContestTypeWeekly    = 1 // 周赛
	ContestTypeMonthly   = 2 // 月赛
	ContestTypeChallenge = 3 // 挑战赛
	// 可以添加更多类型...
)

// GetContestTypePrefix 根据竞赛类型获取前缀
func GetContestTypePrefix(contestType int) string {
	switch contestType {
	case ContestTypeWeekly:
		return "WS" // Weekly Series
	case ContestTypeMonthly:
		return "MS" // Monthly Series
	case ContestTypeChallenge:
		return "CS" // Challenge Series
	default:
		return "OT" // Other
	}
}

// GenerateContestCode 生成竞赛代码
// 格式: 竞赛类型缩写+年份+随机4位数字，如WS20240001
func GenerateContestCode(contestType int) string {
	// 获取竞赛类型前缀
	prefix := GetContestTypePrefix(contestType)

	// 获取当前年份
	year := time.Now().Year()

	// 生成4位随机数
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(9000) + 1000 // 1000-9999之间的随机数

	// 组合成竞赛代码
	return fmt.Sprintf("%s%d%04d", prefix, year, randomNum)
}
