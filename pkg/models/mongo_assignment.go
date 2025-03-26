package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Assignment 作业模型
type Assignment struct {
	ID               primitive.ObjectID  `json:"_id,omitempty" bson:"_id,omitempty"`
	AssignmentCode   string              `json:"assignment_code" bson:"assignment_code"`       // 作业代码：ASG+年月+序号
	Title            string              `json:"title" bson:"title"`                           // 作业标题
	Description      string              `json:"description" bson:"description"`               // 作业描述
	CourseID         string              `json:"course_id" bson:"course_id"`                   // 关联课程ID
	CourseCode       string              `json:"course_code" bson:"course_code"`               // 关联课程代码
	CourseName       string              `json:"course_name" bson:"course_name"`               // 课程名称(冗余字段)
	ClassIDs         []ClassInfo         `json:"class_ids" bson:"class_ids"`                   // 关联班级ID列表
	CreatorID        string              `json:"creator_id" bson:"creator_id"`                 // 创建者ID
	CreatorName      string              `json:"creator_name" bson:"creator_name"`             // 创建者姓名
	AssignmentType   int                 `json:"assignment_type" bson:"assignment_type"`       // 作业类型：0-普通作业，1-OJ作业
	StartTime        int64               `json:"start_time" bson:"start_time"`                 // 开始时间
	EndTime          int64               `json:"end_time" bson:"end_time"`                     // 截止时间
	MaxResubmitCount int                 `json:"max_resubmit_count" bson:"max_resubmit_count"` // 允许重做的最大次数
	Members          map[string][]Member `json:"members" bson:"members"`                       // 成员列表，按角色分类
	Stats            AssignmentStats     `json:"stats" bson:"stats"`                           // 统计信息
	OJProblems       []string            `json:"oj_problems" bson:"oj_problems"`               // 关联的OJ题目ID列表
	Status           int                 `json:"status" bson:"status"`                         // 作业状态：0-已归档，1-进行中
	Ctime            int64               `json:"ctime,omitempty" bson:"ctime,omitempty"`       // 创建时间
	Mtime            int64               `json:"mtime,omitempty" bson:"mtime,omitempty"`       // 修改时间
}

// ClassInfo 班级信息
type ClassInfo struct {
	ClassID   string `json:"class_id" bson:"class_id"`     // 班级ID
	ClassName string `json:"class_name" bson:"class_name"` // 班级名称
}

// Member 成员信息
type Member struct {
	UserID     string      `json:"user_id" bson:"user_id"`       // 用户ID
	UserName   string      `json:"user_name" bson:"user_name"`   // 用户名称
	Submission *Submission `json:"submission" bson:"submission"` // 提交信息
}

// Submission 提交信息
type Submission struct {
	Status      int    `json:"status" bson:"status"`             // 提交状态：0-未提交，1-已提交未批改，2-已批改，3-被打回
	TextContent string `json:"text_content" bson:"text_content"` // 作业文本内容
	FileURL     string `json:"file_url" bson:"file_url"`         // 文件URL
	Score       int    `json:"score" bson:"score"`               // 分数
	Comment     string `json:"comment" bson:"comment"`           // 评语
	GraderID    string `json:"grader_id" bson:"grader_id"`       // 批改人ID
	GraderName  string `json:"grader_name" bson:"grader_name"`   // 批改人姓名
	SubmitTime  int64  `json:"submit_time" bson:"submit_time"`   // 提交时间
	GradeTime   int64  `json:"grade_time" bson:"grade_time"`     // 批改时间
}

// AssignmentStats 作业统计信息
type AssignmentStats struct {
	TotalStudents  int `json:"total_students" bson:"total_students"`   // 总学生数
	SubmittedCount int `json:"submitted_count" bson:"submitted_count"` // 已提交数
	GradedCount    int `json:"graded_count" bson:"graded_count"`       // 已批改数
	AverageScore   int `json:"average_score" bson:"average_score"`     // 平均分
}
