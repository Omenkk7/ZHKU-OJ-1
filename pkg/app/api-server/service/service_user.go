/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 20:16
@Name: service_user.go
@Description:
*/

package service

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// UpdateUser 通过id更新数据
func (s *Service) UpdateUser(res *utils.Result, reqUser *models.User) {
	lg := utils.GetDefaultLogger()
	//selector是筛选条件，update是要更新的内容
	selector := bson.M{
		"_id": reqUser.ID,
	}
	//更改的密码需要加密
	hashPassword, _ := utils.HashPassword(reqUser.Password) //hash加密
	reqUser.Password = hashPassword
	//动态构造bson
	update, err := utils.GenerateUpdateBson(reqUser)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		res.Fail(utils.ConstructingBsonErr)
		return
	}
	lg.Infof("更新_id:%s\nselector:%s\nbson:%s", reqUser.ID, selector, update)
	err = s.dao.UpdateUser(context.Background(), selector, update)
	if err != nil {
		lg.Info(utils.UpdateErr, err)
		res.Fail(utils.UpdateErr)
		return
	}
	res.Success(utils.UpdateSuccess, "")
	return
}

// DeleteUser 通过id删除用户
func (s *Service) DeleteUser(res *utils.Result, id string) {
	lg := utils.GetDefaultLogger()
	//1.删之前先查询是否有该条数据
	objectId, _ := primitive.ObjectIDFromHex(id)
	selector := bson.M{
		"_id": objectId,
	}
	lg.Infof("查询_id: %s", objectId)

	//2.调用Dao查询user
	daoUser, err := s.dao.GetOneUser(context.Background(), selector)
	if err != nil {
		lg.Info(utils.UserNotExistErr, err)
		res.Fail(utils.UserNotExistErr)
		return
	}
	//3.不能删管理员
	if daoUser.Role == utils.StatusAdmin {
		res.Fail(utils.DeleteAdminErr)
		return
	}

	//4.删除
	lg.Infof("删除id%s", id)
	err = s.dao.DeleteUser(context.Background(), selector)
	if err != nil {
		lg.Info(utils.DeleteErr, err)
		res.Fail(utils.DeleteErr)
		return
	}
	lg.Info(utils.DeleteSuccess)
	res.Success(utils.DeleteSuccess, "")
	return
}

// GetOneUser 构造query条件查用户
func (s *Service) GetOneUser(res *utils.Result, reqUser *models.User) {
	lg := utils.GetDefaultLogger()
	//构造bson
	query, err := bson.Marshal(reqUser)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		res.Fail(utils.ConstructingBsonErr)
		return
	}
	lg.Infof("查询条件: %s", query)
	//调用Dao查询user
	daoUser, err := s.dao.GetOneUser(context.Background(), query)
	if err != nil || daoUser == nil {
		lg.Info(utils.UserNotExistErr, err)
		res.Fail(utils.UserNotExistErr)
		return
	}
	//响应
	res.Success(utils.SelectSuccess, daoUser)
	return
}

// GetUserList 查一堆数据
func (s *Service) GetUserList(comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	items, err = s.dao.GetUserList(context.Background(), comQuery)
	return
}

// PostUser 注册
// TODO 校验邮箱，电话，密码的格式
func (s *Service) PostUser(res *utils.Result, reqUser *dto.ReqPostUser) {
	lg := utils.GetDefaultLogger()
	//1.判空
	if reqUser.Username == "" || reqUser.Password == "" {
		lg.Info(utils.CountOrPasswordNullErr)
		res.Fail(utils.CountOrPasswordNullErr)
		return
	}
	//2.检查用户名是否存在
	query := bson.M{
		"username": reqUser.Username,
	}
	daoUser, _ := s.dao.GetOneUser(context.Background(), query)
	if daoUser != nil {
		lg.Infof("用户名%s已被注册", reqUser.Username)
		res.Fail(utils.RegisteredErr)
		return
	}
	//3.未被注册，一切正常
	hashPassword, _ := utils.HashPassword(reqUser.Password) //hash加密
	daoUser = &models.User{
		Username: reqUser.Username,
		Password: hashPassword,
		Email:    reqUser.Email,
		Phone:    reqUser.Phone,
		Status:   utils.StatusNormal,
		Role:     utils.StatusUser,
		Ctime:    time.Now().Unix(),
		Mtime:    time.Now().Unix(),
	}
	_, err := s.dao.CreateUser(context.Background(), daoUser)
	if err != nil {
		lg.Info(utils.ServerErr, err)
		res.Fail(utils.ServerErr)
		return
	}
	res.Success(utils.RegisterSuccess, "")
	return
}

// UserLogin 登录
func (s *Service) UserLogin(res *utils.Result, reqUser *dto.ReqPostLoginUser) {
	lg := utils.GetDefaultLogger()
	//1.判空
	if reqUser.Username == "" || reqUser.Password == "" {
		lg.Info(utils.CountOrPasswordNullErr)
		res.Fail(utils.CountOrPasswordNullErr)
		return
	}
	//2.查询用户是否存在
	query := bson.M{
		"username": reqUser.Username,
	}
	daoUser, err := s.dao.GetOneUser(context.Background(), query)
	if err != nil {
		lg.Info("查询用户异常: %v", err)
		res.Fail(utils.CountErr)
		return
	}
	//3.判断密码是否正确
	if !utils.CheckPasswordHash(reqUser.Password, daoUser.Password) {
		lg.Info(utils.PasswordErr)
		res.Fail(utils.PasswordErr)
		return
	}
	//4.判断账号是否被封禁
	if daoUser.Status == utils.StatusBanned {
		lg.Info(utils.CountBanedErr)
		res.Fail(utils.CountBanedErr)
		return
	}
	//5.登录成功，生成jwt
	jwt, err := utils.GenerateStringToken(daoUser) //jwt令牌
	if err != nil {
		lg.Info(utils.ConstructingJWTErr, err)
		res.Fail(utils.ConstructingJWTErr)
		return
	}
	lg.Info("jwt:", jwt)
	res.Success(utils.LoginSuccess, jwt)
	return
}
