package errors

import "errors"

const (
	UpdateFailed     = "更新失败！"
	CanotDeleteAdmin = "不能删除管理员!"
	DeleteFail       = "删除失败!"

	ConstructingBsonException = "构造Bson异常!"
	ServerAbnormal            = "服务器异常!"

	CountOrPasswordCannotBeNull = "账号或密码不能为空！"
	CountHasBeenRegistered      = "用户名已被注册!"
	UserNotExist                = "用户不存在！"
	CountException              = "账号错误!"
	PasswordException           = "密码错误！"
	CountBaned                  = "账号被封禁！"
	NOPerrmission               = "没有权限!"

	ConstructingJWTException = "生成令牌失败！"
	JwtFail                  = "令牌无效！"
	JwtEmpty                 = "令牌为空！"
)

func New(s string) error {
	return errors.New(s)
}
