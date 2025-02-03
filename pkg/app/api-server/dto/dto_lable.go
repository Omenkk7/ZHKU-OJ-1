package dto

// ReqLabel 由于增删改查时，前端都传这四个参数，所以直接命名为ReqLabel
type ReqLabel struct {
	ID      string `json:"_id,omitempty" bson:"_id,omitempty"`
	Creator string `json:"creator,omitempty" bson:"creator,omitempty"`
	Status  int32  `json:"status,omitempty" bson:"status,omitempty"`
	Name    string `json:"name,omitempty" bson:"name,omitempty"`
	Ctime   int64  `json:"ctime,omitempty" bson:"ctime,omitempty"`
	Mtime   int64  `json:"mtime,omitempty" bson:"mtime,omitempty"`
}
