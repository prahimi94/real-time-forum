package controller

import (
	"encoding/json"
	"fmt"
	errorManagementControllers "forum/modules/errorManagement/controllers"
	forumManagementModels "forum/modules/forumManagement/models"
	userManagementControllers "forum/modules/userManagement/controllers"
	userManagementModels "forum/modules/userManagement/models"

	"forum/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebsocketMsg struct {
	Type      string                        `json:"type"`
	Message   forumManagementModels.Message `json:"message"`
	Sender    string                        `json:"sender"`
	Recipient string                        `json:"recipient"`
	Users     []userManagementModels.User   `json:"users"`
}

var Broadcast = make(chan WebsocketMsg) // Broadcast channel

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	defer conn.Close()

	// Get myUsername from userid related to session token
	cookie, err := r.Cookie("session_token")
	if err == nil && cookie != nil && cookie.Value != "" {
		myUserID, myUsername, err := userManagementModels.GetUserIDFromCookie(r)
		if err != nil {
			fmt.Println("Error getting username:", err)
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		// Add the connection and username to the OnlineUsers map
		utils.Mutex.Lock()
		userManagementControllers.OnlineUsers[conn] = myUsername
		userManagementControllers.UpdateOnlineUsers()
		utils.Mutex.Unlock()

		var chatID int // Declare chatID outside the loop
		fmt.Println(userManagementControllers.OnlineUsers)
		for {
			var socketmsg WebsocketMsg

			cookie, err := r.Cookie("session_token")
			if err != nil || (cookie != nil && cookie.Value == "") {
				utils.Mutex.Lock()
				delete(userManagementControllers.OnlineUsers, conn)
				
				socketmsg.Type = "fetch_all_users"
				Broadcast <- socketmsg
				//userManagementControllers.UpdateOnlineUsers()
				utils.Mutex.Unlock()
				break
			}
			var msgData struct {
				Type      string    `json:"type"`
				Content   string    `json:"content"`
				Sender    string    `json:"sender"`
				Recipient string    `json:"recipient"`
				Timestamp time.Time `json:"timestamp"`
			}
			err = conn.ReadJSON(&msgData)
			fmt.Println("the error: ",err)
			if err != nil {
				// Remove the connection from the clients map on disconnect
				utils.Mutex.Lock()
				delete(userManagementControllers.OnlineUsers, conn)
				fmt.Println("connection is closed")
				//fmt.Println(userManagementControllers.OnlineUsers)
				socketmsg.Type = "fetch_all_users"
				Broadcast <- socketmsg
				utils.Mutex.Unlock()
				break
			}

			// Parse incoming message as JSON
			/* 			if err := json.Unmarshal(message, &msgData); err != nil {
				fmt.Printf("Invalid message format: %s | Error: %v\n", string(message), err)
				continue
			} */

			// Handle "private_chat" message type
			if msgData.Type == "private_chat" {
				recipientUsername := msgData.Recipient

				// Get recipient user ID
				recipientUserID, err := userManagementModels.GetUserIDByUsername(recipientUsername)
				if err != nil {
					fmt.Println("Error getting recipient user ID:", err)
					errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
					continue
				}

				// Check if chat exists, if not create it and add chat members
				chatID, err = forumManagementModels.CheckChatExists(myUserID, recipientUserID)
				if err != nil {
					fmt.Println("Error checking chat existence:", err)
					errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
					continue
				}
				// If no chatID exists (0), InsertChat
				if chatID == 0 {
					chat := &forumManagementModels.Chat{ID: chatID, Type: "private"}
					chatID, err = forumManagementModels.InsertChat(chat, myUserID, recipientUserID, nil)
					if err != nil {
						fmt.Println("Error creating or retrieving chat:", err)
						errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
						continue
					}
				}

				fmt.Printf("Chat initialized between %s and %s (Chat ID: %d)\n", myUsername, recipientUsername, chatID)
				continue
			}

			sanitizedMsg := utils.SanitizeInput(msgData.Content)
			// Ignore empty messages
			if sanitizedMsg == "" {
				continue
			}

			// If chatID exists, go directly to InsertMsg
			if chatID != 0 {
				msg := &forumManagementModels.Message{
					ChatID:    chatID, // Use the chat ID from the "private_chat" logic
					Content:   sanitizedMsg,
					Status:    "enable",
					CreatedBy: myUserID,
					CreatedAt: time.Now(), // Ensure CreatedAt is set
					UpdatedBy: &myUserID,
				}
				_, err = forumManagementModels.InsertMsg(msg, nil)
				if err != nil {
					fmt.Println("Error inserting message into database:", msg, err)
					errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
					continue
				}
			}
			socketmsg.Recipient = msgData.Recipient
			socketmsg.Sender = msgData.Sender
			socketmsg.Message.Content = sanitizedMsg
			socketmsg.Message.CreatedAt = msgData.Timestamp
			socketmsg.Type = "message_content"
			Broadcast <- socketmsg
		}
	}
}

