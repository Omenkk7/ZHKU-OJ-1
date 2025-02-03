/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 19:17
@Name: dao.go
@Description:
*/

package dao

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
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

// GenerateUpdateBson 接收结构体，生成包含非空字段的 $set BSON
func (d *Dao) GenerateUpdateBson(model interface{}) (bson.M, error) {
	updateFields := bson.M{}
	val := reflect.ValueOf(model)

	// 处理指针类型的结构体
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// 确保传入的是结构体
	if val.Kind() != reflect.Struct {
		return nil, errors.New("input must be a struct")
	}

	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// 解析 bson 标签
		bsonTag := field.Tag.Get("bson")
		bsonField := d.parseBsonTag(bsonTag, field.Name)

		//不把ID转为bson,因为他是主键，修改会报错！！！
		if bsonField == mongoID {
			continue
		}

		// 处理指针类型字段
		if fieldVal.Kind() == reflect.Ptr && !fieldVal.IsNil() {
			updateFields[bsonField] = fieldVal.Elem().Interface()
			continue
		}

		// 处理值类型字段（非零值时包含）
		if fieldVal.Kind() != reflect.Ptr && !fieldVal.IsZero() {
			updateFields[bsonField] = fieldVal.Interface()
		}
	}

	if len(updateFields) == 0 {
		return nil, errors.New("no fields to update")
	}

	return bson.M{mongoSet: updateFields}, nil
}

// parseBsonTag 解析 bson 标签，返回字段名称
func (d *Dao) parseBsonTag(tag string, defaultName string) string {
	if tag == "" {
		return strings.ToLower(defaultName)
	}

	// 分割标签选项，取第一个部分作为字段名
	parts := strings.SplitN(tag, ",", 2)
	if parts[0] == "" {
		return strings.ToLower(defaultName)
	}

	return parts[0]
}
