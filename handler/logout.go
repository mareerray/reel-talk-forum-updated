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

    cookie, err := r.Cookie("session_token")
    if err != nil {
        respondWithError(w, "No active session", "general", http.StatusUnauthorized)
        return
    }

    _, err = DB.Exec("DELETE FROM sessions WHERE session_token = $1", cookie.Value)
    if err != nil {
        log.Printf("Failed to delete session: %v", err)
        respondWithError(w, "Error during logout", "general", http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "session_token",
        Value:    "",
        Path:     "/",
        Expires:  time.Unix(0, 0),
        MaxAge:   -1,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteLaxMode,
    })

    json.NewEncoder(w).Encode(model.Response{
        Success: true,
        Message: "Logout successful",
    })
}

