package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Result 统一响应格式
type Result struct {
	Code int32       `json:"code"` //状态码
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (r *Result) Success(msg string, data interface{}) {
	r.Code = 0
	r.Msg = msg
	r.Data = data
}

func (r *Result) Fail(msg string) {
	r.Code = -1
	r.Msg = msg
	r.Data = "null"
}

func (r *Result) Response(c *gin.Context, res *Result) {
	if res.Code == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  res.Msg,
			"data": res.Data,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": -1,
			"msg":  res.Msg,
			"data": res.Data,
		})
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
