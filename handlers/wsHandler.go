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
	fmt.Println("hhh", messageData.Cookie)
	for {
		GetActiveUsers(w, user)

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
		fmt.Println("hhh", messageData.Cookie)
		if messageData.Cookie != "" {
			a := models.DeleteSession(messageData.Cookie)
			if a == nil {
				fmt.Println("delete session")
				break
			}
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

	var temp []*websocket.Conn
	if len(OnlineConnections.Clients[user.Nickname]) == 1 {
		OnlineConnections.Mutex.Lock()
		delete(OnlineConnections.Clients, user.Nickname)
		OnlineConnections.Mutex.Unlock()
	} else if messageData.Cookie != "" {
		OnlineConnections.Mutex.Lock()
		delete(OnlineConnections.Clients, user.Nickname)
		OnlineConnections.Mutex.Unlock()
		messageData.Cookie = ""
		messageData.Receiver = ""
	} else if messageData.Message == "logout" {
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

	GetActiveUsers(w, user)
}

func GetActiveUsers(w http.ResponseWriter, user models.RegisterRequest) {
	var onlineUsers []string
	for username := range OnlineConnections.Clients {
		onlineUsers = append(onlineUsers, username)
	}

	onlineMap := make(map[string]bool)
	for _, id := range onlineUsers {
		onlineMap[id] = true
	}

	users, err := models.GetAllUsers(onlineMap)
	if err != nil {
		log.Printf("Failed to retrieve all users: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"sender":   user.Nickname,
		"users":    users,
		"receiver": messageData.Receiver,
	}

	// Marshal the user data to JSON
	message, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshalling user data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, connections := range OnlineConnections.Clients {
		for _, conn := range connections {
			if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Println("Error sending active users message:", err)
			}
		}
	}
}
