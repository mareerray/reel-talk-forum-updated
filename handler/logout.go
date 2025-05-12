package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/model"
	"time"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get session cookie
	cookie, err := r.Cookie("session_token")
	if err != nil {
		respondWithError(w, "No active session", "general", http.StatusUnauthorized)
		return
	}

	// Begin transaction
	tx, err := DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		respondWithError(w, "Error during logout", "general", http.StatusInternalServerError)
		return
	}

	// Delete session from sessions table
	_, err = tx.Exec("DELETE FROM sessions WHERE session_token = ?", cookie.Value)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to delete session: %v", err)
		respondWithError(w, "Error during logout", "general", http.StatusInternalServerError)
		return
	}

	// Clear session from users table
	_, err = tx.Exec(
		"UPDATE users SET session_token = NULL, session_expiry = NULL WHERE session_token = ?",
		cookie.Value,
	)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to clear user session: %v", err)
		respondWithError(w, "Error during logout", "general", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		respondWithError(w, "Error during logout", "general", http.StatusInternalServerError)
		return
	}

	// Expire the cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode, // <-- Required for cross-origin
	})

	json.NewEncoder(w).Encode(model.Response{
		Success: true,
		Message: "Logout successful",
	})
}
