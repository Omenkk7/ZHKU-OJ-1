/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 22:48
@Name: auth.go
@Description:
*/

package middleware

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"zhku-oj-server/pkg/dao"
	"zhku-oj-server/pkg/utils"
)

// JWTInterceptor 拦截器 通用
func JWTInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg := utils.GetDefaultLogger()
		lg.Info("拦截请求......")

		// 从请求头中获取token
		token := c.Request.Header.Get(utils.JwtTokenHeaderKey) //TODO 这个key好像是前端设置的，每次请求都自动携带

		// 检查token是否存在
		if token == "" {
			lg.Info(utils.JwtEmptyErr)
			utils.BadRequest(c, utils.New(utils.JwtEmptyErr))
			c.Abort()
			return
		}

		//校验jwt是否有效
		jwtClaims, err := utils.ParseToken(token, utils.JwtTokenSecretKey)
		if err != nil {
			lg.Info(utils.JwtFailErr, err)
			utils.BadRequest(c, utils.New(utils.JwtFailErr))
			c.Abort()
			return
		}

		//检查权限，判断是否为管理员
		d := dao.NewDao()
		query := bson.M{
			"username": jwtClaims.Username, //旧jwt中的旧username数据
		}
		user, _ := d.GetOneUser(context.Background(), query) //user已为最新状态
		if user.Role == utils.StatusUser {
			utils.BadRequest(c, utils.New(utils.NOPermissionErr))
			c.Abort()
			return
		}

		//校验成功，解析并拿到jwt的用户数据，存进gin.Context，可通过c.Get("user")重新获得数据
		c.Set("User", user)
		c.Next()
	}
}
