package service

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/net/context"
	"log"
	"strconv"
	"time"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/models"
	"zhku-oj-server/pkg/utils"
)

// CreateClass 创建班级
func (s *Service) CreateClass(ctx context.Context, req *dto.CreateClassRequest, creatorID string) (*dto.ClassResponse, error) {
	// 如果没有提供班级代码，则自动生成
	if req.ClassCode == "" {
		// 生成班级代码
		req.ClassCode = utils.GenerateClassCode(req.Department)

		// 检查班级代码是否已存在
		existingClass, err := s.dao.GetClassByCode(ctx, req.ClassCode)
		if err != nil {
			return nil, err
		}

		// 如果班级代码已存在，则重新生成
		for existingClass != nil {
			req.ClassCode = utils.GenerateClassCode(req.Department)
			existingClass, err = s.dao.GetClassByCode(ctx, req.ClassCode)
			if err != nil {
				return nil, err
			}
		}
	} else {
		// 如果提供了班级代码，检查是否已存在
		existingClass, err := s.dao.GetClassByCode(ctx, req.ClassCode)
		if err != nil {
			return nil, err
		}
		if existingClass != nil {
			return nil, errors.New("班级代码已存在")
		}
	}

	// 获取创建者信息
	// 添加日志，输出creatorID的值
	log.Printf("创建班级的用户ID: %s", creatorID)

	// 如果从上下文中获取不到userId，尝试从JWT中获取
	if creatorID == "" {
		// 尝试从上下文中获取完整的JWT声明
		if claims, ok := ctx.Value("claims").(jwt.MapClaims); ok {
			if id, ok := claims["userId"].(string); ok {
				creatorID = id
				log.Printf("从JWT中获取到用户ID: %s", creatorID)
			}
		}
	}

	// 创建班级对象
	class := &models.Class{
		ClassCode:    req.ClassCode,
		Name:         req.Name,
		Department:   req.Department,
		Description:  req.Description,
		Creator:      creatorID,
		Courses:      req.Courses,
		Members:      make(map[string][]models.ClassMember),
		StudentCount: 0,
		Status:       1, // 正常状态
	}

	// 创建班级
	id, err := s.dao.CreateClass(ctx, class)
	if err != nil {
		return nil, err
	}

	// 获取创建后的班级信息
	createdClass, err := s.dao.GetClassByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// 转换为响应DTO
	return s.convertClassToResponse(createdClass), nil
}

// UpdateClass 更新班级
func (s *Service) UpdateClass(ctx context.Context, classID string, req *dto.UpdateClassRequest) (*dto.ClassResponse, error) {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("班级不存在")
	}

	// 构建更新内容
	update := bson.M{
		"$set": bson.M{},
	}

	// 只更新非空字段
	setFields := update["$set"].(bson.M)
	if req.Name != "" {
		setFields["name"] = req.Name
	}
	if req.Department != "" {
		setFields["department"] = req.Department
	}
	if req.Description != "" {
		setFields["description"] = req.Description
	}

	// 如果提供了课程信息，则更新课程
	if req.Courses != nil {
		setFields["courses"] = req.Courses
	}

	// 执行更新
	err = s.dao.UpdateClass(ctx, classID, update)
	if err != nil {
		return nil, err
	}

	// 获取更新后的班级信息
	updatedClass, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return nil, err
	}

	// 转换为响应DTO
	return s.convertClassToResponse(updatedClass), nil
}

// DeleteClass 删除班级
func (s *Service) DeleteClass(ctx context.Context, classID string) error {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("班级不存在")
	}

	// 检查班级是否有学生
	if class.StudentCount > 0 {
		return errors.New("班级中还有学生，无法删除")
	}

	// 执行删除
	return s.dao.DeleteClass(ctx, classID)
}

// ArchiveClass 归档班级
func (s *Service) ArchiveClass(ctx context.Context, classID string) error {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("班级不存在")
	}

	// 执行归档
	return s.dao.ArchiveClass(ctx, classID)
}

// GetClassByID 根据ID获取班级
func (s *Service) GetClassByID(ctx context.Context, classID string) (*dto.ClassResponse, error) {
	// 获取班级信息
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("班级不存在")
	}

	// 转换为响应DTO
	return s.convertClassToResponse(class), nil
}

