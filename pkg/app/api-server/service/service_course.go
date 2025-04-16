package service

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"math/rand"
	"strconv"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// CreateCourse 创建课程
func (s *Service) CreateCourse(c *gin.Context, req *dto.CreateCourseRequest, creatorID string) (*dto.CourseResponse, error) {
	// 检查课程代码是否已存在
	existCourse, err := s.dao.GetCourseByCode(c, req.CourseCode)
	if err != nil {
		return nil, err
	}
	if existCourse != nil {
		return nil, errors.New("课程代码已存在")
	}

	// 创建课程对象
	course := &models.Course{
		CourseCode:   req.CourseCode,
		Name:         req.Name,
		Department:   req.Department,
		Description:  req.Description,
		Creator:      creatorID,
		Members:      make(map[string][]models.CourseMember),
		StudentCount: 0,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Status:       utils.CourseStatusOpen, // 默认为开放状态
	}

	// 初始化成员映射
	course.Members = map[string][]models.CourseMember{
		"1": {{
			UserID:   creatorID,
			UserName: "管理员", // 这里应该从用户服务获取用户名
			Status:   1,
			JoinTime: time.Now().Unix(),
		}},
		"2": {},
		"3": {},
		"4": {},
	}

	// 保存到数据库
	courseID, err := s.dao.CreateCourse(c, course)
	if err != nil {
		return nil, err
	}

	// 获取创建后的课程信息
	newCourse, err := s.dao.GetCourseByID(c, courseID)
	if err != nil {
		return nil, err
	}

	// 转换为响应对象
	return s.convertCourseToResponse(newCourse), nil
}

// UpdateCourse 更新课程
func (s *Service) UpdateCourse(c *gin.Context, courseID string, req *dto.UpdateCourseRequest) (*dto.CourseResponse, error) {
	// 获取课程信息
	course, err := s.dao.GetCourseByID(c, courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("课程不存在")
	}

	// 构建更新内容
	update := bson.M{"$set": bson.M{}}
	updateSet := update["$set"].(bson.M)

	// 更新字段（只更新非空字段）
	if req.Name != "" {
		updateSet["name"] = req.Name
	}
	if req.Department != "" {
		updateSet["department"] = req.Department
	}
	if req.Description != "" {
		updateSet["description"] = req.Description
	}
	if req.StartTime > 0 {
		updateSet["start_time"] = req.StartTime
	}
	if req.EndTime > 0 {
		updateSet["end_time"] = req.EndTime
	}
	if req.Status > -1 { // 0是有效值，所以用-1判断
		updateSet["status"] = req.Status
	}

	// 如果没有要更新的字段，直接返回
	if len(updateSet) == 0 {
		return s.convertCourseToResponse(course), nil
	}

	// 更新数据库
	err = s.dao.UpdateCourse(c, courseID, update)
	if err != nil {
		return nil, err
	}

	// 获取更新后的课程信息
	updatedCourse, err := s.dao.GetCourseByID(c, courseID)
	if err != nil {
		return nil, err
	}

	// 转换为响应对象
	return s.convertCourseToResponse(updatedCourse), nil
}

// DeleteCourse 删除课程
func (s *Service) DeleteCourse(c *gin.Context, courseID string) error {
	// 获取课程信息
	course, err := s.dao.GetCourseByID(c, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("课程不存在")
	}

	// TODO: 检查课程是否有关联的班级，如果有则不允许删除
	// 这部分需要与班级管理模块集成

	// 删除课程
	return s.dao.DeleteCourse(c, courseID)
}

// ArchiveCourse 归档课程
func (s *Service) ArchiveCourse(c *gin.Context, courseID string) error {
	// 获取课程信息
	course, err := s.dao.GetCourseByID(c, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("课程不存在")
	}

	// 归档课程（将状态设为已结课）
	return s.dao.ArchiveCourse(c, courseID)
}

