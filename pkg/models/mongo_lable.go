package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Label struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Creator string             `json:"creator,omitempty" bson:"creator,omitempty"`
	Status  int32              `json:"status,omitempty" bson:"status,omitempty"`
	Ctime   int64              `json:"ctime,omitempty" bson:"ctime,omitempty"`
	Mtime   int64              `json:"mtime,omitempty" bson:"mtime,omitempty"`
}
