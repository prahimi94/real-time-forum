package models

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/db"
	"forum/utils"
	"log"
	"net/http"
	"time"
)

// User struct represents the user data model
type Session struct {
	ID           int       `json:"id"`
	SessionToken string    `json:"session_token"`
	UserId       int       `json:"user_id"`
	CreatedAt    time.Time `json:"created_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func InsertSession(session *Session) (*Session, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Generate UUID for the user if not already set
	if session.SessionToken == "" {
		uuidSessionTokenid, err := utils.GenerateUuid()
		if err != nil {
			return nil, err
		}
		session.SessionToken = uuidSessionTokenid
	}

	// Set session expiration time to 1 hour
	session.ExpiresAt = time.Now().Add(1 * time.Hour)

	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return &Session{}, err
	}

	updateQuery := `UPDATE sessions SET expires_at = CURRENT_TIMESTAMP WHERE user_id = ? AND expires_at > CURRENT_TIMESTAMP;`
	_, updateErr := tx.Exec(updateQuery, session.UserId)
	if updateErr != nil {
		tx.Rollback()
		return nil, updateErr
	}

	insertQuery := `INSERT INTO sessions (session_token, user_id, expires_at) VALUES (?, ?, ?);`
	_, insertErr := tx.Exec(insertQuery, session.SessionToken, session.UserId, session.ExpiresAt)
	if insertErr != nil {
		tx.Rollback()
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return nil, sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return nil, insertErr
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback() // Rollback on error
		return nil, err
	}

	return session, nil
}

func SelectSession(sessionToken string) (User, time.Time, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	var user User
	var expirationTime time.Time
	err := db.QueryRow(`SELECT 
							u.id as user_id, u.type as user_type, u.firstname as user_firstname, u.lastname as user_lastname, u.age as user_age, u.gender as user_gender, u.username as username, u.email as user_email, IFNULL(u.profile_photo, '') as profile_photo,
							expires_at 
						FROM sessions s
							INNER JOIN users u
								ON s.user_id = u.id
						WHERE session_token = ?`, sessionToken).Scan(&user.ID, &user.Type, &user.Firstname, &user.Lastname, &user.Age, &user.Gender, &user.Username, &user.Email, &user.ProfilePhoto, &expirationTime)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			// Handle other database errors
			return User{}, time.Time{}, errors.New("sql: no rows in result set")
		} else {
			// Handle other database errors
			return User{}, time.Time{}, errors.New("database error")
		}
	}

	return user, expirationTime, nil
}

func DeleteSession(sessionToken string) error {

	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes
	_, err := db.Exec(`UPDATE sessions
					SET expires_at = CURRENT_TIMESTAMP
					WHERE session_token = ?;`, sessionToken)
	if err != nil {
		// Handle other database errors
		log.Fatal(err)
		return errors.New("database error")
	}

	return nil

}

// IsSessionActive checks if a session is active based on the session token
func IsSessionActive(sessionToken string) (bool, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	var expiresAt time.Time

	// Query the database for the session's expiration time
	err := db.QueryRow(`SELECT expires_at FROM sessions WHERE session_token = ?`, sessionToken).Scan(&expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// No session found for the given token
			return false, errors.New("session not found")
		}
		// Handle other database errors
		return false, err
	}

	// Check if the session is still active
	if expiresAt.After(time.Now()) {
		return true, nil // Session is active
	}

	return false, nil // Session is expired
}

func GetUserIDFromCookie(r *http.Request) (int, string, error) {
	// Retrieve user data (e.g., from session or database)
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		// Return an error if the session token is not found
		return 0, "", fmt.Errorf("error retrieving session token: %v", err)
	}

	// Fetch the user from the database using the session token
	user, _, err := SelectSession(sessionToken.Value)
	if err != nil {
		return 0, "", fmt.Errorf("error retrieving session: %v", err)
	}

	myUserID := user.ID         // Get the ID field from the User struct
	myUsername := user.Username // Use the Username field from the User struct

	return myUserID, myUsername, nil
}
