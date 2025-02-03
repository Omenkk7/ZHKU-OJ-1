package server

import (
	"github.com/gin-gonic/gin"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/utils"
)

// PostLabel 增加一个标签
func (s *Server) PostLabel(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("添加标签......")
	//把参数解析到postLabel
	var reqLabel *dto.ReqLabel
	if err := c.BindJSON(&reqLabel); err != nil {
		return
	}
	//调用service_label层
	err := s.svc.PostLabel(reqLabel)
	//返回结果
	if err != nil {
		lg.Errorf("getUserList: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c)
}

// GetSomeLabel 查一堆标签
func (s *Server) GetSomeLabel(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("查一堆标签......")
	query := c.Request.URL.Query()
	cq := utils.BuildCommonQuery(utils.Query(query))
	lg.Info("get", cq)
	//调用service_label层
	res, err := s.svc.GetLabelList(cq)
	//返回结果
	if err != nil {
		lg.Errorf("getUserList: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, res)
}

// GetOneLabel 条件查询标签
func (s *Server) GetOneLabel(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	//把参数解析结构体到label
	var label *dto.ReqLabel
	if err := c.BindJSON(&label); err != nil {
		return
	}
	lg.Println("条件查询标签......")
	//调用service_label层
	res, err := s.svc.GetOneLabel(label)
	//返回结果
	if err != nil {
		lg.Errorf("getUserList: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, res)
}

// PutLabel 修改标签
func (s *Server) PutLabel(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("通过id改一个用户......")

	//把参数解析结构体到reqLabel
	var reqLabel *dto.ReqLabel
	if err := c.BindJSON(&reqLabel); err != nil {
		return
	}

	//调用service_label层
	res, err := s.svc.UpdateLabel(reqLabel)
	//返回结果
	if err != nil {
		lg.Errorf("getUserList: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, res)
}

// DeleteLabel 删除标签
func (s *Server) DeleteLabel(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("删标签......")
	id := c.Param("id")
	//调用service_label层
	res, err := s.svc.DeleteLabel(id)
	//返回结果
	if err != nil {
		lg.Errorf("getUserList: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, res)
}
