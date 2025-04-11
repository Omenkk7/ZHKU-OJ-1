package dto

import "zhku-oj-server/pkg/models"

// 班级相关DTO

// CreateClassRequest 创建班级请求
type CreateClassRequest struct {
	ClassCode   string              `json:"class_code" binding:"required"` // 班级代码
	Name        string              `json:"name" binding:"required"`       // 班级名称
	Department  string              `json:"department" binding:"required"` // 所属院系
	Description string              `json:"description"`                   // 班级描述
	Courses     []models.CourseInfo `json:"courses"`                       // 关联课程
}

// UpdateClassRequest 更新班级请求
type UpdateClassRequest struct {
	Name        string              `json:"name"`        // 班级名称
	Department  string              `json:"department"`  // 所属院系
	Description string              `json:"description"` // 班级描述
	Courses     []models.CourseInfo `json:"courses"`     // 关联课程
}

// ClassResponse
type ClassResponse struct {
	ID           string                           `json:"id"`
	ClassCode    string                           `json:"class_code"`
	Name         string                           `json:"name"`
	Department   string                           `json:"department"`
	Description  string                           `json:"description"`
	Creator      string                           `json:"creator"`
	Courses      []models.CourseInfo              `json:"courses"`
	Members      map[string][]ClassMemberResponse `json:"members,omitempty"`
	StudentCount int                              `json:"student_count"`
	Status       int                              `json:"status"`
	CreateTime   int64                            `json:"create_time"`
	UpdateTime   int64                            `json:"update_time"`
}

// AddStudentRequest 添加学生请求
type AddStudentRequest struct {
	StudentID     string `json:"student_id" binding:"required"`     // 学生ID
	StudentName   string `json:"student_name" binding:"required"`   // 学生姓名
	StudentNumber string `json:"student_number" binding:"required"` // 学号
	JoinType      int    `json:"join_type" binding:"required"`      // 加入方式
}

// BatchAddStudentsRequest 批量添加学生请求
type BatchAddStudentsRequest struct {
	Students []AddStudentRequest `json:"students" binding:"required"` // 学生列表
}

// ClassStudentResponse 班级学生响应
type ClassStudentResponse struct {
	ID            string `json:"id"`             // ID
	ClassID       string `json:"class_id"`       // 班级ID
	StudentID     string `json:"student_id"`     // 学生ID
	StudentName   string `json:"student_name"`   // 学生姓名
	StudentNumber string `json:"student_number"` // 学号
	JoinType      int    `json:"join_type"`      // 加入方式
	Status        int    `json:"status"`         // 状态
	JoinTime      int64  `json:"join_time"`      // 加入时间
	LeaveTime     *int64 `json:"leave_time"`     // 离开时间
	CreateTime    int64  `json:"create_time"`    // 创建时间
	UpdateTime    int64  `json:"update_time"`    // 更新时间
}

// JoinClassRequest 加入班级请求
type JoinClassRequest struct {
	ClassCode     string `json:"class_code" binding:"required"`     // 班级代码
	StudentID     string `json:"student_id" binding:"required"`     // 学生ID
	StudentName   string `json:"student_name" binding:"required"`   // 学生姓名
	StudentNumber string `json:"student_number" binding:"required"` // 学号
	RequestMsg    string `json:"request_msg"`                       // 申请信息
}

// ReviewJoinRequest 审核加入申请请求
type ReviewJoinRequest struct {
	Status    int    `json:"status" binding:"required"` // 状态：1-通过，2-拒绝
	ReviewMsg string `json:"review_msg"`                // 审核意见
}

// JoinRequestResponse 加入申请响应
type JoinRequestResponse struct {
	ID            string `json:"id"`             // ID
	ClassID       string `json:"class_id"`       // 班级ID
	StudentID     string `json:"student_id"`     // 学生ID
	StudentName   string `json:"student_name"`   // 学生姓名
	StudentNumber string `json:"student_number"` // 学号
	RequestMsg    string `json:"request_msg"`    // 申请信息
	Status        int    `json:"status"`         // 状态
	ReviewerID    string `json:"reviewer_id"`    // 审核人ID
	ReviewMsg     string `json:"review_msg"`     // 审核意见
	ReviewTime    *int64 `json:"review_time"`    // 审核时间
	CreateTime    int64  `json:"create_time"`    // 创建时间
	UpdateTime    int64  `json:"update_time"`    // 更新时间
}

// AddCourseRequest 添加课程请求
type AddCourseRequest struct {
	CourseID   string `json:"course_id" binding:"required"`   // 课程ID
	CourseName string `json:"course_name" binding:"required"` // 课程名称
	Status     int    `json:"status" binding:"required"`      // 状态
}

// UpdateCourseStatusRequest 更新课程状态请求
type UpdateCourseStatusRequest struct {
	Status int `json:"status" binding:"required"` // 状态
}

// AddClassMemberRequest 添加班级成员请求
type AddClassMemberRequest struct {
	UserID   string `json:"user_id" binding:"required"`   // 用户ID
	UserName string `json:"user_name" binding:"required"` // 用户名称
	Role     int    `json:"role" binding:"required"`      // 角色：1-管理员，2-教师，3-助教
}

// RemoveClassMemberRequest 移除班级成员请求
type RemoveClassMemberRequest struct {
	UserID string `json:"user_id" binding:"required"` // 用户ID
	Role   int    `json:"role" binding:"required"`    // 角色：1-管理员，2-教师，3-助教
}

// ClassMemberResponse 班级成员响应
type ClassMemberResponse struct {
	UserID   string `json:"user_id"`   // 用户ID
	UserName string `json:"user_name"` // 用户名称
}

// ClassMembersResponse 班级成员列表响应
type ClassMembersResponse struct {
	Admins     []ClassMemberResponse `json:"admins"`     // 管理员列表
	Teachers   []ClassMemberResponse `json:"teachers"`   // 教师列表
	Assistants []ClassMemberResponse `json:"assistants"` // 助教列表
}
