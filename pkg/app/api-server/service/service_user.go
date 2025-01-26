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
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

func (s *Service) checkUsernameExists(username string) (bool, error) {
	var query = map[string]interface{}{
		"name": username,
	}
	m, err := s.dao.GetOneUser(context.Background(), query)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
	}
	if m != nil {
		return true, nil
	}
	return false, nil

}

func (s *Service) CheckPostUserParams(user *dto.ReqPostUser) error {
	// 重名校验
	ok, err := s.checkUsernameExists(user.Name)
	if err != nil {
		return err
	}
	if ok {
		return utils.ErrUsernameExisted
	}
	return nil
}

func (s *Service) CheckUserLoginParams(user *dto.ReqPostLoginUser) error {
	if user.Name == "" || user.Password == "" {
		return errors.New("username or password is empty")
	}
	return nil
}

func (s *Service) BuildPostUser(u *dto.ReqPostUser) (*models.User, error) {
	passwordHash, err := utils.Encoder.Encode([]byte(u.Password))
	return &models.User{
		Name:     u.Name,
		Email:    u.Email,
		Password: string(passwordHash),
		Phone:    u.Phone,
		Role:     utils.DefaultRole,
		Ctime:    time.Now().Unix(),
		Mtime:    time.Now().Unix(),
	}, err

}

func (s *Service) CreateUser(user *models.User) (string, error) {
	return s.dao.CreateUser(context.Background(), user)
}

func (s *Service) UpdateUser(user *models.User) error {
	return nil
}

func (s *Service) DeleteUser(user *models.User) error {
	return nil
}

func (s *Service) GetUserList(comQuery *utils.CommonQuery) (items *utils.RespPageQuery, err error) {
	items, err = s.dao.GetUserList(context.Background(), comQuery)
	return
}
func (s *Service) GetOneUser(query interface{}) (user *models.User, err error) {
	user, err = s.dao.GetOneUser(context.Background(), query)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) UserLogin(loginUser *dto.ReqPostLoginUser) (token string, err error) {
	lg := utils.GetDefaultLogger()
	userModel, err := s.GetOneUser(map[string]interface{}{
		"name": loginUser.Name,
	})
	if err != nil {
		return
	}
	if err = utils.Encoder.Match([]byte(userModel.Password), []byte(loginUser.Password)); err != nil {
		lg.Errorf("Match: %v", err)
		return
	}
	token, err = utils.GenerateStringToken(userModel.ID.Hex(), userModel.Name, userModel.Role)
	if err != nil {
		lg.Errorf("GenerateToken: %v", err)
		return
	}
	return
}
