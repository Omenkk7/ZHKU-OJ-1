package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Submit struct {
	ID         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Code       string             `json:"code,omitempty" bson:"code,omitempty"`
	Evaluation string             `json:"evaluation,omitempty" bson:"evaluation,omitempty"`
	Language   string             `json:"language,omitempty" bson:"language,omitempty"`
	ProblemId  string             `json:"problem_id,omitempty" bson:"problem_id,omitempty"`
	Status     int32              `json:"status,omitempty" bson:"status,omitempty"`
	Ctime      int64              `json:"ctime,omitempty" bson:"ctime,omitempty"`
	Mtime      int64              `json:"mtime,omitempty" bson:"mtime,omitempty"`
}
