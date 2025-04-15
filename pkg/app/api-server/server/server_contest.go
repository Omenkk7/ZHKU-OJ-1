package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"zhku-oj-server/pkg/app/api-server/dto"

	"zhku-oj-server/pkg/utils"
)

// CreateContest 创建竞赛
func (s *Server) CreateContest(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 获取用户名
	userName, exists := c.Get("userName")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求
	var req dto.CreateContestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		lg := utils.GetDefaultLogger()
		lg.Errorf("解析创建竞赛请求失败: %v", err)
		utils.BadRequest(c, err)
		return
	}

	// 添加日志记录请求内容
	lg := utils.GetDefaultLogger()
	lg.Infof("创建竞赛请求内容: %+v", req)

	// 确保 access_type 字段正确设置
	lg.Infof("竞赛访问类型: %d", req.AccessType)

	// 创建竞赛
	id, err := s.svc.CreateContest(c, &req, userID.(string), userName.(string))
	if err != nil {
		lg.Errorf("创建竞赛失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, gin.H{"id": id})
}

// GetContestDetail 获取竞赛详情
func (s *Server) GetContestDetail(c *gin.Context) {
	// 获取竞赛ID
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, errors.New("竞赛ID不能为空"))
		return
	}

	// 获取竞赛详情
	contest, err := s.svc.GetContestByID(c, id)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, contest)
}

// GetContestList 获取竞赛列表
func (s *Server) GetContestList(c *gin.Context) {
	// 获取用户信息
	var userID string
	var userRole int

	// 从上下文中获取用户信息
	if contextUser, exists := c.Get("contextUser"); exists {
		if user, ok := contextUser.(*utils.ContextUser); ok {
			userID = user.ID
			userRole = user.Role
		}
	}

	// 解析请求参数
	var req dto.GetContestListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取竞赛列表
	resp, err := s.svc.GetContestList(c, &req, userID, userRole)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// UpdateContest 更新竞赛
func (s *Server) UpdateContest(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 从URL路径参数获取竞赛ID
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, errors.New("竞赛ID不能为空"))
		return
	}

	// 解析请求
	var req dto.UpdateContestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 更新竞赛
	err := s.svc.UpdateContest(c, id, &req, userID.(string))
	if err != nil {
		lg := utils.GetDefaultLogger()
		lg.Errorf("更新竞赛错误: %v", err)
		switch err.Error() {
		case "mongo: no documents in result":
			utils.FailedResponse(c, http.StatusNotFound, errors.New("竞赛不存在或已被删除"))
		case "无效的竞赛ID":
			utils.BadRequest(c, err)
		case "无权限更新此竞赛":
			utils.FailedResponse(c, http.StatusForbidden, err)
		default:
			utils.FailedResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, nil)
}

// DeleteContest 删除竞赛
func (s *Server) DeleteContest(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 获取竞赛ID
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, errors.New("竞赛ID不能为空"))
		return
	}

	// 删除竞赛
	err := s.svc.DeleteContest(c, id, userID.(string))
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// ArchiveContest 归档竞赛
func (s *Server) ArchiveContest(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 获取竞赛ID
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, errors.New("竞赛ID不能为空"))
		return
	}

	// 归档竞赛
	err := s.svc.ArchiveContest(c, id, userID.(string))
	if err != nil {
		switch err.Error() {
		case "mongo: no documents in result":
			utils.FailedResponse(c, http.StatusNotFound, errors.New("竞赛不存在或已被删除"))
		case "无权限归档此竞赛":
			utils.FailedResponse(c, http.StatusForbidden, err)
		default:
			utils.FailedResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, nil)
}

