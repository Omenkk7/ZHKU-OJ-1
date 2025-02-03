package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"zhku-oj-server/pkg/utils"
)

// 通用CURD  要传结构体、表名
// 写成接口的模式

// CreateOne 插入1条数据
func (d *Dao) CreateOne(ctx context.Context, tableName string, model interface{}) (id string, err error) {
	res, err := d.mongo.InsertOne(ctx, tableName, model)
	if err != nil {
		return "", err
	}
	if res.InsertedID != nil {
		return res.InsertedID.(primitive.ObjectID).Hex(), nil
	}
	return "", nil
}

// GetSome 列表查询
func (d *Dao) GetSome(ctx context.Context, tableName string, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	lg := utils.GetDefaultLogger()
	items = &utils.RespPageQuery{
		Items: make([]*map[string]interface{}, 0),
	}
	//query := bson.M{}
	opts := utils.BuildMongoOptions(comQuery)
	lg.Println(opts)
	lg.Println(comQuery.Filters)
	err = d.mongo.FindSome(ctx, tableName, comQuery.Filters, items, opts)
	return
}

// GetOne 查询一条
func (d *Dao) GetOne(ctx context.Context, tableName string, model interface{}, query interface{}) (interface{}, error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("查询条件：", query)
	//按query条件查询
	err := d.mongo.FindOne(ctx, tableName, query, &model)
	return model, err
}

// DeleteOne 删除一条
func (d *Dao) DeleteOne(ctx context.Context, tableName string, selector bson.M) (err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("删除：", selector)
	_, err = d.mongo.Remove(ctx, tableName, selector)
	return
}

// Update 更新
func (d *Dao) Update(ctx context.Context, tableName string, selector bson.M, update bson.M) (err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("修改：", selector, update)
	_, err = d.mongo.Upsert(ctx, tableName, selector, update)
	return
}
