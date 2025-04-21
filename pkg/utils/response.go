package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Result 统一响应格式

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
		return //要return，不然data[0]可能会报空指针
	}
	c.JSON(http.StatusOK, data[0])
}

func FailedResponse(c *gin.Context, status int, err error) {
	// TODO 如果是乱码，应该是解码和编码不一致导致的
	// TODO 解决方法：需要统一UTF-8；；前端要检查Content-Type: application/json; charset=utf-8
	c.JSON(status, gin.H{
		//Message: err.Error(), //TODO:这种写法无法将中文的错误原因返回前端，改成下一行
		Message: err,
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
