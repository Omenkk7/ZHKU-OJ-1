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
func (s *Service) UpdateUser(user *models.User) (res *utils.Result) {
	lg := utils.GetDefaultLogger()
	//selector是筛选条件，update是要更新的内容
	selector := bson.M{
		"_id": user.ID,
	}
	//更改的密码需要加密
	hashPassword, _ := utils.HashPassword(user.Password) //hash加密
	user.Password = hashPassword
	//动态构造bson
	update, err := utils.GenerateUpdateBson(user)
	if err != nil {
		lg.Info(utils.ConstructingBsonException, err)
		return res.Fail(utils.ConstructingBsonException)
	}
	lg.Infof("更新_id:%s\nselector:%s\nbson:%s", user.ID, selector, update)
	err = s.dao.UpdateUser(context.Background(), selector, update)
	if err != nil {
		lg.Info(utils.UpdateFailed, err)
		return res.Fail(utils.UpdateFailed)
	}
	return res.Success(utils.UpdateSuccess, "")
}

// DeleteUser 通过id删除用户
func (s *Service) DeleteUser(id string) (res *utils.Result) {
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
		lg.Info(utils.UserNotExist, err)
		return res.Fail(utils.UserNotExist)
	}
	//3.不能删管理员
	if daoUser.Role == utils.StatusAdmin {
		return res.Fail(utils.CanotDeleteAdmin)
	}

	//4.删除
	lg.Infof("删除id%s", id)
	err = s.dao.DeleteUser(context.Background(), selector)
	if err != nil {
		lg.Info(utils.DeleteFail, err)
		return res.Fail(utils.DeleteFail)
	}
	lg.Info(utils.DeleteSuccess)
	return res.Success(utils.DeleteSuccess, "")
}

// GetOneUser 构造query条件查用户
func (s *Service) GetOneUser(user *models.User) (res *utils.Result) {
	lg := utils.GetDefaultLogger()
	//构造bson
	query, err := bson.Marshal(user)
	if err != nil {
		lg.Info(utils.ConstructingBsonException, err)
		return res.Fail(utils.ConstructingBsonException)
	}
	lg.Infof("查询条件: %s", query)
	//调用Dao查询user
	user, err = s.dao.GetOneUser(context.Background(), query)
	if err != nil || user == nil {
		lg.Info(utils.UserNotExist, err)
		return res.Fail(utils.UserNotExist)
	}
	//响应
	return res.Success(utils.SelectSuccess, user)
}

// GetUserList 查一堆数据
func (s *Service) GetUserList(comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	items, err = s.dao.GetUserList(context.Background(), comQuery)
	return
}

// PostUser 注册
// TODO 校验邮箱，电话，密码的格式
func (s *Service) PostUser(postUser *dto.ReqPostUser) (res *utils.Result) {
	lg := utils.GetDefaultLogger()
	//1.判空
	if postUser.Username == "" || postUser.Password == "" {
		lg.Info(utils.CountOrPasswordCannotBeNull)
		return res.Fail(utils.CountOrPasswordCannotBeNull)
	}
	//2.检查用户名是否存在
	query := bson.M{
		"username": postUser.Username,
	}
	user, _ := s.dao.GetOneUser(context.Background(), query)
	if user != nil {
		lg.Infof("用户名%s已被注册", postUser.Username)
		return res.Fail(utils.CountHasBeenRegistered)
	}
	//3.未被注册，一切正常
	hashPassword, _ := utils.HashPassword(postUser.Password) //hash加密
	user = &models.User{
		Username: postUser.Username,
		Password: hashPassword,
		Email:    postUser.Email,
		Phone:    postUser.Phone,
		Status:   utils.StatusNormal,
		Role:     utils.StatusUser,
		Ctime:    time.Now().Unix(),
		Mtime:    time.Now().Unix(),
	}
	_, err := s.dao.CreateUser(context.Background(), user)
	if err != nil {
		lg.Info(utils.ServerAbnormal, err)
		return res.Fail(utils.ServerAbnormal)
	}
	return res.Success(utils.RegisterSuccess, "")
}

// UserLogin 登录
func (s *Service) UserLogin(loginUser *dto.ReqPostLoginUser) (res *utils.Result) {
	lg := utils.GetDefaultLogger()
	//1.判空
	if loginUser.Username == "" || loginUser.Password == "" {
		lg.Info(utils.CountOrPasswordCannotBeNull)
		return res.Fail(utils.CountOrPasswordCannotBeNull)
	}
	//2.查询用户是否存在
	query := bson.M{
		"username": loginUser.Username,
	}
	user, err := s.dao.GetOneUser(context.Background(), query)
	if err != nil {
		lg.Info("查询用户异常: %v", err)
		return res.Fail(utils.CountException)
	}
	//3.判断密码是否正确
	if !utils.CheckPasswordHash(loginUser.Password, user.Password) {
		lg.Info(utils.PasswordException)
		return res.Fail(utils.PasswordException)
	}
	//4.判断账号是否被封禁
	if user.Status == utils.StatusBanned {
		lg.Info(utils.CountBaned)
		return res.Fail(utils.CountBaned)
	}
	//5.登录成功，生成jwt
	jwt, err := utils.GenerateStringToken(user) //jwt令牌
	if err != nil {
		lg.Info(utils.ConstructingJWTException, err)
		return res.Fail(utils.ConstructingJWTException)
	}
	lg.Info("jwt:", jwt)
	return res.Success(utils.LoginSuccess, jwt)
}
