package dto

type ReqSubmit struct {
	Code      string `json:"code,omitempty" bson:"code,omitempty"`
	Language  string `json:"language,omitempty" bson:"language,omitempty"`
	ProblemId string `json:"problem_id,omitempty" bson:"problem_id,omitempty"`
	Status    int32  `json:"status,omitempty" bson:"status,omitempty"`
	Ctime     int64  `json:"ctime,omitempty" bson:"ctime,omitempty"`
	Mtime     int64  `json:"mtime,omitempty" bson:"mtime,omitempty"`
}
