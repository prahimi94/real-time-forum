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
	"sync"
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
	Typing    bool                          `json:"typing"` // New field for typing status
}

var Broadcast = make(chan WebsocketMsg) // Broadcast channel
var Mutex = &sync.Mutex{}               // Protect OnlineUsers map

func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
	defer conn.Close()

	// Get myUsername from userid related to session token
	cookie, err := r.Cookie("session_token")
	if err == nil && cookie != nil && cookie.Value != "" {
		myUserID, myUsername, err := userManagementModels.GetUserIDFromCookie(r)
		if err != nil {
			errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
			return
		}

		// Add the connection and username to the OnlineUsers map
		Mutex.Lock()
		userManagementControllers.OnlineUsers[conn] = myUsername
		userManagementControllers.UpdateOnlineUsers()
		broadcastUserOnlineStatus(myUserID, true)
		Mutex.Unlock()

		var chatID int // Declare chatID outside the loop

		for {

			var socketmsg WebsocketMsg

			cookie, err := r.Cookie("session_token")
			if err != nil || (cookie != nil && cookie.Value == "") {
				Mutex.Lock()
				defer Mutex.Unlock()
				conn.Close()
				delete(userManagementControllers.OnlineUsers, conn)
				userManagementControllers.UpdateOnlineUsers()
				broadcastUserOnlineStatus(myUserID, false)
				break
			}

			var msgData struct {
				Type      string    `json:"type"`
				Content   string    `json:"content"`
				Sender    string    `json:"sender"`
				Recipient string    `json:"recipient"`
				Timestamp time.Time `json:"timestamp"`
				Typing    bool      `json:"typing"` // New field for typing status
			}

			err = conn.ReadJSON(&msgData)
			if err != nil {
				// Lock the mutex to safely modify the OnlineUsers map
				Mutex.Lock()

				// Close the WebSocket connection
				conn.Close()

				// Remove the connection from the OnlineUsers map
				delete(userManagementControllers.OnlineUsers, conn)
				userManagementControllers.UpdateOnlineUsers()
				broadcastUserOnlineStatus(myUserID, false)
				Mutex.Unlock()
				// Exit the loop to stop processing messages for this connection
				break
			}

			for _, username := range userManagementControllers.OnlineUsers {
				if username == msgData.Recipient {
					// Check if both sender and recipient are online
					senderOnline := false
					recipientOnline := false

					for _, onlineUsername := range userManagementControllers.OnlineUsers {
						if onlineUsername == msgData.Sender {
							senderOnline = true
						}
						if onlineUsername == msgData.Recipient {
							recipientOnline = true
						}
						// Break early if both are found
						if senderOnline && recipientOnline {
							break
						}
					}

					// If either sender or recipient is not online, skip processing and do not save the message
					if !senderOnline || !recipientOnline {
						fmt.Printf("Skipping message. Sender online: %v, Recipient online: %v\n", senderOnline, recipientOnline)
						continue
					}

					// Handle "typing" message type
					if msgData.Type == "typing" {
						//fmt.Println("Typing message received: ", msgData)
						socketmsg.Type = "typing"
						socketmsg.Sender = msgData.Sender
						socketmsg.Recipient = msgData.Recipient
						socketmsg.Typing = msgData.Typing
						Broadcast <- socketmsg // Notify the recipient about typing status
						continue
					}

					// Handle "private_chat" message type
					//if msgData.Type == "private_chat" {
					recipientUsername := msgData.Recipient

					// Get recipient user ID
					recipientUserID, err := userManagementModels.GetUserIDByUsername(recipientUsername)
					if err != nil {
						errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
						continue
					}

					// Check if chat exists, if not create it and add chat members
					chatID, err = forumManagementModels.CheckChatExists(myUserID, recipientUserID)
					if err != nil {
						errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
						continue
					}
					// If no chatID exists (0), InsertChat
					if chatID == 0 {
						chat := &forumManagementModels.Chat{ID: chatID, Type: "private"}
						chatID, err = forumManagementModels.InsertChat(chat, myUserID, recipientUserID, nil)
						if err != nil {
							errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
							continue
						}
					}
					/* socketmsg.Type = "private_chat_confirmed"
					socketmsg.Recipient = msgData.Recipient
					socketmsg.Sender = msgData.Sender
					Broadcast <- socketmsg
					continue */
					//}

					sanitizedMsg := utils.SanitizeInput(msgData.Content)
					// Ignore empty messages
					if sanitizedMsg == "" {
						continue
					}

					// If chatID exists, go directly to InsertMsg
					if chatID != 0 {
						msg := &forumManagementModels.Message{
							ChatID:            chatID,
							Content:           sanitizedMsg,
							Status:            "enable",
							CreatedBy:         myUserID,
							CreatedAt:         time.Now(),
							UpdatedBy:         &myUserID,
							CreatedByUsername: msgData.Sender,
						}
						_, err = forumManagementModels.InsertMsg(msg, nil)
						if err != nil {
							errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
							continue
						}
					}
					socketmsg.Type = "private_chat"
					socketmsg.Message.Content = sanitizedMsg
					socketmsg.Sender = msgData.Sender
					socketmsg.Recipient = msgData.Recipient
					socketmsg.Message.CreatedAt = msgData.Timestamp
					socketmsg.Message.CreatedByUsername = msgData.Sender
					Broadcast <- socketmsg
				}
			}
		}
	}
}

