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

// PostUser 注册 /user
func (s *Server) PostUser(c *gin.Context) {
	// 打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("注册......")

	// 获取操作者角色
	var operatorRole int32 = 4 // 默认为学生角色

	// 从上下文中获取用户信息
	if contextUser, exists := c.Get("contextUser"); exists {
		if user, ok := contextUser.(*utils.ContextUser); ok {
			operatorRole = int32(user.Role)
			lg.Infof("操作者角色: %d", operatorRole)
		}
	}

	// 把参数解析到postUser
	var postUser *dto.ReqPostUser
	if err := c.BindJSON(&postUser); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 如果未指定角色，默认为学生角色
	if postUser.Role == 0 {
		postUser.Role = 4 // 默认为学生
	}

	// 调用service_user层
	id, err := s.svc.PostUser(postUser, operatorRole)

	// 返回结果
	if err != nil {
		lg.Errorf("注册用户失败: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, gin.H{"id": id})
}

// TODO 个人信息管理已解决？ 这个方法可以解析jwt，直接查询到个人信息
func (s *Server) GetInfor(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	// 从请求头中获取token
	token := c.Request.Header.Get("Authorization") //TODO 这个key好像是前端设置的，每次请求都自动携带

	lg.Infoln("token:", token)

	// 检查token是否存在
	if token == "" {
		lg.Info(utils.JwtEmptyErr)
		utils.BadRequest(c, utils.New(utils.JwtEmptyErr))
		c.Abort()
		return
	}

	//校验jwt是否有效
	jwtClaims, err := utils.ParseToken(token, utils.JwtTokenSecretKey)
	if err != nil {
		lg.Info(utils.JwtFailErr, err)
		utils.BadRequest(c, utils.New(utils.JwtFailErr))
		c.Abort()
		return
	}

	//调用service_user层
	res, err := s.svc.GetInfor(jwtClaims.ID)
	//返回结果
	if err != nil {
		lg.Errorf("getOneUser: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, res)
}

// Login 登录 /login
func (s *Server) Login(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("登录......")
	//把参数解析到结构体loginUser
	var loginUser *dto.ReqPostLoginUser
	if err := c.BindJSON(&loginUser); err != nil {
		utils.BadRequest(c, err)
		return
	}
	//调用service_user层
	response, err := s.svc.UserLogin(loginUser)
	//返回结果
	if err != nil {
		lg.Errorf("登录失败: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, response)
}

// GetOneUser 构造query条件查用户 /:id
func (s *Server) GetOneUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	id := c.Param("id")
	lg.Infof("id: %s", id)
	//把参数解析到结构体user
	reqUser := &dto.ReqUser{ID: id}
	lg.Println("条件查询用户......")
	//调用service_user层
	res, err := s.svc.GetOneUser(reqUser)
	//返回结果
	if err != nil {
		lg.Errorf("getOneUser: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, res)
}

// GetSomeUser 查一堆用户
func (s *Server) GetSomeUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("查一堆用户......")
	query := c.Request.URL.Query()
	cq := utils.BuildCommonQuery(utils.Query(query))
	lg.Info("get", cq)
	//调用service_user层
	res, err := s.svc.GetUserList(cq)
	//返回结果
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

	// 获取操作者角色
	var operatorRole int32 = 4 // 默认为学生角色

	// 从上下文中获取用户信息
	if contextUser, exists := c.Get("contextUser"); exists {
		if user, ok := contextUser.(*utils.ContextUser); ok {
			operatorRole = int32(user.Role)
			lg.Infof("操作者角色: %d", operatorRole)
		}
	}

	// 把参数解析到结构体user
	var reqUser *dto.ReqUser
	if err := c.BindJSON(&reqUser); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service_user层
	res, err := s.svc.UpdateUser(reqUser, operatorRole)

	// 返回结果
	if err != nil {
		lg.Errorf("更新用户失败: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, res)
}

// DeleteUser 通过id删一个用户 /:id
func (s *Server) DeleteUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("通过id删一个用户......")

	// 获取操作者角色
	var operatorRole int32 = 4 // 默认为学生角色

	// 从上下文中获取用户信息
	if contextUser, exists := c.Get("contextUser"); exists {
		if user, ok := contextUser.(*utils.ContextUser); ok {
			operatorRole = int32(user.Role)
			lg.Infof("操作者角色: %d", operatorRole)
		}
	}

	// 获取要删除的用户ID
	id := c.Param("id")

	// 调用service_user层
	res, err := s.svc.DeleteUser(id, operatorRole)

	// 返回结果
	if err != nil {
		lg.Errorf("删除用户失败: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, res)
}
