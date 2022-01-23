package websocket

import (
	"encoding/json"
	"log"
	"net/http"

	mdls "customer-support/models"
	mdb "customer-support/mongo"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var clients = make(map[string]*websocket.Conn) // connected clients
var broadcast = make(chan map[string]interface{})
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

	// contacts := make([]mdls.Contact, 0)

	// for uuid := range clients {
	// 	if uuid != id {
	// 		contact := mdls.Contact{UUID: uuid}
	// 		contacts = append(contacts, contact)
	// 	}
	// }

	response := mdls.JoinedResponse{Type: "joined", UUID: id}
	ws.WriteJSON(response)
	//sendMessageConnect(id)

	for {
		var msg map[string]interface{}
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error conn: %v", err)
			delete(clients, id)
			break
		}
		// Send the newly received message to the broadcast channel
		msg["UUID"] = id
		broadcast <- msg
	}

}

// func sendMessageConnect(uuidNew string) {
// 	msg := mdls.Contact{Type: "connectContact", UUID: uuidNew}
// 	err := sendBroadcastWithoutMe(msg, uuidNew)
// 	if err != nil {
// 		log.Printf("\nError sending new connection %s", err.Error())
// 	}
// }

// func sendMessageDisconnect(uuidOld string) {
// 	msg := mdls.Contact{Type: "disconectContact", UUID: uuidOld}
// 	err := sendBroadcast(msg)
// 	if err != nil {
// 		log.Printf("\nError sending disconnect %s", err.Error())
// 	}
// }

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		result := make(map[string]interface{})
		uuid := msg["UUID"].(string)
		client := clients[uuid]

		if msg["type"] == "getContacts" {

			resultDB, errM := mdb.GetActiveIncidents("acunqueroc")
			if errM != nil {
				log.Printf("error mongo: %v", errM)
				continue
			}

			result["type"] = "contacts"
			result["payload"] = resultDB

			err := client.WriteJSON(result)
			if err != nil {
				log.Printf("error msg: %v", err)
				client.Close()
				delete(clients, uuid)
			}

		} else if msg["type"] == "newMessage" {

			jsonStr, errJ := json.Marshal(msg["payload"])
			if errJ != nil {
				log.Printf("Error decoding json: %v", errJ)
				continue
			}

			var payload mdls.Message
			json.Unmarshal(jsonStr, &payload)
			msg, errDB := SaveMessage(payload)
			if errDB != nil {
				log.Printf("Error saving message: %v", errDB)
				continue
			}

			result["type"] = "newMessageOk"
			result["payload"] = msg

			err := client.WriteJSON(result)
			if err != nil {
				log.Printf("error msg: %v", err)
				client.Close()
				delete(clients, uuid)
			}

		} else if msg["type"] == "getCustomerMessages" {

			resultDB, errM := mdb.GetActiveIncidents("acunqueroc")
			if errM != nil {
				log.Printf("error mongo: %v", errM)
				continue
			}

			result["type"] = "contacts"
			result["payload"] = resultDB

			err := client.WriteJSON(result)
			if err != nil {
				log.Printf("error msg: %v", err)
				client.Close()
				delete(clients, uuid)
			}
		}

		// // Send it out to every client that is currently connected
		// for uuid, client := range clients {
		// 	err := client.WriteJSON(result)
		// 	if err != nil {
		// 		log.Printf("error msg: %v", err)
		// 		client.Close()
		// 		delete(clients, uuid)
		// 	}
		// }
	}
}

// func sendBroadcast(msg interface{}) error {
// 	for _, client := range clients {
// 		err := client.WriteJSON(msg)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func sendBroadcastWithoutMe(msg interface{}, uuidMe string) error {
// 	for uuid, client := range clients {
// 		if uuid != uuidMe {
// 			err := client.WriteJSON(msg)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 	}
// 	return nil
// }

func getNewUUID() string {
	id := uuid.New()
	return id.String()
}