// GetClassByCode 根据班级代码获取班级
func (s *Service) GetClassByCode(ctx context.Context, classCode string) (*dto.ClassResponse, error) {
	// 获取班级信息
	class, err := s.dao.GetClassByCode(ctx, classCode)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("班级不存在")
	}

	// 转换为响应DTO
	return s.convertClassToResponse(class), nil
}

// AddStudentToClass 添加学生到班级
func (s *Service) AddStudentToClass(ctx context.Context, classID string, req *dto.AddStudentRequest) (*dto.ClassStudentResponse, error) {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("班级不存在")
	}

	// 检查学生是否已在班级中
	inClass, err := s.dao.CheckStudentInClass(ctx, classID, req.StudentID)
	if err != nil {
		return nil, err
	}
	if inClass {
		return nil, errors.New("学生已在班级中")
	}

	// 添加学生到班级
	classStudent := &models.ClassStudent{
		ClassID:       classID,
		StudentID:     req.StudentID,
		StudentName:   req.StudentName,
		StudentNumber: req.StudentNumber,
		JoinType:      1,
		Status:        1, // 正常状态
		JoinTime:      time.Now().Unix(),
		LeaveTime:     nil,
	}

	id, err := s.dao.AddStudentToClass(ctx, classStudent)
	if err != nil {
		return nil, err
	}

	// 更新班级学生数量
	err = s.dao.UpdateClassStudentCount(ctx, classID)
	if err != nil {
		return nil, err
	}

	// 转换为响应DTO
	return &dto.ClassStudentResponse{
		ID:            id,
		ClassID:       classStudent.ClassID,
		StudentID:     classStudent.StudentID,
		StudentName:   classStudent.StudentName,
		StudentNumber: classStudent.StudentNumber,
		JoinType:      classStudent.JoinType,
		Status:        classStudent.Status,
		JoinTime:      classStudent.JoinTime,
		LeaveTime:     classStudent.LeaveTime,
		CreateTime:    classStudent.Ctime,
		UpdateTime:    classStudent.Mtime,
	}, nil
}

// BatchAddStudentsToClass 批量添加学生到班级
func (s *Service) BatchAddStudentsToClass(ctx context.Context, classID string, req *dto.BatchAddStudentsRequest) error {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("班级不存在")
	}

	// 准备批量添加的学生数据
	classStudents := make([]models.ClassStudent, 0, len(req.Students))
	now := time.Now().Unix()

	for _, student := range req.Students {
		// 检查学生是否已在班级中
		inClass, err := s.dao.CheckStudentInClass(ctx, classID, student.StudentID)
		if err != nil {
			return err
		}
		if inClass {
			continue // 跳过已在班级中的学生
		}

		classStudents = append(classStudents, models.ClassStudent{
			ClassID:       classID,
			StudentID:     student.StudentID,
			StudentName:   student.StudentName,
			StudentNumber: student.StudentNumber,
			JoinType:      1,
			Status:        1, // 正常状态
			JoinTime:      now,
			LeaveTime:     nil,
			Ctime:         now,
			Mtime:         now,
		})
	}

	// 如果没有需要添加的学生，直接返回
	if len(classStudents) == 0 {
		return nil
	}

	// 批量添加学生
	err = s.dao.BatchAddStudentsToClass(ctx, classStudents)
	if err != nil {
		return err
	}

	// 更新班级学生数量
	return s.dao.UpdateClassStudentCount(ctx, classID)
}

// RemoveStudentFromClass 从班级移除学生
func (s *Service) RemoveStudentFromClass(ctx context.Context, classID, studentID string) error {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("班级不存在")
	}

	// 检查学生是否在班级中
	inClass, err := s.dao.CheckStudentInClass(ctx, classID, studentID)
	if err != nil {
		return err
	}
	if !inClass {
		return errors.New("学生不在班级中")
	}

	// 移除学生
	err = s.dao.RemoveStudentFromClass(ctx, classID, studentID)
	if err != nil {
		return err
	}

	// 更新班级学生数量
	return s.dao.UpdateClassStudentCount(ctx, classID)
}

