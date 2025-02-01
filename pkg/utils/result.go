package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Result 统一响应格式
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
		SuccessResponse(c, res)
	} else {
		//TODO 统一响应格式，返回状态码-1，data为空字符串？
		BadRequest(c, errors.New(res.Msg))
	}
}

//原utils/response

func BadRequest(c *gin.Context, err error) {
	FailedResponse(c, http.StatusBadRequest, err)
}

func Unauthorized(c *gin.Context, err error) {
	FailedResponse(c, http.StatusUnauthorized, err)
}

func Forbidden(c *gin.Context, err error) {
	FailedResponse(c, http.StatusForbidden, err)
}

func SuccessResponse(c *gin.Context, data ...interface{}) {
	if len(data) == 0 {
		c.JSON(http.StatusOK, gin.H{})
	}
	c.JSON(http.StatusOK, data[0])
}

func FailedResponse(c *gin.Context, status int, err error) {
	c.JSON(status, gin.H{
		Message: err.Error(),
	})
}

func FailStringResponse(c *gin.Context, status int, err error) {
	c.String(status, "FailString: %s", err.Error())
}

func FailedResponseWithData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		Message: data,
	})
}
