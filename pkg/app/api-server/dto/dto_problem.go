package dto

import "go.mongodb.org/mongo-driver/bson/primitive"

type ReqProblem struct {
	ID            primitive.ObjectID `json:"_id,omitempty"`
	Creator       string             `json:"creator,omitempty"`
	Description   string             `json:"description,omitempty"`
	InputExample  string             `json:"input_example,omitempty"`
	OutputExample string             `json:"output_example,omitempty"`
	Status        int32              `json:"status,omitempty"`
	SubmitNum     int32              `json:"submit_num,omitempty"`
	Difficulty    string             `json:"difficulty,omitempty"`
	Labels        []string           `json:"labels,omitempty"`
	PassNum       int32              `json:"pass_num,omitempty"`
	Title         string             `json:"title,omitempty"`
	URL           string             `json:"url,omitempty"`
}