// GetCourseByID 根据ID获取课程
func (s *Service) GetCourseByID(c *gin.Context, courseID string) (*dto.CourseResponse, error) {
	// 获取课程信息
	course, err := s.dao.GetCourseByID(c, courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("课程不存在")
	}

	// 转换为响应对象
	return s.convertCourseToResponse(course), nil
}

// GetCourseByCode 根据课程代码获取课程
func (s *Service) GetCourseByCode(c *gin.Context, courseCode string) (*dto.CourseResponse, error) {
	// 获取课程信息
	course, err := s.dao.GetCourseByCode(c, courseCode)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("课程不存在")
	}

	// 转换为响应对象
	return s.convertCourseToResponse(course), nil
}

// GetCourseList 获取课程列表
func (s *Service) GetCourseList(c *gin.Context, page, pageSize int, filters, sorts map[string]interface{}) (*utils.RespPageQuery, error) {
	// 构建查询条件
	comQuery := utils.BuildQueryParamsFromValues(page, pageSize, filters, sorts)

	// 查询课程列表
	return s.dao.GetCourseList(c, comQuery)
}

// GetUserCourses 获取用户参与的课程列表
func (s *Service) GetUserCourses(c *gin.Context, userID string, page, pageSize int, filters, sorts map[string]interface{}) (*utils.RespPageQuery, error) {
	// 构建查询条件
	comQuery := utils.BuildQueryParamsFromValues(page, pageSize, filters, sorts)

	// 查询用户参与的课程列表
	return s.dao.GetCoursesByUserID(c, userID, comQuery)
}

// AddCourseMember 添加课程成员
func (s *Service) AddCourseMember(c *gin.Context, courseID string, req *dto.AddCourseMemberRequest) error {
	// 获取课程信息
	course, err := s.dao.GetCourseByID(c, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("课程不存在")
	}

	// 检查角色是否有效
	if req.Role < 1 || req.Role > 4 {
		return errors.New("无效的角色")
	}

	// 检查用户是否已经是该角色的成员
	roleStr := strconv.Itoa(req.Role)
	if members, ok := course.Members[roleStr]; ok {
		for _, member := range members {
			if member.UserID == req.UserID && member.Status == 1 {
				return errors.New("用户已经是该角色的成员")
			}
		}
	}

	// 创建成员对象
	member := models.CourseMember{
		UserID:        req.UserID,
		UserName:      req.UserName,
		Email:         req.Email,
		StudentNumber: req.StudentNumber,
		Status:        1,
		JoinTime:      time.Now().Unix(),
	}

	// 添加成员
	err = s.dao.AddCourseMember(c, courseID, member, req.Role)
	if err != nil {
		return err
	}

	// 如果添加的是学生，更新学生数量
	if req.Role == utils.RoleStudent {
		return s.dao.UpdateCourseStudentCount(c, courseID)
	}

	return nil
}

// RemoveCourseMember 移除课程成员
func (s *Service) RemoveCourseMember(c *gin.Context, courseID string, userID string, role int) error {
	// 获取课程信息
	course, err := s.dao.GetCourseByID(c, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("课程不存在")
	}

	// 检查角色是否有效
	if role < 1 || role > 4 {
		return errors.New("无效的角色")
	}

	// 移除成员
	err = s.dao.RemoveCourseMember(c, courseID, userID, role)
	if err != nil {
		return err
	}

	// 如果移除的是学生，更新学生数量
	if role == utils.RoleStudent {
		return s.dao.UpdateCourseStudentCount(c, courseID)
	}

	return nil
}

// UpdateCourseMemberStatus 更新课程成员状态
func (s *Service) UpdateCourseMemberStatus(c *gin.Context, courseID string, userID string, role int, status int) error {
	// 获取课程信息
	course, err := s.dao.GetCourseByID(c, courseID)
	if err != nil {
		return err
	}
	if course == nil {
		return errors.New("课程不存在")
	}

	// 检查角色是否有效
	if role < 1 || role > 4 {
		return errors.New("无效的角色")
	}

	// 更新成员状态
	err = s.dao.UpdateCourseMemberStatus(c, courseID, userID, role, status)
	if err != nil {
		return err
	}

	// 如果更新的是学生，更新学生数量
	if role == utils.RoleStudent {
		return s.dao.UpdateCourseStudentCount(c, courseID)
	}

	return nil
}

