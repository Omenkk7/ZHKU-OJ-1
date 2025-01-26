/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 19:46
@Name: service.go
@Description:
*/

package service

import (
	"context"
	"zhku-oj-server/pkg/dao"
)

type Service struct {
	dao *dao.Dao
}

func NewService() *Service {
	return &Service{
		dao: dao.NewDao(),
	}
}

func (s *Service) Close(ctx context.Context) {
	s.dao.Close(ctx)
}
