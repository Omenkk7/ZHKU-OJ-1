package dto

type ReqSubmit struct {
	Code      string `json:"code,omitempty" bson:"code,omitempty"`
	Language  string `json:"language,omitempty" bson:"language,omitempty"`
	ProblemId string `json:"problem_id,omitempty" bson:"problem_id,omitempty"`
	UserId    string `json:"user_id,omitempty" bson:"user_id,omitempty"`
}
