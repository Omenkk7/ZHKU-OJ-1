package service

import (
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// Submit 提交代码
func (s *Service) Submit(reqSubmit *dto.ReqSubmit) (submitId string, err error) {
	lg := utils.GetDefaultLogger()
	//1.代码不能为空
	if reqSubmit.Code == "" {
		lg.Info(utils.CodeCannotBeNull)
		return "", errors.New(utils.CodeCannotBeNull)
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
	submitId, err = s.dao.Submit(context.Background(), dtoSubmit)
	if err != nil {
		lg.Info(utils.ServerErr, err)
		return "", errors.New(utils.ServerErr)
	}
	return submitId, nil
}

// GetOneSubmit 查一个数据
func (s *Service) GetOneSubmit(reqSubmit *dto.ReqSubmit) (daoSubmit *models.Submit, err error) {
	lg := utils.GetDefaultLogger()
	//构造bson
	query, err := bson.Marshal(reqSubmit)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		return nil, errors.New(utils.ConstructingBsonErr)
	}
	lg.Infof("查询条件: %s", query)
	//调用Dao查询daoSubmit
	daoSubmit, err = s.dao.GetOneSubmit(context.Background(), query)
	if daoSubmit == nil {
		lg.Info(utils.SubmitNotExist, err)
		return nil, errors.New(utils.SubmitNotExist)
	}
	return daoSubmit, nil
}
