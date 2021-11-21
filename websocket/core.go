package websocket

import (
	"log"
	"net/http"

	mdls "customer-support/models"
	mng "customer-support/mongo"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var clients = make(map[string]*websocket.Conn) // connected clients
var broadcast = make(chan mdls.Message)
var upgrader = websocket.Upgrader{}

func Run() {
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
	go handleMessages()

	log.Println("http server started on :8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	id := getNewUUID()
	clients[id] = ws

	contacts := make([]mdls.Contact, 0)

	for uuid := range clients {
		if uuid != id {
			contact := mdls.Contact{UUID: uuid}
			contacts = append(contacts, contact)
		}
	}

	response := mdls.JoinedResponse{Type: "joined", FromUUID: id, Contacts: contacts}
	ws.WriteJSON(response)
	sendMessageConnect(id)

	for {
		var msg mdls.Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error conn: %v", err)
			delete(clients, id)
			sendMessageDisconnect(id)
			break
		}
		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}

}

func sendMessageConnect(uuidNew string) {
	msg := mdls.Contact{Type: "connectContact", UUID: uuidNew}
	err := sendBroadcastWithoutMe(msg, uuidNew)
	if err != nil {
		log.Printf("\nError sending new connection %s", err.Error())
	}
}

func sendMessageDisconnect(uuidOld string) {
	msg := mdls.Contact{Type: "disconectContact", UUID: uuidOld}
	err := sendBroadcast(msg)
	if err != nil {
		log.Printf("\nError sending disconnect %s", err.Error())
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for uuid, client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error msg: %v", err)
				client.Close()
				delete(clients, uuid)
			}
		}
		mng.SaveMessage(msg)
	}
}

func sendBroadcast(msg interface{}) error {
	for _, client := range clients {
		err := client.WriteJSON(msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func sendBroadcastWithoutMe(msg interface{}, uuidMe string) error {
	for uuid, client := range clients {
		if uuid != uuidMe {
			err := client.WriteJSON(msg)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getNewUUID() string {
	id := uuid.New()
	return id.String()
}
