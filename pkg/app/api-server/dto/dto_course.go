package dto

// 课程相关DTO

// CreateCourseRequest 创建课程请求
type CreateCourseRequest struct {
	CourseCode  string `json:"course_code" binding:"required"` // 课程代码
	Name        string `json:"name" binding:"required"`        // 课程名称
	Department  string `json:"department" binding:"required"`  // 所属院系
	Description string `json:"description"`                    // 课程描述
	StartTime   int64  `json:"start_time" binding:"required"`  // 开课时间
	EndTime     int64  `json:"end_time" binding:"required"`    // 结课时间
}

// UpdateCourseRequest 更新课程请求
type UpdateCourseRequest struct {
	Name        string `json:"name"`        // 课程名称
	Department  string `json:"department"`  // 所属院系
	Description string `json:"description"` // 课程描述
	StartTime   int64  `json:"start_time"`  // 开课时间
	EndTime     int64  `json:"end_time"`    // 结课时间
	Status      int    `json:"status"`      // 状态
}

// CourseResponse 课程响应
type CourseResponse struct {
	ID           string                            `json:"id"`
	CourseCode   string                            `json:"course_code"`
	Name         string                            `json:"name"`
	Department   string                            `json:"department"`
	Description  string                            `json:"description"`
	Creator      string                            `json:"creator"`
	Members      map[string][]CourseMemberResponse `json:"members,omitempty"`
	StudentCount int                               `json:"student_count"`
	StartTime    int64                             `json:"start_time"`
	EndTime      int64                             `json:"end_time"`
	Status       int                               `json:"status"`
	CreateTime   int64                             `json:"create_time"`
	UpdateTime   int64                             `json:"update_time"`
}

// CourseMemberResponse 课程成员响应
type CourseMemberResponse struct {
	UserID        string `json:"user_id"`                  // 用户ID
	UserName      string `json:"user_name"`                // 用户名称
	Email         string `json:"email,omitempty"`          // 邮箱
	StudentNumber string `json:"student_number,omitempty"` // 学号（学生特有）
	ClassID       string `json:"class_id,omitempty"`       // 关联班级ID
	Status        int    `json:"status"`                   // 状态
	JoinTime      int64  `json:"join_time"`                // 加入时间
}

// AddCourseMemberRequest 添加课程成员请求
type AddCourseMemberRequest struct {
	UserID        string `json:"user_id" binding:"required"`   // 用户ID
	UserName      string `json:"user_name" binding:"required"` // 用户名称
	Email         string `json:"email,omitempty"`              // 邮箱
	StudentNumber string `json:"student_number,omitempty"`     // 学号（学生特有）
	Role          int    `json:"role" binding:"required"`      // 角色：1-管理员，2-教师，3-助教，4-学生
}

// RemoveCourseMemberRequest 移除课程成员请求
type RemoveCourseMemberRequest struct {
	UserID string `json:"user_id" binding:"required"` // 用户ID
	Role   int    `json:"role" binding:"required"`    // 角色
}

// JoinCourseRequest 加入课程请求
type JoinCourseRequest struct {
	CourseCode    string `json:"course_code" binding:"required"`    // 课程代码
	StudentID     string `json:"student_id" binding:"required"`     // 学生ID
	StudentName   string `json:"student_name" binding:"required"`   // 学生姓名
	StudentNumber string `json:"student_number" binding:"required"` // 学号
	RequestMsg    string `json:"request_msg"`                       // 申请信息
}

// ReviewJoinCourseRequest 审核加入课程申请请求
type ReviewJoinCourseRequest struct {
	Status    int    `json:"status" binding:"required"` // 状态：1-通过，2-拒绝
	ReviewMsg string `json:"review_msg"`                // 审核意见
}

// JoinCourseRequestResponse 加入课程申请响应
type JoinCourseRequestResponse struct {
	ID            string `json:"id"`             // ID
	CourseID      string `json:"course_id"`      // 课程ID
	CourseCode    string `json:"course_code"`    // 课程代码
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

// CourseMembersResponse 课程成员列表响应
type CourseMembersResponse struct {
	Admins     []CourseMemberResponse `json:"admins"`     // 管理员列表
	Teachers   []CourseMemberResponse `json:"teachers"`   // 教师列表
	Assistants []CourseMemberResponse `json:"assistants"` // 助教列表
	Students   []CourseMemberResponse `json:"students"`   // 学生列表
}