// GetCourseMembers 获取课程成员
func (s *Service) GetCourseMembers(c *gin.Context, courseID string) (*dto.CourseMembersResponse, error) {
	// 获取课程信息
	course, err := s.dao.GetCourseByID(c, courseID)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("课程不存在")
	}

	// 构建响应对象
	resp := &dto.CourseMembersResponse{
		Admins:     make([]dto.CourseMemberResponse, 0),
		Teachers:   make([]dto.CourseMemberResponse, 0),
		Assistants: make([]dto.CourseMemberResponse, 0),
		Students:   make([]dto.CourseMemberResponse, 0),
	}

	// 填充管理员
	if admins, ok := course.Members["1"]; ok {
		for _, admin := range admins {
			if admin.Status == 1 { // 只返回状态为正常的成员
				resp.Admins = append(resp.Admins, s.convertMemberToResponse(admin))
			}
		}
	}

	// 填充教师
	if teachers, ok := course.Members["2"]; ok {
		for _, teacher := range teachers {
			if teacher.Status == 1 {
				resp.Teachers = append(resp.Teachers, s.convertMemberToResponse(teacher))
			}
		}
	}

	// 填充助教
	if assistants, ok := course.Members["3"]; ok {
		for _, assistant := range assistants {
			if assistant.Status == 1 {
				resp.Assistants = append(resp.Assistants, s.convertMemberToResponse(assistant))
			}
		}
	}

	// 填充学生
	if students, ok := course.Members["4"]; ok {
		for _, student := range students {
			if student.Status == 1 {
				resp.Students = append(resp.Students, s.convertMemberToResponse(student))
			}
		}
	}

	return resp, nil
}

// CreateJoinCourseRequest 创建加入课程申请
func (s *Service) CreateJoinCourseRequest(c *gin.Context, req *dto.JoinCourseRequest) (*dto.JoinCourseRequestResponse, error) {
	// 获取课程信息
	course, err := s.dao.GetCourseByCode(c, req.CourseCode)
	if err != nil {
		return nil, err
	}
	if course == nil {
		return nil, errors.New("课程不存在")
	}

	// 检查课程状态
	if course.Status != utils.CourseStatusOpen {
		return nil, errors.New("课程已结课，不能申请加入")
	}

	// 检查学生是否已经是课程成员
	if students, ok := course.Members["4"]; ok {
		for _, student := range students {
			if student.UserID == req.StudentID && student.Status == 1 {
				return nil, errors.New("您已经是该课程的学生")
			}
		}
	}

	// 检查是否已有待审核的申请
	filters := map[string]interface{}{
		"course_id":  course.ID.Hex(),
		"student_id": req.StudentID,
		"status":     utils.CourseJoinStatusPending,
	}
	// 使用BuildQueryParamsFromValues替代直接创建CommonQuery
	comQuery := utils.BuildQueryParamsFromValues(1, 1, filters, nil)

	joinRequests, err := s.dao.GetJoinCourseRequestList(c, comQuery)
	if err != nil {
		return nil, err
	}

	if joinRequests.Total > 0 {
		return nil, errors.New("您已提交过申请，请等待审核")
	}

	// 创建申请对象
	joinRequest := &models.CourseJoinRequest{
		CourseID:      course.ID.Hex(),
		CourseCode:    req.CourseCode,
		StudentID:     req.StudentID,
		StudentName:   req.StudentName,
		StudentNumber: req.StudentNumber,
		RequestMsg:    req.RequestMsg,
		Status:        utils.CourseJoinStatusPending,
	}

	// 保存到数据库
	requestID, err := s.dao.CreateJoinCourseRequest(c, joinRequest)
	if err != nil {
		return nil, err
	}

	// 获取创建后的申请信息
	newRequest, err := s.dao.GetJoinCourseRequestByID(c, requestID)
	if err != nil {
		return nil, err
	}

	// 转换为响应对象
	return s.convertJoinRequestToResponse(newRequest), nil
}

