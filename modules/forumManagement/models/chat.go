package models

import (
	"database/sql"
	"fmt"
	"forum/db"
	"forum/utils"
	"log"
	"time"
)

// Chat represents the "chats" table
type Chat struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	Type      string    `json:"type"` // "private" or "group"
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy int       `json:"created_by"`
}

// ChatMember represents the "chat_members" table
type ChatMember struct {
	ID        int       `json:"id"`
	ChatID    int       `json:"chat_id"`
	UserID    int       `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Message represents the "messages" table
type Message struct {
	ID        int        `json:"id"`
	ChatID    int        `json:"chat_id"`
	Content   string     `json:"content"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	CreatedBy int        `json:"created_by"`
	UpdatedAt *time.Time `json:"updated_at"`
	UpdatedBy *int       `json:"updated_by"`
}

// MessageFile represents the "message_files" table
type MessageFile struct {
	ID               int        `json:"id"`
	ChatID           int        `json:"chat_id"`
	MessageID        int        `json:"message_id"`
	FileUploadedName string     `json:"file_uploaded_name"`
	FileRealName     string     `json:"file_real_name"`
	Status           string     `json:"status"`
	CreatedAt        time.Time  `json:"created_at"`
	CreatedBy        int        `json:"created_by"`
	UpdatedAt        *time.Time `json:"updated_at"`
	UpdatedBy        *int       `json:"updated_by"`
}

func CheckChatExists(user1ID, user2ID int) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close() //

	var chatID int

	// Check if a private chat already exists between the two users
	query := `
        SELECT c.id
        FROM chats c
        JOIN chat_members cm1 ON c.id = cm1.chat_id
        JOIN chat_members cm2 ON c.id = cm2.chat_id
        WHERE c.type = 'private' AND cm1.user_id = ? AND cm2.user_id = ?
    `
	err := db.QueryRow(query, user1ID, user2ID).Scan(&chatID)
	if err == nil {
		// Chat already exists
		return chatID, nil
	} else if err != sql.ErrNoRows {
		// Unexpected error
		return 0, fmt.Errorf("failed to check for existing chat: %w", err)
	}

	// No chat exists
	return 0, nil
}

func InsertChat(chat *Chat, user1ID, user2ID int, uploadedFiles map[string]string) (int, error) {
	chatID, err := CheckChatExists(user1ID, user2ID)
	if err != nil {
		return 0, err
	}
	if chatID != 0 {
		return chatID, nil
	}

	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	chat.UUID, err = utils.GenerateUuid()
	if err != nil {
		tx.Rollback() // Rollback if UUID generation fails
		return -1, err
	}

	insertQuery := `INSERT INTO chats (uuid, type, created_by) VALUES (?, ?, ?);`
	result, insertErr := tx.Exec(insertQuery, chat.UUID, "private", user1ID) // Assuming user1ID initiated chat
	if insertErr != nil {
		return -1, insertErr
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
		return -1, err
	}

	insertChatMember1Err := InsertChatMember(int(lastInsertID), user1ID, tx)
	if insertChatMember1Err != nil {
		tx.Rollback()
		return -1, insertChatMember1Err
	}
	insertChatMember2Err := InsertChatMember(int(lastInsertID), user2ID, tx)
	if insertChatMember2Err != nil {
		tx.Rollback()
		return -1, insertChatMember2Err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return -1, err
	}

	return int(lastInsertID), nil
}

func InsertChatMember(chatID, userID int, tx *sql.Tx) error {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	insertChatMemberQuery := `INSERT INTO chat_members (chat_id, user_id) VALUES (?, ?);`
	_, err := tx.Exec(insertChatMemberQuery, chatID, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func InsertMsg(msg *Message, uploadedFiles map[string]string) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return -1, err
	}

	insertMsgQuery := `INSERT INTO messages (chat_id, content, status, created_at, created_by, updated_at, updated_by) VALUES (?, ?, ?, ?, ?, ?, ?);`
	result, insertMsgQueryErr := tx.Exec(insertMsgQuery, msg.ChatID, msg.Content, msg.Status, msg.CreatedAt, msg.CreatedBy, `CURRENT_TIMESTAMP`, nil)
	if insertMsgQueryErr != nil {
		tx.Rollback()
		return -1, insertMsgQueryErr
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	insertMsgFilesErr := InsertMsgFiles(msg.ChatID, int(lastInsertID), uploadedFiles, msg.CreatedBy, tx)
	if insertMsgFilesErr != nil {
		tx.Rollback()
		return -1, insertMsgFilesErr
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return -1, err
	}

	return int(lastInsertID), nil
}

func InsertMsgFiles(chatID int, msgID int, uploadedFiles map[string]string, user_id int, tx *sql.Tx) error {
	if len(uploadedFiles) > 0 {
		query := `INSERT INTO message_files (chat_id, message_id, file_real_name, file_uploaded_name, created_by) VALUES `
		values := make([]any, 0, len(uploadedFiles)*3)

		for i := 0; i < len(uploadedFiles); i++ {
			if i > 0 {
				query += ", "
			}
			query += "(?, ?, ?, ?, ?)"
			for key, value := range uploadedFiles {
				values = append(values, chatID, msgID, key, value, user_id)
			}
		}
		query += ";"

		// Execute the bulk insert query
		_, err := tx.Exec(query, values...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func ReadAllMsgs(chatID int) ([]Message, error) {
	db := db.OpenDBConnection()
	defer db.Close()

	var messages []Message

	query := `SELECT id, chat_id, content, status, FROM messages WHERE chat_id = ?;`
	rows, err := db.Query(query, chatID)
	if err != nil {
		return nil, fmt.Errorf("failed to read messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.ChatID, &message.Content, &message.Status); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}

// ReadTenMsgs at a time for a given chat, with pagination support
func ReadTenMsgs(chatID int, offset int) ([]Message, error) {
	db := db.OpenDBConnection()
	defer db.Close()

	var messages []Message

	// Query to fetch messages with pagination
	query := `
		SELECT id, chat_id, content, status, created_at, created_by, updated_at, updated_by
		FROM messages
		WHERE chat_id = ?
		ORDER BY created_at DESC
		LIMIT 10 OFFSET ?;
	`
	rows, err := db.Query(query, chatID, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to read messages: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var message Message
		if err := rows.Scan(
			&message.ID,
			&message.ChatID,
			&message.Content,
			&message.Status,
			&message.CreatedAt,
			&message.CreatedBy,
			&message.UpdatedAt,
			&message.UpdatedBy,
		); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	return messages, nil
}
