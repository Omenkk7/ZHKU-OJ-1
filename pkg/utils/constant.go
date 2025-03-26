/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/26 21:17
@Name: constant.go
@Description:
*/

package utils

import "time"

const (
	LocalJudge     = 0
	RemoteJudge    = 1
	ReadTaskTime   = 5
	LocalJudgeCfg  = "LocalJudge"
	RemoteJudgeCfg = "RemoteJudge"

	//成功
	UpdateSuccess   = "更新成功"
	DeleteSuccess   = "删除成功"
	SelectSuccess   = "查询成功"
	RegisterSuccess = "注册成功！"
	LoginSuccess    = "登录成功"
	CreateSuccess   = "创建成功！"

	//账号状态
	StatusBanned = 0
	StatusNormal = 1
	StatusUser   = 0
	StatusAdmin  = 1

	//题目状态
	StatusPrivate = 0
	StatusPublic  = 1

	OneDaySeconds  = 86400
	OneHourSeconds = 3600
	TenMinutes     = 600
	DateLayout     = "20060102"

	Logger        = "logger"
	DefaultLogger = "default_logger"
	TraceID       = "trace_id"

	Message = "message"

	// DefaultNum
	DefaultPageNum   = 1
	PageNum          = "_pageNum"
	DefaultPageSize  = 20
	PageSize         = "_pageSize"
	Direction        = "_direction"
	DefaultSort      = "_id"
	Sort             = "_sort"
	Desc             = "desc"
	Asc              = "asc"
	Descending       = -1
	Ascending        = 1
	DefaultDirection = Ascending

	// Token jwt token
	DefaultRole              = 1  //user
	RoleAdmin                = 99 //user
	JwtTokenSecretKey        = "zkoj10086"
	JwtTokenHeaderKey        = "X-Auth-Token"
	RememberEffectiveTime    = time.Hour * time.Duration(24*14)
	NotRememberEffectiveTime = time.Hour * time.Duration(2)

	ClassTable            = "classes"             // 班级表
	ClassStudentTable     = "class_students"      // 班级学生关系表
	ClassJoinRequestTable = "class_join_requests" // 班级加入申请表

	CourseTable            = "courses"              // 课程表
	CourseJoinRequestTable = "course_join_requests" // 课程加入申请表

	CourseRoleAdmin     = 1 // 管理员
	CourseRoleTeacher   = 2 // 教师
	CourseRoleAssistant = 3 // 助教
	CourseRoleStudent   = 4 // 学生

	CourseStatusOpen   = 1 // 开放
	CourseStatusClosed = 0 // 已结课

	CourseJoinStatusPending  = 0 // 待审核
	CourseJoinStatusApproved = 1 // 已通过
	CourseJoinStatusRejected = 2 // 已拒绝

	AssignmentTable = "assignments" // 作业表

)
