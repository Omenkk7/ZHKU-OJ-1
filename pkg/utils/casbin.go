package utils

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"github.com/gin-gonic/gin"
	"log"
	"path/filepath"

	"sync"
)

// 角色常量定义
const (
	RBACRoleSuperAdmin = "super_admin" // 超级管理员
	RBACRoleAdmin      = "admin"       // 管理员
	RBACRoleTeacher    = "teacher"     // 教师
	RBACRoleAssistant  = "assistant"   // 助教
	RBACRoleStudent    = "student"     // 学生
)

// 操作常量定义
const (
	ActCreate  = "create"  // 创建
	ActRead    = "read"    // 读取
	ActUpdate  = "update"  // 更新
	ActDelete  = "delete"  // 删除
	ActArchive = "archive" // 归档
	ActAddMem  = "add_mem" // 添加成员
	ActRemMem  = "rem_mem" // 移除成员
	ActReview  = "review"  // 审核
	ActAll     = "*"       // 所有操作
)

// 资源常量定义
const (
	ObjClass        = "class"         // 班级
	ObjClassMember  = "class_member"  // 班级成员
	ObjClassStudent = "class_student" // 班级学生
	ObjClassCourse  = "class_course"  // 班级课程
	ObjJoinRequest  = "join_request"  // 加入申请
	ObjAll          = "*"             // 所有对象
	ObjUser         = "user"          // 用户
	ObjAdmin        = "admin"         // 管理员
)

var (
	enforcer *casbin.Enforcer
	once     sync.Once
)

// InitCasbin 初始化Casbin
func InitCasbin(configDir string) (*casbin.Enforcer, error) {
	var err error
	once.Do(func() {
		// 加载模型
		modelPath := filepath.Join(configDir, "model.conf") // 注意这里的文件名
		m, err := model.NewModelFromFile(modelPath)
		if err != nil {
			log.Printf("加载模型文件失败: %v", err)
			return
		}

		// 使用文件适配器
		adapter := fileadapter.NewAdapter(filepath.Join(configDir, "policy.csv"))

		// 创建enforcer
		enforcer, err = casbin.NewEnforcer(m, adapter)
		if err != nil {
			log.Printf("创建enforcer失败: %v", err)
			return
		}

		// 加载策略
		err = enforcer.LoadPolicy()
		if err != nil {
			log.Printf("加载策略失败: %v", err)
			return
		}

		// 添加默认策略
		addDefaultPolicies(enforcer)

		log.Println("Casbin初始化成功")
	})

	if err != nil {
		return nil, err
	}

	return enforcer, nil
}

// GetEnforcer 获取Casbin enforcer实例
func GetEnforcer() *casbin.Enforcer {
	return enforcer
}

// addDefaultPolicies 添加默认策略
func addDefaultPolicies(e *casbin.Enforcer) {
	// 清除所有现有策略
	e.ClearPolicy()

	// 超级管理员拥有所有权限
	e.AddPolicy(RBACRoleSuperAdmin, ObjAll, ActAll)

	// 管理员拥有所有权限，但不能管理其他管理员
	e.AddPolicy(RBACRoleAdmin, ObjAll, ActAll)
	// 管理员不能创建管理员
	e.AddPolicy(RBACRoleAdmin, ObjAdmin, ActCreate)
	// 管理员不能删除管理员
	e.AddPolicy(RBACRoleAdmin, ObjAdmin, ActDelete)

	// 教师权限
	e.AddPolicy(RBACRoleTeacher, ObjClass, ActRead)
	e.AddPolicy(RBACRoleTeacher, ObjClass, ActCreate)
	e.AddPolicy(RBACRoleTeacher, ObjClass, ActUpdate)
	e.AddPolicy(RBACRoleTeacher, ObjClassMember, ActRead)
	e.AddPolicy(RBACRoleTeacher, ObjClassStudent, ActRead)
	e.AddPolicy(RBACRoleTeacher, ObjClassStudent, ActAddMem)
	e.AddPolicy(RBACRoleTeacher, ObjClassStudent, ActRemMem)
	e.AddPolicy(RBACRoleTeacher, ObjClassCourse, ActRead)
	e.AddPolicy(RBACRoleTeacher, ObjClassCourse, ActUpdate)
	e.AddPolicy(RBACRoleTeacher, ObjJoinRequest, ActRead)
	e.AddPolicy(RBACRoleTeacher, ObjJoinRequest, ActReview)

	// 助教权限
	e.AddPolicy(RBACRoleAssistant, ObjClass, ActRead)
	e.AddPolicy(RBACRoleAssistant, ObjClassMember, ActRead)
	e.AddPolicy(RBACRoleAssistant, ObjClassStudent, ActRead)
	e.AddPolicy(RBACRoleAssistant, ObjClassStudent, ActAddMem)
	e.AddPolicy(RBACRoleAssistant, ObjJoinRequest, ActRead)
	e.AddPolicy(RBACRoleAssistant, ObjJoinRequest, ActReview)

	// 学生权限
	e.AddPolicy(RBACRoleStudent, ObjClass, ActRead)
	e.AddPolicy(RBACRoleStudent, ObjClassMember, ActRead)

	// 保存策略
	e.SavePolicy()
}