// GetClassStudents 获取班级学生列表
func (s *Service) GetClassStudents(ctx context.Context, classID string, page, pageSize int, filters map[string]interface{}, sorts map[string]interface{}) (*utils.RespPageQuery, error) {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("班级不存在")
	}

	comQuery := utils.BuildQueryParamsFromValues(page, pageSize, filters, sorts)

	// 执行查询
	return s.dao.GetClassStudents(ctx, classID, comQuery)
}

// CreateJoinRequest 创建加入班级申请
func (s *Service) CreateJoinRequest(ctx context.Context, req *dto.JoinClassRequest) (*dto.JoinRequestResponse, error) {
	// 检查班级是否存在
	class, err := s.dao.GetClassByCode(ctx, req.ClassCode)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("班级不存在")
	}

	// 检查学生是否已在班级中
	inClass, err := s.dao.CheckStudentInClass(ctx, class.ID.Hex(), req.StudentID)
	if err != nil {
		return nil, err
	}
	if inClass {
		return nil, errors.New("学生已在班级中")
	}

	// 创建加入申请
	joinRequest := &models.ClassJoinRequest{
		ClassID:       class.ID.Hex(),
		StudentID:     req.StudentID,
		StudentName:   req.StudentName,
		StudentNumber: req.StudentNumber,
		RequestMsg:    req.RequestMsg,
		Status:        0, // 待审核
		ReviewerID:    "",
		ReviewMsg:     "",
		ReviewTime:    nil,
	}

	id, err := s.dao.CreateJoinRequest(ctx, joinRequest)
	if err != nil {
		return nil, err
	}

	// 转换为响应DTO
	return &dto.JoinRequestResponse{
		ID:            id,
		ClassID:       joinRequest.ClassID,
		StudentID:     joinRequest.StudentID,
		StudentName:   joinRequest.StudentName,
		StudentNumber: joinRequest.StudentNumber,
		RequestMsg:    joinRequest.RequestMsg,
		Status:        joinRequest.Status,
		ReviewerID:    joinRequest.ReviewerID,
		ReviewMsg:     joinRequest.ReviewMsg,
		ReviewTime:    joinRequest.ReviewTime,
		CreateTime:    joinRequest.Ctime,
		UpdateTime:    joinRequest.Mtime,
	}, nil
}

