package controller

import (
	"encoding/json"
	errorManagementControllers "forum/modules/errorManagement/controllers"
	forumManagementModels "forum/modules/forumManagement/models"
	"net/http"
	"strings"
)

// CreateChatHandler handles the creation of a new chat
func CreateChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	var chat forumManagementModels.Chat
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract user IDs from the request
	var userIDs []int
	err = json.NewDecoder(r.Body).Decode(&userIDs)
	if err != nil {
		http.Error(w, "Invalid user IDs", http.StatusBadRequest)
		return
	}

	chatID, err := forumManagementModels.InsertChat(&chat, userIDs)
	if err != nil {
		http.Error(w, "Failed to create chat", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"chat_id": chatID})
}

// GetChatMessagesHandler retrieves chat messages for a specific user
func GetChatMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the username from the URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	username := pathParts[3]

	// Retrieve chat messages for the user
	messages, err := forumManagementModels.GetMessagesByUsername(username)
	if err != nil {
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

// GetMessagesHandler retrieves all messages for a specific chat
func GetMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the username from the URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	username := pathParts[3]

	messages, err := forumManagementModels.GetMessagesByUsername(username)
	if err != nil {
		http.Error(w, "Failed to retrieve messages", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(messages)
}

// AddMessageHandler handles adding a new message to a chat
func AddMessageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var message forumManagementModels.Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	messageID, err := forumManagementModels.AddMessage(&message)
	if err != nil {
		http.Error(w, "Failed to add message", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"message_id": messageID})
}

// AddMessageFileHandler handles adding a file to a message
func AddMessageFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var file forumManagementModels.MessageFile
	err := json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fileID, err := forumManagementModels.AddMessageFile(&file)
	if err != nil {
		http.Error(w, "Failed to add message file", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"file_id": fileID})
}
