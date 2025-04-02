package controller

import (
	"encoding/json"
	"fmt"
	userManagementModels "forum/modules/userManagement/models"
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

var OnlineUsers = make(map[*websocket.Conn]string) // Map of online users (connected to WS) to usernames
var Broadcast = make(chan []byte)                  // Broadcast channel
var Mutex = &sync.Mutex{}                          // Protect OnlineUsers map

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	defer conn.Close()

	// Get myUsername from userid related to session token
	_, myUsername, err := userManagementModels.GetUserIDFromCookie(r)
	if err != nil {
		fmt.Println("Error getting username:", err)
		return
	}

	// Add the connection and username to the OnlineUsers map
	Mutex.Lock()
	OnlineUsers[conn] = myUsername
	UpdateOnlineUsers()
	fmt.Printf("Online: User %s connected. Current OnlineUsers: %v\n", myUsername, OnlineUsers)
	Mutex.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			// Remove the connection from the clients map on disconnect
			Mutex.Lock()
			delete(OnlineUsers, conn)
			fmt.Printf("Offline: User %s disconnected. Current OnlineUsers: %v\n", myUsername, OnlineUsers)
			UpdateOnlineUsers()
			Mutex.Unlock()
			break
		}
		
		// Ignore empty messages
		if len(message) == 0 {
			continue
		}

		// Add timestamp and username to the message
		timestamp := time.Now().Format("2006-01-02 15:04:05")
		formattedMessage := fmt.Sprintf("[%s] %s: %s", timestamp, myUsername, string(message))

		Broadcast <- []byte(formattedMessage)
	}
}

func HandleMessages() {
	for {
		// Grab the next message from the Broadcast channel
		message := <-Broadcast

		// Send the message to all online users
		Mutex.Lock()
		fmt.Println("Broadcasting message:", string(message), " | OnlineUsers:", OnlineUsers)
		for client := range OnlineUsers {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(OnlineUsers, client)
				UpdateOnlineUsers()
			}
		}
		Mutex.Unlock()
	}
}

// Helper function to broadcast the list of online users
func UpdateOnlineUsers() {
	usernames := make([]string, 0, len(OnlineUsers))
	for _, username := range OnlineUsers {
		usernames = append(usernames, username)
	}

	// Encode the list of usernames as JSON
	userListJSON, err := json.Marshal(usernames)
	if err != nil {
		fmt.Println("Error encoding online users:", err)
		return
	}

	// Send the list to all online clients
	for client := range OnlineUsers {
		err := client.WriteMessage(websocket.TextMessage, userListJSON)
		if err != nil {
			client.Close()
			delete(OnlineUsers, client)
		}
	}
}

func OnlineUsersHandler(w http.ResponseWriter, r *http.Request) {
	Mutex.Lock()
	defer Mutex.Unlock()

	// Collect usernames of online users
	usernames := make([]string, 0, len(OnlineUsers))
	for _, username := range OnlineUsers {
		usernames = append(usernames, username)
	}

	// Respond with the list of usernames
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(usernames); err != nil {
		http.Error(w, "Failed to encode online users", http.StatusInternalServerError)
	}
}
