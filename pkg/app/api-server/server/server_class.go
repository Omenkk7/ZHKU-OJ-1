package server

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net/http"
	"strconv"
	"zhku-oj-server/pkg/app/api-server/dto"
	"zhku-oj-server/pkg/utils"
)

/*
// createClass 创建班级
func (s *Server) createClass(c *gin.Context) {
	// contextUser中获取ID
	var userID string
	if contextUser, exists := c.Get("contextUser"); exists {
		if user, ok := contextUser.(*utils.ContextUser); ok {
			userID = user.ID
		}
	}

	// 解析请求参数
	var req dto.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层创建班级
	resp, err := s.svc.CreateClass(c, &req, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}*/
// createClass 创建班级
func (s *Server) createClass(c *gin.Context) {
	// contextUser中获取ID和角色
	var userID string
	var userRole int
	if contextUser, exists := c.Get("contextUser"); exists {
		if user, ok := contextUser.(*utils.ContextUser); ok {
			userID = user.ID
			userRole = user.Role
			// 添加日志输出用户角色
			log.Printf("创建班级 - 用户ID: %s, 角色ID: %d", userID, userRole)
		}
	}

	// 使用Casbin检查权限
	role := utils.GetRoleName(userRole) // 确保使用正确的函数名
	log.Printf("用户角色名称: %s", role)      // 添加日志输出角色名称

	if !utils.CheckPermission(role, utils.ObjClass, utils.ActCreate) {
		log.Printf("权限检查失败 - 角色: %s, 对象: %s, 操作: %s", role, utils.ObjClass, utils.ActCreate)
		utils.Forbidden(c, utils.ErrNoPermission)
		return
	}

	// 解析请求参数
	var req dto.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层创建班级
	resp, err := s.svc.CreateClass(c, &req, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

/*// updateClass 更新班级
func (s *Server) updateClass(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 检查权限（只有管理员和教师可以更新班级）
	hasPermission, err := s.svc.CheckUserClassPermission(c, classID, userID, []int{1, 2})
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}
	if !hasPermission {
		utils.FailedResponse(c, http.StatusForbidden, utils.NoPermissionErr)
		return
	}

	// 解析请求参数
	var req dto.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层更新班级
	resp, err := s.svc.UpdateClass(c, classID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}*/
// updateClass 更新班级
func (s *Server) updateClass(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 获取当前用户ID和角色
	userID := c.GetString("userId")

	// 获取用户在班级中的角色
	userRoles, err := s.svc.GetUserRolesInClass(c, classID, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 使用Casbin检查权限
	if !utils.CheckClassPermission(c, userRoles, classID, utils.ActUpdate) {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 解析请求参数
	var req dto.UpdateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层更新班级
	resp, err := s.svc.UpdateClass(c, classID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

/*// deleteClass 删除班级
func (s *Server) deleteClass(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 检查权限（只有管理员可以删除班级）
	hasPermission, err := s.svc.CheckUserClassPermission(c, classID, userID, []int{1})
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}
	if !hasPermission {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 调用service层删除班级
	err = s.svc.DeleteClass(c, classID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}*/
// deleteClass 删除班级
func (s *Server) deleteClass(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 获取用户在班级中的角色
	userRoles, err := s.svc.GetUserRolesInClass(c, classID, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 使用Casbin检查权限
	if !utils.CheckClassPermission(c, userRoles, classID, utils.ActDelete) {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 调用service层删除班级
	err = s.svc.DeleteClass(c, classID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

/*// archiveClass 归档班级
func (s *Server) archiveClass(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 调用service层归档班级
	err := s.svc.ArchiveClass(c, classID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}*/
// archiveClass 归档班级
func (s *Server) archiveClass(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 获取用户在班级中的角色
	userRoles, err := s.svc.GetUserRolesInClass(c, classID, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 使用Casbin检查权限
	if !utils.CheckClassPermission(c, userRoles, classID, utils.ActArchive) {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 调用service层归档班级
	err = s.svc.ArchiveClass(c, classID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// getClassByID 根据ID获取班级
func (s *Server) getClassByID(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 调用service层获取班级
	resp, err := s.svc.GetClassByID(c, classID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// getClassByCode 根据班级代码获取班级
func (s *Server) getClassByCode(c *gin.Context) {
	// 获取班级代码
	classCode := c.Param("code")

	// 调用service层获取班级
	resp, err := s.svc.GetClassByCode(c, classCode)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// getClassList 获取班级列表
func (s *Server) getClassList(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 解析过滤条件
	filters := make(map[string]interface{})

	// 班级名称模糊查询
	if name := c.Query("name"); name != "" {
		filters["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	// 班级代码模糊查询
	if code := c.Query("class_code"); code != "" {
		filters["class_code"] = bson.M{"$regex": code, "$options": "i"}
	}

	// 所属院系
	if department := c.Query("department"); department != "" {
		filters["department"] = department
	}

	// 状态
	if status := c.Query("status"); status != "" {
		statusInt, _ := strconv.Atoi(status)
		filters["status"] = statusInt
	}

	// 创建者ID
	if creator := c.Query("creator"); creator != "" {
		filters["creator"] = creator
	}

	// 排序条件
	sorts := make(map[string]interface{})
	if sortField := c.Query("sortField"); sortField != "" {
		sortOrder := 1 // 默认升序
		if c.Query("sortOrder") == "desc" {
			sortOrder = -1
		}
		sorts[sortField] = sortOrder
	} else {
		// 默认按创建时间降序
		sorts["ctime"] = -1
	}

	// 调用service层获取班级列表
	resp, err := s.svc.GetClassList(c, page, pageSize, filters, sorts)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

/*// addStudentToClass 添加学生到班级
func (s *Server) addStudentToClass(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 检查权限（管理员、教师和助教可以添加学生）
	hasPermission, err := s.svc.CheckUserClassPermission(c, classID, userID, []int{1, 2, 3})
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}
	if !hasPermission {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 解析请求参数
	var req dto.AddStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层添加学生
	resp, err := s.svc.AddStudentToClass(c, classID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}*/
// addStudentToClass 添加学生到班级
func (s *Server) addStudentToClass(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 获取用户在班级中的角色
	userRoles, err := s.svc.GetUserRolesInClass(c, classID, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 使用Casbin检查权限
	if !utils.CheckClassPermission(c, userRoles, classID, utils.ActAddMem) {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 解析请求参数
	var req dto.AddStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层添加学生
	resp, err := s.svc.AddStudentToClass(c, classID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

/*// batchAddStudentsToClass 批量添加学生到班级
func (s *Server) batchAddStudentsToClass(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 解析请求参数
	var req dto.BatchAddStudentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层批量添加学生
	err := s.svc.BatchAddStudentsToClass(c, classID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}*/
// batchAddStudentsToClass 批量添加学生到班级
func (s *Server) batchAddStudentsToClass(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 获取用户在班级中的角色
	userRoles, err := s.svc.GetUserRolesInClass(c, classID, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 使用Casbin检查权限
	if !utils.CheckClassPermission(c, userRoles, classID, utils.ActAddMem) {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 解析请求参数
	var req dto.BatchAddStudentsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层批量添加学生
	err = s.svc.BatchAddStudentsToClass(c, classID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

/*// removeStudentFromClass 从班级移除学生
func (s *Server) removeStudentFromClass(c *gin.Context) {
	// 获取班级ID和学生ID
	classID := c.Param("id")
	studentID := c.Param("studentId")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 检查权限（管理员和教师可以移除学生）
	hasPermission, err := s.svc.CheckUserClassPermission(c, classID, userID, []int{1, 2})
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}
	if !hasPermission {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 调用service层移除学生
	err = s.svc.RemoveStudentFromClass(c, classID, studentID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}*/
// removeStudentFromClass 从班级移除学生
func (s *Server) removeStudentFromClass(c *gin.Context) {
	// 获取班级ID和学生ID
	classID := c.Param("id")
	studentID := c.Param("studentId")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 获取用户在班级中的角色
	userRoles, err := s.svc.GetUserRolesInClass(c, classID, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 使用Casbin检查权限
	if !utils.CheckClassPermission(c, userRoles, classID, utils.ActRemMem) {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 调用service层移除学生
	err = s.svc.RemoveStudentFromClass(c, classID, studentID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// getClassStudents 获取班级学生列表
func (s *Server) getClassStudents(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 解析过滤条件
	filters := make(map[string]interface{})

	// 学生姓名模糊查询
	if name := c.Query("student_name"); name != "" {
		filters["student_name"] = bson.M{"$regex": name, "$options": "i"}
	}

	// 学号模糊查询
	if number := c.Query("student_number"); number != "" {
		filters["student_number"] = bson.M{"$regex": number, "$options": "i"}
	}

	// 状态
	if status := c.Query("status"); status != "" {
		statusInt, _ := strconv.Atoi(status)
		filters["status"] = statusInt
	}

	// 加入方式
	if joinType := c.Query("join_type"); joinType != "" {
		joinTypeInt, _ := strconv.Atoi(joinType)
		filters["join_type"] = joinTypeInt
	}

	// 排序条件
	sorts := make(map[string]interface{})
	if sortField := c.Query("sortField"); sortField != "" {
		sortOrder := 1 // 默认升序
		if c.Query("sortOrder") == "desc" {
			sortOrder = -1
		}
		sorts[sortField] = sortOrder
	} else {
		// 默认按加入时间降序
		sorts["join_time"] = -1
	}

	// 调用service层获取班级学生列表
	resp, err := s.svc.GetClassStudents(c, classID, page, pageSize, filters, sorts)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// addCourseToClass 添加课程到班级
func (s *Server) addCourseToClass(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 解析请求参数
	var req dto.AddCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层添加课程
	err := s.svc.AddCourseToClass(c, classID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

/*// removeCourseFromClass 从班级移除课程
func (s *Server) removeCourseFromClass(c *gin.Context) {
	// 获取班级ID和课程ID
	classID := c.Param("id")
	courseID := c.Param("courseId")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 检查权限（只有管理员和教师可以移除课程）
	hasPermission, err := s.svc.CheckUserClassPermission(c, classID, userID, []int{1, 2})
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}
	if !hasPermission {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 调用service层移除课程
	err = s.svc.RemoveCourseFromClass(c, classID, courseID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}*/
// removeCourseFromClass 从班级移除课程
func (s *Server) removeCourseFromClass(c *gin.Context) {
	// 获取班级ID和课程ID
	classID := c.Param("id")
	courseID := c.Param("courseId")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 获取用户在班级中的角色
	userRoles, err := s.svc.GetUserRolesInClass(c, classID, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 使用Casbin检查权限
	if !utils.CheckClassPermission(c, userRoles, classID, utils.ActUpdate) {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 调用service层移除课程
	err = s.svc.RemoveCourseFromClass(c, classID, courseID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

/*// updateCourseStatusInClass 更新班级中课程的状态
func (s *Server) updateCourseStatusInClass(c *gin.Context) {
	// 获取班级ID和课程ID
	classID := c.Param("id")
	courseID := c.Param("courseId")

	// 解析请求参数
	var req dto.UpdateCourseStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层更新课程状态
	err := s.svc.UpdateCourseStatusInClass(c, classID, courseID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}*/
// updateCourseStatusInClass 更新班级中课程的状态
func (s *Server) updateCourseStatusInClass(c *gin.Context) {
	// 获取班级ID和课程ID
	classID := c.Param("id")
	courseID := c.Param("courseId")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 获取用户在班级中的角色
	userRoles, err := s.svc.GetUserRolesInClass(c, classID, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 使用Casbin检查权限
	if !utils.CheckClassPermission(c, userRoles, classID, utils.ActUpdate) {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 解析请求参数
	var req dto.UpdateCourseStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层更新课程状态
	err = s.svc.UpdateCourseStatusInClass(c, classID, courseID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// getClassesByCourseID 获取绑定了指定课程的班级列表
func (s *Server) getClassesByCourseID(c *gin.Context) {
	// 获取课程ID
	courseID := c.Param("courseId")

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 解析过滤条件
	filters := make(map[string]interface{})

	// 班级名称模糊查询
	if name := c.Query("name"); name != "" {
		filters["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	// 班级代码模糊查询
	if code := c.Query("class_code"); code != "" {
		filters["class_code"] = bson.M{"$regex": code, "$options": "i"}
	}

	// 所属院系
	if department := c.Query("department"); department != "" {
		filters["department"] = department
	}

	// 状态
	if status := c.Query("status"); status != "" {
		statusInt, _ := strconv.Atoi(status)
		filters["status"] = statusInt
	}

	// 排序条件
	sorts := make(map[string]interface{})
	if sortField := c.Query("sortField"); sortField != "" {
		sortOrder := 1 // 默认升序
		if c.Query("sortOrder") == "desc" {
			sortOrder = -1
		}
		sorts[sortField] = sortOrder
	} else {
		// 默认按创建时间降序
		sorts["ctime"] = -1
	}

	// 调用service层获取班级列表
	resp, err := s.svc.GetClassesByCourseID(c, courseID, page, pageSize, filters, sorts)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// createJoinRequest 创建加入班级申请
func (s *Server) createJoinRequest(c *gin.Context) {
	// 解析请求参数
	var req dto.JoinClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层创建加入申请
	resp, err := s.svc.CreateJoinRequest(c, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

/*// reviewJoinRequest 审核加入班级申请
func (s *Server) reviewJoinRequest(c *gin.Context) {
	// 获取申请ID
	requestID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 获取申请信息
	joinRequest, err := s.svc.GetJoinRequestByID(c, requestID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 检查权限（只有管理员、教师和助教可以审核申请）
	hasPermission, err := s.svc.CheckUserClassPermission(c, joinRequest.ClassID, userID, []int{1, 2, 3})
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}
	if !hasPermission {
		utils.FailedResponse(c, http.StatusForbidden, utils.NoPermissionErr)
		return
	}

	// 解析请求参数
	var req dto.ReviewJoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层审核申请
	resp, err := s.svc.ReviewJoinRequest(c, requestID, &req, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}*/
// reviewJoinRequest 审核加入班级申请
func (s *Server) reviewJoinRequest(c *gin.Context) {
	// 获取申请ID
	requestID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 获取申请信息
	joinRequest, err := s.svc.GetJoinRequestByID(c, requestID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 获取用户在班级中的角色
	userRoles, err := s.svc.GetUserRolesInClass(c, joinRequest.ClassID, userID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// 使用Casbin检查权限
	if !utils.CheckClassPermission(c, userRoles, joinRequest.ClassID, utils.ActReview) {
		utils.FailedResponse(c, http.StatusForbidden, utils.ErrNoPermission)
		return
	}

	// 解析请求参数
	var req dto.ReviewJoinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 调用service层审核申请
	resp, err := s.svc.ReviewJoinRequest(c, requestID, &req, userID)
	if err != nil {
		// 根据错误类型返回不同的HTTP状态码
		switch err.Error() {
		case "申请不存在":
			utils.FailedResponse(c, http.StatusNotFound, err)
		case "申请已处理":
			// 对于"申请已处理"的情况，返回400错误而不是500
			utils.BadRequest(c, err)
		default:
			utils.FailedResponse(c, http.StatusInternalServerError, err)
		}
		return
	}

	utils.SuccessResponse(c, resp)
}

// getJoinRequestList 获取加入申请列表
func (s *Server) getJoinRequestList(c *gin.Context) {
	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 解析过滤条件
	filters := make(map[string]interface{})

	// 班级ID
	if classID := c.Query("class_id"); classID != "" {
		filters["class_id"] = classID
	}

	// 学生ID
	if studentID := c.Query("student_id"); studentID != "" {
		filters["student_id"] = studentID
	}

	// 学生姓名模糊查询
	if name := c.Query("student_name"); name != "" {
		filters["student_name"] = bson.M{"$regex": name, "$options": "i"}
	}

	// 学号模糊查询
	if number := c.Query("student_number"); number != "" {
		filters["student_number"] = bson.M{"$regex": number, "$options": "i"}
	}

	// 状态
	if status := c.Query("status"); status != "" {
		statusInt, _ := strconv.Atoi(status)
		filters["status"] = statusInt
	}

	// 排序条件
	sorts := make(map[string]interface{})
	if sortField := c.Query("sortField"); sortField != "" {
		sortOrder := 1 // 默认升序
		if c.Query("sortOrder") == "desc" {
			sortOrder = -1
		}
		sorts[sortField] = sortOrder
	} else {
		// 默认按创建时间降序
		sorts["ctime"] = -1
	}

	// 调用service层获取申请列表
	resp, err := s.svc.GetJoinRequestList(c, page, pageSize, filters, sorts)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// getStudentClasses 获取学生所在的班级列表
func (s *Server) getStudentClasses(c *gin.Context) {
	// 获取学生ID
	studentID := c.Param("studentId")

	// 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	// 解析过滤条件
	filters := make(map[string]interface{})

	// 班级名称模糊查询
	if name := c.Query("name"); name != "" {
		filters["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	// 班级代码模糊查询
	if code := c.Query("class_code"); code != "" {
		filters["class_code"] = bson.M{"$regex": code, "$options": "i"}
	}

	// 所属院系
	if department := c.Query("department"); department != "" {
		filters["department"] = department
	}

	// 状态
	if status := c.Query("status"); status != "" {
		statusInt, _ := strconv.Atoi(status)
		filters["status"] = statusInt
	}

	// 排序条件
	sorts := make(map[string]interface{})
	if sortField := c.Query("sortField"); sortField != "" {
		sortOrder := 1 // 默认升序
		if c.Query("sortOrder") == "desc" {
			sortOrder = -1
		}
		sorts[sortField] = sortOrder
	} else {
		// 默认按创建时间降序
		sorts["ctime"] = -1
	}

	// 调用service层获取学生班级列表
	resp, err := s.svc.GetStudentClasses(c, studentID, page, pageSize, filters, sorts)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}

// addClassMember 添加班级成员
func (s *Server) addClassMember(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 解析请求参数
	var req dto.AddClassMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 检查权限（只有管理员和教师可以添加成员）
	hasPermission, err := s.svc.CheckUserClassPermission(c, classID, userID, []int{1, 2})
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}
	if !hasPermission {
		utils.FailedResponse(c, http.StatusForbidden, utils.NoPermissionErr)
		return
	}

	// 调用service层添加班级成员
	err = s.svc.AddClassMember(c, classID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// removeClassMember 移除班级成员
func (s *Server) removeClassMember(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 解析请求参数
	var req dto.RemoveClassMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err)
		return
	}

	// 检查权限（只有管理员可以移除成员）
	hasPermission, err := s.svc.CheckUserClassPermission(c, classID, userID, []int{1})
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}
	if !hasPermission {
		utils.FailedResponse(c, http.StatusForbidden, utils.NoPermissionErr)
		return
	}

	// 调用service层移除班级成员
	err = s.svc.RemoveClassMember(c, classID, &req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, nil)
}

// getClassMembers 获取班级成员列表
func (s *Server) getClassMembers(c *gin.Context) {
	// 获取班级ID
	classID := c.Param("id")

	// 获取当前用户ID
	userID := c.GetString("userId")

	// 检查权限（班级成员可以查看成员列表）
	hasPermission, err := s.svc.CheckUserClassPermission(c, classID, userID, []int{1, 2, 3, 4})
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}
	if !hasPermission {
		utils.FailedResponse(c, http.StatusForbidden, utils.NoPermissionErr)
		return
	}

	// 调用service层获取班级成员列表
	resp, err := s.svc.GetClassMembers(c, classID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, resp)
}
