/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 19:17
@Name: dao.go
@Description:
*/

package dao

import (
	"context"
	"zhku-oj-server/conf"
	"zhku-oj-server/pkg/utils"
)

type Dao struct {
	mongo *utils.MongoDB
}

func NewDao() *Dao {
	mongoCfg := conf.Config.Mongo
	return &Dao{
		mongo: utils.NewMongoDB(mongoCfg.Uri, mongoCfg.DbName),
	}
}
func (d *Dao) Close(ctx context.Context) {
	_ = d.mongo.Close(ctx)
}
