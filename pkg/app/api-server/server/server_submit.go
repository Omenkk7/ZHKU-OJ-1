package server

import (
	"github.com/gin-gonic/gin"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/utils"
)

func (s *Server) Submit(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("提交代码......")
	//把参数解析到reqSubmit
	var reqSubmit *dto.ReqSubmit
	if err := c.BindJSON(&reqSubmit); err != nil {
		return
	}
	//调用service_submit层
	err := s.svc.Submit(reqSubmit)
	//返回结果
	if err != nil {
		lg.Errorf("Submit: %v", err)
		utils.BadRequest(c, err)
		return
	}
	utils.SuccessResponse(c)
}
