package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type BaseTask struct {
	Code      string `json:"code,omitempty" bson:"code,omitempty"`
	Language  string `json:"language,omitempty" bson:"language,omitempty"`
	ProblemId string `json:"problem_id,omitempty" bson:"problem_id,omitempty"`
	UserId    string `json:"user_id,omitempty" bson:"user_id,omitempty"`
}

type LocalTask struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` //这个用于记录submit_id
	BaseTask                    //嵌入
}
