package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

func (d *Dao) CreateProblem(ctx context.Context, problem *models.Problem) (id string, err error) {
	res, err := d.mongo.InsertOne(ctx, ProblemTable, problem)
	if err != nil {
		return "", err
	}
	if res.InsertedID != nil {
		return res.InsertedID.(primitive.ObjectID).Hex(), nil
	}
	return "", nil
}

func (d *Dao) GetProblemList(ctx context.Context, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
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

func (d *Dao) GetOneProblem(ctx context.Context, query interface{}) (problem *models.Problem, err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("查询条件：", query)
	//按query条件查询
	err = d.mongo.FindOne(ctx, ProblemTable, query, &problem)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return
}

func (d *Dao) DeleteProblem(ctx context.Context, selector bson.M) (objectID primitive.ObjectID, err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("删除：", selector)
	_, err = d.mongo.Remove(ctx, ProblemTable, selector)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return primitive.NilObjectID, nil
}

func (d *Dao) UpdateProblem(ctx context.Context, selector bson.M, update bson.M) (objectID primitive.ObjectID, err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("修改：", selector, update)
	_, err = d.mongo.Upsert(ctx, ProblemTable, selector, update)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return primitive.NilObjectID, nil
}
