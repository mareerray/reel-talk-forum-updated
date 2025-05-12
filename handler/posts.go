// In handler/posts.go
package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/model"
	"time"
)

// GetPostsHandler handles requests for fetching all posts
func GetPostsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Query the database for all posts with username
	rows, err := DB.Query(`
		SELECT p.id, p.user_id, p.title, p.content, p.categories, 
			p.created_at, p.updated_at, u.nickname 
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		log.Printf("Error querying posts: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to hold the posts
	var posts []model.Post

	// Iterate through the rows and populate the posts slice
	for rows.Next() {
		var post model.Post
		var createdAt, updatedAt string

		if err := rows.Scan(
			&post.ID, &post.UserID, &post.Title, &post.Content,
			&post.Categories, &createdAt, &updatedAt, &post.Nickname,
		); err != nil {
			log.Printf("Error scanning post row: %v", err)
			continue
		}

		// Parse the time strings
		post.CreatedAt, _ = time.Parse("02-01-2006 15:04:05", createdAt)
		post.UpdatedAt, _ = time.Parse("02-01-2006 15:04:05", updatedAt)

		posts = append(posts, post)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		log.Printf("Error iterating post rows: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Set the content type header
	w.Header().Set("Content-Type", "application/json")

	// Encode the posts as JSON and send the response
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		log.Printf("Error encoding posts to JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
