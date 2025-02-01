/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 18:24
@Name: handler_user.go
@Description:
*/

package server

import (
	"github.com/gin-gonic/gin"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// PostUser 注册 /user
func (s *Server) PostUser(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("注册......")
	//把参数解析到postUser
	var postUser *dto.ReqPostUser
	if err := c.BindJSON(&postUser); err != nil {
		return
	}
	//调用service_user层
	res := s.svc.PostUser(postUser)
	//返回结果
	res.Response(c, res)
	return
}

// Login 登录 /login
func (s *Server) Login(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("登录......")
	//把参数解析到结构体loginUser
	var loginUser *dto.ReqPostLoginUser
	if err := c.BindJSON(&loginUser); err != nil {
		return
	}
	//调用service_user层
	res := s.svc.UserLogin(loginUser)
	//返回结果
	res.Response(c, res)
	return
}

// GetOneUser 构造query条件查用户 /:id
func (s *Server) GetOneUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	var user *models.User
	if err := c.BindJSON(&user); err != nil {
		return
	}
	lg.Println("条件查询用户......")
	//调用service_user层
	res := s.svc.GetOneUser(user)
	//返回结果
	res.Response(c, res)
	return
}

// GetSomeUser 查一堆用户 /
func (s *Server) GetSomeUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("查一堆用户......")
	query := c.Request.URL.Query()
	cq := utils.BuildCommonQuery(utils.Query(query))
	lg.Info("get", cq)
	res, err := s.svc.GetUserList(cq)
	if err != nil {
		lg.Errorf("getUserList: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, res)
}

// PutUser 通过id改一个用户 /:id
func (s *Server) PutUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("通过id改一个用户......")

	//把参数解析到结构体loginUser
	var user *models.User
	if err := c.BindJSON(&user); err != nil {
		return
	}

	//调用service_user层
	res := s.svc.UpdateUser(user)
	//返回结果
	res.Response(c, res)
	return
}

// DeleteUser 通过id删一个用户 /:id
func (s *Server) DeleteUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("删用户......")
	id := c.Param("id")
	//调用service_user层
	res := s.svc.DeleteUser(id)
	//返回结果
	res.Response(c, res)
	return
}
