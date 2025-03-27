package controller

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan []byte)            // Broadcast channel
var mutex = &sync.Mutex{}                    // Protect clients map

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	mutex.Lock()
	clients[conn] = true
	mutex.Unlock()
	/*
		// Retrieve session token from cookie
		sessionToken, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println("Error retrieving session token:", err)
			return
		}

		// Fetch the user_id from the session table using the session token
		user, _, err := userManagementModels.SelectSession(sessionToken.Value)
		if err != nil {
			fmt.Println("Error retrieving session:", err)
			return
		}

		// Fetch the username from the users table using the user_id
		userID, err := userManagementModels.ReadUserByID(user.ID)
		if err != nil {
			fmt.Println("Error retrieving user:", err)
			return
		}

		username := userID.Username */

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			mutex.Unlock()
			break
		}

		// Add timestamp and username to the message
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMessage := fmt.Sprintf("[%s] %s: %s", timestamp /* username */, "user", string(message))

		broadcast <- []byte(formattedMessage)
	}
}

func HandleMessages() {
	for {
		// Grab the next message from the broadcast channel
		message := <-broadcast

		// Send the message to all connected clients
		mutex.Lock()
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mutex.Unlock()
	}
}
