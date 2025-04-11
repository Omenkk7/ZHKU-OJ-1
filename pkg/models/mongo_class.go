package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassMember struct {
	UserID   string `json:"user_id" bson:"user_id"`     // 用户ID
	UserName string `json:"user_name" bson:"user_name"` // 用户名称
}

// Class 班级模型
type Class struct {
	ID           primitive.ObjectID       `json:"_id,omitempty" bson:"_id,omitempty"`
	ClassCode    string                   `json:"class_code" bson:"class_code"`           // 班级代码：院系缩写+年份+随机4位数字
	Name         string                   `json:"name" bson:"name"`                       // 班级名称
	Department   string                   `json:"department" bson:"department"`           // 所属院系
	Description  string                   `json:"description" bson:"description"`         // 班级描述
	Creator      string                   `json:"creator" bson:"creator"`                 // 创建者ID
	Courses      []CourseInfo             `json:"courses" bson:"courses"`                 // 关联课程
	Members      map[string][]ClassMember `json:"members" bson:"members"`                 // 班级成员权限映射
	StudentCount int                      `json:"student_count" bson:"student_count"`     // 学生数量（冗余字段）
	Status       int                      `json:"status" bson:"status"`                   // 状态：1-正常，0-已归档
	Ctime        int64                    `json:"ctime,omitempty" bson:"ctime,omitempty"` // 创建时间
	Mtime        int64                    `json:"mtime,omitempty" bson:"mtime,omitempty"` // 修改时间
}

// CourseInfo 课程信息（嵌入到Class中）
type CourseInfo struct {
	CourseID   string `json:"course_id" bson:"course_id"`     // 关联课程ID
	CourseName string `json:"course_name" bson:"course_name"` // 课程名称
	BindTime   int64  `json:"bind_time" bson:"bind_time"`     // 绑定时间
	Status     int    `json:"status" bson:"status"`           // 状态：1-进行中，0-已结课
}

// ClassStudent 班级学生关联模型
type ClassStudent struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClassID       string             `json:"class_id" bson:"class_id"`               // 班级ID
	StudentID     string             `json:"student_id" bson:"student_id"`           // 学生ID
	StudentName   string             `json:"student_name" bson:"student_name"`       // 学生姓名（冗余字段）
	StudentNumber string             `json:"student_number" bson:"student_number"`   // 学号（冗余字段）
	JoinType      int                `json:"join_type" bson:"join_type"`             // 加入方式：1-管理员添加，2-自主申请
	Status        int                `json:"status" bson:"status"`                   // 状态：1-正常，0-已移出
	JoinTime      int64              `json:"join_time" bson:"join_time"`             // 加入时间
	LeaveTime     *int64             `json:"leave_time" bson:"leave_time"`           // 离开时间，如果未离开则为null
	Ctime         int64              `json:"ctime,omitempty" bson:"ctime,omitempty"` // 创建时间
	Mtime         int64              `json:"mtime,omitempty" bson:"mtime,omitempty"` // 修改时间
}

// ClassJoinRequest 班级加入申请
type ClassJoinRequest struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClassID       string             `json:"class_id" bson:"class_id"`               // 班级ID
	StudentID     string             `json:"student_id" bson:"student_id"`           // 学生ID
	StudentName   string             `json:"student_name" bson:"student_name"`       // 学生姓名（冗余字段）
	StudentNumber string             `json:"student_number" bson:"student_number"`   // 学号（冗余字段）
	RequestMsg    string             `json:"request_msg" bson:"request_msg"`         // 申请信息
	Status        int                `json:"status" bson:"status"`                   // 状态：0-待审核，1-已通过，2-已拒绝
	ReviewerID    string             `json:"reviewer_id" bson:"reviewer_id"`         // 审核人ID
	ReviewMsg     string             `json:"review_msg" bson:"review_msg"`           // 审核意见
	ReviewTime    *int64             `json:"review_time" bson:"review_time"`         // 审核时间
	Ctime         int64              `json:"ctime,omitempty" bson:"ctime,omitempty"` // 创建时间
	Mtime         int64              `json:"mtime,omitempty" bson:"mtime,omitempty"` // 修改时间
}
