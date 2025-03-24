package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Course 课程模型
type Course struct {
	ID           primitive.ObjectID        `json:"_id,omitempty" bson:"_id,omitempty"`
	CourseCode   string                    `json:"course_code" bson:"course_code"`         // 课程代码：院系缩写+年份+随机4位数字
	Name         string                    `json:"name" bson:"name"`                       // 课程名称
	Department   string                    `json:"department" bson:"department"`           // 所属院系
	Description  string                    `json:"description" bson:"description"`         // 课程描述
	Creator      string                    `json:"creator" bson:"creator"`                 // 创建者ID
	Members      map[string][]CourseMember `json:"members" bson:"members"`                 // 课程成员权限映射
	StudentCount int                       `json:"student_count" bson:"student_count"`     // 学生数量（冗余字段，便于查询）
	StartTime    int64                     `json:"start_time" bson:"start_time"`           // 开课时间
	EndTime      int64                     `json:"end_time" bson:"end_time"`               // 结课时间
	Status       int                       `json:"status" bson:"status"`                   // 状态：1-开放，0-已结课
	Ctime        int64                     `json:"ctime,omitempty" bson:"ctime,omitempty"` // 创建时间
	Mtime        int64                     `json:"mtime,omitempty" bson:"mtime,omitempty"` // 修改时间
}

// CourseMember 课程成员
type CourseMember struct {
	UserID        string `json:"user_id" bson:"user_id"`                                   // 用户ID
	UserName      string `json:"user_name" bson:"user_name"`                               // 用户名称
	Email         string `json:"email,omitempty" bson:"email,omitempty"`                   // 邮箱
	StudentNumber string `json:"student_number,omitempty" bson:"student_number,omitempty"` // 学号（学生特有）
	ClassID       string `json:"class_id,omitempty" bson:"class_id,omitempty"`             // 关联班级ID（如果是通过班级绑定加入）
	Status        int    `json:"status" bson:"status"`                                     // 状态：1-正常，0-已移出
	JoinTime      int64  `json:"join_time" bson:"join_time"`                               // 加入时间
}

// CourseJoinRequest 课程加入申请（选修课场景）
type CourseJoinRequest struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CourseID      string             `json:"course_id" bson:"course_id"`             // 课程ID
	CourseCode    string             `json:"course_code" bson:"course_code"`         // 课程代码（冗余字段，便于查询）
	StudentID     string             `json:"student_id" bson:"student_id"`           // 申请学生ID
	StudentName   string             `json:"student_name" bson:"student_name"`       // 学生姓名（冗余字段）
	StudentNumber string             `json:"student_number" bson:"student_number"`   // 学号（冗余字段）
	RequestMsg    string             `json:"request_msg" bson:"request_msg"`         // 申请理由
	Status        int                `json:"status" bson:"status"`                   // 状态：0-待审核，1-已通过，2-已拒绝
	ReviewerID    string             `json:"reviewer_id" bson:"reviewer_id"`         // 审核人ID
	ReviewMsg     string             `json:"review_msg" bson:"review_msg"`           // 审核意见
	ReviewTime    *int64             `json:"review_time" bson:"review_time"`         // 审核时间
	Ctime         int64              `json:"ctime,omitempty" bson:"ctime,omitempty"` // 创建时间
	Mtime         int64              `json:"mtime,omitempty" bson:"mtime,omitempty"` // 修改时间
}
