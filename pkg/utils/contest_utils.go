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
)

// 竞赛参赛者状态常量
const (
	ContestParticipantStatusPending  = 0 // 待审核
	ContestParticipantStatusApproved = 1 // 已通过
	ContestParticipantStatusRejected = 2 // 已拒绝

	// 角色常量
	ContestRoleAdmin = 1 // 管理员

	// 竞赛状态常量
	ContestStatusDeleted  = 0 // 已删除
	ContestStatusNotStart = 1 // 未开始
	ContestStatusRunning  = 2 // 进行中
	ContestStatusEnded    = 3 // 已结束
	ContestStatusArchived = 4 // 已归档

	// 参赛规则常量
	ContestAccessPrivate = 2 // 私有
	ContestAccessPublic  = 1 // 公开
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
