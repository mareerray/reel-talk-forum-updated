package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/model"

	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword compares a hashed password with a plain text password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

func ValidateInputs(DB *sql.DB, nickname, email, password string) error {
	nickname = strings.TrimSpace(nickname)
	email = strings.TrimSpace(email)

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if nickname == "" || email == "" || password == "" {
		return ValidationError{Field: "general", Message: "all fields are required"}
	}

	if !emailRegex.MatchString(email) {
		return ValidationError{Field: "email", Message: "invalid email format"}
	}

	if len(nickname) < 5 || len(nickname) > 15 {
		return ValidationError{Field: "nickname", Message: "nickname must be between 5 and 30 characters long"}
	}

	if err := ValidatePassword(password); err != nil {
		return ValidationError{Field: "password", Message: err.Error()}
	}

	if !isValidNickname(nickname) {
		return ValidationError{Field: "nickname", Message: "nickname can only contain letters, numbers, underscores, and dashes"}
	}

	nicknameAvailable, err := NicknameNotTaken(DB, nickname)
	if err != nil {
		return fmt.Errorf("error checking nickname availability: %w", err)
	}
	if !nicknameAvailable {
		return ValidationError{Field: "nickname", Message: "nickname already taken"}
	}

	emailAvailable, err := EmailNotTaken(DB, email)
	if err != nil {
		return fmt.Errorf("error checking email availability: %w", err)
	}
	if !emailAvailable {
		return ValidationError{Field: "email", Message: "email already registered"}
	}

	return nil
}

func isValidNickname(nickname string) bool {
	for _, char := range nickname {
		if !(isValidCharacter(char) || char == '_' || char == '-') {
			return false
		}
	}
	return true
}

func isValidCharacter(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

func EmailNotTaken(DB *sql.DB, email string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking email: %w", err)
	}
	return count == 0, nil
}

func NicknameNotTaken(DB *sql.DB, nickname string) (bool, error) {

	if DB == nil {
		return false, fmt.Errorf("database connection is nil123")
	}
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM users WHERE nickname = ?", nickname).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking nickname: %w", err)
	}
	return count == 0, nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	lowercase := regexp.MustCompile(`[a-z]`)
	uppercase := regexp.MustCompile(`[A-Z]`)
	digit := regexp.MustCompile(`[0-9]`)
	specialChar := regexp.MustCompile(`[@$!%*?&]`)

	if !lowercase.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !uppercase.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !digit.MatchString(password) {
		return errors.New("password must contain at least one digit")
	}
	if !specialChar.MatchString(password) {
		return errors.New("password must contain at least one special character (@, $, !, %, *, ?, &)")
	}

	return nil
}

func ValidateSession(w http.ResponseWriter, r *http.Request) (bool, model.User, string, error) {
	var sessionToken string

	// 1. Check Authorization header first
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		sessionToken = strings.TrimPrefix(authHeader, "Bearer ")
	} else if cookie, err := r.Cookie("session_token"); err == nil { // 2. Fallback to cookie
		sessionToken = cookie.Value
	}
	log.Println("Validating session token:", sessionToken)

	if sessionToken == "" {
		return false, model.User{}, "", nil
	}
	user, expirationTime, selectError := SelectSession(sessionToken)
	if selectError != nil {
		if selectError.Error() == "sql: no rows in result set" {
			DeleteSession(sessionToken)
			return false, model.User{}, "", nil
		} else {
			return false, model.User{}, "", selectError
		}
	}
	log.Printf("[DEBUG] CurrentUserID: %d | UUID: %s", user.ID, user.UUID)


	// Check if the cookie has expired
	if time.Now().After(expirationTime) {
		DeleteSession(sessionToken)
		return false, model.User{}, "", nil
	}
	return true, user, sessionToken, nil
}

func SelectSession(sessionToken string) (model.User, time.Time, error) {

	log.Printf("Validating session token: %s", sessionToken)

	var user model.User
	var expirationTime time.Time

	// Adjusted SQL to match database schema
	// utils.go
	err := DB.QueryRow(`
	SELECT u.id, u.uuid, u.nickname, s.session_expiry 
	FROM sessions s
	INNER JOIN users u ON s.user_id = u.id
	WHERE s.session_token = ? 
	AND s.is_active = true
	AND s.session_expiry > ?`,
		sessionToken, time.Now().Format(time.RFC3339),
	).Scan(&user.ID, &user.UUID, &user.Nickname, &expirationTime)

	if err != nil {
		log.Printf("Database error: %v", err) // Add this line
		return model.User{}, time.Time{}, err
	}
	return user, expirationTime, nil
}

func DeleteSession(sessionToken string) error {

	db := OpenDBConnection()
	defer db.Close() // Close the connection after the function finishes
	_, err := db.Exec(`DELETE FROM sessions WHERE session_token = ?`, sessionToken)

	if err != nil {
		// Handle other database errors
		log.Fatal(err)
		return errors.New("database error")
	}

	return nil

}
