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
	"zhku-oj-server/pkg/utils/errors"
	"zhku-oj-server/pkg/utils/myUtils"
	"zhku-oj-server/pkg/utils/success"
)

// UpdateUser 通过id更新数据
func (s *Service) UpdateUser(user *models.User) (res *models.Result) {
	lg := utils.GetDefaultLogger()
	selector := bson.M{
		"_id": user.ID,
	}
	hashPassword, _ := myUtils.HashPassword(user.Password) //hash加密
	update := bson.D{
		{"$set", bson.D{
			{"password", hashPassword},
			{"email", user.Email},
			{"phone", user.Phone},
			{"username", user.Username},
			{"nickname", user.Nickname},
			{"mtime", time.Now().Unix()},
		}},
	}
	lg.Infof("更新_id:%s\nselector:%s\n", user.ID, selector)
	err := s.dao.UpdateUser(context.Background(), selector, update)
	if err != nil {
		lg.Info(errors.UpdateFailed, err)
		return res.Fail(errors.UpdateFailed)
	}
	return res.Success(success.UpdateSuccess, "")
}

// DeleteUser 通过id删除用户
func (s *Service) DeleteUser(id string) (res *models.Result) {
	lg := utils.GetDefaultLogger()
	//1.删之前先查询是否有该条数据
	objectId, err := primitive.ObjectIDFromHex(id)
	selector := bson.M{
		"_id": objectId,
	}
	lg.Infof("查询_id: %s", objectId)

	//2.调用Dao查询user
	user, err := s.dao.GetOneUser(context.Background(), selector)
	if err != nil {
		lg.Info(errors.UserNotExist, err)
		return res.Fail(errors.UserNotExist)
	}
	//3.不能删管理员
	if user.Role == utils.StatusAdmin {
		return res.Fail(errors.CanotDeleteAdmin)
	}

	//4.删除
	lg.Infof("删除id%s", id)
	err = s.dao.DeleteUser(context.Background(), selector)
	if err != nil {
		lg.Info(errors.DeleteFail, err)
		return res.Fail(errors.DeleteFail)
	}
	lg.Info(success.DeleteSuccess)
	return res.Success(success.DeleteSuccess, "")
}

// GetOneUser 构造query条件查用户
func (s *Service) GetOneUser(user *models.User) (res *models.Result) {
	lg := utils.GetDefaultLogger()
	//构造bson
	query, err := bson.Marshal(user)
	if err != nil {
		lg.Info(errors.ConstructingBsonException, err)
		return res.Fail(errors.ConstructingBsonException)
	}
	lg.Infof("查询条件: %s", query)
	//调用Dao查询user
	user, err = s.dao.GetOneUser(context.Background(), query)
	if err != nil || user == nil {
		lg.Info(errors.UserNotExist, err)
		return res.Fail(errors.UserNotExist)
	}
	//响应
	return res.Success(success.SelectSuccess, user)
}

// GetUserList 查一堆数据
func (s *Service) GetUserList(comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	items, err = s.dao.GetUserList(context.Background(), comQuery)
	return
}

// PostUser 注册
// TODO 校验邮箱，电话，密码的格式
func (s *Service) PostUser(postUser *dto.ReqPostUser) (res *models.Result) {
	lg := utils.GetDefaultLogger()
	//1.判空
	if postUser.Username == "" || postUser.Password == "" {
		lg.Info(errors.CountOrPasswordCannotBeNull)
		return res.Fail(errors.CountOrPasswordCannotBeNull)
	}
	//2.检查用户名是否存在
	query := bson.M{
		"username": postUser.Username,
	}
	user, _ := s.dao.GetOneUser(context.Background(), query)
	if user != nil {
		lg.Infof("用户名%s已被注册", postUser.Username)
		return res.Fail(errors.CountHasBeenRegistered)
	}
	//3.未被注册，一切正常
	hashPassword, _ := myUtils.HashPassword(postUser.Password) //hash加密
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
		lg.Info(errors.ServerAbnormal, err)
		return res.Fail(errors.ServerAbnormal)
	}
	return res.Success(success.RegisterSuccess, "")
}

// UserLogin 登录
func (s *Service) UserLogin(loginUser *dto.ReqPostLoginUser) (res *models.Result) {
	lg := utils.GetDefaultLogger()
	//1.判空
	if loginUser.Username == "" || loginUser.Password == "" {
		lg.Info(errors.CountOrPasswordCannotBeNull)
		return res.Fail(errors.CountOrPasswordCannotBeNull)
	}
	//2.查询用户是否存在
	query := bson.M{
		"username": loginUser.Username,
	}
	user, err := s.dao.GetOneUser(context.Background(), query)
	if err != nil {
		lg.Info("查询用户异常: %v", err)
		return res.Fail(errors.CountException)
	}
	//3.判断密码是否正确
	if !myUtils.CheckPasswordHash(loginUser.Password, user.Password) {
		lg.Info(errors.PasswordException)
		return res.Fail(errors.PasswordException)
	}
	//4.判断账号是否被封禁
	if user.Status == utils.StatusBanned {
		lg.Info(errors.CountBaned)
		return res.Fail(errors.CountBaned)
	}
	//5.成功查询
	jwt, err := myUtils.GenerateToken(user, utils.JwtTokenSecretKey, time.Hour*24) //jwt令牌
	if err != nil {
		lg.Info(errors.ConstructingJWTException, err)
		return res.Fail(errors.ConstructingJWTException)
	}
	lg.Info("jwt:", jwt)
	return res.Success(success.LoginSuccess, jwt)
}
