package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Submit struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Code       string             `json:"code,omitempty" bson:"code,omitempty"`
	Evaluation interface{}        `json:"evaluation,omitempty" bson:"evaluation,omitempty"`
	Language   string             `json:"language,omitempty" bson:"language,omitempty"`
	UserId     string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
	ProblemId  string             `json:"problem_id,omitempty" bson:"problem_id,omitempty"`
	Status     int32              `json:"status,omitempty" bson:"status,omitempty"`
	Stime      int64              `json:"stime,omitempty" bson:"stime,omitempty"`
}
