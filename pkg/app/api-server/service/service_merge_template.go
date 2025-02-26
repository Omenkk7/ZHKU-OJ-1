package service

import (
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

func (s *Service) MergeTemplate(choice int, task interface{}) (code interface{}) {
	lg := utils.GetDefaultLogger()
	switch choice {
	case utils.LocalJudge:
		t := task.(models.LocalTask)
		lg.Infof("题目id：%v语言为：%v", t.ProblemId, t.Language)
		//TODO 查询题目对应的模板template,然后合并
		//假设已经合并了完整代码，以下实例为java的两数之和
		return "import java.util.Scanner; \n class a {\n   public static void main(String[] args){\n Scanner sc=new Scanner(System.in);\n System.out.print(sc.nextInt()+sc.nextInt());\n }\n \n}"
	}
	return nil
}
