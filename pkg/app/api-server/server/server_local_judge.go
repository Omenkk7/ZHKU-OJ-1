package server

import (
	"fmt"
	"time"
	"zhku-oj-server/pkg/app/api-server/service"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

type LocalJudge struct {
	svc *service.Service
}

func (lj *LocalJudge) RunJudge() error {
	lj.work()
	//TODO error返回什么？
	return nil
}

func (lj *LocalJudge) NewJudge() *LocalJudge {
	return &LocalJudge{
		svc: service.NewService(),
	}
}

func (lj *LocalJudge) work() {
	lg := utils.GetDefaultLogger()
	lg.Info("已启动本地判题..........")

	taskChan := make(chan interface{}, 100)

	//生产者和消费者，开启两个协程
	go lj.producer(taskChan)
	go lj.consumer(taskChan)
}

func (lj *LocalJudge) producer(taskChan chan interface{}) {
	lg := utils.GetDefaultLogger()
	var maxTimeStamp int64 = 0
	//循环读取列表
	for {
		//生产者
		results := lj.svc.ReadTask(&maxTimeStamp)
		if len(results) != 0 {
			lg.Infof("新查到%d条数据待判题\n", len(results))
			for _, result := range results {
				lg.Infof("语言: %s, 代码: %s\n", result.Language, result.Code)
				//丢进任务管道
				taskChan <- models.LocalTask{
					ID: result.ID,
					BaseTask: models.BaseTask{
						Language:  result.Language,
						Code:      result.Code,
						UserId:    result.UserId,
						ProblemId: result.ProblemId,
					},
				}
			}
		}
		time.Sleep(utils.ReadTaskTime * time.Second)
	}
}
func (lj *LocalJudge) consumer(taskChan chan interface{}) {
	lg := utils.GetDefaultLogger()
	i := 1
	for {
		//用于隔开任务日志输出，便于测试观察
		fmt.Println()
		fmt.Println()

		//从阻塞队列中拿出任务
		t := <-taskChan
		task := t.(models.LocalTask)

		//合并模板，拿到合并后的完整代码
		c, err := lj.svc.MergeTemplate(task)
		if err != nil || c == nil {
			continue
		}
		task.Code = c.(string)

		//TODO 完善InvokeSandbox
		//接下来开启http调用go-oj进行判题
		lj.svc.InvokeSandbox(task.Language, task.Code)

		//TODO 判断结果，把结果写回mongo

		lg.Infof("已完成任务%d：语言: %s, 代码: %s\n", i, task.Language, task.Code)
		time.Sleep(5 * time.Second)
		i++
	}
}
