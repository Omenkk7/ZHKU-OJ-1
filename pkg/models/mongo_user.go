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
	Name     string             `json:"name,omitempty" bson:"name,omitempty"`
	Role     int                `json:"role,omitempty" bson:"role,omitempty"`
	Email    string             `json:"email,omitempty" bson:"email,omitempty"`
	Password string             `json:"password,omitempty" bson:"password,omitempty"`
	Phone    string             `json:"phone,omitempty" bson:"phone,omitempty"`
	Ctime    int64              `json:"ctime,omitempty" bson:"ctime,omitempty"`
	Mtime    int64              `json:"mtime,omitempty" bson:"mtime,omitempty"`
}
