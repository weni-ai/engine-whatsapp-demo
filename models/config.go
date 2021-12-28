package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Config struct {
	ID    primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Token string             `json:"token" bson:"token"`
}
