package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JoinedResponse struct {
	Type     string    `json:"type"`
	FromUUID string    `json:"fromUUID"`
	Contacts []Contact `json:"contactsUUID"`
}

type Contact struct {
	Type string `json:"type"`
	UUID string `json:"uuid"`
}

type Audit struct {
	UserIdCreated  string    `json:"userIdCreated" bson:"userIdCreated"`
	UserIdModified string    `json:"userIdModified" bson:"userIdModified"`
	Created        time.Time `json:"created" bson:"created"`
	Modified       time.Time `json:"modified" bson:"modified"`
}

type Conversation struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ClientId string             `json:"clientId" bson:"clientId"`
	Status   int                `json:"status" bson:"status"`
	Audit    Audit              `json:"audit" bson:"audit"`
}

type Incident struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AgentId        string             `json:"agentId" bson:"agentId"`
	ConversationId string             `json:"incidentId" bson:"incidentId"`
	Status         int                `json:"status" bson:"status"`
	Audit          Audit              `json:"audit" bson:"audit"`
}

type Message struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Type           string             `json:"type,omitempty" bson:"type,omitempty"`
	AgentId        string             `json:"agentId" bson:"agentId"`
	ClientId       string             `json:"clientId" bson:"clientId"`
	IncidentId     string             `json:"incidentId" bson:"incidentId"`
	DateTime       time.Time          `json:"dateTime" bson:"dateTime"`
	FromMe         bool               `json:"fromMe" bson:"fromMe"`
	FromUUID       string             `json:"fromUUID,omitempty" bson:"fromUUID,omitempty"`
	ToUUID         string             `json:"toUUID,omitempty" bson:"toUUID,omitempty"`
	MessageContent MessageContent     `json:"messageContent" bson:"messageContent"`
}

type MessageContent struct {
	MessageType      int    `json:"messageType" bson:"messageType"`
	Text             string `json:"text" bson:"text"`
	MediaUrl         string `json:"mediaUrl,omitempty" bson:"mediaUrl,omitempty"`
	MediaType        int    `json:"mediaType,omitempty" bson:"mediaType,omitempty"`
	MediaContentType string `json:"mediaContentType,omitempty" bson:"mediaContentType,omitempty"`
	Caption          string `json:"caption,omitempty" bson:"caption,omitempty"`
}
