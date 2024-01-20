package controllers

import (
	"encoding/json"
	"github.com/JMURv/e-commerce/gateway/pkg/models"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocketChatMessage struct {
	Type        string         `json:"type"`
	UserID      uint           `json:"userID"`
	ReceiverID  uint           `json:"receiverID"`
	RoomID      uint           `json:"roomID"`
	MessageData models.Message `json:"messageData"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var clients = make(map[uint]*websocket.Conn)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	reader(conn)
}

func reader(conn *websocket.Conn) {
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var receivedMessage WebSocketChatMessage
		err = json.Unmarshal(p, &receivedMessage)
		if err != nil {
			log.Println("Error decoding JSON:", err)
			continue
		}

		clients[receivedMessage.UserID] = conn

		switch receivedMessage.Type {
		case "create":
			log.Println("Get create message type")

			newMessage, err := receivedMessage.MessageData.CreateMessage()
			if err != nil {
				log.Printf("Error while creating message: %v", err)
				continue
			}

			response, err := json.Marshal(newMessage)
			if err != nil {
				log.Printf("Error while encoding new message: %v", err)
				continue
			}

			broadcast(receivedMessage.UserID, receivedMessage.ReceiverID, response)
		case "edit":
			log.Println("Get edit message type")

			updatedMessage, err := models.UpdateMessage(receivedMessage.MessageData.ID, &receivedMessage.MessageData)
			if err != nil {
				log.Printf("Error while updating message: %v", err)
				continue
			}

			response, err := json.Marshal(updatedMessage)
			if err != nil {
				log.Printf("Error while encoding updated message: %v", err)
				continue
			}

			broadcast(receivedMessage.UserID, receivedMessage.ReceiverID, response)
		case "delete":
			log.Println("Get delete message type")

			deletedMessage, err := models.DeleteMessage(receivedMessage.MessageData.ID)
			if err != nil {
				log.Printf("Error while deleting message: %v", err)
				continue
			}

			response, err := json.Marshal(deletedMessage)
			if err != nil {
				log.Printf("Error while encoding deleting message: %v", err)
				continue
			}

			broadcast(receivedMessage.UserID, receivedMessage.ReceiverID, response)
		}
	}
}

func broadcast(senderID, receiverID uint, message []byte) {
	for userID, conn := range clients {
		if userID == receiverID || userID == senderID {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
