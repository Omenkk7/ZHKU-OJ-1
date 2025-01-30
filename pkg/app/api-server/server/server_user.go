/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 18:24
@Name: handler_user.go
@Description:
*/

package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/app/api-server/service"
	"zhku-oj-server/pkg/dao"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
	"zhku-oj-server/pkg/utils/myUtils"
)

// 拦截器
func JWTInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg := utils.GetDefaultLogger()
		lg.Info("拦截请求......")

		// 从请求头中获取token
		token := c.Request.Header.Get("Authorization") //TODO 这个key好像是前端设置的，每次请求都自动携带

		// 检查token是否存在
		if token == "" {
			lg.Info("令牌为空")
			utils.SuccessResponse(c, map[string]interface{}{
				"status": utils.StatusFail,
				"mes":    "令牌为空",
				"data":   "",
			})
			c.Abort()
			return
		}

		//校验jwt是否有效
		jwtClaims, err := myUtils.ParseToken(token, utils.JwtTokenSecretKey)
		if err != nil {
			lg.Info("令牌无效：", err)
			utils.SuccessResponse(c, map[string]interface{}{
				"status": utils.StatusFail,
				"mes":    err.Error(),
				"data":   "",
			})
			c.Abort()
			return
		}

		//检查权限，判断是否为管理员
		dao := dao.NewDao()
		query := bson.M{
			"username": jwtClaims.User.Username, //旧jwt中的旧username数据
		}
		user, _ := dao.GetOneUser(context.Background(), query) //user已为最新状态
		if user.Role == utils.StatusUser {
			utils.SuccessResponse(c, map[string]interface{}{
				"status": utils.StatusFail,
				"mes":    "没有权限，拒绝访问",
				"data":   "",
			})
			c.Abort()
			return
		}

		//校验成功，解析并拿到jwt的用户数据，存进gin.Context，可通过c.Get("user")重新获得数据
		c.Set("User", user)
		c.Next()
	}
}

// 路由器
func (s *Server) RegisterUser(g *gin.RouterGroup) {
	userGroup := g.Group("/user")
	{
		userGroup.POST("/", s.postUser)   // 注册
		userGroup.POST("/login", s.login) // 登录
	}

	// 对以下路由应用JWT拦截器，并且只有管理员可以执行以下操作
	securedGroup := userGroup.Group("/")
	securedGroup.Use(JWTInterceptor())
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
	status, mes := service.PostUser(postUser)
	//返回结果
	utils.SuccessResponse(c, map[string]interface{}{
		"status": status,
		"mes":    mes,
		"data":   "",
	})
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
	status, token, mes := service.UserLogin(loginUser)
	//返回结果
	utils.SuccessResponse(c, map[string]interface{}{
		"status": status,
		"mes":    mes,
		"data":   token,
	})
	return
}

// 通过id查用户
func (s *Server) getOneUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	id := c.Param("id")
	lg.Println("通过id查用户......", id)
	//调用service_user层
	status, mes, data := service.GetUserById(id)
	//返回结果
	utils.SuccessResponse(c, map[string]interface{}{
		"status": status,
		"mes":    mes,
		"data":   data,
	})
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
	status, mes, data := service.UpdateUser(user)
	//返回结果
	utils.SuccessResponse(c, map[string]interface{}{
		"status": status,
		"mes":    mes,
		"data":   data,
	})
	return
}

// 通过id删一个用户
func (s *Server) deleteUser(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("删用户......")
	id := c.Param("id")
	//调用service_user层
	status, mes, data := service.DeleteUser(id)
	//返回结果
	utils.SuccessResponse(c, map[string]interface{}{
		"status": status,
		"mes":    mes,
		"data":   data,
	})
	return
}

//原登录逻辑
/*func (s *Server) login(c *gin.Context) {
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

	//判断账号或密码是否为空
	if err := s.svc.CheckUserLoginParams(loginUser); err != nil {
		lg.Errorf("CheckUserLoginParams: %v", err)
		utils.BadRequest(c, err)
		return
	}
	//账号密码都不为空，查询账号密是否正确
	token, err := s.svc.UserLogin(loginUser)
	if err != nil {
		lg.Errorf("UserLogin: %v", err)
		utils.BadRequest(c, err)
		return
	}
	//账号密码正确，返回token，登录成功
	utils.SuccessResponse(c, map[string]interface{}{
		"token": token,
	})
}*/

//原注册逻辑
/*func (s *Server) postUser(c *gin.Context) {
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
*/
