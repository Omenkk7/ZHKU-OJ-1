/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 21:13
@Name: dao_user.go
@Description:
*/

package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

func (d *Dao) CreateUser(ctx context.Context, user *models.User) (id string, err error) {
	res, err := d.mongo.InsertOne(ctx, userTable, user)
	if err != nil {
		return "", err
	}
	if res.InsertedID != nil {
		return res.InsertedID.(primitive.ObjectID).Hex(), nil
	}
	return "", nil
}

func (d *Dao) GetUserList(ctx context.Context, comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	lg := utils.GetDefaultLogger()
	items = &utils.RespPageQuery{
		Items: make([]*map[string]interface{}, 0),
	}
	//query := bson.M{}
	opts := utils.BuildMongoOptions(comQuery)
	lg.Println(opts)
	lg.Println(comQuery.Filters)
	err = d.mongo.FindSome(ctx, userTable, comQuery.Filters, items, opts)
	return
}

func (d *Dao) GetOneUser(ctx context.Context, query interface{}) (user *models.User, err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("查询条件：", query)
	//按query条件查询
	err = d.mongo.FindOne(ctx, userTable, query, &user)
	return user, err
}

func (d *Dao) DeleteUser(ctx context.Context, selector bson.M) (err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("删除：", selector)
	_, err = d.mongo.Remove(ctx, userTable, selector)
	return err
}

func (d *Dao) UpdateUser(ctx context.Context, selector bson.M, update bson.D) (err error) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Println("修改：", selector, update)
	/*d.mongo.UpdateOne(ctx, userTable, selector, update)*/
	collection := mongo.Collection{}
	collection.UpdateOne(ctx, userTable, update)
	return
}
