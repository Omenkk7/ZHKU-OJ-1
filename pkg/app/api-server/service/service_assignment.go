package service

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"log"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// CreateAssignment 创建作业
func (s *Service) CreateAssignment(ctx context.Context, req *dto.CreateAssignmentReq, userID string, userName string) (string, error) {
	course, err := s.dao.GetCourseByID(ctx, req.CourseID)
	if err != nil {
		return "", err
	}
	if course == nil {
		return "", errors.New("课程不存在")
	}

	// 检查用户是否有权限创建作业
	roles, err := s.dao.CheckUserCourseRole(ctx, req.CourseID, userID)
	if err != nil {
		return "", err
	}
	hasPermission := false
	for _, role := range roles {
		if role == 1 || role == 2 { // 管理员或教师
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return "", errors.New("没有权限创建作业")
	}

	// 生成作业代码并确保唯一性
	var assignmentCode string
	for {
		// 使用工具函数生成作业代码
		assignmentCode = utils.GenerateAssignmentCode()

		// 检查作业代码是否已存在
		existingAssignment, err := s.dao.GetAssignmentByCode(ctx, assignmentCode)
		if err != nil {
			return "", errors.New("检查作业代码时出错")
		}

		// 如果不存在，则使用这个代码
		if existingAssignment == nil {
			break
		}

		// 如果存在，则重新生成
		log.Printf("作业代码 %s 已存在，重新生成", assignmentCode)
	}

	// 创建作业对象
	assignment := &models.Assignment{
		AssignmentCode:   assignmentCode,
		Title:            req.Title,
		Description:      req.Description,
		CourseID:         req.CourseID,
		CourseCode:       course.CourseCode,
		CourseName:       course.Name,
		ClassIDs:         req.ClassIDs,
		CreatorID:        userID,
		CreatorName:      userName,
		AssignmentType:   req.AssignmentType,
		StartTime:        req.StartTime,
		EndTime:          req.EndTime,
		MaxResubmitCount: req.MaxResubmitCount,
		Members:          make(map[string][]models.Member),
		Stats: models.AssignmentStats{
			TotalStudents:  0,
			SubmittedCount: 0,
			GradedCount:    0,
			AverageScore:   0,
		},
		OJProblems: req.OJProblems,
		Status:     1, // 进行中
	}

	// 初始化成员列表
	assignment.Members["1"] = []models.Member{} // 管理员
	assignment.Members["2"] = []models.Member{} // 教师
	assignment.Members["3"] = []models.Member{} // 助教
	assignment.Members["4"] = []models.Member{} // 学生

	// 添加创建者为管理员
	creatorMember := models.Member{
		UserID:   userID,
		UserName: userName,
	}
	assignment.Members["1"] = append(assignment.Members["1"], creatorMember)

	// 如果有关联班级，添加班级学生
	if len(req.ClassIDs) > 0 {
		for _, classInfo := range req.ClassIDs {
			// 获取班级学生
			comQuery := &utils.CommonQuery{
				Filters: map[string]any{
					"class_id": classInfo.ClassID,
					"status":   1, // 状态为正常
				},
			}
			result, err := s.dao.GetClassStudents(ctx, classInfo.ClassID, comQuery)
			if err != nil {
				log.Printf("获取班级学生失败: %v", err)
				continue
			}

			// 添加学生到作业
			for _, item := range result.Items {
				studentMap := *item
				studentID := studentMap["student_id"].(string)
				studentName := studentMap["student_name"].(string)

				// 检查学生是否已存在
				exists := false
				for _, member := range assignment.Members["4"] {
					if member.UserID == studentID {
						exists = true
						break
					}
				}

				if !exists {
					studentMember := models.Member{
						UserID:   studentID,
						UserName: studentName,
					}
					assignment.Members["4"] = append(assignment.Members["4"], studentMember)
				}
			}
		}
	}

	// 更新统计信息
	assignment.Stats.TotalStudents = len(assignment.Members["4"])

	// 保存到数据库
	return s.dao.CreateAssignment(ctx, assignment)
}

// UpdateAssignment 更新作业
func (s *Service) UpdateAssignment(ctx context.Context, id string, req *dto.UpdateAssignmentReq, userID string) error {
	// 获取作业信息
	assignment, err := s.dao.GetAssignmentByID(ctx, id)
	if err != nil {
		return err
	}
	if assignment == nil {
		return errors.New("作业不存在")
	}

	// 检查用户是否有权限更新作业
	roles, err := s.dao.CheckUserCourseRole(ctx, assignment.CourseID, userID)
	if err != nil {
		return err
	}
	hasPermission := false
	for _, role := range roles {
		if role == 1 || role == 2 || (role == 3 && assignment.CreatorID == userID) { // 管理员、教师或创建者助教
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return errors.New("没有权限更新作业")
	}

	// 构建更新内容
	update := bson.M{
		"$set": bson.M{},
	}

	// 更新字段
	updateSet := update["$set"].(bson.M)
	if req.Title != "" {
		updateSet["title"] = req.Title
	}
	if req.Description != "" {
		updateSet["description"] = req.Description
	}
	if len(req.ClassIDs) > 0 {
		updateSet["class_ids"] = req.ClassIDs
	}
	if req.StartTime > 0 {
		updateSet["start_time"] = req.StartTime
	}
	if req.EndTime > 0 {
		updateSet["end_time"] = req.EndTime
	}
	if req.MaxResubmitCount >= 0 {
		updateSet["max_resubmit_count"] = req.MaxResubmitCount
	}
	if len(req.OJProblems) > 0 {
		updateSet["oj_problems"] = req.OJProblems
	}
	if req.Status >= 0 {
		updateSet["status"] = req.Status
	}

	// 如果没有更新内容，直接返回
	if len(updateSet) == 0 {
		return nil
	}

	// 更新数据库
	return s.dao.UpdateAssignment(ctx, id, update)
}

// DeleteAssignment 删除作业
func (s *Service) DeleteAssignment(ctx context.Context, id string, userID string) error {
	// 获取作业信息
	assignment, err := s.dao.GetAssignmentByID(ctx, id)
	if err != nil {
		return err
	}
	if assignment == nil {
		return errors.New("作业不存在")
	}

	// 检查用户是否有权限删除作业
	roles, err := s.dao.CheckUserCourseRole(ctx, assignment.CourseID, userID)
	if err != nil {
		return err
	}
	hasPermission := false
	for _, role := range roles {
		if role == 1 || (role == 2 && assignment.CreatorID == userID) { // 管理员或创建者教师
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return errors.New("没有权限删除作业")
	}

	// 删除作业
	return s.dao.DeleteAssignment(ctx, id)
}

// ArchiveAssignment 归档作业
func (s *Service) ArchiveAssignment(ctx context.Context, id string, userID string) error {
	// 获取作业信息
	assignment, err := s.dao.GetAssignmentByID(ctx, id)
	if err != nil {
		return err
	}
	if assignment == nil {
		return errors.New("作业不存在")
	}

	// 检查用户是否有权限归档作业
	roles, err := s.dao.CheckUserCourseRole(ctx, assignment.CourseID, userID)
	if err != nil {
		return err
	}
	hasPermission := false
	for _, role := range roles {
		if role == 1 || role == 2 || (role == 3 && assignment.CreatorID == userID) { // 管理员、教师或创建者助教
			hasPermission = true
			break
		}
	}
	if !hasPermission {
		return errors.New("没有权限归档作业")
	}

	// 归档作业
	return s.dao.ArchiveAssignment(ctx, id)
}

// GetAssignmentDetail 获取作业详情
func (s *Service) GetAssignmentDetail(ctx context.Context, id string, userID string) (*dto.AssignmentDetailResp, error) {
	// 获取作业信息
	assignment, err := s.dao.GetAssignmentByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, errors.New("作业不存在")
	}

	// 检查用户是否有权限查看作业
	hasPermission := false

	// 检查用户是否是作业的成员
	for _, members := range assignment.Members {
		for _, member := range members {
			if member.UserID == userID {
				hasPermission = true
				break
			}
		}
		if hasPermission {
			break
		}
	}

	// 如果不是作业成员，检查用户在课程中的角色
	if !hasPermission {
		roles, err := s.dao.CheckUserCourseRole(ctx, assignment.CourseID, userID)
		if err != nil {
			return nil, err
		}
		for _, role := range roles {
			if role >= 1 && role <= 4 { // 任何角色都可以查看
				hasPermission = true
				break
			}
		}
	}

	if !hasPermission {
		return nil, errors.New("没有权限查看作业")
	}

	// 构建响应
	resp := &dto.AssignmentDetailResp{
		ID:               assignment.ID.Hex(),
		AssignmentCode:   assignment.AssignmentCode,
		Title:            assignment.Title,
		Description:      assignment.Description,
		CourseID:         assignment.CourseID,
		CourseCode:       assignment.CourseCode,
		CourseName:       assignment.CourseName,
		ClassIDs:         assignment.ClassIDs,
		CreatorID:        assignment.CreatorID,
		CreatorName:      assignment.CreatorName,
		AssignmentType:   assignment.AssignmentType,
		StartTime:        assignment.StartTime,
		EndTime:          assignment.EndTime,
		MaxResubmitCount: assignment.MaxResubmitCount,
		OJProblems:       assignment.OJProblems,
		Status:           assignment.Status,
		Stats:            assignment.Stats,
		Ctime:            assignment.Ctime,
		Mtime:            assignment.Mtime,
	}

	return resp, nil
}

// GetAssignmentList 获取作业列表
func (s *Service) GetAssignmentList(ctx context.Context, req *dto.GetAssignmentListReq, userID string) (*utils.RespPageQuery, error) {
	// 构建筛选条件
	filters := make(map[string]interface{})

	// 添加筛选条件
	if req.CourseID != "" {
		filters["course_id"] = req.CourseID
	}
	if req.ClassID != "" {
		filters["class_ids.class_id"] = req.ClassID
	}
	// 修改状态判断逻辑，只有当 status 参数有效时才添加到过滤条件
	if req.Status != nil {
		filters["status"] = *req.Status
	}
	// 修改类型判断逻辑，只有当 type 参数有效时才添加到过滤条件
	if req.Type != nil {
		filters["assignment_type"] = *req.Type
	}
	if req.CreatorID != "" {
		filters["creator_id"] = req.CreatorID
	}
	if req.Keyword != "" {
		filters["$or"] = []bson.M{
			{"title": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"description": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"assignment_code": bson.M{"$regex": req.Keyword, "$options": "i"}},
		}
	}

	// 构建排序条件
	sorts := map[string]interface{}{
		"ctime": -1, // 默认按创建时间倒序
	}

	// 使用 BuildQueryParamsFromValues 构建查询参数
	comQuery := utils.BuildQueryParamsFromValues(req.Page, req.PageSize, filters, sorts)

	// 获取作业列表
	return s.dao.GetAssignmentList(ctx, comQuery)
}

// GetStudentAssignments 获取学生的作业列表
func (s *Service) GetStudentAssignments(ctx context.Context, req *dto.GetStudentAssignmentsReq, studentID string) (*utils.RespPageQuery, error) {
	// 构建筛选条件
	filters := make(map[string]interface{})

	// 添加筛选条件
	if req.CourseID != "" {
		filters["course_id"] = req.CourseID
	}

	// 构建排序条件
	sorts := map[string]interface{}{
		"ctime": -1, // 默认按创建时间倒序
	}

	// 使用 BuildQueryParamsFromValues 构建查询参数
	comQuery := utils.BuildQueryParamsFromValues(req.Page, req.PageSize, filters, sorts)

	// 调用dao层方法获取学生作业列表
	return s.dao.GetStudentAssignments(ctx, studentID, comQuery, req.Status)
}

// SubmitAssignment 提交作业
func (s *Service) SubmitAssignment(ctx context.Context, assignmentID string, req *dto.SubmitAssignmentReq, studentID string) error {
	// 获取作业信息
	assignment, err := s.dao.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		return err
	}
	if assignment == nil {
		return errors.New("作业不存在")
	}

	// 检查作业是否已结束
	now := time.Now().Unix()
	if now > assignment.EndTime {
		return errors.New("作业已截止，无法提交")
	}

	// 检查学生是否在作业成员中
	studentExists := false
	var studentSubmission *models.Submission
	for _, member := range assignment.Members["4"] {
		if member.UserID == studentID {
			studentExists = true
			studentSubmission = member.Submission
			break
		}
	}

	// 添加更严格的检查和日志
	if !studentExists {
		log.Printf("用户 %s 尝试提交作业 %s，但不是该作业的学生成员", studentID, assignmentID)
		return errors.New("您不是该作业的学生成员，无法提交")
	}

	// 检查是否允许重做
	if studentSubmission != nil && studentSubmission.Status == 2 { // 已批改
		if assignment.MaxResubmitCount <= 0 {
			return errors.New("该作业不允许重做")
		}

		// TODO: 检查重做次数，在submission中添加重做次数字段
	}

	// 创建提交信息
	submission := &models.Submission{
		Status:      1, // 已提交未批改
		TextContent: req.TextContent,
		FileURL:     req.FileURL,
		SubmitTime:  now,
	}

	// 提交作业
	return s.dao.SubmitAssignment(ctx, assignmentID, studentID, submission)
}

// GradeAssignment 批改作业
func (s *Service) GradeAssignment(ctx context.Context, req *dto.GradeAssignmentReq, graderID string, graderName string) error {
	// 检查用户是否有权限批改作业
	hasPermission, err := s.CheckAssignmentPermission(ctx, req.AssignmentID, graderID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("没有权限批改作业")
	}

	// 获取作业信息（为了后续检查学生是否在作业成员中）
	assignment, err := s.dao.GetAssignmentByID(ctx, req.AssignmentID)
	if err != nil {
		return err
	}
	if assignment == nil {
		return errors.New("作业不存在")
	}

	// 检查学生是否在作业成员中
	studentExists := false
	var studentSubmission *models.Submission
	for _, member := range assignment.Members["4"] {
		if member.UserID == req.StudentID {
			studentExists = true
			studentSubmission = member.Submission
			break
		}
	}
	if !studentExists {
		return errors.New("该学生不是作业成员")
	}

	// 检查学生是否已提交作业
	if studentSubmission == nil || studentSubmission.Status == 0 {
		return errors.New("该学生尚未提交作业")
	}

	// 批改作业
	return s.dao.GradeAssignment(ctx, req.AssignmentID, req.StudentID, graderID, graderName, req.Score, req.Comment)
}

// RejectAssignment 打回作业
func (s *Service) RejectAssignment(ctx context.Context, req *dto.RejectAssignmentReq, graderID string, graderName string) error {
	// 检查用户是否有权限打回作业
	hasPermission, err := s.CheckAssignmentPermission(ctx, req.AssignmentID, graderID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("没有权限打回作业")
	}

	// 获取作业信息（为了后续检查学生是否在作业成员中）
	assignment, err := s.dao.GetAssignmentByID(ctx, req.AssignmentID)
	if err != nil {
		return err
	}
	if assignment == nil {
		return errors.New("作业不存在")
	}

	// 检查学生是否在作业成员中
	studentExists := false
	var studentSubmission *models.Submission
	for _, member := range assignment.Members["4"] {
		if member.UserID == req.StudentID {
			studentExists = true
			studentSubmission = member.Submission
			break
		}
	}
	if !studentExists {
		return errors.New("该学生不是作业成员")
	}

	// 检查学生是否已提交作业
	if studentSubmission == nil || studentSubmission.Status == 0 {
		return errors.New("该学生尚未提交作业")
	}

	// 打回作业
	return s.dao.RejectAssignment(ctx, req.AssignmentID, req.StudentID, graderID, graderName, req.Comment)
}

// GetStudentSubmission 获取学生的作业提交信息
func (s *Service) GetStudentSubmission(ctx context.Context, req *dto.GetStudentSubmissionReq, userID string) (*dto.StudentSubmissionResp, error) {
	// 获取作业信息
	assignment, err := s.dao.GetAssignmentByID(ctx, req.AssignmentID)
	if err != nil {
		return nil, err
	}
	if assignment == nil {
		return nil, errors.New("作业不存在")
	}

	// 检查权限
	if userID != req.StudentID {
		// 非本人查看，检查是否有权限
		roles, err := s.dao.CheckUserCourseRole(ctx, assignment.CourseID, userID)
		if err != nil {
			return nil, err
		}
		hasPermission := false
		for _, role := range roles {
			if role == 1 || role == 2 || role == 3 { // 管理员、教师或助教
				hasPermission = true
				break
			}
		}
		if !hasPermission {
			return nil, errors.New("没有权限查看其他学生的提交")
		}
	}

	// 获取学生提交信息
	submission, err := s.dao.GetStudentSubmission(ctx, req.AssignmentID, req.StudentID)
	if err != nil {
		return nil, err
	}

	// 如果没有提交信息，返回空的提交信息
	if submission == nil {
		return &dto.StudentSubmissionResp{
			Status: 0, // 未提交
		}, nil
	}

	// 获取批改人姓名
	graderName := submission.GraderName
	if submission.GraderID != "" && graderName == "" {
		// 使用 GetOneUser 方法替代 GetUserByID
		objID, err := primitive.ObjectIDFromHex(submission.GraderID)
		if err == nil {
			selector := bson.M{
				"_id": objID,
			}
			user, err := s.dao.GetOneUser(ctx, selector)
			if err == nil && user != nil {
				// 使用 Nickname 字段替代 Username
				graderName = user.Nickname
			}
		}
	}

	// 返回提交信息
	return &dto.StudentSubmissionResp{
		Status:      submission.Status,
		TextContent: submission.TextContent,
		FileURL:     submission.FileURL,
		Score:       submission.Score,
		Comment:     submission.Comment,
		GraderName:  graderName,
		SubmitTime:  submission.SubmitTime,
		GradeTime:   submission.GradeTime,
	}, nil
}

// ExportAssignmentGrades 导出作业成绩
func (s *Service) ExportAssignmentGrades(ctx context.Context, assignmentID string, userID string) ([]map[string]interface{}, error) {
	// 检查用户是否有权限导出成绩
	hasPermission, err := s.CheckAssignmentPermission(ctx, assignmentID, userID)
	if err != nil {
		return nil, err
	}
	if !hasPermission {
		return nil, errors.New("没有权限导出成绩")
	}

	// 导出成绩
	return s.dao.ExportAssignmentGrades(ctx, assignmentID)
}

// CheckAssignmentPermission 检查用户是否有权限操作作业，管理员（系统管理员或作业管理员）可以操作所有作业，教师只能操作自己负责的作业（在作业的教师列表中）
func (s *Service) CheckAssignmentPermission(ctx context.Context, assignmentID string, userID string) (bool, error) {
	// 获取作业信息
	assignment, err := s.dao.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		return false, err
	}
	if assignment == nil {
		return false, errors.New("作业不存在")
	}

	// 检查用户是否有权限操作作业
	hasPermission := false

	// 检查是否是系统管理员
	objID, err := primitive.ObjectIDFromHex(userID)
	if err == nil {
		selector := bson.M{
			"_id": objID,
		}
		user, err := s.dao.GetOneUser(ctx, selector)
		if err == nil && user != nil && user.Role == 1 { // 系统管理员
			hasPermission = true
			return hasPermission, nil
		}
	}

	// 检查是否是作业管理员
	for _, member := range assignment.Members["1"] {
		if member.UserID == userID {
			hasPermission = true
			return hasPermission, nil
		}
	}

	// 如果不是管理员，检查是否是作业教师
	for _, member := range assignment.Members["2"] {
		if member.UserID == userID {
			hasPermission = true
			return hasPermission, nil
		}
	}

	// 检查用户在课程中的角色（兼容旧逻辑）
	if !hasPermission {
		roles, err := s.dao.CheckUserCourseRole(ctx, assignment.CourseID, userID)
		if err != nil {
			return false, err
		}

		for _, role := range roles {
			if role == 1 { // 课程管理员
				hasPermission = true
				break
			}
		}
	}

	return hasPermission, nil
}

// AddStudentsToAssignmentDirect 批量添加学生到作业
func (s *Service) AddStudentsToAssignmentDirect(ctx context.Context, assignmentID string, students []models.Member, userID string) error {
	// 检查用户是否有权限添加学生
	hasPermission, err := s.CheckAssignmentPermission(ctx, assignmentID, userID)
	if err != nil {
		return err
	}
	if !hasPermission {
		return errors.New("没有权限添加学生")
	}

	// 获取作业信息（为了后续过滤已存在的学生）
	assignment, err := s.dao.GetAssignmentByID(ctx, assignmentID)
	if err != nil {
		return err
	}
	if assignment == nil {
		return errors.New("作业不存在")
	}

	// 过滤已存在的学生
	var newStudents []models.Member
	for _, student := range students {
		// 检查学生是否已存在
		exists := false
		for _, member := range assignment.Members["4"] {
			if member.UserID == student.UserID {
				exists = true
				break
			}
		}
		if !exists {
			newStudents = append(newStudents, student)
		}
	}

	if len(newStudents) == 0 {
		return nil
	}

	// 添加学生到作业
	return s.dao.AddStudentsToAssignment(ctx, assignmentID, newStudents)
}
