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
	Role     int32  `json:"role"`     // 用户角色：1-管理员、2-教师、3-助教、4-学生
	Nickname string `json:"nickname"` // 昵称
	Class    string `json:"class"`    // 班级
	Sid      string `json:"sid"`      // 学号

}

type ReqPostLoginUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ReqUser struct {
	ID       string   `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string   `json:"username,omitempty" bson:"username,omitempty"`
	Password string   `json:"password,omitempty" bson:"password,omitempty"`
	Phone    string   `json:"phone,omitempty" bson:"phone,omitempty"`
	Email    string   `json:"email,omitempty" bson:"email,omitempty"`
	Group    []string `json:"group,omitempty" bson:"group,omitempty"`
	Nickname string   `json:"nickname,omitempty" bson:"nickname,omitempty"`
	Role     int32    `json:"role,omitempty" bson:"role,omitempty"`
	Class    string   `json:"class,omitempty" bson:"class,omitempty"`
	Sid      string   `json:"sid,omitempty" bson:"sid,omitempty"`
	Status   int32    `json:"status,omitempty" bson:"status,omitempty"`
	Ctime    int64    `json:"ctime,omitempty" bson:"ctime,omitempty"`
	Mtime    int64    `json:"mtime,omitempty" bson:"mtime,omitempty"`
}
