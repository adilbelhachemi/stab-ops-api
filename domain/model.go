package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Operator struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	FirstName string             `json:"firstName,omitempty" bson:"first_name"`
	LastName  string             `json:"lastName,omitempty" bson:"last_name"`
	Role      string             `json:"role,omitempty" bson:"role"`
	Actions   []Action           `json:"actions" bson:"actions"`
	Password  string             `json:"password,omitempty" bson:"password"`
}

type Action struct {
	Type      string    `json:"type" bson:"type"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"` // primitive.DateTime
}

type OperatorFilter struct {
	Action   string    `json:"action,omitempty"`
	FromDate time.Time `json:"from,omitempty"`
	ToDate   time.Time `json:"to,omitempty"`
	Current  bool      `json:"current,omitempty"`
	Password string    `json:"password,omitempty"`
	Role     string    `json:"role,omitempty"`
}

type OperatorSigninRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
