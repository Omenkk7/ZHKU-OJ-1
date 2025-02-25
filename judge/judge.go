package judge

import (
	"time"
	"zhku-oj-server/pkg/utils"
)

func StartJudge() {
	lg := utils.GetDefaultLogger()

	//循环读取列表
	for {
		results := ReadTask()
		lg.Infof("共查到%d条数据待判题\n", len(results))
		for _, result := range results {
			lg.Infof("语言: %s, 代码: %s\n,提交时间：%s\n", result.Language, result.Code, result.Stime)
		}
		time.Sleep(ReadTaskTime * time.Second)
		//TODO 采用管道+消费者--生产者模式??

	}
}
