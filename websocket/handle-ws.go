package websocket

import (
	"log"
	"net/http"

	mdls "customer-support/models"

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

func sendMessageConnect(newUUID string) {
	for uuid, client := range clients {
		if uuid != newUUID {
			msg := mdls.Contact{Type: "connectContact", UUID: newUUID}
			client.WriteJSON(msg)
		}
	}
}

func sendMessageDisconnect(oldUUID string) {
	for _, client := range clients {
		msg := mdls.Contact{Type: "disconectContact", UUID: oldUUID}
		client.WriteJSON(msg)
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for UUID, client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error msg: %v", err)
				client.Close()
				delete(clients, UUID)
			}
		}
	}
}

func getNewUUID() string {
	id := uuid.New()
	return id.String()
}
