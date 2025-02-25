package judge

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/net/context"
	"log"
	"time"
	"zhku-oj-server/pkg/utils"
)

type Submit struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Code       string             `json:"code,omitempty" bson:"code,omitempty"`
	Evaluation interface{}        `json:"evaluation,omitempty" bson:"evaluation,omitempty"`
	Language   string             `json:"language,omitempty" bson:"language,omitempty"`
	UserId     string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	ProblemId  string             `json:"problem_id,omitempty" bson:"problem_id,omitempty"`
	Status     int32              `json:"status,omitempty" bson:"status,omitempty"`
	Stime      int64              `json:"stime,omitempty" bson:"stime,omitempty"`
}

// 用于自动读取任务
func ReadTask() {
	lg := utils.GetDefaultLogger()
	lg.Info("读任务......")

	//TODO 直接用原有的Dao报空指针，找了一个多小时都不知道是什么的问题，先直接用原生mongo实现该功能

	// 1. 连接 MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())

	// 2. 获取集合
	collection := client.Database("zkoj").Collection("submits")

	// 3. 查询所有文档
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)

	// 4. 解码结果
	var results []Submit
	if err = cursor.All(ctx, &results); err != nil {
		log.Fatal(err)
	}

	// 5. 处理结果
	lg.Infof("找到了%d条数据\n", len(results))
	lg.Infof("------------排序前---------------")
	for _, result := range results {
		lg.Infof("语言: %s, 代码: %s\n,提交时间：%s\n", result.Language, result.Code, result.Stime)
	}

	//6.待判题的代码，按时间顺序排序，排好序再进行判题
	for i := 0; i < len(results); i++ {
		min := results[i].Stime
		minIndex := i
		for j := i; j < len(results); j++ {
			if min > results[j].Stime {
				min = results[j].Stime
				minIndex = j
			}
		}
		temp := results[minIndex]
		results[minIndex] = results[i]
		results[i] = temp
	}

	lg.Infof("------------排序后---------------")
	for _, result := range results {
		lg.Infof("语言: %s, 代码: %s\n,提交时间：%s\n", result.Language, result.Code, result.Stime)
	}
}
