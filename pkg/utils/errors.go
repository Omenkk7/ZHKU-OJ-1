/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/26 22:45
@Name: errors.go
@Description:
*/

package utils

import "errors"

var (
	ErrUsernameExisted = errors.New("用户名已被使用")
	ErrUserNotFound    = errors.New("用户名不存在")
)
