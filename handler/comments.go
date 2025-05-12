package handler

import (
	"encoding/json"
	"net/http"
	"real-time-forum/model"
	"strconv"
	"strings"
)

func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get token from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing authorization header", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Get user ID from session
	var userID int
	var nickname string
	err := DB.QueryRow(`
		SELECT s.user_id, u.nickname 
        FROM sessions s
        JOIN users u ON s.user_id = u.id
        WHERE s.session_token = ? 
        AND s.is_active = true 
        AND s.session_expiry > datetime('now')`, token).Scan(&userID, &nickname)

	if err != nil {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	var comment struct {
		PostID  int    `json:"post_id"`
		Content string `json:"content"`
	}

	// Decode JSON
	err = json.NewDecoder(r.Body).Decode(&comment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	// input validation for the comment content and post ID
	if comment.PostID <= 0 {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}
	if len(comment.Content) == 0 || len(comment.Content) > 200 {
		http.Error(w, "Comment must be between 1 and 200 characters", http.StatusBadRequest)
		return
	}

	// Execute SQL with the userID from the session
	_, err = DB.Exec(`INSERT INTO comments 
        (user_id, user_name, post_id, content) 
        VALUES (?, ?, ?, ?)`,
		userID, nickname, comment.PostID, comment.Content)

	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func GetCommentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	postID := strings.TrimPrefix(r.URL.Path, "/api/comments/")
	id, err := strconv.Atoi(postID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid post ID"})
		return
	}

	rows, err := DB.Query(`
        SELECT user_name, content, created_at 
        FROM comments 
        WHERE post_id = ? 
        ORDER BY created_at DESC`, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
		return
	}
	defer rows.Close()

	comments := []model.Comment{}
	for rows.Next() {
		var c model.Comment
		if err := rows.Scan(&c.UserName, &c.Content, &c.CreatedAt); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to read comments"})
			return
		}
		comments = append(comments, c)
	}

	// Return empty array if no comments exist
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}
