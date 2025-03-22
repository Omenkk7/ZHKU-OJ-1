package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

func (d *Dao) Submit(ctx context.Context, submit *models.Submit) (submitId string, err error) {
	res, err := d.mongo.InsertOne(ctx, SubmitTable, submit)
	if err != nil {
		return "", err
	}
	if res.InsertedID != nil {
		return res.InsertedID.(primitive.ObjectID).Hex(), nil
	}
	if objectId, ok := res.InsertedID.(primitive.ObjectID); ok {
		return objectId.Hex(), nil
	}
	return "", nil
}

func (d *Dao) UpdateSubmit(ctx context.Context, selector bson.M, update bson.M) (objectID primitive.ObjectID, err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("修改：", selector, update)
	_, err = d.mongo.Upsert(ctx, SubmitTable, selector, update)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return primitive.NilObjectID, nil
}

func (d *Dao) GetOneSubmit(ctx context.Context, query interface{}) (submit *models.Submit, err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("查询条件：", query)
	//按query条件查询
	err = d.mongo.FindOne(ctx, SubmitTable, query, &submit)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return submit, nil
}