func HandleMessages() {
	for {
		// Grab the next message from the Broadcast channel
		message := <-Broadcast

		Mutex.Lock()

		for client, username := range userManagementControllers.OnlineUsers {

			// Retrieve the userID of the disconnected user
			disconnectedUserID, err := userManagementModels.GetUserIDByUsername(username)
			if err != nil {
				fmt.Println("Error retrieving userID for disconnected user:", err)
			}

			if message.Type == "typing" && username == message.Recipient {
				err := client.WriteJSON(message)
				if err != nil {
					client.Close()
					delete(userManagementControllers.OnlineUsers, client)
					userManagementControllers.UpdateOnlineUsers()
					broadcastUserOnlineStatus(disconnectedUserID, false)
				}
			} else if message.Type == "private_chat" && (username == message.Recipient || username == message.Sender) {
				err := client.WriteJSON(message)
				if err != nil {
					client.Close()
					delete(userManagementControllers.OnlineUsers, client)
					userManagementControllers.UpdateOnlineUsers()
					broadcastUserOnlineStatus(disconnectedUserID, false)
				}
			} else {
				err := client.WriteJSON(message)
				if err != nil {
					client.Close()
					delete(userManagementControllers.OnlineUsers, client)
					userManagementControllers.UpdateOnlineUsers()
					broadcastUserOnlineStatus(disconnectedUserID, false)

				}
			}
		}
		Mutex.Unlock()
	}
}

/*
	 func OnlineUsersHandler(w http.ResponseWriter, r *http.Request) {
		Mutex.Lock()
		defer Mutex.Unlock()

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
*/
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
	/* 	w.Header().Set("Content-Type", "application/json")
	   	json.NewEncoder(w).Encode(map[string]int{"chatID": chatID}) */
	res := utils.Result{
		Success: true,
		Message: "ChatID fetched successfully",
		Data:    chatID,
	}
	utils.ReturnJson(w, res)
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
		users[i].IsOnline = false // Default to false
		for _, username := range userManagementControllers.OnlineUsers {
			if username == user.Username {
				users[i].IsOnline = true
				continue
			}
		}
	}

	// Return the users as JSON
	/* 	w.Header().Set("Content-Type", "application/json")
	   	var socketmsg WebsocketMsg
	   	socketmsg.Type = "show_all_users"
	   	socketmsg.Users = users
	   	Broadcast <- socketmsg
	   	json.NewEncoder(w).Encode(map[string]bool{"success": true}) */

	res := utils.Result{
		Success: true,
		Message: "Users fetched successfully",
		Data:    users,
	}
	utils.ReturnJson(w, res)

}

func broadcastUserOnlineStatus(userID int, isOnline bool) {
	// Broadcast the updated user list to all clients
	users, err := userManagementModels.ReadAllChatUsers(userID)

	for i, user := range users {
		for _, username := range userManagementControllers.OnlineUsers {
			if username == user.Username {
				users[i].IsOnline = isOnline
				continue
			}
		}
	}
	if err != nil {
		fmt.Println("Error fetching users: ", err)
	}
	var socketmsg WebsocketMsg
	socketmsg.Type = "fetch_all_users"
	socketmsg.Users = users
	Broadcast <- socketmsg
	//fmt.Println("broadcastUserOnlineStatus: ", socketmsg.Users)
}
