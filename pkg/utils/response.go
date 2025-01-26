/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/24 11:36
@Name: response.go
@Description:
*/

package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

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
