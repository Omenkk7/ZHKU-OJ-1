package server

import (
	"github.com/gin-gonic/gin"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
	"zhku-oj-server/pkg/utils/middleware"
)

// RegisterLabel 路由器 ——————label_manager
func (s *Server) RegisterLabel(g *gin.RouterGroup) {
	//查询操作不拦截
	userGroup := g.Group("/user")
	{
		userGroup.GET("/", s.GetSomeLabel)   //查一堆
		userGroup.GET("/:id", s.GetOneLabel) //查一个
	}
	// 对以下路由应用JWT拦截器，并且只有管理员可以执行以下操作
	securedGroup := userGroup.Group("/label").Use(middleware.JWTInterceptor())
	{
		securedGroup.PUT("/:id", s.PutLabel)       //改一个
		securedGroup.DELETE("/:id", s.DeleteLabel) //删一个
		securedGroup.POST("/", s.PostLabel)        //增一个
	}
}

// PostLabel 增加一个标签
func (s *Server) PostLabel(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("添加标签......")
	//把参数解析到postLabel
	var postLabel *dto.ReqLabel
	if err := c.BindJSON(&postLabel); err != nil {
		return
	}
	//调用service_label层
	res := s.svc.PostLabel(postLabel)
	//返回结果
	res.Response(c, res)
	return
}

// GetSomeLabel 查一堆标签
func (s *Server) GetSomeLabel(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("查一堆标签......")
	query := c.Request.URL.Query()
	cq := utils.BuildCommonQuery(utils.Query(query))
	lg.Info("get", cq)
	res, err := s.svc.GetLabelList(cq)
	if err != nil {
		lg.Errorf("getLabelList: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, res)
}

// GetOneLabel 条件查询标签
func (s *Server) GetOneLabel(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	var label *models.Label
	if err := c.BindJSON(&label); err != nil {
		return
	}
	lg.Println("条件查询标签......")
	//调用service_label层
	res := s.svc.GetOneLabel(label)
	//返回结果
	res.Response(c, res)
	return
}

// PutLabel 修改标签
func (s *Server) PutLabel(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("通过id改一个用户......")

	//把参数解析到结构体loginUser
	var label *dto.ReqLabel
	if err := c.BindJSON(&label); err != nil {
		return
	}

	//调用service_label层
	res := s.svc.UpdateLabel(label)
	//返回结果
	res.Response(c, res)
	return
}

// DeleteLabel 删除标签
func (s *Server) DeleteLabel(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("删标签......")
	id := c.Param("id")
	//调用service_user层
	res := s.svc.DeleteLabel(id)
	//返回结果
	res.Response(c, res)
	return
}
