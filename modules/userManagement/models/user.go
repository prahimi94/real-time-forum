package models

import (
	"database/sql"
	"errors"
	"fmt"
	"forum/db"
	"forum/utils"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User struct represents the user data model
type User struct {
	ID           int        `json:"id"`
	UUID         string     `json:"uuid"`
	Type         string     `json:"type"`
	Username     string     `json:"username"`
	Firstname    string     `json:"firstname"`
	Lastname     string     `json:"lastname"`
	Gender       string     `json:"gender"`
	Age          int        `json:"age"`
	Email        string     `json:"email"`
	Password     string     `json:"password"`
	ProfilePhoto string     `json:"profile_photo"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
	UpdatedBy    *int       `json:"updated_by"`
}

func InsertUser(user *User) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Generate UUID for the user if not already set
	if user.UUID == "" {
		uuid, err := utils.GenerateUuid()
		if err != nil {
			return -1, err
		}
		user.UUID = uuid
	}

	var existingEmail string
	var existingUsername string
	emailCheckQuery := `SELECT email, username FROM users WHERE email = ? OR username = ? LIMIT 1;`
	err := db.QueryRow(emailCheckQuery, user.Email, user.Username).Scan(&existingEmail, &existingUsername)
	if err == nil {
		if existingEmail == user.Email {
			return -1, errors.New("duplicateEmail")
		}
		if existingUsername == user.Username {
			return -1, errors.New("duplicateUsername")
		}
	}

	insertQuery := `INSERT INTO users (uuid, firstname, lastname, gender, age, username, email, password) VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
	result, insertErr := db.Exec(insertQuery, user.UUID, user.Firstname, user.Lastname, user.Gender, user.Age, user.Username, user.Email, user.Password)
	if insertErr != nil {
		fmt.Println(insertErr)
		fmt.Println(user.Gender)
		// Check if the error is a SQLite constraint violation (duplicate entry)
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // 19 = UNIQUE constraint failed (SQLite error code)
				return -1, errors.New("user with this email or username already exists")
			}
		}
		return -1, insertErr // Other DB errors
	}

	// Retrieve the last inserted ID
	userId, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
		return -1, err
	}

	return int(userId), nil
}

func UpdateUser(user *User) error {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	updateUser := `UPDATE users
					SET firstname = ?,
						lastname = ?,
						gender = ?,
						age = ?,
						profile_photo = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := db.Exec(updateUser, user.Firstname, user.Lastname, user.Gender, user.Age, user.ProfilePhoto, user.ID, user.ID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func AuthenticateUser(username, password string) (bool, int, error) {
	// Open SQLite database
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query to retrieve the hashed password stored in the database for the given username
	var userId int
	var storedHashedPassword string
	err := db.QueryRow("SELECT id, password FROM users WHERE username = ? or email = ?", username, username).Scan(&userId, &storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			// Username not found
			return false, -1, errors.New("username not found")
		}
		// Handle other database errors
		log.Fatal(err)
	}

	// Compare the entered password with the stored hashed password using bcrypt
	err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
	if err != nil {
		// Password is incorrect
		return false, -1, errors.New("password is incorrect")
	}

	// Successful login if no errors occurred
	return true, userId, nil
}

func ReadAllUsers() ([]User, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT u.id as user_id, u.firstname as user_firstname, u.lastname as user_lastname, u.username as user_username, u.email as user_email, 
		IFNULL(u.profile_photo, '') as profile_photo, u.status as user_status, u.created_at as user_created_at, 
		u.updated_at as user_updated_at, u.updated_by as user_updated_by
		FROM users u
		WHERE u.status != 'delete'
		AND u.type != 'admin'
		ORDER BY u.id desc;
    `)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		// Scan the post and user data
		err := rows.Scan(
			&user.ID, &user.Firstname, &user.Lastname, &user.Username, &user.Email,
			&user.ProfilePhoto, &user.Status, &user.CreatedAt,
			&user.UpdatedAt, &user.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		users = append(users, user)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return users, nil
}

func ReadAllChatUsers(user_id int) ([]User, error) {
	db := db.OpenDBConnection()
	defer db.Close()

	rows, selectError := db.Query(`
        SELECT u.id as user_id, u.firstname as user_firstname, u.lastname as user_lastname, u.username as user_username, u.email as user_email, 
		IFNULL(u.profile_photo, '') as profile_photo, u.status as user_status, u.created_at as user_created_at, 
		u.updated_at as user_updated_at, u.updated_by as user_updated_by
			FROM users u
			LEFT JOIN chat_members cm 
				ON u.id = cm.user_id
				AND cm.status != 'delete'
			LEFT JOIN chats c
				ON cm.chat_id = c.id
				AND c.status != 'delete'
				AND c.id in (
					SELECT chat_id
					FROM chat_members
					WHERE user_id = ?
					AND status != 'delete'
				)
			WHERE u.status != 'delete'
			AND u.type != 'admin'
			GROUP BY u.id
			ORDER BY MAX(c.updated_at) desc, u.username;
    `, user_id)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User

		err := rows.Scan(
			&user.ID, &user.Firstname, &user.Lastname, &user.Username, &user.Email,
			&user.ProfilePhoto, &user.Status, &user.CreatedAt,
			&user.UpdatedAt, &user.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return users, nil
}

func ReadUserByID(user_id int) (User, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT u.id as user_id, u.firstname as user_firstname, u.lastname as user_lastname, u.username as user_username, u.email as user_email, 
		IFNULL(u.profile_photo, '') as profile_photo, u.status as user_status, u.created_at as user_created_at, 
		u.updated_at as user_updated_at, u.updated_by as user_updated_by
		FROM users u
		WHERE u.status != 'delete'
		AND u.type != 'admin'
		AND u.id = ?
		ORDER BY u.id desc;
    `, user_id)
	if selectError != nil {
		return User{}, selectError
	}
	defer rows.Close()

	var user User

	for rows.Next() {
		// Scan the post and user data
		err := rows.Scan(
			&user.ID, &user.Firstname, &user.Lastname, &user.Username, &user.Email,
			&user.ProfilePhoto, &user.Status, &user.CreatedAt,
			&user.UpdatedAt, &user.UpdatedBy,
		)
		if err != nil {
			return User{}, fmt.Errorf("error scanning row: %v", err)
		}

	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return User{}, fmt.Errorf("row iteration error: %v", err)
	}

	return user, nil
}

func UpdateStatusUser(user_id int, status string, login_user_id int) error {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	updateQuery := `UPDATE users
					SET status = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := db.Exec(updateQuery, status, login_user_id, user_id)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func GetUserIDByUsername(username string) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	var userID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return -1, errors.New("user not found")
		}
		return -1, err // Other DB errors
	}

	return userID, nil
}
