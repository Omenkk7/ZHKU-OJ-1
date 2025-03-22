package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReqSubmit struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Code      string             `json:"code,omitempty" bson:"code,omitempty"`
	Language  string             `json:"language,omitempty" bson:"language,omitempty"`
	ProblemId string             `json:"problem_id,omitempty" bson:"problem_id,omitempty"`
	UserId    string             `json:"user_id,omitempty" bson:"user_id,omitempty"`
}
