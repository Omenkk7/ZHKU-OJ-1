package server

import (
	"log"
	"os"
	"zhku-oj-server/pkg/app/api-server/service"
	"zhku-oj-server/pkg/utils"
)

type BaseTask struct {
	Code      string `json:"code,omitempty" bson:"code,omitempty"`
	Language  string `json:"language,omitempty" bson:"language,omitempty"`
	ProblemId string `json:"problem_id,omitempty" bson:"problem_id,omitempty"`
	UserId    string `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

type Judge struct {
	svc    *service.Service
	bt     *BaseTask
	lj     *LocalJudge
	rj     *RemoteJudge
	choice int
}

func (j Judge) Work() {
	switch j.choice {
	case utils.LocalJudge:
		j.lj.work()
	case utils.RemoteJudge:
		j.rj.work()
	}
}

func NewJudge() *Judge {
	//通过读取配置文件，决定启动哪个判题
	configPath := "conf/config.yaml"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 如果当前目录没有conf，尝试上级目录
		configPath = "../conf/config.yaml"
	}

	cfg, err := utils.LoadConfig(configPath)
	if err != nil {
		log.Printf("加载配置文件失败: %v", err)
		// 返回默认配置的Judge
		return &Judge{
			choice: utils.LocalJudge,
			svc:    service.NewService(),
		}
	}
	judgeCfg := cfg.GetJudgeConfig()

	switch judgeCfg.Type {
	case utils.LocalJudgeCfg:
		return &Judge{
			choice: utils.LocalJudge,
			svc:    service.NewService(),
		}
	case utils.RemoteJudgeCfg:
		return &Judge{
			choice: utils.RemoteJudge,
			svc:    service.NewService(),
		}
	}
	/*
		。。。。启动其他的服务
	*/
	return nil
}

func (j Judge) RunJudge() error {
	switch j.choice {
	case utils.LocalJudge:
		j.lj = j.lj.NewJudge()
		return j.lj.RunJudge()
	case utils.RemoteJudge:
		j.rj = j.rj.NewJudge()
		return j.rj.RunJudge()
	}
	/*
		。。。。其他的服务
	*/
	return nil
}

func (j Judge) Producer(taskChan chan interface{}) {
	switch j.choice {
	case utils.LocalJudge:
		j.lj.producer(taskChan)
		return
	case utils.RemoteJudge:
		j.rj.producer(taskChan)
		return
	}
}

func (j Judge) Consumer(taskChan chan interface{}) {
	switch j.choice {
	case utils.LocalJudge:
		j.lj.consumer(taskChan)
		return
	case utils.RemoteJudge:
		j.rj.consumer(taskChan)
		return
	}
}
