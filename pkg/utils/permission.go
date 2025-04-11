package utils

import "errors"

// 定义角色常量
const (
	ClassRoleAdmin = 1 // 管理员
	RoleTeacher    = 2 // 教师
	RoleAssistant  = 3 // 助教
	RoleStudent    = 4 // 学生
)

// 权限相关错误
var (
	ErrNoPermission = errors.New("没有权限执行此操作")
)

// CheckRolePermission 检查用户角色是否有权限
func CheckRolePermission(userRoles []int, requiredRoles []int) bool {
	if len(userRoles) == 0 || len(requiredRoles) == 0 {
		return false
	}

	// 检查用户角色是否包含任一所需角色
	for _, userRole := range userRoles {
		for _, requiredRole := range requiredRoles {
			if userRole == requiredRole {
				return true
			}
		}
	}

	return false
}
