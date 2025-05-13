package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"real-time-forum/model"
	"real-time-forum/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles user registration requests
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Received registration request")

	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		respondWithError(w, "Invalid request format", "general", http.StatusBadRequest)
		return
	}

	log.Printf("Registration request for nickname: %s, email: %s", req.Username, req.Email)

	if req.FirstName == "" || req.LastName == "" || req.Gender == "" || req.Age == "" {
		respondWithError(w, "All fields are required", "general", http.StatusBadRequest)
		return
	}

	if req.Gender != "Male" && req.Gender != "Female" && req.Gender != "Other" {
		respondWithError(w, "Gender must be Male, Female, or Other", "gender", http.StatusBadRequest)
		return
	}

	if err := utils.ValidateInputs(DB, req.Username, req.Email, req.Password); err != nil {
		if valErr, ok := err.(utils.ValidationError); ok {
			log.Printf("Validation error: %v", valErr)
			respondWithError(w, valErr.Message, valErr.Field, http.StatusBadRequest)
			return
		}
		log.Printf("Validation error: %v", err)
		respondWithError(w, "An error occurred during validation", "general", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		respondWithError(w, "Error processing your request", "general", http.StatusInternalServerError)
		return
	}

	// Generate UUID for user
	userUUID, err := utils.GenerateUuid()
	if err != nil {
		log.Printf("UUID generation failed: %v", err)
		respondWithError(w, "Error creating user account", "general", http.StatusInternalServerError)
		return
	}

	_, err = DB.Exec(
		"INSERT INTO users (first_name, last_name, nickname, gender, age, email, password_hash, uuid) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		req.FirstName, req.LastName, req.Username, req.Gender, req.Age, req.Email, string(hashedPassword), userUUID,
	)
	if err != nil {
		log.Printf("Error adding user: %v", err)
		respondWithError(w, "Error creating user account", "general", http.StatusInternalServerError)
		return
	}

	log.Println("User registered successfully")
	json.NewEncoder(w).Encode(model.Response{
		Success: true,
		Message: "Registration successful",
	})
}

// LoginHandler handles user login requests
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, model.Response{
			Success: false,
			Message: "Method not allowed",
		})
		return
	}

	log.Println("Received login request")

	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, model.Response{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	var userID int
	var nickname, hashedPassword, userUUID string
	var query string

	if req.LoginType == "email" {
		query = "SELECT id, nickname, password_hash, uuid FROM users WHERE email = ?"
	} else {
		query = "SELECT id, nickname, password_hash, uuid FROM users WHERE nickname = ?"
	}

	err := DB.QueryRow(query, req.Identifier).Scan(&userID, &nickname, &hashedPassword, &userUUID)
	if err != nil {
		log.Printf("Login error: %v", err)
		respondWithError(w, "Invalid credentials", "login-general", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
	if err != nil {
		log.Printf("Password verification failed: %v", err)
		respondWithError(w, "Invalid credentials", "login-general", http.StatusUnauthorized)
		return
	}

	sessionToken, err := utils.GenerateUuid()
	if err != nil {
		log.Printf("Failed to generate session token: %v", err)
		respondWithError(w, "Error creating session", "login-general", http.StatusInternalServerError)
		return
	}
	expiresAt := time.Now().Add(24 * time.Hour)
	expiryStr := expiresAt.Format(time.RFC3339)

	// Begin transaction
	tx, err := DB.Begin()
	if err != nil {
		log.Printf("Failed to begin transaction: %v", err)
		respondWithError(w, "Error creating session", "login-general", http.StatusInternalServerError)
		return
	}

	// Delete existing sessions for this user
	_, delErr := tx.Exec("DELETE FROM sessions WHERE user_id = ?", userID)
	if delErr != nil {
		tx.Rollback()
		log.Printf("Failed to clear existing sessions: %v", delErr)
		respondWithError(w, "Login conflict", "login-general", http.StatusConflict)
		return
	}

	// Insert into sessions table
	_, err = tx.Exec(
		"INSERT INTO sessions (id, user_id, is_active, session_token, session_expiry) VALUES (?, ?, ?, ?, ?)",
		sessionToken, userID, true, sessionToken, expiryStr,
	)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to create session record: %v", err)
		respondWithError(w, "Error creating session", "login-general", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		respondWithError(w, "Error creating session", "login-general", http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Path:     "/",
		MaxAge:   86400, // 24h
		HttpOnly: true,
		Secure:   false, // Set to true if using HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	log.Println("User logged in successfully")
	json.NewEncoder(w).Encode(model.Response{
		Success: true,
		Message: "Login successful",
		Token:   sessionToken,
	})
	
	// After successful SELECT query:
	log.Printf(
		"Login success - UserID: %d, Nickname: %s, UUID: %s", 
		userID, 
		nickname, 
		userUUID,
	)

	// After session creation:
	log.Printf(
		"Session created - Token: %s → UserID: %d", 
		sessionToken, 
		userID,
	)
}

// Helper function to respond with error
func respondWithError(w http.ResponseWriter, message, field string, statusCode int) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(model.Response{
		Success: false,
		Message: message,
		Field:   field,
	})
}

type contextKey string

const userContextKey contextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Reuse ValidateSession logic here
		valid, user, _, _ := utils.ValidateSession(w, r)
		if !valid {
			respondJSON(w, http.StatusUnauthorized, map[string]any{
				"success": false,
				"message": "Authentication required",
			})
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
