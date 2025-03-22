package dao

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

func (d *Dao) Submit(ctx context.Context, submit *models.Submit) (id string, err error) {
	res, err := d.mongo.InsertOne(ctx, SubmitTable, submit)
	if err != nil {
		return "", err
	}
	if res.InsertedID != nil {
		return res.InsertedID.(primitive.ObjectID).Hex(), nil
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
