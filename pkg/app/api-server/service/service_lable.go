package service

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// PostLabel 新建标签
func (s *Service) PostLabel(reqLabel *dto.ReqLabel) (err error) {
	lg := utils.GetDefaultLogger()
	//1.标签名称和作者不能为空
	if reqLabel.Name == "" || reqLabel.Creator == "" {
		lg.Info(utils.LabelOrAuthorCannotBeNull)
		return errors.New(utils.LabelOrAuthorCannotBeNull)
	}
	//2.检查标签是否存在
	query := bson.M{
		"name": reqLabel.Name,
	}
	dtoLabel, _ := s.dao.GetOneLabel(context.Background(), query)
	if dtoLabel != nil {
		lg.Infof("标签%s已存在", dtoLabel.Name)
		return errors.New(utils.LabelIsExist)
	}
	//3.标签不存在，一切正常，构建入库模型
	dtoLabel = &models.Label{
		Name:    reqLabel.Name,
		Creator: reqLabel.Creator,
		Status:  reqLabel.Status,
		Ctime:   time.Now().Unix(),
		Mtime:   time.Now().Unix(),
	}
	//4.入库
	_, err = s.dao.CreateLabel(context.Background(), dtoLabel)
	if err != nil {
		lg.Info(utils.ServerErr, err)
		return err
	}
	return
}

// GetLabelList  查一堆数据
func (s *Service) GetLabelList(comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	items, err = s.dao.GetLabelList(context.Background(), comQuery)
	return
}

// GetOneLabel 查一个数据
func (s *Service) GetOneLabel(reqLabel *dto.ReqLabel) (daoLabel *models.Label, err error) {
	lg := utils.GetDefaultLogger()
	//构造bson
	query, err := bson.Marshal(reqLabel)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		return nil, errors.New(utils.ConstructingBsonErr)
	}
	lg.Infof("查询条件: %s", query)
	//调用Dao查询daoLabel
	daoLabel, err = s.dao.GetOneLabel(context.Background(), query)
	if daoLabel == nil {
		lg.Info(utils.LabelNotExist, err)
		return nil, errors.New(utils.LabelNotExist)
	}
	return daoLabel, nil
}

// UpdateLabel 更新数据
func (s *Service) UpdateLabel(reqLabel *dto.ReqLabel) (id string, err error) {
	lg := utils.GetDefaultLogger()
	//selector是筛选条件，update是要更新的内容
	objectId, _ := primitive.ObjectIDFromHex(reqLabel.ID)
	selector := bson.M{
		"_id": objectId,
	}
	//动态构造bson
	update, err := s.dao.GenerateUpdateBson(reqLabel)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		return "", errors.New(utils.ConstructingBsonErr)
	}
	lg.Infof("更新_id:%s\nselector:%s\n", reqLabel.ID, selector)
	_, err = s.dao.UpdateLabel(context.Background(), selector, update)
	if err != nil {
		lg.Info(utils.UpdateErr, err)
		return "", errors.New(utils.UpdateErr)
	}
	return reqLabel.ID, nil
}

// DeleteLabel 通过id删除标签
func (s *Service) DeleteLabel(id string) (_ string, err error) { //返回的id类型应该为string，不然全是0
	lg := utils.GetDefaultLogger()
	//1.删之前先查询是否有该条数据
	objectId, _ := primitive.ObjectIDFromHex(id)
	selector := bson.M{
		"_id": objectId,
	}
	lg.Infof("查询_id: %s", objectId)

	//2.调用Dao查询daoLabel
	_, err = s.dao.GetOneLabel(context.Background(), selector)
	if err != nil {
		lg.Info(utils.LabelNotExist, err)
		return "", errors.New(utils.LabelNotExist)
	}

	//3.删除
	lg.Infof("删除id%s", id)
	_, err = s.dao.DeleteLabel(context.Background(), selector)
	if err != nil {
		lg.Info(utils.DeleteErr, err)
		return "", errors.New(utils.DeleteErr)
	}
	lg.Info(utils.DeleteSuccess)
	return id, nil
}
