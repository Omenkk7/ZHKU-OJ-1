package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

func (d *Dao) CreateLabel(ctx context.Context, label *models.Label) (id string, err error) {
	res, err := d.mongo.InsertOne(ctx, labelTable, label)
	if err != nil {
		return "", err
	}
	if res.InsertedID != nil {
		return res.InsertedID.(primitive.ObjectID).Hex(), nil
	}
	return "", nil
}

func (d *Dao) GetLabelList(ctx context.Context, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	lg := utils.GetDefaultLogger()
	items = &utils.RespPageQuery{
		Items: make([]*map[string]interface{}, 0),
	}
	//query := bson.M{}
	opts := utils.BuildMongoOptions(comQuery)
	lg.Println(opts)
	lg.Println(comQuery.Filters)
	err = d.mongo.FindSome(ctx, ProblemTable, comQuery.Filters, items, opts)
	return
}

func (d *Dao) GetOneLabel(ctx context.Context, query interface{}) (label *models.Label, err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("查询条件：", query)
	//按query条件查询
	err = d.mongo.FindOne(ctx, labelTable, query, &label)
	return
}

func (d *Dao) DeleteLabel(ctx context.Context, selector bson.M) (err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("删除：", selector)
	_, err = d.mongo.Remove(ctx, labelTable, selector)
	return
}

func (d *Dao) UpdateLabel(ctx context.Context, selector bson.M, update bson.M) (err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("修改：", selector, update)
	_, err = d.mongo.Upsert(ctx, labelTable, selector, update)
	return
}
