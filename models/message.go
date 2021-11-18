package models

type JoinedResponse struct {
	Type     string    `json:"type"`
	FromUUID string    `json:"fromUUID"`
	Contacts []Contact `json:"contactsUUID"`
}

type Contact struct {
	Type string `json:"type"`
	UUID string `json:"uuid"`
}

type Message struct {
	//Email       string `json:"email"`
	Type     string `json:"type"`
	Username string `json:"username"`
	FromUUID string `json:"fromUUID,omitempty"`
	DestUUID string `json:"destUUID,omitempty"`
	Message  string `json:"message"`
}
