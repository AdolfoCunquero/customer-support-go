package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Incident struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AgentId        string             `json:"agentId" bson:"agentId"`
	ConversationId string             `json:"conversationId" bson:"conversationId"`
	CustomerId     string             `json:"customerId" bson:"customerId"`
	Status         int                `json:"status" bson:"status"`
	Audit          Audit              `json:"audit" bson:"audit"`
}

type ActiveIncident struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	AgentId        string             `bson:"agentId" json:"agentId"`
	ConversationId string             `bson:"conversationId" json:"conversationId"`
	CustomerId     string             `bson:"customerId" json:"customerId"`
	Status         int16              `bson:"status" json:"status"`
	CustomerInfo   Customer           `bson:"customerInfo" json:"customerInfo,omitempty"`
}

type ActiveConversationInc struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AgentId        string             `json:"agentId" bson:"agentId"`
	ConversationId string             `json:"conversationId" bson:"conversationId"`
	CustomerId     string             `json:"customerId" bson:"customerId"`
	Status         int                `json:"status" bson:"status"`
	Messages       []Message          `json:"messages" bson:"messages"`
}
