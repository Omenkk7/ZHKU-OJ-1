package judge

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"zhku-oj-server/pkg/utils"
)

type Task struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` //这个用于记录submit_id
	Code      string             `json:"code,omitempty" bson:"code,omitempty"`
	Language  string             `json:"language,omitempty" bson:"language,omitempty"`
	ProblemId string             `json:"problem_id,omitempty" bson:"problem_id,omitempty"`
	UserId    string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

func StartJudge() {
	lg := utils.GetDefaultLogger()
	lg.Info("已启动判题..........")

	taskChan := make(chan Task, 100)
	var maxTimeStamp int64 = 0

	//消费者
	go work(taskChan)
	//循环读取列表
	for {
		//生产者
		results := ReadTask(&maxTimeStamp)
		if len(results) != 0 {
			lg.Infof("新查到%d条数据待判题\n", len(results))
			for _, result := range results {
				lg.Infof("语言: %s, 代码: %s\n", result.Language, result.Code)
				//丢进任务管道
				taskChan <- Task{
					ID:        result.ID,
					Language:  result.Language,
					Code:      result.Code,
					UserId:    result.UserId,
					ProblemId: result.ProblemId,
				}
			}
		}
		time.Sleep(ReadTaskTime * time.Second)
		//TODO 采用管道+消费者--生产者模式??
	}
}

func work(taskChan chan Task) {
	lg := utils.GetDefaultLogger()
	i := 1
	for {
		task := <-taskChan
		//TODO 添加合并模板，判题的逻辑

		lg.Infof("已完成任务%d：语言: %s, 代码: %s\n", i, task.Language, task.Code)
		time.Sleep(5 * time.Second)
		i++
	}
}
