package handler

import (
    "encoding/json"
    "net/http"
    "strings"
)

func CreatePostHandler(w http.ResponseWriter, r *http.Request) {
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
    err := DB.QueryRow(`
        SELECT user_id FROM sessions
        WHERE session_token = ?
        AND is_active = true
        AND session_expiry > datetime('now')`, token).Scan(&userID)
    
    if err != nil {
        http.Error(w, "Invalid session", http.StatusUnauthorized)
        return
    }
    
    // Parse request
    var req struct {
        Title      string `json:"title"`
        Content    string `json:"content"`
        Categories string `json:"categories"` // Single category with emoji
    }
    
    err = json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }
    
    // Validate input
    if req.Title == "" || len(req.Title) > 100 {
        http.Error(w, "Title must be between 1 and 100 characters", http.StatusBadRequest)
        return
    }
    
    if req.Content == "" {
        http.Error(w, "Content cannot be empty", http.StatusBadRequest)
        return
    }
    
    // Validate category
    if req.Categories == "" {
        http.Error(w, "A category must be selected", http.StatusBadRequest)
        return
    }
    
    // Verify the category exists in the database
    parts := strings.Split(req.Categories, " ")
    if len(parts) < 2 {
        http.Error(w, "Invalid category format", http.StatusBadRequest)
        return
    }
    
    categoryName := strings.TrimSpace(parts[0])
    var exists bool
    err = DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE name = ?)", categoryName).Scan(&exists)
    if err != nil || !exists {
        http.Error(w, "Invalid category: "+categoryName, http.StatusBadRequest)
        return
    }
    
    // Insert the post
    result, err := DB.Exec(`
        INSERT INTO posts (user_id, title, content, categories)
        VALUES (?, ?, ?, ?)`,
        userID, req.Title, req.Content, req.Categories)
    
    if err != nil {
        http.Error(w, "Failed to create post: "+err.Error(), http.StatusInternalServerError)
        return
    }
    
    postID, err := result.LastInsertId()
    if err != nil {
        http.Error(w, "Failed to get post ID", http.StatusInternalServerError)
        return
    }
    
    // Return success response
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status": "success",
        "post_id": postID,
    })
}

