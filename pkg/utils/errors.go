package utils

import "errors"

const (
	UpdateErr      = "更新失败！"
	DeleteAdminErr = "不能删除管理员!"
	DeleteErr      = "删除失败!"

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
	BadJudge         = 3
)

func New(s string) error {
	return errors.New(s)
}
