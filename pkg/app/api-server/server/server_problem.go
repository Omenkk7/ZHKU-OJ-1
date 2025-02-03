package server

import (
	"github.com/gin-gonic/gin"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// GetSomeProblem 查一堆题目
func (s *Server) GetSomeProblem(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("查一堆题目......")
	query := c.Request.URL.Query()
	cq := utils.BuildCommonQuery(utils.Query(query))
	lg.Info("get", cq)
	res, err := s.svc.GetProblemList(cq)
	if err != nil {
		lg.Errorf("getProbelmList: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c, res)
}

// GetOneProblem 条件查询查一个题目
func (s *Server) GetOneProblem(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	var problem *models.Problem
	if err := c.BindJSON(&problem); err != nil {
		return
	}
	lg.Println("条件查询题目......")
	//调用service_problem层
	s.svc.GetOneProblem(s.res, problem)
	//返回结果
	s.res.Response(c, s.res)
	return
}

// PutProblem 通过id改题目
func (s *Server) PutProblem(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("通过id改一个题目......")

	//把参数解析到结构体ReqProblem
	var problem *dto.ReqProblem
	if err := c.BindJSON(&problem); err != nil {
		return
	}

	//调用service_problem层
	s.svc.UpdateProblem(s.res, problem)
	//返回结果
	s.res.Response(c, s.res)
	return
}

// DeleteProblem 通过id删除题目
func (s *Server) DeleteProblem(c *gin.Context) {
	lg := utils.GetDefaultLogger()
	lg.Info("删除题目......")
	id := c.Param("id")
	//调用service_problem层
	s.svc.DeleteProblem(s.res, id)
	//返回结果
	s.res.Response(c, s.res)
	return
}

// PostProblem 添加题目
func (s *Server) PostProblem(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("添加题目......")
	//把参数解析到ReqProblem
	var reqProblem *dto.ReqProblem
	if err := c.BindJSON(&reqProblem); err != nil {
		return
	}
	//调用service_problem层
	s.svc.PostProblem(s.res, reqProblem)
	//返回结果
	s.res.Response(c, s.res)
	return
}
