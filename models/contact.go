package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Contact struct {
	ID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	URN     string             `json:"urn,omitempty" bson:"urn,omitempty"`
	Name    string             `json:"name,omitempty" bson:"name,omitempty"`
	Channel primitive.ObjectID `json:"channel,omitempty" bson:"channel,omitempty"`
}
