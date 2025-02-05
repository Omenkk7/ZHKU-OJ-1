/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 22:48
@Name: auth.go
@Description:
*/

package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"strings"
	"zhku-oj-server/pkg/utils"
)

var (
	whiteList = map[string]string{
		"/api/v1/user/":      "POST",
		"/api/v1/user/login": "POST",
	}
)

func withInWhiteList(url *url.URL, method string) bool {
	target := whiteList
	queryUrl := strings.Split(fmt.Sprint(url), "?")[0]
	if _, ok := target[queryUrl]; ok {
		if target[queryUrl] == method {
			return true
		}
		return false
	}
	return false
}

// JWTMiddleware 拦截器 通用
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lg := utils.GetDefaultLogger()
		lg.Info("拦截请求......")

		// TODO:// 配置白名单和黑名单
		//if withInWhiteList(c.Request.URL, c.Request.Method) {
		//	c.Next()
		//	return
		//}

		// 从请求头中获取token
		token := c.Request.Header.Get("Authorization") //TODO 这个key好像是前端设置的，每次请求都自动携带

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
		/*dao := dao.NewDao()
		query := bson.M{
			"username": jwtClaims.Username, //旧jwt中的旧username数据
		}
		user, _ := dao.GetOneUser(context.Background(), query) //user已为最新状态*/
		//TODO 权限控制，目前的逻辑可能会造成数据不一致
		if jwtClaims.Role == utils.StatusUser {
			utils.BadRequest(c, utils.New(utils.NOPermissionErr))
			c.Abort()
			return
		}

		//校验成功，解析并拿到jwt的用户数据，存进gin.Context，可通过c.Get("user")重新获得数据
		contextUser := &utils.ContextUser{
			ID:       jwtClaims.ID.Hex(),
			Username: jwtClaims.Username,
			Role:     int(jwtClaims.Role),
		}
		c.Set("contextUser", contextUser)
		c.Next()
	}
}
