package models

import (
	"forum/db"
	userManagementModels "forum/modules/userManagement/models"
	"log"
	"time"
)

type Comment struct {
	ID               int                       `json:"id"`
	PostId           int                       `json:"post_id"`
	Description      string                    `json:"description"`
	UserId           int                       `json:"user_id"`
	Status           string                    `json:"status"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        *time.Time                `json:"updated_at"`
	UpdatedBy        *int                      `json:"updated_by"`
	IsLikedByUser    bool                      `json:"liked"`
	IsDislikedByUser bool                      `json:"disliked"`
	NumberOfLikes    int                       `json:"number_of_likes"`
	NumberOfDislikes int                       `json:"number_of_dislikes"`
	Post             Post                      `json:"post"`
	User             userManagementModels.User `json:"user"`
}

func InsertComment(postId int, userId int, description string) (int, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	insertQuery := `INSERT INTO comments (post_id, description, user_id) VALUES (?, ?, ?);`
	result, insertErr := db.Exec(insertQuery, postId, description, userId)
	if insertErr != nil {
		// Check if the error is a SQLite constraint violation
		return -1, insertErr
	}
	// Retrieve the last inserted ID
	lastInsertID, errFind := result.LastInsertId()
	if errFind != nil {
		log.Fatal(errFind)
		return -1, errFind
	}
	return int(lastInsertID), nil
}

func UpdateComment(comment *Comment, user_id int, newDescription string) error {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Start a transaction for atomicity
	updateQuery := `UPDATE comments
					SET description = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := db.Exec(updateQuery, newDescription, user_id, comment.ID)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func UpdateCommentStatus(id int, status string, user_id int) error {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	updateQuery := `UPDATE comments
					SET status = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := db.Exec(updateQuery, status, user_id, id)
	if updateErr != nil {
		return updateErr
	}

	return nil
}

func ReadAllComments() ([]Comment, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	var comments []Comment
	selectQuery := `
		SELECT 
			p.id AS post_id, p.uuid AS post_uuid, p.title AS post_title, p.description AS post_description, 
			p.status AS post_status, p.created_at AS post_created_at, p.updated_at AS post_updated_at, p.updated_by AS post_updated_by,
			c.id AS comment_id, c.post_id AS comment_post_id ,c.description AS comment_description,c.user_id AS comment_user_id, 
			c.status AS comment_status, c.created_at AS comment_created_at, c.updated_at AS comment_updated_at, c.updated_by AS comment_updated_by,
			u.id AS user_id, u.uuid AS user_uuid, u.username AS user_username, u.name AS user_name, u.type AS user_type, u.email AS user_email, IFNULL(u.profile_photo, '') as user_profile_photo, 
			u.status AS user_status, u.created_at AS user_created_at, u.updated_at AS user_updated_at, u.updated_by AS user_updated_by
		FROM comments c
		INNER JOIN posts p ON c.post_id = p.id AND p.status != 'delete' AND c.status != 'delete'
		INNER JOIN users u ON c.user_id = u.id AND u.status != 'delete'
		ORDER BY c.id desc
	`
	// Query the records
	rows, selectError := db.Query(selectQuery)

	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close() // Ensure rows are closed after processing

	// Iterate over rows and populate the slice
	for rows.Next() {
		var comment Comment
		var user userManagementModels.User
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.UUID,
			&post.Title,
			&post.Description,
			&post.Status,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.UpdatedBy,

			&comment.ID,
			&comment.PostId,
			&comment.Description,
			&comment.UserId,
			&comment.Status,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.UpdatedBy,

			&user.ID,
			&user.UUID,
			&user.Username,
			&user.Name,
			&user.Type,
			&user.Email,
			&user.ProfilePhoto,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.UpdatedBy,
		)

		if err != nil {
			return nil, err
		}
		comment.Post = post
		comment.User = user

		comments = append(comments, comment)
	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func ReadCommentsFromUserId(userId int) ([]Comment, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	var comments []Comment

	// Updated query to join comments with posts
	selectQuery := `
		SELECT 
			p.id AS post_id, p.uuid AS post_uuid, p.title AS post_title, p.description AS post_description, 
			p.status AS post_status, p.created_at AS post_created_at, p.updated_at AS post_updated_at, p.updated_by AS post_updated_by,
			c.id AS comment_id, c.user_id AS comment_user_id, c.description AS comment_description, 
			c.status AS comment_status, c.created_at AS comment_created_at, c.updated_at AS comment_updated_at, c.updated_by AS comment_updated_by
		FROM comments c
		INNER JOIN posts p ON c.post_id = p.i
		WHERE c.status != 'delete' AND p.status != 'delete' AND c.user_id = ?
		ORDER BY c.id desc;
	`

	rows, selectError := db.Query(selectQuery, userId) // Query the database
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close() // Ensure rows are closed after processing

	// Iterate over rows and populate the slice
	for rows.Next() {
		var comment Comment
		var post Post

		err := rows.Scan(
			// Map post fields
			&post.ID,
			&post.UUID,
			&post.Title,
			&post.Description,
			&post.Status,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.UpdatedBy,

			// Map comment fields
			&comment.ID,
			&comment.UserId,
			&comment.Description,
			&comment.Status,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.UpdatedBy,
		)

		if err != nil {
			return nil, err
		}

		// Assign the post to the comment
		comment.Post = post

		comments = append(comments, comment)
	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func ReadAllCommentsForPost(postId int) ([]Comment, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	var comments []Comment
	commentMap := make(map[int]*Comment)
	// Updated query to join comments with posts
	selectQuery := `
		SELECT 
			u.id AS user_id, u.uuid AS user_uuid, u.username AS user_username, u.name AS user_name, u.type AS user_type, u.email AS user_email, IFNULL(u.profile_photo, '') as user_profile_photo, 
			u.status AS user_status, u.created_at AS user_created_at, u.updated_at AS user_updated_at, u.updated_by AS user_updated_by,
			c.id AS comment_id, c.post_id as comment_post_id, c.user_id AS comment_user_id, c.description AS comment_description, 
			c.status AS comment_status, c.created_at AS comment_created_at, c.updated_at AS comment_updated_at, c.updated_by AS comment_updated_by,
			COALESCE(cl.type, '')
		FROM comments c
			INNER JOIN users u
				ON c.user_id = u.id AND c.status != 'delete' AND u.status != 'delete' AND c.post_id = ?
			LEFT JOIN comment_likes cl
				ON c.id = cl.comment_id AND cl.status != 'delete'
		ORDER BY c.id desc;
	`
	rows, selectError := db.Query(selectQuery, postId) // Query the database
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close() // Ensure rows are closed after processing

	// Iterate over rows and populate the slice
	for rows.Next() {
		var comment Comment
		var user userManagementModels.User
		var Type string
		err := rows.Scan(
			// Map post fields
			&user.ID,
			&user.UUID,
			&user.Username,
			&user.Name,
			&user.Type,
			&user.Email,
			&user.ProfilePhoto,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.UpdatedBy,

			// Map comment fields
			&comment.ID,
			&comment.PostId,
			&comment.UserId,
			&comment.Description,
			&comment.Status,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.UpdatedBy,

			&Type,
		)
		comment.User = user
		if err != nil {
			return nil, err
		}

		if existingComment, found := commentMap[comment.ID]; found {
			if Type == "like" {
				existingComment.NumberOfLikes++
			} else if Type == "dislike" {
				existingComment.NumberOfDislikes++
			}
		} else {
			if Type == "like" {
				comment.NumberOfLikes++
			} else if Type == "dislike" {
				comment.NumberOfDislikes++
			}

			commentMap[comment.ID] = &comment
		}

	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Convert the map of comments into a slice
	for _, comment := range commentMap {
		comments = append(comments, *comment)
	}

	return comments, nil
}

func ReadAllCommentsForPostByUserID(postId int, userID int) ([]Comment, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	var comments []Comment
	commentMap := make(map[int]*Comment)
	// Updated query to join comments with posts
	selectQuery := `
		SELECT 
			u.id AS user_id, u.uuid AS user_uuid, u.username AS user_username, u.name AS user_name, u.type AS user_type, u.email AS user_email, IFNULL(u.profile_photo, '') as user_profile_photo, 
			u.status AS user_status, u.created_at AS user_created_at, u.updated_at AS user_updated_at, u.updated_by AS user_updated_by,
			c.id AS comment_id, c.post_id as comment_post_id ,c.user_id AS comment_user_id, c.description AS comment_description, 
			c.status AS comment_status, c.created_at AS comment_created_at, c.updated_at AS comment_updated_at, c.updated_by AS comment_updated_by,
			(SELECT COUNT(DISTINCT id) from comment_likes WHERE comment_id = c.id AND status != 'delete' AND type = 'like') AS number_of_likes,
			(SELECT COUNT(DISTINCT id) from comment_likes WHERE comment_id = c.id AND status != 'delete' AND type = 'dislike') AS number_of_dislikes,
			CASE 
                WHEN EXISTS (SELECT 1 FROM comment_likes WHERE comment_id = c.id AND status != 'delete' AND type = 'like' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_liked_by_user,
            CASE 
                WHEN EXISTS (SELECT 1 FROM comment_likes WHERE comment_id = c.id AND status != 'delete' AND type = 'dislike' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_disliked_by_user
		FROM comments c
			INNER JOIN users u
				ON c.user_id = u.id AND c.status != 'delete' AND u.status != 'delete' AND c.post_id = ?	
		ORDER BY c.id desc;
	`
	rows, selectError := db.Query(selectQuery, userID, userID, postId) // Query the database
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close() // Ensure rows are closed after processing
	// Iterate over rows and populate the slice
	for rows.Next() {
		var comment Comment
		var user userManagementModels.User

		err := rows.Scan(
			// Map post fields
			&user.ID,
			&user.UUID,
			&user.Username,
			&user.Name,
			&user.Type,
			&user.Email,
			&user.ProfilePhoto,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.UpdatedBy,

			// Map comment fields
			&comment.ID,
			&comment.PostId,
			&comment.UserId,
			&comment.Description,
			&comment.Status,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.UpdatedBy,

			&comment.NumberOfLikes, &comment.NumberOfDislikes,
			&comment.IsLikedByUser, &comment.IsDislikedByUser,
		)
		comment.User = user
		if err != nil {
			return nil, err
		}

		_, found := commentMap[comment.ID]
		if !found {
			commentMap[comment.ID] = &comment
		}

	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}
	// Convert the map of comments into a slice
	for _, comment := range commentMap {
		comments = append(comments, *comment)
	}

	return comments, nil
}

// func ReadAllCommentsForPostLikedByUser(postId int, userId int) ([]Comment, error) {
// 	db := db.OpenDBConnection()
// 	defer db.Close() // Close the connection after the function finishes
//
// 	var comments []Comment
//
// 	// Updated query to join comments with posts
// 	selectQuery := `
// 		SELECT
// 			u.id AS user_id, u.uuid AS user_uuid, u.username AS user_username, u.name AS user_name, u.type AS user_type, u.email AS user_email, IFNULL(u.profile_photo, '') as user_profile_photo,
// 			u.status AS user_status, u.created_at AS user_created_at, u.updated_at AS user_updated_at, u.updated_by AS user_updated_by,
// 			c.id AS comment_id, c.user_id AS comment_user_id, c.description AS comment_description,
// 			c.status AS comment_status, c.created_at AS comment_created_at, c.updated_at AS comment_updated_at, c.updated_by AS comment_updated_by,
// 			count(CASE WHEN cl.type = 'like' THEN 1 END) as likes_count, count(CASE WHEN cl.type = 'dislike' THEN 1 END) as dislikes_count
// 		FROM comments c
// 			INNER JOIN users u
// 				ON c.user_id = u.id AND c.status != 'delete' AND u.status != 'delete' AND c.post_id = ?
// 			INNER JOIN comment_likes cl
// 				ON c.id = cl.comment_id AND cl.status != 'delete'
// 		GROUP BY cl.comment_id;
// 	`
// 	rows, selectError := db.Query(selectQuery, postId) // Query the database
// 	if selectError != nil {
// 		return nil, selectError
// 	}
// 	defer rows.Close() // Ensure rows are closed after processing
//
// 	// Iterate over rows and populate the slice
// 	for rows.Next() {
// 		var comment Comment
// 		var user userManagementModels.User
//
// 		err := rows.Scan(
// 			// Map post fields
// 			&user.ID,
// 			&user.UUID,
// 			&user.Username,
// 			&user.Name,
// 			&user.Type,
// 			&user.Email,
// 			&user.ProfilePhoto,
// 			&user.Status,
// 			&user.CreatedAt,
// 			&user.UpdatedAt,
// 			&user.UpdatedBy,
//
// 			// Map comment fields
// 			&comment.ID,
// 			&comment.UserId,
// 			&comment.Description,
// 			&comment.Status,
// 			&comment.CreatedAt,
// 			&comment.UpdatedAt,
// 			&comment.UpdatedBy,
// 			&comment.NumberOfLikes,
// 			&comment.NumberOfDislikes,
// 		)
//
// 		if err != nil {
// 			return nil, err
// 		}
//
// 		// Assign the post to the comment
// 		comment.User = user
//
// 		comments = append(comments, comment)
// 	}
//
// 	// Check for any errors during the iteration
// 	if err := rows.Err(); err != nil {
// 		return nil, err
// 	}
//
// 	return comments, nil
// }

func ReadAllCommentsOfUserForPost(postId int, userId int) ([]Comment, error) {
	db := db.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	var comments []Comment
	selectQuery := `
		SELECT 
			p.id AS post_id, p.uuid AS post_uuid, p.title AS post_title, p.description AS post_description, 
			p.status AS post_status, p.created_at AS post_created_at, p.updated_at AS post_updated_at, p.updated_by AS post_updated_by,
			c.id AS comment_id, c.user_id AS comment_user_id, c.description AS comment_description, 
			c.status AS comment_status, c.created_at AS comment_created_at, c.updated_at AS comment_updated_at, c.updated_by AS comment_updated_by,
			u.id AS user_id, u.uuid AS user_uuid, u.username AS user_username, u.name AS user_name, u.type AS user_type, u.email AS user_email, IFNULL(u.profile_photo, '') as user_profile_photo, 
			u.status AS user_status, u.created_at AS user_created_at, u.updated_at AS user_updated_at, u.updated_by AS user_updated_by
		FROM comments c
		INNER JOIN posts p ON c.post_id = p.id AND p.status != 'delete' AND c.status != 'delete' AND p.id = ?
		INNER JOIN users u ON c.user_id = u.id AND u.status != 'delete'
		where u.id = ?
		ORDER BY c.id desc
	`
	// Query the records
	rows, selectError := db.Query(selectQuery, postId, userId)

	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close() // Ensure rows are closed after processing

	// Iterate over rows and populate the slice
	for rows.Next() {
		var comment Comment
		var user userManagementModels.User
		var post Post
		err := rows.Scan(
			&post.ID,
			&post.UUID,
			&post.Title,
			&post.Description,
			&post.Status,
			&post.CreatedAt,
			&post.UpdatedAt,
			&post.UpdatedBy,

			&comment.ID,
			&comment.UserId,
			&comment.Description,
			&comment.Status,
			&comment.CreatedAt,
			&comment.UpdatedAt,
			&comment.UpdatedBy,

			&user.ID,
			&user.UUID,
			&user.Username,
			&user.Name,
			&user.Type,
			&user.Email,
			&user.ProfilePhoto,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.UpdatedBy,
		)

		if err != nil {
			return nil, err
		}
		comment.Post = post
		comment.User = user

		comments = append(comments, comment)
	}

	// Check for any errors during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
