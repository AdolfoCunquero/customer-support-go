package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conversation struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CustomerId string             `json:"customerId" bson:"customerId"`
	Status     int                `json:"status" bson:"status"`
	Audit      Audit              `json:"audit" bson:"audit"`
}
