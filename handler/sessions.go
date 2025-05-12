package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/model"
	"real-time-forum/utils"
	"strings"
	"time"
)

func CheckUserLoggedIn(r *http.Request) (bool, int) {
	sessionToken, err := r.Cookie("session_token")
	if err != nil {
		return false, 0
	}

	var userID int
	var expiresAt time.Time

	query := `
    SELECT user_id, session_expiry 
    FROM sessions 
    WHERE session_token = ? 
    AND session_token = (
        SELECT session_token 
        FROM sessions AS s2 
        WHERE s2.user_id = sessions.user_id 
        ORDER BY session_expiry DESC 
        LIMIT 1
    )
	`

	err = DB.QueryRow(query, sessionToken.Value).Scan(&userID, &expiresAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, 0
		}
		log.Printf("Error checking session: %v", err)
		return false, 0
	}

	if time.Now().After(expiresAt) {
		// Session has expired
		DeleteSession(sessionToken.Value)
		return false, 0
	}

	return true, userID
}

func InsertSession(session *model.Session) (*model.Session, error) {
	db := utils.OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes

	// Generate UUID for the user if not already set
	if session.SessionToken == "" {
		uuidSessionTokenid, err := utils.GenerateUuid()
		if err != nil {
			return nil, err
		}
		session.SessionToken = uuidSessionTokenid
	}

	// Set session expiration time
	session.ExpiresAt = time.Now().Add(12 * time.Hour)

	// Start a transaction for atomicity
	tx, err := db.Begin()
	if err != nil {
		return &model.Session{}, err
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

func CreateSession(w http.ResponseWriter, userID int) error {
	// First, invalidate any existing session for this user
	_, err := DB.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if err != nil {
		log.Printf("Error deleting old sessions: %v", err)
	}
	token, cookie, err := generateSessionToken()
	if err != nil {
		return err
	}

	query := `
    INSERT INTO sessions 
    (user_id, session_token, session_expiry) 
    VALUES (?, ?, ?)
	`

	_, err = DB.Exec(query, userID, token, cookie.Expires)
	if err != nil {
		return fmt.Errorf("failed to insert session: %v", err)
	}

	http.SetCookie(w, cookie)
	return nil
}

// -- Non-Global Functions : Only happens in this package server -- //

func DeleteSession(token string) {
	_, err := DB.Exec("DELETE FROM sessions WHERE session_token = ?", token)
	if err != nil {
		log.Printf("Error deleting session: %v", err)
	}
}

func generateSessionToken() (string, *http.Cookie, error) {
    token, err := utils.GenerateUuid()
    if err != nil {
        return "", nil, err
    }

    expiration := time.Now().Add(24 * time.Hour)
    cookie := &http.Cookie{
        Name:     "session_token",
        Value:    token,
        Expires:  expiration,
        Path:     "/",
        HttpOnly: true,      // Prevent XSS
        Secure:   true,     // Allow HTTP for local development
        SameSite: http.SameSiteLaxMode,  // Prevent CSRF
        MaxAge:   86400,     // 24h in seconds (backup for Expires)
    }
    return token, cookie, nil
}


func ValidateSessionHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get token from header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing authorization header", http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Check session validity
	var isValid bool
	err := DB.QueryRow(`
        SELECT EXISTS(
            SELECT 1 
            FROM sessions 
            WHERE session_token = ? 
            AND is_active = true 
            AND session_expiry > datetime('now')
        )`, token).Scan(&isValid)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !isValid {
		http.Error(w, "Invalid session", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(model.Response{Success: true})

	// Get user ID associated with the session
	var userID int
	err = DB.QueryRow(`
        SELECT user_id 
        FROM sessions 
        WHERE session_token = ?`, token).Scan(&userID)

	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"user_id": userID,
	})
}
