package models

import (
	"database/sql"
	"fmt"
	"forum/db"
	userManagementModels "forum/modules/userManagement/models"
	"log"
	"time"
)

// Category struct represents the user data model
type Category struct {
	ID             int                       `json:"id"`
	Name           string                    `json:"name"`
	Color          string                    `json:"color"`
	Icon           string                    `json:"icon"`
	Status         string                    `json:"status"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      *time.Time                `json:"updated_at"`
	CreatedBy      int                       `json:"created_by"`
	UpdatedBy      *int                      `json:"updated_by"`
	User           userManagementModels.User `json:"user"` // Embedded user data
	PostsCount     *int                      `json:"posts_count"`
	CommentsCount  *int                      `json:"comments_count"`
	PostLikesCount *int                      `json:"post_likes_count"`
}

func InsertCategory(category *Category) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	insertQuery := `INSERT INTO categories (name, color, icon) VALUES (?, ?, ?);`
	result, insertErr := db.Exec(insertQuery, category.Name, category.Color, category.Icon)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := insertErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return -1, sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
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

func UpdateCategory(category *Category, userId int) error {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	updateQuery := `UPDATE categories
					SET name = ?,
						color = ?,
						icon = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := db.Exec(updateQuery, category.Name, category.Color, category.Icon, userId, category.ID)
	if updateErr != nil {
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := updateErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return updateErr
	}

	return nil
}

func UpdateStatuCategory(categoryId int, status string, userId int) error {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	updateQuery := `UPDATE categories
					SET status = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := db.Exec(updateQuery, status, userId, categoryId)
	if updateErr != nil {
		// Check if the error is a SQLite constraint violation
		if sqliteErr, ok := updateErr.(interface{ ErrorCode() int }); ok {
			if sqliteErr.ErrorCode() == 19 { // SQLite constraint violation error code
				return sql.ErrNoRows // Return custom error to indicate a duplicate
			}
		}
		return updateErr
	}

	return nil
}

func AdminReadAllCategories() ([]Category, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT c.id as category_id, c.name as category_name, c.color as category_color, c.icon as category_icon, c.status as category_status, 
               c.created_at as category_created_at, c.created_by as category_created_by, 
               c.updated_at as category_updated_at, c.updated_by as category_updated_by,
               u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email, IFNULL(u.profile_photo, '') as user_profile_photo,
			   (SELECT COUNT(DISTINCT p.id) 
			   	FROM post_categories pc
				INNER JOIN posts p
					ON pc.post_id = p.id
				WHERE p.status != 'delete'
				AND pc.status != 'delete'
				AND pc.category_id = c.id
			   ) as posts_count,
			   (SELECT COUNT(DISTINCT com.id) 
			   	FROM post_categories pc
				INNER JOIN posts p
					ON pc.post_id = p.id
				INNER JOIN comments com
					ON com.post_id = p.id
					AND com.status != 'delete'
				WHERE p.status != 'delete'
				AND pc.status != 'delete'
				AND pc.category_id = c.id
			   ) as comments_count,
			   (SELECT COUNT(DISTINCT pl.id) 
			   	FROM post_categories pc
				INNER JOIN posts p
					ON pc.post_id = p.id
				INNER JOIN post_likes pl
					ON pl.post_id = p.id
					AND pl.status != 'delete'
				WHERE p.status != 'delete'
				AND pc.status != 'delete'
				AND pc.category_id = c.id
			   ) as post_likes_count
        FROM categories c
        INNER JOIN users u ON c.created_by = u.id
        WHERE c.status != 'delete';
    `)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category
		var user userManagementModels.User

		// Scan the data into variables
		err := rows.Scan(
			&category.ID, &category.Name, &category.Color, &category.Icon, &category.Status, &category.CreatedAt, &category.CreatedBy,
			&category.UpdatedAt, &category.UpdatedBy,
			&user.ID, &user.Name, &user.Username, &user.Email, &user.ProfilePhoto,
			&category.PostsCount, &category.CommentsCount, &category.PostLikesCount,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Assign the user to the category
		category.User = user

		// Append category to the categories slice
		categories = append(categories, category)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return categories, nil
}

