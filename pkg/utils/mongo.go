/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/24 01:58
@Name: mongo.go
@Description:
*/

package utils

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDB struct {
	client  *mongo.Client
	retry   int
	db      string
	timeout time.Duration
	cancel  context.CancelFunc
}

func NewMongoDB(uri, database string, timeoutOpt ...time.Duration) *MongoDB {
	client, cancel := initDB(uri)

	// 设置CRUD超时时间
	var timeout time.Duration
	if len(timeoutOpt) > 0 {
		timeout = timeoutOpt[0]
	} else {
		timeout = time.Second * 30
	}

	return &MongoDB{
		client:  client,
		retry:   3,
		db:      database,
		timeout: timeout,
		cancel:  cancel,
	}
}

func initDB(uri string) (*mongo.Client, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	client, err := mongo.Connect(
		ctx,
		options.Client().
			ApplyURI(uri).
			SetMaxPoolSize(100).
			SetConnectTimeout(10*time.Second).
			SetMaxConnIdleTime(10*time.Second),
	)
	if err != nil {
		panic(fmt.Errorf("init mongo failed, uri: %s", uri))
	}
	return client, cancel
}

func (m *MongoDB) Close(ctx context.Context) error {
	m.cancel()
	return m.client.Disconnect(ctx)
}

func BuildMongoOptions(comQuery *CommonQuery) *options.FindOptions {
	opt := &options.FindOptions{}
	if comQuery.PageNum > 0 && comQuery.PageSize > 0 {
		nextPage := int64(1)
		opt.SetSkip((comQuery.PageNum - nextPage) * comQuery.PageSize)
		opt.SetLimit(comQuery.PageSize)
	}
	sort := bson.M{"_id": -1}
	if comQuery.Sort != "" {
		sort = bson.M{comQuery.Sort: comQuery.Direction}
	}
	opt.SetSort(sort)
	return opt
}

func (m *MongoDB) FindOne(ctx context.Context, table string, query interface{}, result interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	err = m.client.Database(m.db).Collection(table).FindOne(ctx, query).Decode(result)
	if err != nil {
		logger := GetLogger(ctx)
		logger.Error("MongoFindOneErr", zap.Error(err), zap.String("table", table), zap.Any("query", query))
		return
	}
	return
}

func (m *MongoDB) Find(ctx context.Context, table string, query interface{}, result interface{}, opt ...*options.FindOptions) (err error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	cursor, err := m.client.Database(m.db).Collection(table).Find(ctx, query, opt...)
	if err != nil {
		logger := GetLogger(ctx)
		logger.Error("MongoFindErr", zap.Error(err), zap.String("table", table), zap.Any("query", query))
		return
	}
	err = cursor.All(ctx, result)
	return
}

func (m *MongoDB) FindSome(ctx context.Context, table string, query interface{}, result *RespPageQuery, opt ...*options.FindOptions) (err error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	lg := GetDefaultLogger()
	lg.Println("MongoFindSome query: ", query)
	cursor, err := m.client.Database(m.db).Collection(table).Find(ctx, query, opt...)
	if err != nil {
		lg.Error("MongoFindErr", err, table, query)
		return
	}
	err = cursor.All(ctx, &result.Items)
	total, err := m.client.Database(m.db).Collection(table).CountDocuments(ctx, query)
	if err != nil {
		lg.Errorf("MongoFindTotalErr err:%v, table:%s, query:%v", err, table, query)
		return
	}
	result.Total = total
	return
}

func (m *MongoDB) InsertOne(ctx context.Context, table string, doc interface{}) (res *mongo.InsertOneResult, err error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	res, err = m.client.Database(m.db).Collection(table).InsertOne(ctx, doc)
	return
}

func (m *MongoDB) Remove(ctx context.Context, table string, selector interface{}) (res *mongo.DeleteResult, err error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	res, err = m.client.Database(m.db).Collection(table).DeleteOne(ctx, selector)
	return
}

func (m *MongoDB) RemoveMany(ctx context.Context, table string, selector interface{}) (res *mongo.DeleteResult, err error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	res, err = m.client.Database(m.db).Collection(table).DeleteMany(ctx, selector)
	return
}

func (m *MongoDB) Upsert(ctx context.Context, table string, selector, doc interface{}) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	// 设置Upsert选项
	res, err := m.client.Database(m.db).Collection(table).UpdateOne(ctx, selector, doc, options.Update().SetUpsert(true))
	if err != nil {
		logger := GetLogger(ctx)
		logger.Error("MongoUpsertErr", zap.Error(err), zap.String("table", table), zap.Any("selector", selector))
		return false, err
	}
	return res.UpsertedCount > 0, nil
}

func (m *MongoDB) Update(ctx context.Context, table string, selector, update interface{}) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	res, err := m.client.Database(m.db).Collection(table).UpdateMany(ctx, selector, update)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

// UpdateOne 更新单个文档
func (m *MongoDB) UpdateOne(ctx context.Context, table string, selector interface{}, update interface{}, opts ...*options.UpdateOptions) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	res, err := m.client.Database(m.db).Collection(table).UpdateOne(ctx, selector, update, opts...)
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

func (m *MongoDB) BulkInsert(ctx context.Context, table string, ordered bool, docs ...interface{}) (ok bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	operations := make([]mongo.WriteModel, 0)
	for _, doc := range docs {
		op := mongo.NewInsertOneModel().SetDocument(doc)
		operations = append(operations, op)
	}
	bulkOption := options.BulkWriteOptions{}
	bulkOption.SetOrdered(ordered)
	res, err := m.client.Database(m.db).Collection(table).BulkWrite(ctx, operations, &bulkOption)
	if err != nil {
		return
	}
	return res.InsertedCount > 0, nil
}

func (m *MongoDB) Aggregate(ctx context.Context, table string, pipeline interface{}, result interface{}) (err error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	cursor, err := m.client.Database(m.db).Collection(table).Aggregate(ctx, pipeline)
	if err != nil {
		return
	}
	defer func() {
		_ = cursor.Close(ctx)
	}()
	err = cursor.All(ctx, &result)
	return
}

// Count 计算符合条件的文档数量
func (m *MongoDB) Count(ctx context.Context, collection string, filter interface{}) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, m.timeout)
	defer cancel()
	return m.client.Database(m.db).Collection(collection).CountDocuments(ctx, filter)
}
