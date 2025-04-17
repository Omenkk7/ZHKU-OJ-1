package service

import (
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// ReadTask 用于自动读取任务
func (s *Service) ReadTask(maxTimeStamp *int64) (results []models.Submit) {
	//1.条件查询，得到任务列表并返回
	bson := bson.M{
		"status": utils.WaitForJudge,
		"stime":  bson.M{"$gt": *maxTimeStamp}, //原来为*maxTimeStamp-1 但如果判题慢，会导致bug重复判题！！
	}
	results, _ = s.dao.GetTaskList(context.Background(), bson)

	//etcd
	//2.按时间顺序排序
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
	//3.为了避免重复读取数据，记录上次数据的最大时间戳，并进行判定
	if len(results) != 0 {
		*maxTimeStamp = results[len(results)-1].Stime
	}
	//返回排好序的结果
	return results
}
