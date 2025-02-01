/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/26 21:17
@Name: constant.go
@Description:
*/

package utils

import "time"

const (

	//账号状态
	StatusBanned = 0
	StatusNormal = 1
	StatusUser   = 0
	StatusAdmin  = 1

	//题目状态
	StatusPrivate = 0
	StatusPublic  = 1

	//响应result状态
	StatusSuccess = "Success"
	StatusFail    = "Fail"

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
	Descending       = -1
	Ascending        = 1
	DefaultDirection = Ascending

	// Token jwt token
	DefaultRole              = 1  //user
	RoleAdmin                = 99 //user
	Token                    = "X-Auth-Token"
	JwtTokenSecretKey        = "zkoj10086"
	JwtTokenHeaderKey        = "X-Auth-Token"
	RememberEffectiveTime    = time.Hour * time.Duration(24*14)
	NotRememberEffectiveTime = time.Hour * time.Duration(2)
)