func ReadAllCategories() ([]Category, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT c.id as category_id, c.name as category_name, c.color as category_color, c.icon as category_icon, c.status as category_status, 
               c.created_at as category_created_at, c.created_by as category_created_by, 
               c.updated_at as category_updated_at, c.updated_by as category_updated_by,
               u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email, IFNULL(u.profile_photo, '') as user_profile_photo
        FROM categories c
        INNER JOIN users u ON c.created_by = u.id
        WHERE c.status != 'delete';
    `)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		var category Category
		var user userManagementModels.User

		// Scan the data into variables
		err := rows.Scan(
			&category.ID, &category.Name, &category.Color, &category.Icon, &category.Status, &category.CreatedAt, &category.CreatedBy,
			&category.UpdatedAt, &category.UpdatedBy,
			&user.ID, &user.Name, &user.Username, &user.Email, &user.ProfilePhoto,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Assign the user to the category
		category.User = user

		// Append category to the categories slice
		categories = append(categories, category)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	return categories, nil
}

func ReadCategoryById(categoryId int) (Category, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT c.id as category_id, c.name as category_name, c.color as category_color, c.icon as category_icon, c.status as category_status, 
               c.created_at as category_created_at, c.created_by as category_created_by, 
               c.updated_at as category_updated_at, c.updated_by as category_updated_by,
               u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email, IFNULL(u.profile_photo, '') as user_profile_photo
        FROM categories c
        INNER JOIN users u ON c.created_by = u.id  -- Fixed the JOIN to use the correct column for user relation
        WHERE c.status != 'delete'
        AND c.id = ?;
    `, categoryId)
	if selectError != nil {
		return Category{}, selectError
	}
	defer rows.Close()

	// Variable to hold the category and user data
	var category Category
	var user userManagementModels.User

	// Scan the result into variables
	if rows.Next() {
		err := rows.Scan(
			&category.ID, &category.Name, &category.Color, &category.Icon, &category.Status, &category.CreatedAt, &category.CreatedBy,
			&category.UpdatedAt, &category.UpdatedBy,
			&user.ID, &user.Name, &user.Username, &user.Email, &user.ProfilePhoto,
		)
		if err != nil {
			return Category{}, fmt.Errorf("error scanning row: %v", err)
		}

		// Assign the user to the category
		category.User = user
	} else {
		// If no category found with the given ID
		return Category{}, fmt.Errorf("category with ID %d not found", categoryId)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return Category{}, fmt.Errorf("row iteration error: %v", err)
	}

	return category, nil
}

func ReadCategoryByName(categoryName string) (Category, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Query the records
	rows, selectError := db.Query(`
        SELECT c.id as category_id, c.name as category_name, c.color as category_color, c.icon as category_icon, c.status as category_status, 
               c.created_at as category_created_at, c.created_by as category_created_by, 
               c.updated_at as category_updated_at, c.updated_by as category_updated_by,
               u.id as user_id, u.name as user_name, u.username as user_username, u.email as user_email, IFNULL(u.profile_photo, '') as user_profile_photo
        FROM categories c
        INNER JOIN users u ON c.created_by = u.id  -- Fixed the JOIN to use the correct column for user relation
        WHERE c.status != 'delete'
        AND c.name = ?;
    `, categoryName)
	if selectError != nil {
		return Category{}, selectError
	}
	defer rows.Close()

	// Variable to hold the category and user data
	var category Category
	var user userManagementModels.User

	// Scan the result into variables
	if rows.Next() {
		err := rows.Scan(
			&category.ID, &category.Name, &category.Color, &category.Icon, &category.Status, &category.CreatedAt, &category.CreatedBy,
			&category.UpdatedAt, &category.UpdatedBy,
			&user.ID, &user.Name, &user.Username, &user.Email, &user.ProfilePhoto,
		)
		if err != nil {
			return Category{}, fmt.Errorf("error scanning row: %v", err)
		}

		// Assign the user to the category
		category.User = user
	} else {
		// If no category found with the given Name
		return Category{}, fmt.Errorf("category with Name %v not found", categoryName)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return Category{}, fmt.Errorf("row iteration error: %v", err)
	}

	return category, nil
}
