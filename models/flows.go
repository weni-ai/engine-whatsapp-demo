package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Flows struct {
	Channel     primitive.ObjectID `json:"channel_uuid,omitempty" bson:"channel_uuid,omitempty"`
	FlowsStarts []Flow             `json:"flows_starts,omitempty" bson:"flows_starts,omitempty"`
}

type Flow struct {
	Name    string `json:"flow_name,omitempty" bson:"flow_name,omitempty"`
	UUID    string `json:"flow_uuid,omitempty" bson:"flow_uuid,omitempty"`
	Keyword string `json:"keyword,omitempty" bson:"keyword,omitempty"`
}
