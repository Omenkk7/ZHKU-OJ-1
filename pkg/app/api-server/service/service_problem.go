package service

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// PostProblem 新建题目
func (s *Service) PostProblem(reqPostProblem *dto.ReqProblem) (res *utils.Result) {
	lg := utils.GetDefaultLogger()
	//1.题目名称和出题人不能为空
	if reqPostProblem.Title == "" || reqPostProblem.Creator == "" {
		lg.Info(utils.ProblemOrAuthorCannotBeNull)
		return res.Fail(utils.ProblemOrAuthorCannotBeNull)
	}
	//2.检查题目是否存在
	query := bson.M{
		"name": reqPostProblem.Title,
	}
	dtoProblem, _ := s.dao.GetOneProblem(context.Background(), query)
	if dtoProblem != nil {
		lg.Infof("题目%s已存在", reqPostProblem.Title)
		return res.Fail(utils.ProblemIsExist)
	}
	//TODO ————————————————————文件存储和入库，要保证事务性————————————————————————
	//3.题目不存在，一切正常，构建入库模型
	//TODO 写一个工具类，动态构建入库模型
	dtoProblem = &models.Problem{
		ID:            primitive.ObjectID{},
		Creator:       reqPostProblem.Creator,
		Description:   reqPostProblem.Description,
		InputExample:  reqPostProblem.InputExample,
		OutputExample: reqPostProblem.OutputExample,
		Status:        reqPostProblem.Status,
		SubmitNum:     reqPostProblem.SubmitNum,
		Difficulty:    reqPostProblem.Difficulty,
		Labels:        reqPostProblem.Labels,
		PassNum:       reqPostProblem.PassNum,
		Title:         reqPostProblem.Title,
		URL:           reqPostProblem.URL,
		Ctime:         time.Now().Unix(),
		Mtime:         time.Now().Unix(),
	}
	//4.入库
	_, err := s.dao.CreateProblem(context.Background(), dtoProblem)
	if err != nil {
		lg.Info(utils.ServerErr, err)
		return res.Fail(utils.ServerErr)
	}
	//TODO —————————————————————————————保证事务一致性———————————————————————————————————————-
	return res.Success(utils.CreateSuccess, "")
}

// GetProblemList  查一堆数据
func (s *Service) GetProblemList(comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	items, err = s.dao.GetProblemList(context.Background(), comQuery)
	return
}

// GetOneProblem 查一个数据
func (s *Service) GetOneProblem(reqProblem *models.Problem) (res *utils.Result) {
	lg := utils.GetDefaultLogger()
	//构造bson
	query, err := bson.Marshal(reqProblem)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		return res.Fail(utils.ConstructingBsonErr)
	}
	lg.Infof("查询条件: %s", query)
	//调用Dao查询daoProblem
	daoProblem, err := s.dao.GetOneProblem(context.Background(), query)
	if err != nil || reqProblem == nil {
		lg.Info(utils.ProblemNotExist, err)
		return res.Fail(utils.ProblemNotExist)
	}
	//响应
	return res.Success(utils.SelectSuccess, daoProblem)
}

// UpdateProblem 更新题目
func (s *Service) UpdateProblem(reqProblem *dto.ReqProblem) (res *utils.Result) {
	lg := utils.GetDefaultLogger()
	//selector是筛选条件，update是要更新的内容
	selector := bson.M{
		"_id": reqProblem.ID,
	}
	//动态构造bson
	update, err := utils.GenerateUpdateBson(reqProblem)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		return res.Fail(utils.ConstructingBsonErr)
	}
	lg.Infof("更新_id:%s\nselector:%s\n", reqProblem.ID, selector)
	err = s.dao.UpdateProblem(context.Background(), selector, update)
	if err != nil {
		lg.Info(utils.UpdateErr, err)
		return res.Fail(utils.UpdateErr)
	}
	return res.Success(utils.UpdateSuccess, "")
}

// DeleteProblem 通过id删除题目
func (s *Service) DeleteProblem(id string) (res *utils.Result) {
	lg := utils.GetDefaultLogger()
	//1.删之前先查询是否有该条数据
	objectId, _ := primitive.ObjectIDFromHex(id)
	selector := bson.M{
		"_id": objectId,
	}
	lg.Infof("查询_id: %s", objectId)

	//2.调用Dao查询daoProblem
	daoProblem, err := s.dao.GetOneProblem(context.Background(), selector)
	if err != nil || daoProblem == nil {
		lg.Info(utils.ProblemNotExist, err)
		return res.Fail(utils.ProblemNotExist)
	}

	//3.删除
	lg.Infof("删除id%s", id)
	err = s.dao.DeleteProblem(context.Background(), selector)
	if err != nil {
		lg.Info(utils.DeleteErr, err)
		return res.Fail(utils.DeleteErr)
	}
	lg.Info(utils.DeleteSuccess)
	return res.Success(utils.DeleteSuccess, "")
}
