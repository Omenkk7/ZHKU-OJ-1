package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"zhku-oj-server/pkg/models"
)

// GetTaskList 列表查询
func (d *Dao) GetTaskList(ctx context.Context, query bson.M) (results []models.Submit, err error) {
	// 执行查询，获取所有数据
	err = d.mongo.Find(ctx, SubmitTable, query, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
