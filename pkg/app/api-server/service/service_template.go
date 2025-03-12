package service

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"strings"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

func (s *Service) MergeTemplate(task interface{}) (code interface{}, err error) {
	lg := utils.GetDefaultLogger()
	t := task.(models.LocalTask)
	lg.Infof("题目id：%v语言为：%v", t.ProblemId, t.Language)

	//1.template集合内查模板，problemID——>模板
	problemId, _ := primitive.ObjectIDFromHex(t.ProblemId)
	query := bson.M{
		"_id": problemId,
	}
	daoProblem, err := s.dao.GetOneProblem(context.Background(), query)
	/*res, err := s.dao.GetOne(context.Background(), dao.ProblemTable, &models.Problem{}, query)*/
	if err != nil {
		lg.Info("find template error:", err)
		return "", err
	}

	//2.拿到模板后，进行代码的合并—————直接用用户提交的代码，替换掉模板里的#function,去除多余的空格
	//字面的\n转为真正的换行\n
	daoProblem.Template = strings.Replace(daoProblem.Template, "\\n", "\n", -1)
	t.Code = strings.Replace(t.Code, "\\n", "\n", -1)
	//替换#function为用户写的函数
	code = strings.TrimSpace(strings.Replace(daoProblem.Template, "#function", t.Code, 1))

	//便于观察模板代码，函数代码，完整代码
	fmt.Println()
	fmt.Println()
	fmt.Println("模板代码：", daoProblem.Template)
	fmt.Println("函数代码:", t.Code)
	fmt.Println("完整代码：", code)

	return code, err
}