// CheckPermission 检查权限
func CheckPermission(sub string, obj string, act string) bool {
	if enforcer == nil {
		return false
	}

	ok, _ := enforcer.Enforce(sub, obj, act)
	return ok
}

// AddRoleForUser 为用户添加角色
func AddRoleForUser(user string, role string) bool {
	if enforcer == nil {
		return false
	}

	ok, _ := enforcer.AddGroupingPolicy(user, role)
	return ok
}

// GetRolesForUser 获取用户的所有角色
func GetRolesForUser(user string) []string {
	if enforcer == nil {
		return []string{}
	}

	roles, _ := enforcer.GetRolesForUser(user)
	return roles
}

// CasbinMiddleware Casbin中间件
func CasbinMiddleware(obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文中获取用户信息
		contextUser, exists := c.Get("contextUser")
		if !exists {
			Unauthorized(c, New("未登录或会话已过期"))
			c.Abort()
			return
		}

		user, ok := contextUser.(*ContextUser)
		if !ok {
			Unauthorized(c, New("用户信息获取失败"))
			c.Abort()
			return
		}

		// 将用户角色转换为Casbin角色
		role := GetRoleName(user.Role)

		// 检查权限
		if !CheckPermission(role, obj, act) {
			Forbidden(c, ErrNoPermission)
			c.Abort()
			return
		}

		c.Next()
	}
}

// 将数字角色ID转换为角色名称
func GetRoleName(roleID int) string {
	switch roleID {
	/*case RoleSuperAdmin:return RBACRoleSuperAdmin*/ //报错
	case ClassRoleAdmin:
		return RBACRoleAdmin
	case RoleTeacher:
		return RBACRoleTeacher
	case RoleAssistant:
		return RBACRoleAssistant
	case RoleStudent:
		return RBACRoleStudent
	default:
		return "unknown"
	}
}

// CheckClassPermission 检查用户对特定班级的权限
func CheckClassPermission(ctx context.Context, userRoles []int, classID string, act string) bool {
	// 如果用户是管理员，直接返回true
	for _, role := range userRoles {
		if role == ClassRoleAdmin {
			return true
		}
	}

	// 将用户角色转换为Casbin角色
	casbinRoles := make([]string, 0, len(userRoles))
	for _, role := range userRoles {
		casbinRoles = append(casbinRoles, GetRoleName(role))
	}

	// 检查每个角色的权限
	for _, role := range casbinRoles {
		// 检查对班级的权限
		if CheckPermission(role, ObjClass, act) {
			return true
		}
	}

	return false
}

// GetUserRoleInClass 获取用户在班级中的角色
func GetUserRoleInClass(userID string, classID string, dao interface{}) ([]int, error) {
	// 这里需要调用dao层的方法获取用户在班级中的角色
	// 由于我们没有直接访问dao的方法，这里提供一个接口
	type ClassRoleChecker interface {
		CheckUserClassRole(ctx context.Context, classID string, userID string) ([]int, error)
	}

	if checker, ok := dao.(ClassRoleChecker); ok {
		return checker.CheckUserClassRole(context.Background(), classID, userID)
	}

	return nil, fmt.Errorf("无法检查用户角色")
}

// ReloadCasbinPolicy 重新加载Casbin策略
func ReloadCasbinPolicy() error {
	if enforcer == nil {
		return fmt.Errorf("enforcer未初始化")
	}

	// 清除现有策略
	enforcer.ClearPolicy()

	// 添加默认策略
	addDefaultPolicies(enforcer)

	// 保存策略
	return enforcer.SavePolicy()
}