// UpdateContestStatus 更新竞赛状态
func (s *Server) UpdateContestStatus(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 获取竞赛ID
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, errors.New("竞赛ID不能为空"))
		return
	}

	// 解析请求
	var req struct {
		Status int `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 更新竞赛状态
	err := s.svc.UpdateContestStatus(c, id, req.Status, userID.(string))
	if err != nil {
		switch err.Error() {
		case "mongo: no documents in result":
			utils.FailedResponse(c, http.StatusNotFound, errors.New("竞赛不存在或已被删除"))
		case "无效的竞赛ID":
			utils.BadRequest(c, err)
		case "无权限更新此竞赛状态":
			utils.FailedResponse(c, http.StatusForbidden, err)
		case "无效的竞赛状态":
			utils.BadRequest(c, err)
		default:
			utils.FailedResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, nil)
}

// AddParticipant 添加参赛者
func (s *Server) AddParticipant(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求
	var req dto.AddParticipantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 添加参赛者
	err := s.svc.AddParticipant(c, &req, userID.(string))
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// BatchAddParticipants 批量添加参赛者
func (s *Server) BatchAddParticipants(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求
	var req dto.BatchAddParticipantsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 批量添加参赛者
	err := s.svc.BatchAddParticipants(c, &req, userID.(string))
	if err != nil {
		switch err.Error() {
		case "竞赛不存在或已被删除":
			utils.FailedResponse(c, http.StatusNotFound, err)
		case "无效的竞赛ID":
			utils.BadRequest(c, err)
		case "只能在未开始的竞赛添加参赛者":
			utils.BadRequest(c, err)
		case "参赛人数将超过上限":
			utils.BadRequest(c, err)
		case "无权限批量添加参赛者":
			utils.FailedResponse(c, http.StatusForbidden, err)
		default:
			utils.FailedResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, nil)
}

// RemoveParticipant 移除参赛者
func (s *Server) RemoveParticipant(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求
	var req dto.RemoveParticipantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 移除参赛者
	err := s.svc.RemoveParticipant(c, &req, userID.(string))
	if err != nil {
		switch err.Error() {
		case "参赛者不存在或已被删除":
			utils.FailedResponse(c, http.StatusNotFound, err)
		case "无效的竞赛ID":
			utils.BadRequest(c, err)
		case "只能在未开始的竞赛移除参赛者":
			utils.BadRequest(c, err)
		case "无权限移除参赛者":
			utils.FailedResponse(c, http.StatusForbidden, err)
		default:
			utils.FailedResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, nil)
}

// GetParticipantList 获取参赛者列表
func (s *Server) GetParticipantList(c *gin.Context) {
	// 解析请求参数
	var req dto.GetParticipantListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取参赛者列表
	resp, err := s.svc.GetParticipantList(c, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// AuditParticipant 审核参赛者
func (s *Server) AuditParticipant(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求
	var req dto.AuditParticipantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 审核参赛者
	err := s.svc.AuditParticipant(c, &req, userID.(string))
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// GetContestRanking 获取竞赛排名
func (s *Server) GetContestRanking(c *gin.Context) {
	// 获取竞赛ID
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, errors.New("竞赛ID不能为空"))
		return
	}

	req := dto.GetContestRankingReq{
		ContestID: id,
	}

	// 解析其他请求参数
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取竞赛排名
	resp, err := s.svc.GetContestRanking(c, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// ExportContestScore 导出竞赛成绩
func (s *Server) ExportContestScore(c *gin.Context) {
	// 获取竞赛ID
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, errors.New("竞赛ID不能为空"))
		return
	}

	// 构建请求
	req := dto.ExportContestScoreReq{
		ContestID: id,
	}

	// 导出竞赛成绩
	data, err := s.svc.ExportContestScore(c, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename=contest_score_"+id+".csv")
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Length", strconv.Itoa(len(data)))

	// 写入响应体
	_, err = c.Writer.Write(data)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, errors.New("写入响应失败"))
		return
	}
}

// ApplyJoinContest 学生申请加入竞赛
func (s *Server) ApplyJoinContest(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 获取用户名
	userName, exists := c.Get("userName")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求
	var req dto.ApplyJoinContestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 添加日志记录
	lg := utils.GetDefaultLogger()
	lg.Infof("学生申请加入竞赛，学生ID: %s, 竞赛ID: %s", userID.(string), req.ContestID)

	// 申请加入竞赛
	err := s.svc.ApplyJoinContest(c, &req, userID.(string), userName.(string))
	if err != nil {
		lg.Errorf("申请加入竞赛失败: %v", err)
		switch err.Error() {
		case "无效的竞赛ID":
			utils.BadRequest(c, err)
		case "竞赛不存在或已被删除":
			utils.FailedResponse(c, http.StatusNotFound, err)
		case "您已经是该竞赛的参赛者", "您已经申请过该竞赛，请等待审核":
			utils.BadRequest(c, err)
		case "只能申请加入私有竞赛", "竞赛已开始或已结束，无法申请加入":
			utils.BadRequest(c, err)
		default:
			utils.FailedResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "申请已提交，请等待审核"})
}
