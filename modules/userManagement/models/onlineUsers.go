package models

import (
	"fmt"
	"forum/db"
	"net/http"
	"strings"
)

func GetActiveSessionUserIDs(r *http.Request) ([]int, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query to get unique user_ids with active sessions
	rows, err := db.Query(`SELECT DISTINCT user_id FROM sessions WHERE expires_at > CURRENT_TIMESTAMP`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get myUserID from session token
	myUserID, _, err := GetUserIDFromCookie(r)
	if err != nil {
		return nil, err
	}

	var userIds []int
	for rows.Next() {
		var userId int
		if err := rows.Scan(&userId); err != nil {
			return nil, err
		}
		// Do not include myUserID
		if userId == myUserID {
			continue
		}
		userIds = append(userIds, userId)
	}

	if len(userIds) == 0 {
		// No active sessions found
		return []int{}, nil
	}

	return userIds, nil
}

func GetActiveSessionUsernames(r *http.Request) ([]string, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	userIds, err := GetActiveSessionUserIDs(r)
	if err != nil {
		return nil, err
	}

	if len(userIds) == 0 {
		// No active user IDs found
		return []string{}, nil
	}

	// Dynamically construct the query with the correct number of placeholders
	placeholders := strings.Repeat("?,", len(userIds))
	placeholders = strings.TrimSuffix(placeholders, ",") // Remove the trailing comma
	query := fmt.Sprintf(`SELECT username FROM users WHERE id IN (%s)`, placeholders)

	// Convert userIds to a slice of interface{} for the query arguments
	args := make([]interface{}, len(userIds))
	for i, id := range userIds {
		args[i] = id
	}

	userRows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer userRows.Close()

	var usernames []string
	for userRows.Next() {
		var username string
		if err := userRows.Scan(&username); err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}

	return usernames, nil
}
