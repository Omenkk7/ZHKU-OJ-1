package dto

import (
	"zhku-oj-server/pkg/models"
)

// CreateAssignmentReq 创建作业请求
type CreateAssignmentReq struct {
	Title            string             `json:"title" binding:"required"`      // 作业标题
	Description      string             `json:"description"`                   // 作业描述
	CourseID         string             `json:"course_id" binding:"required"`  // 关联课程ID
	ClassIDs         []models.ClassInfo `json:"class_ids"`                     // 关联班级ID列表
	AssignmentType   int                `json:"assignment_type"`               // 作业类型：0-普通作业，1-OJ作业
	StartTime        int64              `json:"start_time" binding:"required"` // 开始时间
	EndTime          int64              `json:"end_time" binding:"required"`   // 截止时间
	MaxResubmitCount int                `json:"max_resubmit_count"`            // 允许重做的最大次数
	OJProblems       []string           `json:"oj_problems"`                   // 关联的OJ题目ID列表
}

// UpdateAssignmentReq 更新作业请求
type UpdateAssignmentReq struct {
	Title            string             `json:"title"`              // 作业标题
	Description      string             `json:"description"`        // 作业描述
	ClassIDs         []models.ClassInfo `json:"class_ids"`          // 关联班级ID列表
	StartTime        int64              `json:"start_time"`         // 开始时间
	EndTime          int64              `json:"end_time"`           // 截止时间
	MaxResubmitCount int                `json:"max_resubmit_count"` // 允许重做的最大次数
	OJProblems       []string           `json:"oj_problems"`        // 关联的OJ题目ID列表
	Status           int                `json:"status"`             // 作业状态：0-已归档，1-进行中
}

// SubmitAssignmentReq 提交作业请求
type SubmitAssignmentReq struct {
	TextContent string `json:"text_content"` // 作业文本内容
	FileURL     string `json:"file_url"`     // 文件URL
}

// GradeAssignmentReq 批改作业请求
type GradeAssignmentReq struct {
	AssignmentID string `json:"assignment_id" binding:"required"` // 作业ID
	StudentID    string `json:"student_id" binding:"required"`    // 学生ID
	Score        int    `json:"score" binding:"required"`         // 分数
	Comment      string `json:"comment"`                          // 评语
}

// RejectAssignmentReq 打回作业请求
type RejectAssignmentReq struct {
	AssignmentID string `json:"assignment_id" binding:"required"` // 作业ID
	StudentID    string `json:"student_id" binding:"required"`    // 学生ID
	Comment      string `json:"comment" binding:"required"`       // 评语
}

// GetAssignmentListReq 获取作业列表请求
type GetAssignmentListReq struct {
	CourseID  string `json:"course_id" form:"course_id"`   // 课程ID
	ClassID   string `json:"class_id" form:"class_id"`     // 班级ID
	Status    *int   `json:"status" form:"status"`         // 作业状态
	Type      *int   `json:"type" form:"type"`             // 作业类型
	CreatorID string `json:"creator_id" form:"creator_id"` // 创建者ID
	Keyword   string `json:"keyword" form:"keyword"`       // 关键词
	Page      int    `json:"page" form:"page"`             // 页码
	PageSize  int    `json:"page_size" form:"page_size"`   // 每页数量
}

// GetStudentAssignmentsReq 获取学生作业列表请求
type GetStudentAssignmentsReq struct {
	CourseID string `json:"course_id" form:"course_id"` // 课程ID
	Status   int    `json:"status" form:"status"`       // 提交状态
	Page     int    `json:"page" form:"page"`           // 页码
	PageSize int    `json:"page_size" form:"page_size"` // 每页数量
}

// AssignmentDetailResp 作业详情响应
type AssignmentDetailResp struct {
	ID               string                 `json:"id"`                 // 作业ID
	AssignmentCode   string                 `json:"assignment_code"`    // 作业代码
	Title            string                 `json:"title"`              // 作业标题
	Description      string                 `json:"description"`        // 作业描述
	CourseID         string                 `json:"course_id"`          // 关联课程ID
	CourseCode       string                 `json:"course_code"`        // 关联课程代码
	CourseName       string                 `json:"course_name"`        // 课程名称
	ClassIDs         []models.ClassInfo     `json:"class_ids"`          // 关联班级ID列表
	CreatorID        string                 `json:"creator_id"`         // 创建者ID
	CreatorName      string                 `json:"creator_name"`       // 创建者姓名
	AssignmentType   int                    `json:"assignment_type"`    // 作业类型
	StartTime        int64                  `json:"start_time"`         // 开始时间
	EndTime          int64                  `json:"end_time"`           // 截止时间
	MaxResubmitCount int                    `json:"max_resubmit_count"` // 允许重做的最大次数
	OJProblems       []string               `json:"oj_problems"`        // 关联的OJ题目ID列表
	Status           int                    `json:"status"`             // 作业状态
	Stats            models.AssignmentStats `json:"stats"`              // 统计信息
	Ctime            int64                  `json:"ctime"`              // 创建时间
	Mtime            int64                  `json:"mtime"`              // 修改时间
}

// StudentSubmissionResp 学生提交信息响应
type StudentSubmissionResp struct {
	Status      int    `json:"status"`       // 提交状态
	TextContent string `json:"text_content"` // 作业文本内容
	FileURL     string `json:"file_url"`     // 文件URL
	Score       int    `json:"score"`        // 分数
	Comment     string `json:"comment"`      // 评语
	GraderName  string `json:"grader_name"`  // 批改人姓名
	SubmitTime  int64  `json:"submit_time"`  // 提交时间
	GradeTime   int64  `json:"grade_time"`   // 批改时间
}

type GetStudentSubmissionReq struct {
	AssignmentID string `json:"assignment_id" form:"assignment_id" binding:"required"` // 作业ID
	StudentID    string `json:"student_id" form:"student_id" binding:"required"`       // 学生ID
	CourseID     string `json:"course_id" form:"course_id"`                            // 课程ID（可选）
	Status       int    `json:"status" form:"status"`                                  // 状态（可选）
	Page         int    `json:"page" form:"page"`                                      // 页码
	PageSize     int    `json:"page_size" form:"page_size"`                            // 每页数量
}

// ExportGradesReq 导出成绩请求
type ExportGradesReq struct {
	AssignmentID string `json:"assignment_id" form:"assignment_id" binding:"required"` // 作业ID
}

type AddStudentsToAssignmentReq struct {
	AssignmentID string                   `json:"assignment_id" binding:"required"` // 作业ID
	Students     []AddStudentToAssignment `json:"students" binding:"required"`      // 学生列表
}

// AddStudentToAssignment 添加到作业的学生信息
type AddStudentToAssignment struct {
	StudentID   string `json:"student_id" binding:"required"`   // 学生ID
	StudentName string `json:"student_name" binding:"required"` // 学生姓名
}
