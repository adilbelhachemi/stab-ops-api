package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Operator struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"firstName,omitempty" bson:"first_name"`
	LastName  string             `json:"lastName,omitempty" bson:"last_name"`
	Position  string             `json:"position,omitempty" bson:"position"`
	Actions   []Action           `json:"actions" bson:"actions"`
}

type Action struct {
	Type      string `json:"type" bson:"type"`
	CreatedAt int64  `json:"created_at" bson:"created_at"`
}
