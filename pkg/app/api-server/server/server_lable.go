package server

import (
	"github.com/gin-gonic/gin"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

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
	s.svc.PostLabel(s.res, postLabel)
	//返回结果
	s.res.Response(c, s.res)
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
	//把参数解析结构体到label
	var label *models.Label
	if err := c.BindJSON(&label); err != nil {
		return
	}
	lg.Println("条件查询标签......")
	//调用service_label层
	s.svc.GetOneLabel(s.res, label)
	//返回结果
	s.res.Response(c, s.res)
	return
}

// PutLabel 修改标签
func (s *Server) PutLabel(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("通过id改一个用户......")

	//把参数解析结构体ReqLabel到
	var label *dto.ReqLabel
	if err := c.BindJSON(&label); err != nil {
		return
	}

	//调用service_label层
	s.svc.UpdateLabel(s.res, label)
	//返回结果
	s.res.Response(c, s.res)
	return
}

// DeleteLabel 删除标签
func (s *Server) DeleteLabel(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("删标签......")
	id := c.Param("id")
	//调用service_label层
	s.svc.DeleteLabel(s.res, id)
	//返回结果
	s.res.Response(c, s.res)
	return
}
