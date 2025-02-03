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

func (s *Service) PostLabel(reqLabel *dto.ReqLabel) (data interface{}, err error) {
	lg := utils.GetDefaultLogger()
	//1.标签名称和作者不能为空
	if reqLabel.Name == "" || reqLabel.Creator == "" {
		lg.Info(utils.LabelOrAuthorCannotBeNull)
		return
	}
	//2.检查标签是否存在
	query := bson.M{
		"name": reqLabel.Name,
	}
	dtoLabel, _ := s.dao.GetOneLabel(context.Background(), query)
	if dtoLabel != nil {
		lg.Infof("标签%s已存在", dtoLabel.Name)
		return
	}
	//3.标签不存在，一切正常，构建入库模型
	dtoLabel = &models.Label{
		Name:    reqLabel.Name,
		Creator: reqLabel.Creator,
		Status:  utils.StatusPublic,
		Ctime:   time.Now().Unix(),
		Mtime:   time.Now().Unix(),
	}
	//4.入库
	_, err = s.dao.CreateLabel(context.Background(), dtoLabel)
	if err != nil {
		lg.Info(utils.ServerErr, err)
		return
	}
	return
}

// GetLabelList  查一堆数据
func (s *Service) GetLabelList(comQuery *utils.CommonQuery) (data interface{}, err error) {
	data, err = s.dao.GetLabelList(context.Background(), comQuery)
	return
}

// GetOneLabel 查一个数据
func (s *Service) GetOneLabel(label *models.Label) (data interface{}, err error) {
	lg := utils.GetDefaultLogger()
	//构造bson
	query, err := bson.Marshal(label)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		return
	}
	lg.Infof("查询条件: %s", query)
	//调用Dao查询daoLabel
	daoLabel, err := s.dao.GetOneLabel(context.Background(), query)
	if err != nil || label == nil {
		lg.Info(utils.LabelNotExist, err)
		return
	}
	return daoLabel, nil
}

// UpdateLabel 更新数据
func (s *Service) UpdateLabel(reqLabel *dto.ReqLabel) (data interface{}, err error) {
	lg := utils.GetDefaultLogger()
	//selector是筛选条件，update是要更新的内容
	selector := bson.M{
		"_id": reqLabel.ID,
	}
	//动态构造bson
	update, err := utils.GenerateUpdateBson(reqLabel)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		return
	}
	lg.Infof("更新_id:%s\nselector:%s\n", reqLabel.ID, selector)
	err = s.dao.UpdateLabel(context.Background(), selector, update)
	if err != nil {
		lg.Info(utils.UpdateErr, err)
		return
	}
	return
}

// DeleteLabel 通过id删除标签
func (s *Service) DeleteLabel(id string) (data interface{}, err error) {
	lg := utils.GetDefaultLogger()
	//1.删之前先查询是否有该条数据
	objectId, _ := primitive.ObjectIDFromHex(id)
	selector := bson.M{
		"_id": objectId,
	}
	lg.Infof("查询_id: %s", objectId)

	//2.调用Dao查询daoLabel
	daoLabel, err := s.dao.GetOneLabel(context.Background(), selector)
	if err != nil || daoLabel == nil {
		lg.Info(utils.LabelNotExist, err)
		return
	}

	//3.删除
	lg.Infof("删除id%s", id)
	err = s.dao.DeleteLabel(context.Background(), selector)
	if err != nil {
		lg.Info(utils.DeleteErr, err)
		return
	}
	lg.Info(utils.DeleteSuccess)
	return
}
