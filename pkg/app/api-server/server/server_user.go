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
	"zhku-oj-server/pkg/utils/middleware"
)

// 路由器
func (s *Server) RegisterUser(g *gin.RouterGroup) {
	userGroup := g.Group("/user")
	{
		userGroup.POST("/", s.postUser)   // 注册
		userGroup.POST("/login", s.login) // 登录
	}

	// 对以下路由应用JWT拦截器，并且只有管理员可以执行以下操作
	securedGroup := userGroup.Group("/").Use(middleware.JWTInterceptor())
	{
		securedGroup.GET("/", s.getSomeUser)      //查一堆
		securedGroup.GET("/:id", s.getOneUser)    //查一个
		userGroup.PUT("/:id", s.putUser)          //改一个
		securedGroup.DELETE("/:id", s.deleteUser) //删一个
	}
}

// 注册
func (s *Server) postUser(c *gin.Context) {
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

// 登录
func (s *Server) login(c *gin.Context) {
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

// 构造query条件查用户
func (s *Server) getOneUser(c *gin.Context) {
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

// 查一堆用户
func (s *Server) getSomeUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("查一堆用户......")
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

// 通过id改一个用户
func (s *Server) putUser(c *gin.Context) {
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

// 通过id删一个用户
func (s *Server) deleteUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("删用户......")
	id := c.Param("id")
	//调用service_user层
	res := s.svc.DeleteUser(id)
	//返回结果
	res.Response(c, res)
	return
}