// ReviewJoinRequest 审核加入班级申请
func (s *Service) ReviewJoinRequest(ctx context.Context, requestID string, req *dto.ReviewJoinRequest, reviewerID string) (*dto.JoinRequestResponse, error) {
	// 添加日志记录
	lg := utils.GetDefaultLogger()
	lg.Infof("开始审核班级加入申请，申请ID: %s, 审核人ID: %s, 状态: %d", requestID, reviewerID, req.Status)

	// 验证 requestID 是否为有效的 ObjectID
	_, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		lg.Errorf("无效的申请ID格式: %s, 错误: %v", requestID, err)
		return nil, errors.New("无效的申请ID格式")
	}

	// 获取申请信息
	joinRequest, err := s.dao.GetJoinRequestByID(ctx, requestID)
	if err != nil {
		lg.Errorf("获取申请信息失败: %v", err)
		return nil, err
	}
	if joinRequest == nil {
		lg.Warnf("申请不存在，ID: %s", requestID)
		return nil, errors.New("申请不存在")
	}

	// 检查申请状态
	if joinRequest.Status != 0 {
		lg.Warnf("申请已处理，ID: %s, 当前状态: %d", requestID, joinRequest.Status)
		return nil, errors.New("申请已处理")
	}

	// 更新申请状态
	now := time.Now().Unix()
	update := bson.M{
		"$set": bson.M{
			"status":      req.Status,
			"reviewer_id": reviewerID,
			"review_msg":  req.ReviewMsg,
			"review_time": now,
		},
	}

	lg.Infof("更新申请状态，ID: %s, 状态: %d", requestID, req.Status)
	err = s.dao.UpdateJoinRequest(ctx, requestID, update)
	if err != nil {
		lg.Errorf("更新申请状态失败: %v", err)
		return nil, err
	}

	// 如果审核通过，添加学生到班级
	if req.Status == 1 {
		lg.Infof("审核通过，添加学生到班级，班级ID: %s, 学生ID: %s", joinRequest.ClassID, joinRequest.StudentID)
		classStudent := &models.ClassStudent{
			ClassID:       joinRequest.ClassID,
			StudentID:     joinRequest.StudentID,
			StudentName:   joinRequest.StudentName,
			StudentNumber: joinRequest.StudentNumber,
			JoinType:      2, // 自主申请
			Status:        1, // 正常状态
			JoinTime:      now,
			LeaveTime:     nil,
		}

		_, err = s.dao.AddStudentToClass(ctx, classStudent)
		if err != nil {
			lg.Errorf("添加学生到班级失败: %v", err)
			return nil, err
		}

		// 将学生添加到班级的members字段中（学生角色ID为4）
		member := models.ClassMember{
			UserID:   joinRequest.StudentID,
			UserName: joinRequest.StudentName,
		}
		err = s.dao.AddClassMember(ctx, joinRequest.ClassID, member, 4) // 4表示学生角色
		if err != nil {
			lg.Errorf("添加学生到班级members字段失败: %v", err)
			return nil, err
		}

		// 更新班级学生数量
		lg.Infof("更新班级学生数量，班级ID: %s", joinRequest.ClassID)
		err = s.dao.UpdateClassStudentCount(ctx, joinRequest.ClassID)
		if err != nil {
			lg.Errorf("更新班级学生数量失败: %v", err)
			return nil, err
		}
	}

	// 直接使用已有的 joinRequest 对象构建响应，避免再次查询数据库
	// 更新 joinRequest 对象的相关字段
	joinRequest.Status = req.Status
	joinRequest.ReviewerID = reviewerID
	joinRequest.ReviewMsg = req.ReviewMsg
	reviewTime := now
	joinRequest.ReviewTime = &reviewTime
	joinRequest.Mtime = now

	// 转换为响应DTO
	return &dto.JoinRequestResponse{
		ID:            joinRequest.ID.Hex(),
		ClassID:       joinRequest.ClassID,
		StudentID:     joinRequest.StudentID,
		StudentName:   joinRequest.StudentName,
		StudentNumber: joinRequest.StudentNumber,
		RequestMsg:    joinRequest.RequestMsg,
		Status:        joinRequest.Status,
		ReviewerID:    joinRequest.ReviewerID,
		ReviewMsg:     joinRequest.ReviewMsg,
		ReviewTime:    joinRequest.ReviewTime,
		CreateTime:    joinRequest.Ctime,
		UpdateTime:    joinRequest.Mtime,
	}, nil
}

// GetJoinRequestList 获取加入申请列表
func (s *Service) GetClassList(ctx context.Context, page, pageSize int, filters map[string]interface{}, sorts map[string]interface{}) (*utils.RespPageQuery, error) {
	comQuery := utils.BuildQueryParamsFromValues(page, pageSize, filters, sorts)

	// 执行查询
	return s.dao.GetClassList(ctx, comQuery)
}

// GetStudentClasses 获取学生所在的班级列表
func (s *Service) GetStudentClasses(ctx context.Context, studentID string, page, pageSize int, filters map[string]interface{}, sorts map[string]interface{}) (*utils.RespPageQuery, error) {
	comQuery := utils.BuildQueryParamsFromValues(page, pageSize, filters, sorts)

	// 执行查询
	return s.dao.GetClassesByStudentID(ctx, studentID, comQuery)
}

// AddCourseToClass 添加课程到班级
func (s *Service) AddCourseToClass(ctx context.Context, classID string, req *dto.AddCourseRequest) error {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("班级不存在")
	}

	// 检查课程是否已在班级中
	for _, course := range class.Courses {
		if course.CourseID == req.CourseID {
			return errors.New("课程已在班级中")
		}
	}

	// 添加课程
	courseInfo := models.CourseInfo{
		CourseID:   req.CourseID,
		CourseName: req.CourseName,
		BindTime:   time.Now().Unix(),
		Status:     req.Status,
	}

	return s.dao.AddCoursesToClass(ctx, classID, []models.CourseInfo{courseInfo})
}

