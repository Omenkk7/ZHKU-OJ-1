package utils

import "errors"

const (
	UpdateErr      = "更新失败！"
	DeleteAdminErr = "不能删除管理员!"
	DeleteErr      = "删除失败!"

	QueryErr     = "查询失败!"
	HashErr      = "密码加密失败!"
	CreateErr    = "创建失败!"
	UserExistErr = "用户已存在!"

	ConstructingBsonErr = "构造Bson异常!"
	ServerErr           = "服务器异常!"

	CountOrPasswordNullErr = "账号或密码不能为空！"
	RegisteredErr          = "用户名已被注册!"
	UserNotExistErr        = "用户不存在！"
	CountErr               = "账号错误!"
	PasswordErr            = "密码错误！"
	CountBanedErr          = "账号被封禁！"
	NOPermissionErr        = "没有权限!"

	ConstructingJWTErr = "生成令牌失败！"
	JwtFailErr         = "令牌无效！"
	JwtEmptyErr        = "令牌为空！"

	ProblemOrAuthorCannotBeNull = "题目名称和出题人不能为空!"
	ProblemIsExist              = "题目已存在！"
	ProblemNotExist             = "题目不存在!"

	LabelOrAuthorCannotBeNull = "标签名称和出题人不能为空!"
	LabelIsExist              = "标签已存在！"
	LabelNotExist             = "标签不存在!"

	CodeCannotBeNull = "代码为空！"
	WaitForJudge     = 1
	SuccessJudge     = 2

	// 班级管理相关错误
	ClassNotExistErr     = "班级不存在"
	ClassCodeExistErr    = "班级代码已存在"
	StudentInClassErr    = "学生已在班级中"
	StudentNotInClassErr = "学生不在班级中"
	CourseInClassErr     = "课程已在班级中"
	CourseNotInClassErr  = "课程不在班级中"
	RequestProcessedErr  = "申请已处理"
	ClassHasStudentsErr  = "班级中还有学生，无法删除"

	// 班级成员权限相关错误常量
	NoPermissionErrMsg   = "没有权限执行此操作"
	MemberExistErrMsg    = "成员已存在"
	MemberNotExistErrMsg = "成员不存在"
)

func New(s string) error {
	return errors.New(s)
}

// 预创建的错误变量
var (

	// 班级成员权限相关错误
	NoPermissionErr   = errors.New(NoPermissionErrMsg)
	MemberExistErr    = errors.New(MemberExistErrMsg)
	MemberNotExistErr = errors.New(MemberNotExistErrMsg)

	ErrCourseNotFound        = errors.New("课程不存在")
	ErrCourseCodeExists      = errors.New("课程代码已存在")
	ErrCourseHasClass        = errors.New("课程已绑定班级，无法删除")
	ErrCourseAlreadyArchived = errors.New("课程已归档")
	ErrCourseNotOpen         = errors.New("课程已结课，不能进行此操作")
	ErrMemberAlreadyExists   = errors.New("成员已存在")
	ErrMemberNotExists       = errors.New("成员不存在")
	ErrJoinRequestExists     = errors.New("已提交过申请，请等待审核")
	ErrJoinRequestProcessed  = errors.New("申请已被处理")
)
