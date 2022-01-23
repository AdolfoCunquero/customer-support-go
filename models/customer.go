package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type origin struct {
}

type Customer struct {
	CustomerId primitive.ObjectID `json:"customerId" bson:"_id"`
	Name       string             `json:"name" bson:"name"`
	Origin     origin             `json:"origin" bson:"origin,omitempty"`
	Status     int                `json:"status" bson:"status"`
	Audit      Audit              `json:"audit,omitempty" bson:"audit,omitempty"`
}
