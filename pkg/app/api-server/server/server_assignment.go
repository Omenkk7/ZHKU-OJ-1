package server

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// CreateAssignment 创建作业
func (s *Server) CreateAssignment(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")

	log.Println(userID)
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}
	userName, exists := c.Get("userName")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求参数
	var req dto.CreateAssignmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("参数解析失败: %v", err)
		utils.BadRequest(c, errors.New("参数无效"))
		return
	}

	// 创建作业
	id, err := s.svc.CreateAssignment(c, &req, userID.(string), userName.(string))
	if err != nil {
		log.Printf("创建作业失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, gin.H{"id": id})
}

// UpdateAssignment 更新作业
func (s *Server) UpdateAssignment(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求参数
	var req dto.UpdateAssignmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("参数解析失败: %v", err)
		utils.BadRequest(c, fmt.Errorf("参数无效: %v", err))
		return
	}

	// 获取作业ID
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, errors.New("作业ID参数无效"))
		return
	}

	// 更新作业
	err := s.svc.UpdateAssignment(c, id, &req, userID.(string))
	if err != nil {
		log.Printf("更新作业失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, nil)
}

// DeleteAssignment 删除作业
func (s *Server) DeleteAssignment(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 获取作业ID
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, errors.New("参数无效"))
		return
	}

	// 删除作业
	err := s.svc.DeleteAssignment(c, id, userID.(string))
	if err != nil {
		log.Printf("删除作业失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, nil)
}

// ArchiveAssignment 归档作业
func (s *Server) ArchiveAssignment(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 获取作业ID
	id := c.Query("id")
	if id == "" {
		utils.BadRequest(c, errors.New("参数无效"))
		return
	}

	// 归档作业
	err := s.svc.ArchiveAssignment(c, id, userID.(string))
	if err != nil {
		log.Printf("归档作业失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, nil)
}

// GetAssignmentDetail 获取作业详情
func (s *Server) GetAssignmentDetail(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 获取作业ID - 从路径参数获取
	id := c.Param("id")
	if id == "" {
		utils.BadRequest(c, errors.New("作业ID参数无效"))
		return
	}

	// 获取作业详情
	assignment, err := s.svc.GetAssignmentDetail(c, id, userID.(string))
	if err != nil {
		log.Printf("获取作业详情失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, assignment)
}

// GetAssignmentList 获取作业列表
func (s *Server) GetAssignmentList(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求参数
	var req dto.GetAssignmentListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("参数解析失败: %v", err)
		utils.BadRequest(c, errors.New("参数无效"))
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 修改 status 参数处理逻辑
	statusStr := c.Query("status")
	if statusStr == "" {
		// 当 status 参数为空时，默认设置为 1（进行中）
		defaultStatus := 1
		req.Status = &defaultStatus
	} else if statusStr != "" {
		status, err := strconv.Atoi(statusStr)
		if err == nil {
			req.Status = &status
		}
	}

	// 获取作业列表
	result, err := s.svc.GetAssignmentList(c, &req, userID.(string))
	if err != nil {
		log.Printf("获取作业列表失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, result)
}

// GetStudentAssignments 获取学生作业列表
func (s *Server) GetStudentAssignments(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求参数
	var req dto.GetStudentAssignmentsReq
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("参数解析失败: %v", err)
		utils.BadRequest(c, errors.New("参数无效"))
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取学生作业列表
	result, err := s.svc.GetStudentAssignments(c, &req, userID.(string))
	if err != nil {
		log.Printf("获取学生作业列表失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, result)
}

// SubmitAssignment 提交作业
func (s *Server) SubmitAssignment(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 添加用户角色日志
	log.Printf("用户 %s 尝试提交作业", userID)

	// 解析请求参数
	assignmentID := c.Param("id")
	var req dto.SubmitAssignmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("参数解析失败: %v", err)
		utils.BadRequest(c, errors.New("参数无效"))
		return
	}

	// 提交作业
	err := s.svc.SubmitAssignment(c, assignmentID, &req, userID.(string))
	if err != nil {
		log.Printf("提交作业失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, nil)
}

// GradeAssignment 批改作业
func (s *Server) GradeAssignment(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}
	userName, exists := c.Get("userName")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求参数
	var req dto.GradeAssignmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("参数解析失败: %v", err)
		utils.BadRequest(c, errors.New("参数无效"))
		return
	}

	// 批改作业 - 不再需要从URL获取ID，直接使用req中的AssignmentID
	err := s.svc.GradeAssignment(c, &req, userID.(string), userName.(string))
	if err != nil {
		log.Printf("批改作业失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, nil)
}

// RejectAssignment 打回作业
func (s *Server) RejectAssignment(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}
	userName, exists := c.Get("userName")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求参数
	var req dto.RejectAssignmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("参数解析失败: %v", err)
		utils.BadRequest(c, errors.New("参数无效"))
		return
	}

	// 打回作业
	err := s.svc.RejectAssignment(c, &req, userID.(string), userName.(string))
	if err != nil {
		log.Printf("打回作业失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, nil)
}

// GetStudentSubmission 获取学生提交信息
func (s *Server) GetStudentSubmission(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求参数
	var req dto.GetStudentSubmissionReq
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Printf("参数解析失败: %v", err)
		utils.BadRequest(c, errors.New("参数无效"))
		return
	}

	// 验证必要参数
	if req.AssignmentID == "" || req.StudentID == "" {
		utils.BadRequest(c, errors.New("作业ID和学生ID不能为空"))
		return
	}

	// 设置默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	// 获取学生提交信息
	submission, err := s.svc.GetStudentSubmission(c, &req, userID.(string))
	if err != nil {
		log.Printf("获取作业详情失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, submission)
}

// ExportAssignmentGrades 导出作业成绩
func (s *Server) ExportAssignmentGrades(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 获取参数 - 从路径参数获取
	assignmentID := c.Param("id")
	if assignmentID == "" {
		utils.BadRequest(c, errors.New("参数无效"))
		return
	}

	// 导出成绩
	grades, err := s.svc.ExportAssignmentGrades(c, assignmentID, userID.(string))
	if err != nil {
		log.Printf("导出成绩失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, grades)
}

// AddStudentsToAssignment 添加学生到作业
func (s *Server) AddStudentsToAssignment(c *gin.Context) {
	// 获取用户信息
	userID, exists := c.Get("userId")
	if !exists {
		utils.Unauthorized(c, errors.New("需要登录"))
		return
	}

	// 解析请求参数
	var req dto.AddStudentsToAssignmentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("参数解析失败: %v", err)
		utils.BadRequest(c, errors.New("参数无效"))
		return
	}

	// 转换为服务层需要的格式
	students := make([]models.Member, 0, len(req.Students))
	for _, student := range req.Students {
		students = append(students, models.Member{
			UserID:   student.StudentID,
			UserName: student.StudentName,
		})
	}

	// 添加学生
	err := s.svc.AddStudentsToAssignmentDirect(c, req.AssignmentID, students, userID.(string))
	if err != nil {
		log.Printf("添加学生失败: %v", err)
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 返回成功
	utils.SuccessResponse(c, nil)
}
