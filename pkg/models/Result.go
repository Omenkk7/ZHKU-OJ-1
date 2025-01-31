package models

import (
	"errors"
	"github.com/gin-gonic/gin"
	"zhku-oj-server/pkg/utils"
)

// 统一响应格式
type Result struct {
	Code int32       `json:"code"` //状态码
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *Result) Success(msg string, data interface{}) *Result {
	return &Result{
		Code: 0,
		Msg:  msg,
		Data: data,
	}
}

func (r *Result) Fail(msg string) *Result {
	return &Result{
		Code: -1,
		Msg:  msg,
	}
}

func (r *Result) Response(c *gin.Context, res *Result) {
	if res.Code == 0 {
		utils.SuccessResponse(c, res)
	} else {
		//TODO 统一响应格式，返回状态码-1，data为空字符串？
		utils.BadRequest(c, errors.New(res.Msg))
	}
}