// ReviewJoinCourseRequest 审核加入课程申请
func (s *Service) ReviewJoinCourseRequest(c *gin.Context, requestID string, reviewerID string, req *dto.ReviewJoinCourseRequest) (*dto.JoinCourseRequestResponse, error) {
	// 获取申请信息
	joinRequest, err := s.dao.GetJoinCourseRequestByID(c, requestID)
	if err != nil {
		return nil, err
	}
	if joinRequest == nil {
		return nil, errors.New("申请不存在")
	}

	// 检查申请状态
	if joinRequest.Status != utils.CourseJoinStatusPending {
		return nil, errors.New("申请已被处理")
	}

	// 检查审核状态是否有效
	if req.Status != utils.CourseJoinStatusApproved && req.Status != utils.CourseJoinStatusRejected {
		return nil, errors.New("无效的审核状态")
	}

	// 更新申请状态
	now := time.Now().Unix()
	reviewTime := now
	update := bson.M{
		"$set": bson.M{
			"status":      req.Status,
			"reviewer_id": reviewerID,
			"review_msg":  req.ReviewMsg,
			"review_time": reviewTime,
		},
	}

	err = s.dao.UpdateJoinCourseRequest(c, requestID, update)
	if err != nil {
		return nil, err
	}

	// 如果审核通过，将学生添加到课程
	if req.Status == utils.CourseJoinStatusApproved {
		// 获取课程信息
		course, err := s.dao.GetCourseByID(c, joinRequest.CourseID)
		if err != nil {
			return nil, err
		}
		if course == nil {
			return nil, errors.New("课程不存在")
		}

		// 添加学生到课程
		member := models.CourseMember{
			UserID:        joinRequest.StudentID,
			UserName:      joinRequest.StudentName,
			StudentNumber: joinRequest.StudentNumber,
			Status:        1,
			JoinTime:      now,
		}

		err = s.dao.AddCourseMember(c, joinRequest.CourseID, member, utils.RoleStudent)
		if err != nil {
			return nil, err
		}

		// 更新课程学生数量
		err = s.dao.UpdateCourseStudentCount(c, joinRequest.CourseID)
		if err != nil {
			return nil, err
		}
	}

	// 获取更新后的申请信息
	updatedRequest, err := s.dao.GetJoinCourseRequestByID(c, requestID)
	if err != nil {
		return nil, err
	}

	// 转换为响应对象
	return s.convertJoinRequestToResponse(updatedRequest), nil
}

// GetJoinCourseRequestList 获取加入课程申请列表
func (s *Service) GetJoinCourseRequestList(c *gin.Context, courseID string, page, pageSize int, filters, sorts map[string]interface{}) (*utils.RespPageQuery, error) {
	// 添加课程ID筛选条件
	if filters == nil {
		filters = make(map[string]interface{})
	}
	filters["course_id"] = courseID

	// 构建查询条件
	comQuery := utils.BuildQueryParamsFromValues(page, pageSize, filters, sorts)

	// 查询申请列表
	return s.dao.GetJoinCourseRequestList(c, comQuery)
}

// CheckUserCoursePermission 检查用户在课程中的权限
func (s *Service) CheckUserCoursePermission(c *gin.Context, courseID string, userID string, allowedRoles []int) (bool, error) {
	// 获取用户在课程中的角色
	roles, err := s.dao.CheckUserCourseRole(c, courseID, userID)
	if err != nil {
		return false, err
	}

	// 检查用户是否有允许的角色
	for _, role := range roles {
		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				return true, nil
			}
		}
	}

	return false, nil
}

