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

var messageData models.MessageData

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
	GetActiveUsers(w)
	for {

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

		if messageData.Message == "logout" {
			log.Printf("%s has logged out", user.Nickname)
			break
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
	fmt.Println(OnlineConnections.Clients, len(OnlineConnections.Clients))
	var temp []*websocket.Conn
	if len(OnlineConnections.Clients[user.Nickname]) == 1 {
		OnlineConnections.Mutex.Lock()
		delete(OnlineConnections.Clients, user.Nickname)
		OnlineConnections.Mutex.Unlock()
	} else if messageData.Message == "logout" {
		fmt.Println("KKKK")
		OnlineConnections.Mutex.Lock()
		delete(OnlineConnections.Clients, user.Nickname)
		OnlineConnections.Mutex.Unlock()
		messageData.Message = ""
	} else {
		for _, activeConn := range OnlineConnections.Clients[user.Nickname] {
			if activeConn.RemoteAddr().String() != conn.RemoteAddr().String() {
				temp = append(temp, activeConn)
			}
		}

		if len(temp) > 0 {
			OnlineConnections.Mutex.Lock()
			OnlineConnections.Clients[user.Nickname] = temp
			OnlineConnections.Mutex.Unlock()
		} else {
			OnlineConnections.Mutex.Lock()
			delete(OnlineConnections.Clients, user.Nickname)
			OnlineConnections.Mutex.Unlock()
		}

	}
	fmt.Println(OnlineConnections.Clients, len(OnlineConnections.Clients))

	GetActiveUsers(w)

}

func GetActiveUsers(w http.ResponseWriter) {
	// Convert map keys to a slice
	var onlineUsers []string
	for username := range OnlineConnections.Clients {
		onlineUsers = append(onlineUsers, username)
	}

	// Create a map to check online status
	onlineMap := make(map[string]bool)
	for _, id := range onlineUsers {
		onlineMap[id] = true
	}

	// Fetch all users
	users, err := models.GetAllUsers(onlineMap)
	if err != nil {
		log.Printf("Failed to retrieve all users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Marshal the user data to JSON
	message, err := json.Marshal(users)
	if err != nil {
		log.Println("Error marshalling user data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, connections := range OnlineConnections.Clients {
		for _, conn := range connections {
			// Ensure the connection is still open before sending the message
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error sending active users message:", err)
			}
		}
	}
}
