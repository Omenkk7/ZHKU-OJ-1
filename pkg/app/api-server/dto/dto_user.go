/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/26 21:44
@Name: dto_user.go
@Description:
*/

package dto

type ReqPostUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type ReqPostLoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
