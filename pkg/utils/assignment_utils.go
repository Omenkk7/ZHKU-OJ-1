package utils

import (
	"fmt"
	"math/rand"
	"time"
)

// 作业类型常量
const (
	AssignmentTypeNormal = 0 // 普通作业
	AssignmentTypeOJ     = 1 // OJ作业
)

// 作业状态常量
const (
	AssignmentStatusArchived = 0 // 已归档
	AssignmentStatusActive   = 1 // 进行中
)

// GenerateAssignmentCode 生成作业代码
// 格式：ASG + 年月 + 4位随机数字
// 例如：ASG202504XXXX
func GenerateAssignmentCode() string {
	// 初始化随机数生成器
	rand.Seed(time.Now().UnixNano())

	// 获取当前年月
	now := time.Now()
	yearMonth := now.Format("200601") // 格式为YYYYMM

	// 生成4位随机数
	randomNum := rand.Intn(10000)

	// 组合成作业代码
	return fmt.Sprintf("ASG%s%04d", yearMonth, randomNum)
}
