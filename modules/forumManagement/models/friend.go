package models

import (
	"fmt"
	"forum/db"
	userManagementModels "forum/modules/userManagement/models"
	"log"
	"time"
)

// Post struct represents the user data model
type Friend struct {
	ID           int                         `json:"id"`
	FirstUserId  int                         `json:"first_user_id"`
	SecondUserId int                         `json:"second_user_id"`
	Status       string                      `json:"status"`
	CreatedAt    time.Time                   `json:"created_at"`
	Createdby    int                         `json:"created_by"`
	UpdatedAt    *time.Time                  `json:"updated_at"`
	UpdatedBy    *int                        `json:"updated_by"`
	Users        []userManagementModels.User `json:"users"` // Embedded users data
}

func ReadFriendsByUserId(userId int) ([]Friend, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT 
			first_user.id as first_user_id, first_user.name as first_user_name, first_user.username as first_user_username, first_user.email as first_user_email, first_user.profile_photo as first_user_profile_photo,
			second_user.id as second_user_id, second_user.name as second_user_name, second_user.username as second_user_username, second_user.email as second_user_email, second_user.profile_photo as second_user_profile_photo,
			f.status
		FROM friends f
			INNER JOIN users first_user
				ON f.first_user_id = u.id
			INNER JOIN users second_user
				ON f.second_user_id = u.id
				AND (first_user.id = ? OR second_user.id = ?)
		WHERE f.status != 'delete';
    `, userId, userId)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var friends []Friend
	// Map to track posts by their ID to avoid duplicates

	for rows.Next() {
		var friend Friend
		var first_user userManagementModels.User
		var second_user userManagementModels.User

		// Scan the post and user data
		err := rows.Scan(
			&first_user.ID, &first_user.Name, &first_user.Username, &first_user.Email, &first_user.ProfilePhoto,
			&second_user.ID, &second_user.Name, &second_user.Username, &second_user.Email, &second_user.ProfilePhoto,
			&friend.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		friendUsers := []userManagementModels.User{first_user, second_user}
		friend.Users = friendUsers

		friends = append(friends, friend)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return friends, nil
}

func InsertFriend(firstUserId int, secondUserId int) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	insertQuery := `INSERT INTO friends (first_user_id, second_user_id, created_by) VALUES (?, ?, ?);`
	result, insertErr := db.Exec(insertQuery, firstUserId, secondUserId, firstUserId)
	if insertErr != nil {
		return -1, insertErr
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	return int(lastInsertID), nil
}

func UpdateStatusFriend(firend_id int, status string, user_id int) error {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	updateQuery := `UPDATE friends
					SET status = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := db.Exec(updateQuery, status, user_id, firend_id)
	if updateErr != nil {
		return updateErr
	}

	return nil
}
