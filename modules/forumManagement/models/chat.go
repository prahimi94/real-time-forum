package models

import (
	"forum/db"
	"forum/utils"
	"log"

	"time"
)

// Chat struct represents the chat data model
type Chat struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Type      string    `json:"type"` // "private" or "group"
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Message struct represents the message data model
type Message struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	ChatID    int       `json:"chat_id"`
	UserID    int       `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatMember struct represents the chat member data model
type ChatMember struct {
	ID     int    `json:"id"`
	ChatID int    `json:"chat_id"`
	UserID int    `json:"user_id"`
	Status string `json:"status"`
}

// MessageFile struct represents the message file data model
type MessageFile struct {
	ID               int    `json:"id"`
	MessageID        int    `json:"message_id"`
	FileUploadedName string `json:"file_uploaded_name"`
	FileRealName     string `json:"file_real_name"`
	Status           string `json:"status"`
}

func InsertChat(chat *Chat, userIDs []int) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close()

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	chat.UUID, err = utils.GenerateUuid()
	if err != nil {
		tx.Rollback() // Rollback if UUID generation fails
		return -1, err
	}

	// Insert the chat
	insertChatQuery := `INSERT INTO chats (uuid, type, status) VALUES (?, ?, ?)`
	result, insertErr := tx.Exec(insertChatQuery, chat.UUID, chat.Type, chat.Status)
	if insertErr != nil {
		tx.Rollback()
		return -1, insertErr
	}

	chatID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return -1, err
	}

	// Add members to the chat
	for _, userID := range userIDs {
		insertMemberQuery := `INSERT INTO chat_members (chat_id, user_id, status) VALUES (?, ?, 'enable')`
		_, err := tx.Exec(insertMemberQuery, chatID, userID)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	newMessage := &Message{
		Content: "Welcome to the chat!",
		ChatID:  int(chatID),
		UserID:  userIDs[0], // Assuming the first user is the sender
	}
	_, addMessageErr := AddMessage(newMessage)
	if addMessageErr != nil {
		tx.Rollback() // Rollback if message insertion fails
		return -1, addMessageErr
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return -1, err
	}

	return int(chatID), nil
}

// GetChatsByUserID retrieves all chats for a specific user
func GetChatsByUserID(userID int) ([]Chat, error) {
	db := db.OpenDBConnection()
	defer db.Close()

	query := `
        SELECT c.id, c.uuid, c.type, c.status, c.created_at
        FROM chats c
        INNER JOIN chat_members cm ON c.id = cm.chat_id
        WHERE cm.user_id = ? AND cm.status = 'enable' AND c.status = 'enable'
    `
	rows, err := db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.ID, &chat.UUID, &chat.Type, &chat.Status, &chat.CreatedAt); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return chats, nil
}

// GetMessagesByUsername retrieves all messages for a specific user
func GetMessagesByUsername(username string) ([]Message, error) {
	db := db.OpenDBConnection()
	defer db.Close()

	query := `
        SELECT m.id, m.content, m.chat_id, m.user_id, m.status, m.created_at
        FROM messages m
        INNER JOIN users u ON m.user_id = u.id
        WHERE u.username = ? AND m.status = 'enable'
        ORDER BY m.created_at ASC
    `
	rows, err := db.Query(query, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.Content, &message.ChatID, &message.UserID, &message.Status, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// AddMessage adds a new message to a chat
func AddMessage(message *Message) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close()

	query := `INSERT INTO messages (content, chat_id, user_id, status) VALUES (?, ?, ?, 'enable')`
	result, err := db.Exec(query, message.Content, message.ChatID, message.UserID)
	if err != nil {
		return -1, err
	}

	messageID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(messageID), nil
}

// AddMessageFile adds a file to a message
func AddMessageFile(file *MessageFile) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close()

	query := `INSERT INTO message_files (message_id, file_uploaded_name, file_real_name, status) VALUES (?, ?, ?, 'enable')`
	result, err := db.Exec(query, file.MessageID, file.FileUploadedName, file.FileRealName)
	if err != nil {
		return -1, err
	}

	fileID, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	return int(fileID), nil
}
