package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	models "handler/DataBase/Models"
	utils "handler/Utils"

	"github.com/gorilla/websocket"
)

type Online struct {
	Clients map[string][]*websocket.Conn
	Mutex   sync.Mutex
}

var OnlineConnections = Online{
	Clients: make(map[string][]*websocket.Conn),
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection:", err)
		return
	}
	defer conn.Close()

	user, ok := utils.GetUserFromSession(r)

	if !ok {
		fmt.Println("Error getting username from session:", err)
		return
	}

	OnlineConnections.Mutex.Lock()
	OnlineConnections.Clients[user.Nickname] = append(OnlineConnections.Clients[user.Nickname], conn)
	OnlineConnections.Mutex.Unlock()

	for {

		var messageData models.MessageData
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		err = json.Unmarshal(msg, &messageData)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		if err := models.CreateMessage(user.Nickname, messageData); err != nil {
			log.Println("Error saving message:", err)
			break
		}

		log.Printf("Received message from %s to %s: %s %s", user.Nickname, messageData.Receiver, messageData.Message, messageData.CreatedAt)

		receiverConnections, ok := OnlineConnections.Clients[messageData.Receiver]

		fmt.Println("remote addres", conn.RemoteAddr().String())
		if ok {
			for _, receiverConnection := range receiverConnections {
				responseMessage := map[string]string{
					"sender":     user.Nickname,
					"content":    messageData.Message,
					"created_at": messageData.CreatedAt,
				}
				message, _ := json.Marshal(responseMessage)
				err := receiverConnection.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Println("error sending message")
				}
			}
		} else {
			log.Printf("Receiver %s is not online", messageData.Receiver)
		}
	}
	var temp []*websocket.Conn
	if len(OnlineConnections.Clients[user.Nickname]) == 1 {
		OnlineConnections.Mutex.Lock()
		delete(OnlineConnections.Clients, user.Nickname)
		OnlineConnections.Mutex.Unlock()
	} else {
		for _, activeConn := range OnlineConnections.Clients[user.Nickname] {
			if activeConn.RemoteAddr().String() != conn.RemoteAddr().String() {
				temp = append(temp, activeConn)
			}
		}
		OnlineConnections.Mutex.Lock()
		OnlineConnections.Clients[user.Nickname] = temp
		OnlineConnections.Mutex.Unlock()

	}
}

func GetActiveUsers(w http.ResponseWriter, r *http.Request) []string {
	OnlineConnections.Mutex.Lock()
	defer OnlineConnections.Mutex.Unlock()

	user, ok := utils.GetUserFromSession(r)

	if !ok {
		fmt.Println("user not found in session")
		w.WriteHeader(http.StatusInternalServerError)
		return nil
	}
	var activeUsers []string
	for username := range OnlineConnections.Clients {
		if user.Nickname != username {
			activeUsers = append(activeUsers, username)
		}
	}
	return activeUsers
}
