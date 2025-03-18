package dao

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"zhku-oj-server/pkg/models"
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
