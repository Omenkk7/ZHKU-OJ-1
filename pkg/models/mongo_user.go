/*
@Author: urmsone urmsone@163.com
@Date: 2025/1/25 20:56
@Name: user.go
@Description:
*/

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Username string             `json:"username,omitempty" bson:"username,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Phone    string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Group    []string           `json:"group,omitempty" bson:"group,omitempty"`
	Nickname string             `json:"nickname,omitempty" bson:"nickname,omitempty"`
	Role     int32              `json:"role,omitempty" bson:"role,omitempty"`
	Class    string             `json:"class,omitempty" bson:"class,omitempty"`
	Sid      string             `json:"sid,omitempty" bson:"sid,omitempty"`
	Status   int32              `json:"status,omitempty" bson:"status,omitempty"`
	Ctime    int64              `json:"ctime,omitempty" bson:"ctime,omitempty"`
	Mtime    int64              `json:"mtime,omitempty" bson:"mtime,omitempty"`
}