func HandleMessages() {
	for {
		// Grab the next message from the Broadcast channel
		message := <-Broadcast

		// Send the message to all online users
		utils.Mutex.Lock()
		// fix to send only sender and reciver for private msg
		for client := range userManagementControllers.OnlineUsers {
			err := client.WriteJSON(message)
			if err != nil {
				client.Close()
				delete(userManagementControllers.OnlineUsers, client)
				//userManagementControllers.UpdateOnlineUsers()
				var socketmsg WebsocketMsg
				socketmsg.Type = "fetch_all_users"
				Broadcast <- socketmsg
			}
		}
		utils.Mutex.Unlock()
	}
}

func OnlineUsersHandler(w http.ResponseWriter, r *http.Request) {
	utils.Mutex.Lock()
	defer utils.Mutex.Unlock()

	// Collect usernames of online users
	usernames := make([]string, 0, len(userManagementControllers.OnlineUsers))
	for _, username := range userManagementControllers.OnlineUsers {
		usernames = append(usernames, username)
	}

	// Respond with the list of usernames
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(usernames); err != nil {
		http.Error(w, "Failed to encode online users", http.StatusInternalServerError)
	}
}

func ChatMsgHandler(w http.ResponseWriter, r *http.Request) {
	// Extract chatID from the URL path
	chatIDStr := r.URL.Path[len("/api/chat-messages/"):]
	if chatIDStr == "" {
		http.Error(w, "Chat ID is required", http.StatusBadRequest)
		return
	}

	// Convert chatID to int
	chatID, err := strconv.Atoi(chatIDStr)
	if err != nil {
		http.Error(w, "Invalid Chat ID", http.StatusBadRequest)
		return
	}

	// Retrieve userID from the session or cookie
	userID, _, err := userManagementModels.GetUserIDFromCookie(r)
	if err != nil {
		http.Error(w, "Failed to retrieve user ID", http.StatusUnauthorized)
		return
	}

	// Read all messages for the given chat ID
	messages, err := forumManagementModels.ReadAllMsgs(chatID, userID)
	if err != nil {
		http.Error(w, "Failed to read messages", http.StatusInternalServerError)
		return
	}

	// Respond with the messages in JSON format
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Failed to encode messages", http.StatusInternalServerError)
	}
}

func GetChatIDHandler(w http.ResponseWriter, r *http.Request) {
	// Parse sender and recipient from the request
	sender := r.URL.Query().Get("sender")
	recipient := r.URL.Query().Get("recipient")

	if sender == "" || recipient == "" {
		http.Error(w, "Sender and recipient are required", http.StatusBadRequest)
		return
	}

	// Get sender and recipient user IDs
	senderID, err := userManagementModels.GetUserIDByUsername(sender)
	if err != nil {
		http.Error(w, "Invalid sender username", http.StatusBadRequest)
		return
	}

	recipientID, err := userManagementModels.GetUserIDByUsername(recipient)
	if err != nil {
		http.Error(w, "Invalid recipient username", http.StatusBadRequest)
		return
	}

	// Query the database for the chat ID
	chatID, err := forumManagementModels.CheckChatExists(senderID, recipientID)
	if err != nil {
		http.Error(w, "Failed to retrieve chat ID", http.StatusInternalServerError)
		return
	}

	// Respond with the chat ID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"chatID": chatID})
}

func GetAllChatUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the logged-in user's information
	loginStatus, loginUser, _, checkLoginError := userManagementControllers.CheckLogin(w, r)
	if checkLoginError != nil {
		http.Error(w, "Failed to check login status", http.StatusInternalServerError)
		return
	}
	if !loginStatus {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	// Pass the logged-in user's ID to ReadAllChatUsers
	users, err := userManagementModels.ReadAllChatUsers(loginUser.ID)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	for i, user := range users {
		for _, username := range userManagementControllers.OnlineUsers {
			if username == user.Username {
				users[i].IsOnline = true
				continue
			}
		}
	}

	// Return the users as JSON
	w.Header().Set("Content-Type", "application/json")
	var socketmsg WebsocketMsg
	socketmsg.Type = "show_all_users"
	socketmsg.Users = users
	Broadcast <- socketmsg
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}
