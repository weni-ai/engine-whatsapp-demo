package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Channel struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UUID  string             `json:"uuid,omitempty"`
	Name  string             `json:"name,omitempty"`
	Token string             `json:"token,omitempty"`
}