// GenerateCourseCode 生成课程代码
func (s *Service) GenerateCourseCode(c *gin.Context, department string) (string, error) {
	// 获取院系缩写（假设传入的是全名，需要转换为缩写）
	deptMap := map[string]string{
		"计算机学院": "CS",
		"数学学院":  "MA",
		"物理学院":  "PH",
		"化学学院":  "CH",
		// 可以根据需要添加更多映射
	}

	deptCode, ok := deptMap[department]
	if !ok {
		// 如果没有找到映射，使用默认缩写或者直接使用传入的值
		if len(department) >= 2 {
			deptCode = department[:2]
		} else {
			deptCode = department
		}
	}

	// 获取当前年份
	year := time.Now().Format("2006")

	// 使用更现代的随机数生成方式
	// 创建一个新的随机数生成器，使用当前时间作为种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNum := r.Intn(9000) + 1000 // 生成1000-9999之间的随机数

	// 组合课程代码
	courseCode := fmt.Sprintf("%s%s%d", deptCode, year, randomNum)

	// 检查课程代码是否已存在
	existCourse, err := s.dao.GetCourseByCode(c, courseCode)
	if err != nil {
		return "", err
	}
	if existCourse != nil {
		// 如果已存在，递归调用重新生成
		return s.GenerateCourseCode(c, department)
	}

	return courseCode, nil
}

// 辅助方法：转换课程对象为响应对象
func (s *Service) convertCourseToResponse(course *models.Course) *dto.CourseResponse {
	if course == nil {
		return nil
	}

	resp := &dto.CourseResponse{
		ID:           course.ID.Hex(),
		CourseCode:   course.CourseCode,
		Name:         course.Name,
		Department:   course.Department,
		Description:  course.Description,
		Creator:      course.Creator,
		StudentCount: course.StudentCount,
		StartTime:    course.StartTime,
		EndTime:      course.EndTime,
		Status:       course.Status,
		CreateTime:   course.Ctime,
		UpdateTime:   course.Mtime,
	}

	// 转换成员信息
	resp.Members = make(map[string][]dto.CourseMemberResponse)
	for roleStr, members := range course.Members {
		respMembers := make([]dto.CourseMemberResponse, 0)
		for _, member := range members {
			if member.Status == 1 { // 只返回状态为正常的成员
				respMembers = append(respMembers, s.convertMemberToResponse(member))
			}
		}
		resp.Members[roleStr] = respMembers
	}

	return resp
}

// 辅助方法：转换成员对象为响应对象
func (s *Service) convertMemberToResponse(member models.CourseMember) dto.CourseMemberResponse {
	return dto.CourseMemberResponse{
		UserID:        member.UserID,
		UserName:      member.UserName,
		Email:         member.Email,
		StudentNumber: member.StudentNumber,
		ClassID:       member.ClassID,
		Status:        member.Status,
		JoinTime:      member.JoinTime,
	}
}

// GetJoinCourseRequestByID 根据ID获取加入课程申请
func (s *Service) GetJoinCourseRequestByID(c *gin.Context, requestID string) (*dto.JoinCourseRequestResponse, error) {
	// 获取申请信息
	joinRequest, err := s.dao.GetJoinCourseRequestByID(c, requestID)
	if err != nil {
		return nil, err
	}
	if joinRequest == nil {
		return nil, errors.New("申请不存在")
	}

	// 转换为响应对象
	return s.convertJoinRequestToResponse(joinRequest), nil
}

// 辅助方法：转换加入申请对象为响应对象
func (s *Service) convertJoinRequestToResponse(request *models.CourseJoinRequest) *dto.JoinCourseRequestResponse {
	if request == nil {
		return nil
	}

	return &dto.JoinCourseRequestResponse{
		ID:            request.ID.Hex(),
		CourseID:      request.CourseID,
		CourseCode:    request.CourseCode,
		StudentID:     request.StudentID,
		StudentName:   request.StudentName,
		StudentNumber: request.StudentNumber,
		RequestMsg:    request.RequestMsg,
		Status:        request.Status,
		ReviewerID:    request.ReviewerID,
		ReviewMsg:     request.ReviewMsg,
		ReviewTime:    request.ReviewTime,
		CreateTime:    request.Ctime,
		UpdateTime:    request.Mtime,
	}
}
