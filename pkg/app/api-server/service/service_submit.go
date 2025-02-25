package service

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// Submit 提交代码
func (s *Service) Submit(reqSubmit *dto.ReqSubmit) (err error) {
	lg := utils.GetDefaultLogger()
	//1.代码不能为空
	if reqSubmit.Code == "" {
		lg.Info(utils.CodeCannotBeNull)
		return errors.New(utils.CodeCannotBeNull)
	}
	//2.一切正常，构建入库模型
	dtoSubmit := &models.Submit{
		Code:       reqSubmit.Code,
		Language:   reqSubmit.Language,
		Status:     utils.WaitForJudge,
		ProblemId:  reqSubmit.ProblemId,
		UserId:     reqSubmit.UserId,
		Evaluation: "nil",
		Stime:      time.Now().Unix(),
	}
	//3.入库
	_, err = s.dao.Submit(context.Background(), dtoSubmit)
	if err != nil {
		lg.Info(utils.ServerErr, err)
		return errors.New(utils.ServerErr)
	}
	return
}
