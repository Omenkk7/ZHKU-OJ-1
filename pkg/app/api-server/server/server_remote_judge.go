package server

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"zhku-oj-server/pkg/app/api-server/service"
	"zhku-oj-server/pkg/utils"
)

type RemoteJudge struct {
	svc *service.Service
}
type RemoteTask struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` //这个用于记录submit_id
	BaseTask *BaseTask
}

func (rj *RemoteJudge) RunJudge() error {
	//TODO 其他逻辑
	fmt.Print("2222222222222")
	return nil
}

func (rj *RemoteJudge) NewJudge() *RemoteJudge {
	return &RemoteJudge{
		svc: service.NewService(),
	}
}

func (rj *RemoteJudge) producer(chan interface{}) {

}
func (rj *RemoteJudge) consumer(taskChan chan interface{}) {
	lg := utils.GetDefaultLogger()
	i := 1
	for {
		t := <-taskChan
		task := t.(RemoteTask)

		//TODO 添加合并模板，判题的逻辑

		lg.Infof("已完成任务%d：语言: %s, 代码: %s\n", i, task.BaseTask.Language, task.BaseTask.Code)
		time.Sleep(5 * time.Second)
		i++
	}
}

func (rj *RemoteJudge) work() {

}
