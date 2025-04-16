package server

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/utils"
)

// createCourse 创建课程
func (s *Server) createCourse(c *gin.Context) {
	// 获取当前用户ID和角色
	userID := c.GetString("userId")
	userRole, _ := strconv.Atoi(c.GetString("userRole"))

	// 使用Casbin检查权限
	role := utils.GetRoleName(userRole)
	if !utils.CheckPermission(role, utils.ObjCourse, utils.ActCreate) {
		utils.Forbidden(c, utils.ErrNoPermission)
		return
	}

	// 解析请求参数
	var req dto.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 如果没有提供课程代码，自动生成
	if req.CourseCode == "" {
		courseCode, err := s.svc.GenerateCourseCode(c, req.Department)
		if err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, err)
			return
		}
		req.CourseCode = courseCode
	}

	// 调用service层创建课程
	resp, err := s.svc.CreateCourse(c, &req, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// updateCourse 更新课程
func (s *Server) updateCourse(c *gin.Context) {
	// 获取课程ID
	courseID := c.Param("id")

	// 获取当前用户ID和角色
	userID := c.GetString("userId")
	userRole, _ := strconv.Atoi(c.GetString("userRole"))

	// 使用Casbin检查权限
	role := utils.GetRoleName(userRole)
	if !utils.CheckPermission(role, utils.ObjCourse, utils.ActUpdate) {
		utils.Forbidden(c, utils.ErrNoPermission)
		return
	}

	// 检查用户是否有权限更新此课程（管理员可以更新所有课程，教师只能更新自己的课程）
	if userRole != utils.ClassRoleAdmin {
		hasPermission, err := s.svc.CheckUserCoursePermission(c, courseID, userID, []int{utils.RoleTeacher})
		if err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, err)
			return
		}
		if !hasPermission {
			utils.Forbidden(c, utils.ErrNoPermission)
			return
		}
	}

	// 解析请求参数
	var req dto.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层更新课程
	resp, err := s.svc.UpdateCourse(c, courseID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// deleteCourse 删除课程
func (s *Server) deleteCourse(c *gin.Context) {
	// 获取课程ID
	courseID := c.Param("id")

	// 获取当前用户角色
	userRole, _ := strconv.Atoi(c.GetString("userRole"))

	// 使用Casbin检查权限（只有管理员可以删除课程）
	role := utils.GetRoleName(userRole)
	if !utils.CheckPermission(role, utils.ObjCourse, utils.ActDelete) {
		utils.Forbidden(c, utils.ErrNoPermission)
		return
	}

	// 调用service层删除课程
	err := s.svc.DeleteCourse(c, courseID) // 修复：声明变量err
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// archiveCourse 归档课程
func (s *Server) archiveCourse(c *gin.Context) {
	// 获取课程ID
	courseID := c.Param("id")

	// 获取当前用户ID和角色
	userID := c.GetString("userId")
	userRole, _ := strconv.Atoi(c.GetString("userRole"))

	// 使用Casbin检查权限（只有管理员和教师可以归档课程）
	role := utils.GetRoleName(userRole)
	if !utils.CheckPermission(role, utils.ObjCourse, utils.ActUpdate) {
		utils.Forbidden(c, utils.ErrNoPermission)
		return
	}

	// 检查用户是否有权限归档此课程（管理员可以归档所有课程，教师只能归档自己的课程）
	if userRole != utils.ClassRoleAdmin {
		hasPermission, err := s.svc.CheckUserCoursePermission(c, courseID, userID, []int{utils.RoleTeacher})
		if err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, err)
			return
		}
		if !hasPermission {
			utils.Forbidden(c, utils.ErrNoPermission)
			return
		}
	}

	// 调用service层归档课程
	err := s.svc.ArchiveCourse(c, courseID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// getCourse 获取课程详情
func (s *Server) getCourse(c *gin.Context) {
	// 获取课程ID
	courseID := c.Param("id")

	// 获取当前用户ID和角色
	userID := c.GetString("userId")
	userRole, _ := strconv.Atoi(c.GetString("userRole"))

	// 使用Casbin检查权限
	role := utils.GetRoleName(userRole)
	if !utils.CheckPermission(role, utils.ObjCourse, utils.ActRead) {
		utils.Forbidden(c, utils.ErrNoPermission)
		return
	}

	// 调用service层获取课程
	course, err := s.svc.GetCourseByID(c, courseID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 检查用户是否有权限查看此课程（管理员可以查看所有课程，其他角色只能查看自己参与的课程）
	if userRole != utils.ClassRoleAdmin && userID != course.Creator {
		hasPermission, err := s.svc.CheckUserCoursePermission(c, courseID, userID, []int{utils.RoleTeacher, utils.RoleAssistant, utils.RoleStudent})
		if err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, err)
			return
		}
		if !hasPermission {
			utils.Forbidden(c, utils.ErrNoPermission)
			return
		}
	}

	utils.SuccessResponse(c, course)
}

// getCourseByCode 根据课程代码获取课程
func (s *Server) getCourseByCode(c *gin.Context) {
	// 获取课程代码
	courseCode := c.Query("code")
	if courseCode == "" {
		utils.BadRequest(c, errors.New("课程代码不能为空"))
		return
	}

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 调用service层获取课程
	course, err := s.svc.GetCourseByCode(c, courseCode)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 检查权限（管理员可以查看所有课程，其他角色只能查看自己参与的课程）
	if userID != course.Creator {
		hasPermission, err := s.svc.CheckUserCoursePermission(c, course.ID, userID, []int{1, 2, 3, 4})
		if err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, err)
			return
		}
		if !hasPermission {
			utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
			return
		}
	}

	utils.SuccessResponse(c, course)
}

