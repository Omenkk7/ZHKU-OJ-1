package server

import (
	"github.com/gin-gonic/gin"
	"zhku-oj-server/pkg/utils"
)

// Chat 发起对话 /chat-messages
func (s *Server) Chat(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("发起对话......")
	/*//把参数解析到结构体loginUser
	var loginUser *dto.ReqPostLoginUser
	if err := c.BindJSON(&loginUser); err != nil {
		return
	}
	//调用service_user层
	res, err := s.svc.UserLogin(loginUser)
	//返回结果
	if err != nil {
		lg.Errorf("getUserList: %v", err)
		utils.BadRequest(c, err)
		return
	}*/
	utils.SuccessResponse(c, "成功")
}

// GetMessage 查询对话的信息 /message
func (s *Server) GetMessage(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("查询对话的信息......")
	utils.SuccessResponse(c, "成功")
}

// GetConversations 查询对话记录 /conversations
func (s *Server) GetConversations(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("查询对话记录......")
	utils.SuccessResponse(c, "成功")
}

// RenameConversation 重命名对话 /conversations/:conversation_id/name
func (s *Server) RenameConversation(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("重命名对话......")
	utils.SuccessResponse(c, "成功")
}

// DeleteConversation 删除对话 /conversations/:conversation_id
func (s *Server) DeleteConversation(c *gin.Context) {
	//打印日志
	lg := utils.GetDefaultLogger()
	lg.Info("删除对话......")
	utils.SuccessResponse(c, "成功")
}
