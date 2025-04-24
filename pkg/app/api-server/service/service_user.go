/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 20:16
@Name: service_user.go
@Description:
*/

package service

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// UpdateUser 通过id更新数据
func (s *Service) UpdateUser(reqUser *dto.ReqUser, operatorRole int32) (id string, err error) {
	lg := utils.GetDefaultLogger()

	// 获取要更新的用户信息
	objectId, _ := primitive.ObjectIDFromHex(reqUser.ID)
	selector := bson.M{
		"_id": objectId,
	}

	// 查询原用户信息
	originalUser, err := s.dao.GetOneUser(context.Background(), selector)
	if err != nil {
		lg.Info(utils.UserNotExistErr, err)
		return "", errors.New(utils.UserNotExistErr)
	}

	// 如果要修改的是管理员账号，检查操作者是否为超级管理员
	if originalUser.Role == 1 || reqUser.Role == 1 {
		if operatorRole != 0 { // 不是超级管理员
			lg.Info("权限不足，只有超级管理员可以修改管理员账号")
			return "", errors.New("权限不足，只有超级管理员可以修改管理员账号")
		}
	}

	// 更改的密码需要加密
	if reqUser.Password != "" {
		hashPassword, _ := utils.HashPassword(reqUser.Password) // hash加密
		reqUser.Password = hashPassword
	}
	reqUser.Mtime = time.Now().Unix()

	// 动态构造bson
	update, err := s.dao.GenerateUpdateBson(reqUser)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		return "", err
	}

	lg.Infof("更新_id:%s\nselector:%s\nbson:%s", reqUser.ID, selector, update)
	_, err = s.dao.UpdateUser(context.Background(), selector, update)
	if err != nil {
		lg.Info(utils.UpdateErr, err)
		return "", err
	}
	return reqUser.ID, nil
}

// DeleteUser 通过id删除用户
func (s *Service) DeleteUser(id string, operatorRole int32) (_ string, err error) {
	lg := utils.GetDefaultLogger()
	// 1.删之前先查询是否有该条数据
	objectId, _ := primitive.ObjectIDFromHex(id)
	selector := bson.M{
		"_id": objectId,
	}
	lg.Infof("查询_id: %s", objectId)

	// 2.调用Dao查询user
	daoUser, err := s.dao.GetOneUser(context.Background(), selector)
	if err != nil {
		lg.Info(utils.UserNotExistErr, err)
		return "", errors.New(utils.UserNotExistErr)
	}

	// 3.如果要删除的是管理员账号，检查操作者是否为超级管理员
	if daoUser.Role == 1 { // 管理员角色
		if operatorRole != 0 { // 不是超级管理员
			lg.Info("权限不足，只有超级管理员可以删除管理员账号")
			return "", errors.New("权限不足，只有超级管理员可以删除管理员账号")
		}
	}

	// 4.删除
	lg.Infof("删除id%s", id)
	_, err = s.dao.DeleteUser(context.Background(), selector)
	if err != nil {
		lg.Info(utils.DeleteErr, err)
		return "", errors.New(utils.DeleteErr)
	}
	lg.Info(utils.DeleteSuccess)
	return id, nil
}

// GetOneUser 构造query条件查用户
func (s *Service) GetInfor(userId primitive.ObjectID) (user *models.User, err error) {
	lg := utils.GetDefaultLogger()
	//构造bson
	query := bson.M{
		"_id": userId,
	}
	lg.Infof("查询条件: %s", query)
	//调用Dao查询user
	daoUser, err := s.dao.GetOneUser(context.Background(), query)
	if daoUser == nil {
		lg.Info(utils.UserNotExistErr, err)
		return nil, errors.New(utils.UserNotExistErr)
	}
	return daoUser, nil
}

// GetOneUser 构造query条件查用户
func (s *Service) GetOneUser(reqUser *dto.ReqUser) (user *models.User, err error) {
	lg := utils.GetDefaultLogger()
	//构造bson
	query, err := bson.Marshal(reqUser)
	if err != nil {
		lg.Info(utils.ConstructingBsonErr, err)
		return nil, errors.New(utils.ConstructingBsonErr)
	}
	lg.Infof("查询条件: %s", query)
	//调用Dao查询user
	daoUser, err := s.dao.GetOneUser(context.Background(), query)
	if daoUser == nil {
		lg.Info(utils.UserNotExistErr, err)
		return nil, errors.New(utils.UserNotExistErr)
	}
	return daoUser, nil
}