// getCourseList 获取课程列表
func (s *Server) getCourseList(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 获取筛选条件
	name := c.Query("name")
	code := c.Query("code")
	department := c.Query("department")
	status := c.Query("status")

	// 构建筛选条件
	filters := make(map[string]interface{})
	if name != "" {
		filters["name"] = bson.M{"$regex": name, "$options": "i"} // i表示不区分大小写
	}
	if code != "" {
		filters["course_code"] = bson.M{"$regex": code, "$options": "i"}
	}
	if department != "" {
		filters["department"] = department
	}
	if status != "" {
		statusInt, err := strconv.Atoi(status)
		if err == nil {
			filters["status"] = statusInt
		}
	}

	// 获取排序条件
	sortField := c.DefaultQuery("sort_field", "ctime")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	sorts := map[string]interface{}{
		sortField: sortOrder,
	}

	// 获取当前用户ID和角色
	userID := c.GetString("userId")
	userRole, _ := strconv.Atoi(c.GetString("userRole"))

	var resp *utils.RespPageQuery
	var err error

	// 根据角色进行不同的处理
	if userRole == utils.RoleAdmin {
		// 管理员可以查看所有课程
		resp, err = s.svc.GetCourseList(c, page, pageSize, filters, sorts)
	} else {
		// 非管理员只能查看自己参与的课程
		resp, err = s.svc.GetUserCourses(c, userID, page, pageSize, filters, sorts)
	}

	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// addCourseMember 添加课程成员
func (s *Server) addCourseMember(c *gin.Context) {
	// 获取课程ID
	courseID := c.Param("id")

	// 获取当前用户ID和角色
	userID := c.GetString("userId")
	userRole, _ := strconv.Atoi(c.GetString("userRole"))

	// 使用Casbin检查权限
	role := utils.GetRoleName(userRole)
	if !utils.CheckPermission(role, utils.ObjCourseMember, utils.ActAddMem) {
		utils.Forbidden(c, utils.ErrNoPermission)
		return
	}

	// 解析请求参数
	var req dto.AddCourseMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 检查用户是否有权限添加此课程的成员（管理员可以添加所有课程的成员，教师只能添加自己课程的成员）
	if userRole != utils.ClassRoleAdmin {
		hasPermission, err := s.svc.CheckUserCoursePermission(c, courseID, userID, []int{utils.RoleTeacher})
		if err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, err)
			return
		}
		if !hasPermission {
			utils.Forbidden(c, utils.ErrNoPermission)
			return
		}
	}

	// 额外检查：教师只能添加学生和助教，不能添加其他教师或管理员
	if userRole == utils.RoleTeacher && (req.Role == utils.ClassRoleAdmin || req.Role == utils.RoleTeacher) {
		utils.Forbidden(c, errors.New("教师只能添加学生和助教"))
		return
	}

	// 调用service层添加课程成员
	err := s.svc.AddCourseMember(c, courseID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// removeCourseMember 移除课程成员
func (s *Server) removeCourseMember(c *gin.Context) {
	// 获取课程ID和用户ID
	courseID := c.Param("id")
	memberID := c.Param("member_id")
	roleStr := c.Query("role")
	role, _ := strconv.Atoi(roleStr)

	// 获取当前用户ID和角色
	userID := c.GetString("userId")
	userRole, _ := strconv.Atoi(c.GetString("userRole"))

	// 使用Casbin检查权限
	roleName := utils.GetRoleName(userRole)
	if !utils.CheckPermission(roleName, utils.ObjCourseMember, utils.ActRemMem) {
		utils.Forbidden(c, utils.ErrNoPermission)
		return
	}

	// 检查用户是否有权限移除此课程的成员（管理员可以移除所有课程的成员，教师只能移除自己课程的成员）
	if userRole != utils.ClassRoleAdmin {
		hasPermission, err := s.svc.CheckUserCoursePermission(c, courseID, userID, []int{utils.RoleTeacher})
		if err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, err)
			return
		}
		if !hasPermission {
			utils.Forbidden(c, utils.ErrNoPermission)
			return
		}
	}

	// 额外检查：教师只能移除学生和助教，不能移除其他教师或管理员
	if userRole == utils.RoleTeacher && (role == utils.ClassRoleAdmin || role == utils.RoleTeacher) {
		utils.Forbidden(c, errors.New("教师只能移除学生和助教"))
		return
	}

	// 调用service层移除课程成员
	err := s.svc.RemoveCourseMember(c, courseID, memberID, role)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// getCourseMembers 获取课程成员
