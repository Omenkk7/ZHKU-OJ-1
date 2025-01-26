/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 00:18
@Name: interface.go
@Description:
*/

package common

// ServerManagerInterface router管理器接口
type ServerManagerInterface interface {
	RunManager()
	Shutdown()
}

type ServerInterface interface {
	RegisterRoutes()
	Run() error
}
