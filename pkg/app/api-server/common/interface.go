/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 00:18
@Name: interface.go
@Description:
*/

package common

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"zhku-oj-server/pkg/app/api-server/server"
	"zhku-oj-server/pkg/utils"
)

// ServerManagerInterface router管理器接口
type ServerManagerInterface interface {
	RunManager()
	Shutdown()
}

type ServerInterface interface {
	RegisterRoutes()
	Run() error
}

// DaoInterface 接口
type DaoInterface interface {
	CreateOne(ctx context.Context, tableName string, model interface{}) (id string, err error)
	GetOne(ctx context.Context, tableName string, model interface{}, query interface{}) (interface{}, error)
	GetSome(ctx context.Context, tableName string, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error)
	Update(ctx context.Context, tableName string, selector bson.M, update bson.M) (err error)
	DeleteOne(ctx context.Context, tableName string, selector bson.M) (err error)
}

// Judge接口
type JudgeInterface interface {
	RunJudge() error
	NewJudge() server.Judge

	// TODO 留下work()或producer()+consumer()即可
	Work()
	Producer(obj interface{})
	Consumer(obj interface{})
}
type Task interface {
}