// GetUserList 查一堆数据
func (s *Service) GetUserList(comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	items, err = s.dao.GetUserList(context.Background(), comQuery)
	return items, err
}

// TODO 校验邮箱，电话，密码的格式
// PostUser 注册
func (s *Service) PostUser(reqPostUser *dto.ReqPostUser, operatorRole int32) (id string, err error) {
	lg := utils.GetDefaultLogger()
	lg.Info("注册用户......")

	// 检查用户名是否已存在
	existUser, err := s.dao.GetOneUser(context.Background(), bson.M{"username": reqPostUser.Username})
	if err != nil && err != mongo.ErrNoDocuments {
		lg.Info(utils.QueryErr, err)
		return "", errors.New(utils.QueryErr)
	}
	if existUser != nil {
		lg.Info(utils.UserExistErr)
		return "", errors.New(utils.UserExistErr)
	}

	// 验证角色值是否有效
	if reqPostUser.Role < 1 || reqPostUser.Role > 4 {
		lg.Info("无效的用户角色")
		return "", errors.New("无效的用户角色，请选择有效的角色：1-管理员、2-教师、3-助教、4-学生")
	}

	// 如果要创建管理员账号，检查操作者是否为超级管理员
	if reqPostUser.Role == 1 { // 管理员角色
		if operatorRole != 0 { // 不是超级管理员
			lg.Info("权限不足，只有超级管理员可以创建管理员账号")
			return "", errors.New("权限不足，只有超级管理员可以创建管理员账号")
		}
	}

	// 创建用户对象
	user := &models.User{
		ID:       primitive.NewObjectID(),
		Username: reqPostUser.Username,
		Password: reqPostUser.Password,
		Email:    reqPostUser.Email,
		Phone:    reqPostUser.Phone,
		Role:     reqPostUser.Role,
		Nickname: reqPostUser.Nickname,
		Class:    reqPostUser.Class,
		Sid:      reqPostUser.Sid,
		Status:   1, // 默认状态为正常
		Ctime:    time.Now().Unix(),
		Mtime:    time.Now().Unix(),
	}

	// 密码加密
	hashPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		lg.Info(utils.HashErr, err)
		return "", errors.New(utils.HashErr)
	}
	user.Password = hashPassword

	// 调用dao层创建用户
	id, err = s.dao.CreateUser(context.Background(), user)
	if err != nil {
		lg.Info(utils.CreateErr, err)
		return "", errors.New(utils.CreateErr)
	}
	return id, nil
}

// UserLogin用户登录
func (s *Service) UserLogin(reqLoginUser *dto.ReqPostLoginUser) (response *dto.LoginResponse, err error) {
	lg := utils.GetDefaultLogger()
	//1.判空
	if reqLoginUser.Username == "" || reqLoginUser.Password == "" {
		lg.Info(utils.CountOrPasswordNullErr)
		return nil, errors.New(utils.CountOrPasswordNullErr)
	}
	//2.查询用户是否存在
	query := bson.M{
		"username": reqLoginUser.Username,
	}
	daoUser, err := s.dao.GetOneUser(context.Background(), query)
	if err != nil || daoUser == nil {
		lg.Info("查询用户异常: %v", err)
		return nil, errors.New(utils.UserNotExistErr)
	}
	//3.判断密码是否正确
	if !utils.CheckPasswordHash(reqLoginUser.Password, daoUser.Password) {
		lg.Info(utils.PasswordErr)
		return nil, errors.New(utils.PasswordErr)
	}
	//4.判断账号是否被封禁
	if daoUser.Status == utils.StatusBanned {
		lg.Info(utils.CountBanedErr)
		return nil, errors.New(utils.CountBanedErr)
	}
	//5.登录成功，生成jwt
	jwt, err := utils.GenerateStringToken(daoUser) //jwt令牌
	if err != nil {
		lg.Info(utils.ConstructingJWTErr, err)
		return nil, err
	}
	lg.Info("jwt:", jwt)

	// 构建登录响应
	response = &dto.LoginResponse{
		Token:    jwt,
		Username: daoUser.Username,
		Role:     daoUser.Role,
	}

	return response, nil
}
