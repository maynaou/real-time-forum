package handler

import (
	"encoding/json"
	"fmt"
	database "handler/DataBase/Sqlite"
	utils "handler/Utils"
	"log"

	"net/http"
	"sync"

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

	fmt.Println("Username:", user.Nickname)
	if !ok {
		fmt.Println("Error getting username from session:", err)
		return
	}

	OnlineConnections.Mutex.Lock()
	OnlineConnections.Clients[user.Nickname] = append(OnlineConnections.Clients[user.Nickname], conn)
	OnlineConnections.Mutex.Unlock()

	fmt.Println("HHH", OnlineConnections)

	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var messageData struct {
			Type     string `json:"type"`
			Sender   string `json:"sender"`
			Message  string `json:"message"`
			Receiver string `json:"receiver"`
			Time     string `json:"time"`
		}
		err = json.Unmarshal(msg, &messageData)
		fmt.Println(messageData)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}
		fmt.Println("HHHHHH")
		db := database.GetDatabaseInstance()
		if db.DB != nil {
			log.Println("Connexion à la base de données réussie.")
		}

		_, err = db.DB.Exec(`
            INSERT INTO messages (sender, receiver, content, created_at)
            VALUES (?, ?, ?, ?)
        `, user.Nickname, messageData.Receiver, messageData.Message, messageData.Time)

		if err != nil {
			log.Println("Error inserting message into the database:", err)
			break
		}

		log.Printf("Received message from %s to %s: %s", user.Nickname, messageData.Receiver, messageData.Message)

		receiverConnections, ok := OnlineConnections.Clients[messageData.Receiver]
		//check if the user/sender logs out all connections get deleted
		//if the receiver is loged out the message should still be added to the database and not shown for that connection but when he fetches the messages all the
		fmt.Println("remote addres", conn.RemoteAddr().String())
		if ok {
			for _, receiverConnection := range receiverConnections {
				responseMessage := map[string]string{
					"sender":  user.Nickname,
					"message": messageData.Message,
					"time":    messageData.Time,
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

func GetActiveUsers(w http.ResponseWriter, r *http.Request) {
	OnlineConnections.Mutex.Lock()
	defer OnlineConnections.Mutex.Unlock()

	user, ok := utils.GetUserFromSession(r)

	if !ok {
		fmt.Println("user not found in session")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var activeUsers []string
	for username := range OnlineConnections.Clients {
		if user.Nickname != username {
			activeUsers = append(activeUsers, username)
		}
	}

	fmt.Println("hhh", activeUsers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(activeUsers)
}