// RemoveCourseFromClass 从班级移除课程
func (s *Service) RemoveCourseFromClass(ctx context.Context, classID, courseID string) error {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("班级不存在")
	}

	// 检查课程是否在班级中
	courseExists := false
	for _, course := range class.Courses {
		if course.CourseID == courseID {
			courseExists = true
			break
		}
	}

	if !courseExists {
		return errors.New("课程不在班级中")
	}

	// 移除课程
	return s.dao.RemoveCourseFromClass(ctx, classID, courseID)
}

// UpdateCourseStatusInClass 更新班级中课程的状态
func (s *Service) UpdateCourseStatusInClass(ctx context.Context, classID, courseID string, req *dto.UpdateCourseStatusRequest) error {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("班级不存在")
	}

	// 检查课程是否在班级中
	courseExists := false
	for _, course := range class.Courses {
		if course.CourseID == courseID {
			courseExists = true
			break
		}
	}

	if !courseExists {
		return errors.New("课程不在班级中")
	}

	// 更新课程状态
	return s.dao.UpdateCourseStatusInClass(ctx, classID, courseID, req.Status)
}

// GetClassesByCourseID 获取绑定了指定课程的班级列表
func (s *Service) GetClassesByCourseID(ctx context.Context, courseID string, page, pageSize int, filters map[string]interface{}, sorts map[string]interface{}) (*utils.RespPageQuery, error) {
	comQuery := utils.BuildQueryParamsFromValues(page, pageSize, filters, sorts)

	// 执行查询
	return s.dao.GetClassesByCourseID(ctx, courseID, comQuery)
}

// 辅助方法：将Class模型转换为ClassResponse
func (s *Service) convertClassToResponse(class *models.Class) *dto.ClassResponse {
	if class == nil {
		return nil
	}

	// 转换成员信息
	members := make(map[string][]dto.ClassMemberResponse)
	for role, roleMembers := range class.Members {
		memberResponses := make([]dto.ClassMemberResponse, 0, len(roleMembers))
		for _, member := range roleMembers {
			memberResponses = append(memberResponses, dto.ClassMemberResponse{
				UserID:   member.UserID,
				UserName: member.UserName,
			})
		}
		members[role] = memberResponses
	}

	return &dto.ClassResponse{
		ID:           class.ID.Hex(),
		ClassCode:    class.ClassCode,
		Name:         class.Name,
		Department:   class.Department,
		Description:  class.Description,
		Creator:      class.Creator,
		Courses:      class.Courses,
		Members:      members, // 添加成员信息
		StudentCount: class.StudentCount,
		Status:       class.Status,
		CreateTime:   class.Ctime,
		UpdateTime:   class.Mtime,
	}
}

// AddClassMember 添加班级成员
func (s *Service) AddClassMember(ctx context.Context, classID string, req *dto.AddClassMemberRequest) error {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("班级不存在")
	}

	// 创建班级成员
	member := models.ClassMember{
		UserID:   req.UserID,
		UserName: req.UserName,
	}

	// 添加班级成员
	return s.dao.AddClassMember(ctx, classID, member, req.Role)
}

// RemoveClassMember 移除班级成员
func (s *Service) RemoveClassMember(ctx context.Context, classID string, req *dto.RemoveClassMemberRequest) error {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return err
	}
	if class == nil {
		return errors.New("班级不存在")
	}

	// 移除班级成员
	return s.dao.RemoveClassMember(ctx, classID, req.UserID, req.Role)
}

// GetClassMembers 获取班级成员列表
func (s *Service) GetClassMembers(ctx context.Context, classID string) (*dto.ClassMembersResponse, error) {
	// 检查班级是否存在
	class, err := s.dao.GetClassByID(ctx, classID)
	if err != nil {
		return nil, err
	}
	if class == nil {
		return nil, errors.New("班级不存在")
	}

	// 构建响应
	resp := &dto.ClassMembersResponse{
		Admins:     make([]dto.ClassMemberResponse, 0),
		Teachers:   make([]dto.ClassMemberResponse, 0),
		Assistants: make([]dto.ClassMemberResponse, 0),
	}

	// 添加创建者为管理员
	if class.Creator != "" {
		// 获取创建者信息
		objID, err := primitive.ObjectIDFromHex(class.Creator)
		if err == nil {
			selector := bson.M{
				"_id": objID,
			}
			user, err := s.dao.GetOneUser(ctx, selector)
			if err == nil && user != nil {
				resp.Admins = append(resp.Admins, dto.ClassMemberResponse{
					UserID:   class.Creator, // 直接使用 class.Creator 作为 UserID
					UserName: user.Nickname,
				})
			}
		}
	}

	// 添加其他成员
	for roleStr, members := range class.Members {
		role, _ := strconv.Atoi(roleStr)
		for _, member := range members {
			memberResp := dto.ClassMemberResponse{
				UserID:   member.UserID,
				UserName: member.UserName,
			}

			switch role {
			case 1: // 管理员
				// 避免重复添加创建者
				if member.UserID != class.Creator {
					resp.Admins = append(resp.Admins, memberResp)
				}
			case 2: // 教师
				resp.Teachers = append(resp.Teachers, memberResp)
			case 3: // 助教
				resp.Assistants = append(resp.Assistants, memberResp)
			}
		}
	}

	return resp, nil
}

