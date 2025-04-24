package server

import (
	"errors"
	"fmt"
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

// AddProblemsToContest 添加题目到竞赛
func (s *Server) AddProblemsToContest(c *gin.Context) {
	lg := utils.GetDefaultLogger()

	//  获取竞赛ID
	contestID := c.Param("id")
	if contestID == "" {
		lg.Warn("请求缺少竞赛ID")
		utils.BadRequest(c, errors.New("竞赛ID不能为空"))
		return
	}

	//  获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		lg.Warn("无法获取用户ID，需要登录")
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}
	userRoleStr, exists := c.Get("userRole")
	if !exists {
		lg.Warn("无法获取用户角色，需要登录")
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}
	userRole, err := strconv.Atoi(userRoleStr.(string))
	if err != nil {
		lg.Errorf("用户角色转换失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, errors.New("无效的用户角色"))
		return
	}

	//  解析请求体
	var req dto.AddProblemsToContestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Errorf("解析添加题目请求失败: %v", err)
		utils.BadRequest(c, err)
		return
	}

	// 添加额外验证：确保每个题目都有ID和标题
	for i, problem := range req.Problems {
		if problem.ProblemID == "" {
			lg.Errorf("第 %d 个题目缺少题目ID", i+1)
			utils.BadRequest(c, fmt.Errorf("第 %d 个题目缺少题目ID", i+1))
			return
		}
		if problem.Title == "" {
			lg.Errorf("第 %d 个题目缺少题目名称", i+1)
			utils.BadRequest(c, fmt.Errorf("第 %d 个题目缺少题目名称", i+1))
			return
		}
	}

	//  调用 Service 处理
	err = s.svc.AddProblemsToContest(c, contestID, &req, userID.(string), userRole)
	if err != nil {
		lg.Errorf("调用 Service 添加题目失败: %v", err)
		// 根据错误类型返回不同的状态码
		if errors.Is(err, utils.ErrNoPermission) {
			utils.Forbidden(c, err)
		} else if err.Error() == "竞赛不存在" || err.Error() == "无效的竞赛ID格式" || err.Error() == "题目列表不能为空" || err.Error() == "无法向已结束或已归档的竞赛添加题目" { // 匹配一些预期的业务错误
			utils.BadRequest(c, err)
		} else {
			utils.FailedResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	//  返回成功响应
	lg.Infof("成功处理向竞赛 %s 添加题目的请求", contestID)
	utils.SuccessResponse(c, gin.H{"message": "题目添加成功"})
}

// BatchRemoveProblemsFromContest 批量从竞赛移除题目 (新增)
func (s *Server) BatchRemoveProblemsFromContest(c *gin.Context) {
	lg := utils.GetDefaultLogger()

	//  获取竞赛ID
	contestID := c.Param("id")
	if contestID == "" {
		lg.Warn("请求缺少竞赛ID")
		utils.BadRequest(c, errors.New("竞赛ID不能为空"))
		return
	}

	//  获取用户信息
	userIDVal, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}
	userID := userIDVal.(string)

	userRoleVal, exists := c.Get("userRole")
	if !exists {
		utils.Unauthorized(c, errors.New("无法获取用户角色信息"))
		return
	}

	// 转换用户角色
	var userRole int
	switch v := userRoleVal.(type) {
	case int:
		userRole = v
	case float64:
		userRole = int(v)
	case string:
		var err error
		userRole, err = strconv.Atoi(v)
		if err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, errors.New("用户角色格式无效"))
			return
		}
	default:
		utils.FailedResponse(c, http.StatusInternalServerError, errors.New("无法识别的用户角色类型"))
		return
	}

	//  解析请求体
	var req dto.BatchRemoveProblemsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		lg.Errorf("解析批量移除题目请求失败: %v", err)
		utils.BadRequest(c, err) // 返回具体的绑定错误信息
		return
	}

	// 调用 Service 处理
	err := s.svc.BatchRemoveProblemsFromContest(c.Request.Context(), contestID, &req, userID, userRole)
	if err != nil {
		lg.Errorf("Service层从竞赛 %s 移除题目失败: %v", contestID, err)
		// 根据 Service 返回的错误类型进行不同的响应
		if errors.Is(err, utils.ErrNoPermission) {
			utils.Forbidden(c, err)
		} else if err.Error() == "竞赛不存在" {
			// 使用 FailedResponse 替代 NotFound
			utils.FailedResponse(c, http.StatusNotFound, err)
		} else if err.Error() == "无效的竞赛ID" || err.Error() == "要移除的题目ID列表不能为空" {
			utils.BadRequest(c, err)
		} else {
			utils.FailedResponse(c, http.StatusInternalServerError, errors.New("移除题目时发生内部错误"))
		}
		return
	}

	//  成功响应
	utils.SuccessResponse(c)
}