func (s *Server) getCourseMembers(c *gin.Context) {
	// 获取课程ID
	courseID := c.Param("id")

	// 获取当前用户ID和角色
	userID := c.GetString("userId")
	userRole, _ := strconv.Atoi(c.GetString("userRole"))

	// 使用Casbin检查权限
	role := utils.GetRoleName(userRole)
	if !utils.CheckPermission(role, utils.ObjCourseMember, utils.ActRead) {
		utils.Forbidden(c, utils.ErrNoPermission)
		return
	}

	// 检查用户是否有权限查看此课程的成员（管理员可以查看所有课程的成员，其他角色只能查看自己参与的课程的成员）
	if userRole != utils.ClassRoleAdmin {
		hasPermission, err := s.svc.CheckUserCoursePermission(c, courseID, userID, []int{utils.RoleTeacher, utils.RoleAssistant})
		if err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, err)
			return
		}
		if !hasPermission {
			utils.Forbidden(c, utils.ErrNoPermission)
			return
		}
	}

	// 调用service层获取课程成员
	resp, err := s.svc.GetCourseMembers(c, courseID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// createJoinCourseRequest 创建加入课程申请
func (s *Server) createJoinCourseRequest(c *gin.Context) {
	// 解析请求参数
	var req dto.JoinCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层创建加入课程申请
	resp, err := s.svc.CreateJoinCourseRequest(c, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// reviewJoinCourseRequest 审核加入课程申请
func (s *Server) reviewJoinCourseRequest(c *gin.Context) {
	// 获取申请ID
	requestID := c.Param("id")

	// 获取当前用户ID和角色
	userID := c.GetString("userId")
	userRole, _ := strconv.Atoi(c.GetString("userRole"))

	// 使用Casbin检查权限
	role := utils.GetRoleName(userRole)
	if !utils.CheckPermission(role, utils.ObjCourseRequest, utils.ActReview) {
		utils.Forbidden(c, utils.ErrNoPermission)
		return
	}

	// 解析请求参数
	var req dto.ReviewJoinCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 获取加入请求信息，以便检查关联的课程
	joinRequest, err := s.svc.GetJoinCourseRequestByID(c, requestID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}
	if joinRequest == nil {
		utils.BadRequest(c, errors.New("加入请求不存在"))
		return
	}

	// 检查用户是否有权限审核此课程的申请（管理员可以审核所有课程的申请，教师只能审核自己课程的申请）
	if userRole != utils.ClassRoleAdmin {
		hasPermission, err := s.svc.CheckUserCoursePermission(c, joinRequest.CourseID, userID, []int{utils.RoleTeacher})
		if err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, err)
			return
		}
		if !hasPermission {
			utils.Forbidden(c, utils.ErrNoPermission)
			return
		}
	}

	// 调用service层审核加入课程申请
	resp, err := s.svc.ReviewJoinCourseRequest(c, requestID, userID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// getJoinCourseRequestList 获取加入课程申请列表
func (s *Server) getJoinCourseRequestList(c *gin.Context) {
	// 获取课程ID
	courseID := c.Param("id")

	// 获取当前用户ID和角色
	userID := c.GetString("userId")
	userRole, _ := strconv.Atoi(c.GetString("userRole"))

	// 使用Casbin检查权限
	role := utils.GetRoleName(userRole)
	if !utils.CheckPermission(role, utils.ObjCourseRequest, utils.ActRead) {
		utils.Forbidden(c, utils.ErrNoPermission)
		return
	}

	// 检查用户是否有权限查看此课程的申请列表（管理员可以查看所有课程的申请，教师只能查看自己课程的申请）
	if userRole != utils.ClassRoleAdmin {
		hasPermission, err := s.svc.CheckUserCoursePermission(c, courseID, userID, []int{utils.RoleTeacher})
		if err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, err)
			return
		}
		if !hasPermission {
			utils.Forbidden(c, utils.ErrNoPermission)
			return
		}
	}

	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 获取筛选条件
	status := c.Query("status")
	studentName := c.Query("student_name")
	studentNumber := c.Query("student_number")

	// 构建筛选条件
	filters := make(map[string]interface{})
	if status != "" {
		statusInt, err := strconv.Atoi(status)
		if err == nil {
			filters["status"] = statusInt
		}
	}
	if studentName != "" {
		filters["student_name"] = bson.M{"$regex": studentName, "$options": "i"}
	}
	if studentNumber != "" {
		filters["student_number"] = bson.M{"$regex": studentNumber, "$options": "i"}
	}

	// 获取排序条件
	sortField := c.DefaultQuery("sort_field", "ctime")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	sorts := map[string]interface{}{
		sortField: sortOrder,
	}

	// 调用service层获取申请列表
	resp, err := s.svc.GetJoinCourseRequestList(c, courseID, page, pageSize, filters, sorts)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}