//// CheckUserClassPermission 检查用户是否有班级操作权限
//func (s *Service) CheckUserClassPermission(ctx context.Context, classID string, userID string, requiredRoles []int) (bool, error) {
//	// 获取用户在班级中的角色
//	roles, err := s.dao.CheckUserClassRole(ctx, classID, userID)
//	if err != nil {
//		return false, err
//	}
//
//	// 检查用户是否有所需角色
//	for _, role := range roles {
//		for _, requiredRole := range requiredRoles {
//			if role == requiredRole {
//				return true, nil
//			}
//		}
//	}
//
//	return false, nil
//}

// GetJoinRequestByID 根据ID获取加入申请
func (s *Service) GetJoinRequestByID(ctx context.Context, requestID string) (*dto.JoinRequestResponse, error) {
	// 获取申请信息
	joinRequest, err := s.dao.GetJoinRequestByID(ctx, requestID)
	if err != nil {
		return nil, err
	}
	if joinRequest == nil {
		return nil, errors.New("申请不存在")
	}

	// 转换为响应DTO
	return &dto.JoinRequestResponse{
		ID:            joinRequest.ID.Hex(),
		ClassID:       joinRequest.ClassID,
		StudentID:     joinRequest.StudentID,
		StudentName:   joinRequest.StudentName,
		StudentNumber: joinRequest.StudentNumber,
		RequestMsg:    joinRequest.RequestMsg,
		Status:        joinRequest.Status,
		ReviewerID:    joinRequest.ReviewerID,
		ReviewMsg:     joinRequest.ReviewMsg,
		ReviewTime:    joinRequest.ReviewTime,
		CreateTime:    joinRequest.Ctime,
		UpdateTime:    joinRequest.Mtime,
	}, nil
}

// GetJoinRequestList 获取加入申请列表
func (s *Service) GetJoinRequestList(ctx context.Context, page, pageSize int, filters map[string]interface{}, sorts map[string]interface{}) (*utils.RespPageQuery, error) {
	comQuery := utils.BuildQueryParamsFromValues(page, pageSize, filters, sorts)

	// 执行查询
	return s.dao.GetJoinRequestList(ctx, comQuery)
}

// GetUserRolesInClass 获取用户在班级中的角色
func (s *Service) GetUserRolesInClass(ctx context.Context, classID, userID string) ([]int, error) {
	// 检查用户是否是系统管理员
	objID, err := primitive.ObjectIDFromHex(userID)
	if err == nil {
		selector := bson.M{
			"_id": objID,
		}
		user, err := s.dao.GetOneUser(ctx, selector)
		if err == nil && user != nil && user.Role == utils.ClassRoleAdmin {
			return []int{utils.ClassRoleAdmin}, nil
		}
	}

	// 否则查询用户在班级中的角色
	return s.dao.CheckUserClassRole(ctx, classID, userID)
}

// CheckUserClassPermission 检查用户在班级中是否有指定权限
func (s *Service) CheckUserClassPermission(ctx context.Context, classID, userID string, requiredRoles []int) (bool, error) {
	// 获取用户在班级中的角色
	userRoles, err := s.GetUserRolesInClass(ctx, classID, userID)
	if err != nil {
		return false, err
	}

	// 检查用户角色是否包含所需角色
	return utils.CheckRolePermission(userRoles, requiredRoles), nil
}
