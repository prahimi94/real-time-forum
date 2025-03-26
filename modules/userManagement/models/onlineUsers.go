package models

import (
	"database/sql"
	"encoding/json"
	"forum/db"
	errorManagementControllers "forum/modules/errorManagement/controllers"
	"net/http"
)

func GetOnlineUsers() ([]User, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	query := `
        SELECT 
            u.id, 
            u.name, 
            u.username, 
            u.email, 
            IFNULL(u.profile_photo, '') as profile_photo,
            MAX(m.created_at) as last_message_time
        FROM sessions s
        INNER JOIN users u ON s.user_id = u.id
        LEFT JOIN messages m ON u.id = m.user_id
        WHERE s.expires_at > CURRENT_TIMESTAMP
        GROUP BY u.id
        ORDER BY 
            CASE 
                WHEN MAX(m.created_at) IS NOT NULL THEN 1 
                ELSE 2 
            END, 
            MAX(m.created_at) DESC, 
            u.username ASC
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var onlineUsers []User
	for rows.Next() {
		var user User
		var lastMessageTime sql.NullTime // To handle NULL values for users without messages
		err := rows.Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.ProfilePhoto, &lastMessageTime)
		if err != nil {
			return nil, err
		}
		onlineUsers = append(onlineUsers, user)
	}

	return onlineUsers, nil
}

func GetOnlineUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.MethodNotAllowedError)
		return
	}

	onlineUsers, err := GetOnlineUsers()
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}

	// Check if no users are online
	if len(onlineUsers) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "No users are online now.",
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(onlineUsers)
	if err != nil {
		errorManagementControllers.HandleErrorPage(w, r, errorManagementControllers.InternalServerError)
		return
	}
}
