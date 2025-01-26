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
	"zhku-oj-server/pkg/utils"
)

func (s *Server) RegisterUser(g *gin.RouterGroup) {
	userGroup := g.Group("/user")
	{
		userGroup.GET("/", s.getSomeUser)
		userGroup.GET("/:id", s.getOneUser)
		userGroup.POST("/", s.postUser) // 注册
		userGroup.PUT("/:id", s.putUser)
		userGroup.DELETE("/:id", s.deleteUser)
		userGroup.POST("/login", s.login)
	}
}

func (s *Server) getOneUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	id := c.Param("id")
	lg.Println("getUserById", id)
	utils.SuccessResponse(c, map[string]interface{}{})
}
func (s *Server) getSomeUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("user router get start")
	query := c.Request.URL.Query()
	cq := utils.BuildCommonQuery(utils.Query(query))
	lg.Info("get", cq)
	resp, err := s.svc.GetUserList(cq)
	if err != nil {
		lg.Errorf("getUserList: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, resp)
}

func (s *Server) postUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	// 1) 获取传输参数，绑定到dtoUser
	dtoUser := &dto.ReqPostUser{}
	if err := c.ShouldBind(dtoUser); err != nil {
		lg.Errorf("postUser: %v", err)
		utils.BadRequest(c, err)
		return
	}
	lg.Info("postUser", dtoUser)
	// TODO://鉴权，数据粒度的权限检查，也可以抽到一个路由函数中，非侵入式的鉴权
	// 2）参数校验
	if err := s.svc.CheckPostUserParams(dtoUser); err != nil {
		lg.Errorf("postUser: %v", err)
		utils.BadRequest(c, err)
		return
	}
	// 3）构建mongo入库模型
	userModel, err := s.svc.BuildPostUser(dtoUser)
	if err != nil {
		lg.Errorf("postUser: %v", err)
		utils.BadRequest(c, err)
		return
	}
	lg.Infof("builded userModel: %v", userModel)

	// 4）写mongo
	id, err := s.svc.CreateUser(userModel)
	if err != nil {
		lg.Errorf("postUser err: %v", err)
		utils.SuccessResponse(c, map[string]interface{}{
			"success": false,
		})
		return
	}
	// 5）返回
	utils.SuccessResponse(c, map[string]interface{}{
		"success": true,
		"id":      id,
	})
	return
}

func (s *Server) putUser(c *gin.Context) {}

func (s *Server) deleteUser(c *gin.Context) {}

func (s *Server) login(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("user router get start")
	loginUser := &dto.ReqPostLoginUser{}

	if err := c.ShouldBind(loginUser); err != nil {
		lg.Errorf("loginUser: %v", err)
		utils.BadRequest(c, err)
		return
	}
	lg.Infof("loginUser ReqPostLoginUser: %v", loginUser)
	// 参数校验
	if err := s.svc.CheckUserLoginParams(loginUser); err != nil {
		lg.Errorf("CheckUserLoginParams: %v", err)
		utils.BadRequest(c, err)
		return
	}
	token, err := s.svc.UserLogin(loginUser)
	if err != nil {
		lg.Errorf("UserLogin: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, map[string]interface{}{
		"token": token,
	})
}
